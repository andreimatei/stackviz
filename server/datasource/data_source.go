package datasource

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/google/traceviz/server/go/category"
	"github.com/google/traceviz/server/go/color"
	"github.com/google/traceviz/server/go/label"
	"github.com/google/traceviz/server/go/table"
	tvutil "github.com/google/traceviz/server/go/util"
	weightedtree "github.com/google/traceviz/server/go/weighted_tree"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/hashicorp/golang-lru/v2/simplelru"
	pp "github.com/maruel/panicparse/v2/stack"
	"hash/fnv"
	"io"
	"log"
	"math"
	"regexp"
	"stacksviz/ent"
	"stacksviz/ent/processsnapshot"
	"strconv"
	"strings"
)

const (
	stacksTreeQuery = "stacks.tree"
	stacksRawQuery  = "stacks.raw"

	snapshotIDKey            = "snapshot_id"
	pathPrefixKey            = "path_prefix"
	nameKey                  = "name"
	detailsFormatKey         = "detail_format"
	fullNameKey              = "full_name"
	filterKey                = "filter"
	numTotalGoroutinesKey    = "num_total_goroutines"
	numFilteredGoroutinesKey = "num_filtered_goroutines"
	numBucketsKey            = "num_buckets"

	numGoroutinesInBucketKey = "num_gs_in_bucket"
	goroutineIDKey           = "g_id"
	varsKey                  = "vars"
	pcOffsetKey              = "pc_off"
	fileKey                  = "file"
	lineKey                  = "line"
)

// DataSource implements the querydispatcher.dataSource that deals with
// goroutine stacks.
type DataSource struct {
	fetcher StacksFetcher
}

// New builds a DataSource.
func New(fetcher StacksFetcher) *DataSource {
	return &DataSource{fetcher: fetcher}
}

type FramesOfInterest map[int]map[int][]string

// processSnapshot represents a single fetched log trace, along with any metadata it
// requires.
type processSnapshot struct {
	processID        string
	snapshot         *pp.Snapshot
	agg              *pp.Aggregated
	framesOfInterest FramesOfInterest
}

// StacksFetcher describes types capable of fetching stack traces by collection
// name.
type StacksFetcher interface {
	// Fetch fetches the stacks specified by collectionName, returning a
	// LogTrace or an error if a failure is encountered.
	Fetch(ctx context.Context, snapshotID int) (processSnapshot, error)
}

type stacksFetcherImpl struct {
	client *ent.Client
	// lru is a cache mapping from processSnapshot ID to previously-loaded processSnapshot.
	lru simplelru.LRUCache[int, processSnapshot]
}

var _ StacksFetcher = &stacksFetcherImpl{}

// NewStacksFetcher creates a new StacksFetcher that will read collections from
// the specified directory.
func NewStacksFetcher(client *ent.Client) StacksFetcher {
	lru, err := lru.New[int, processSnapshot](100)
	if err != nil {
		panic(err)
	}
	return &stacksFetcherImpl{
		client: client,
		lru:    lru,
	}
}

func getSnapshot(ctx context.Context, id int, client *ent.Client) (*ent.ProcessSnapshot, error) {
	log.Printf("!!! getSnapshot: id: %d", id)
	c, err := client.ProcessSnapshot.
		Query().
		Where(processsnapshot.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying snapshot: %w", err)
	}
	return c, nil
}

func (f *stacksFetcherImpl) Fetch(ctx context.Context, snapshotID int) (processSnapshot, error) {
	// Check the cache first.
	{
		snap, ok := f.lru.Get(snapshotID)
		if ok {
			return snap, nil
		}
	}

	snapRec, err := getSnapshot(ctx, snapshotID, f.client)
	if err != nil {
		return processSnapshot{}, err
	}

	opts := pp.DefaultOpts()
	opts.ParsePC = true
	snap, _, err := pp.ScanSnapshot(strings.NewReader(snapRec.Snapshot), io.Discard, opts)
	if err != nil && err != io.EOF {
		return processSnapshot{}, err
	}
	if snap == nil {
		return processSnapshot{}, fmt.Errorf("failed to parse any stacks")
	}

	agg := snap.Aggregate(pp.AnyValue)

	var fois FramesOfInterest
	if snapRec.FramesOfInterest != "" {
		if err = json.Unmarshal([]byte(snapRec.FramesOfInterest), &fois); err != nil {
			log.Printf("!!! json: %s", snapRec.FramesOfInterest)
			return processSnapshot{}, fmt.Errorf("failed to unmarshal frames of interest: %w", err)
		}
	}

	// !!!
	//for i, foi := range snapRec.FramesOfInterest {
	//	if err = json.Unmarshal([]byte(foi), &fois[i]); err != nil {
	//		return processSnapshot{}, err
	//	}
	//}

	res := processSnapshot{
		snapshot:         snap,
		agg:              agg,
		framesOfInterest: fois,
	}

	f.lru.Add(snapshotID, res)
	return res, nil
}

// treeNode is a node in a trie of stack traces. Each node represents a
// function; its children are other functions called by the node's function in
// one or more stacks.
type treeNode struct {
	function pp.Func
	file     string
	line     int
	pcOffset int64
	vars     [][]string
	// path is the path from the root to this node, represented by hashes of
	// each ancestor's function.
	path     []weightedtree.ScopeID
	children []treeNode
	// numLeafGoroutines counts how many goroutines have this node as their leaf
	// function. This results in the "self magnitude" of the node when rendered
	// as a flame graph - i.e. how much weight it needs to have in addition to
	// the sum of the children's weights.
	numLeafGoroutines int
	numGoroutines     int
}

// scopeID returns the identifier for this node.
func (t *treeNode) scopeID() weightedtree.ScopeID {
	if len(t.path) > 0 {
		return t.path[len(t.path)-1]
	}
	return 0
}

var _ weightedtree.TreeNode = &treeNode{}

// Path is part of the weightedtree.TreeNode interface.
func (t *treeNode) Path() []weightedtree.ScopeID {
	return t.path
}

func (t *treeNode) pathAsStrings() []string {
	path := make([]string, len(t.Path()))
	for i, p := range t.Path() {
		path[i] = strconv.FormatUint(uint64(p), 10)
	}
	return path
}

// Children is part of the weightedtree.TreeNode interface.
func (t *treeNode) Children(ids ...weightedtree.ScopeID) ([]weightedtree.TreeNode, error) {
	res := make([]weightedtree.TreeNode, 0, len(ids))
	for i := range t.children {
		c := &t.children[i]
		add := false
		if len(ids) > 0 {
			for _, id := range ids {
				if c.scopeID() == id {
					add = true
				}
			}
		} else {
			add = true
		}
		if add {
			res = append(res, c)
		}
	}
	return res, nil
}

func (t *treeNode) prettyPrint() {
	t.prettyPrintInner(0)
}

func (t *treeNode) prettyPrintInner(indent int) {
	var sb strings.Builder
	for i := 0; i < indent; i++ {
		sb.WriteRune('\t')
	}
	fmt.Printf("%s(%d) %s (%s:%d) (%v)\n", sb.String(), t.numLeafGoroutines, t.function.Complete, t.file, t.line, t.path)
	for i := range t.children {
		t.children[i].prettyPrintInner(indent + 1)
	}
}

// findChild finds the child of t for a call at file:line. If such a child
// doesn't exist, returns nil.
func (t *treeNode) findChild(file string, line int) *treeNode {
	for i := range t.children {
		c := &t.children[i]
		if c.file == file && c.line == line {
			return c
		}
	}
	return nil
}

// addStack adds the stack to the tree rooted at t, creating new nodes for calls
// that don't yet exist.
func (t *treeNode) addStack(stack []frame) {
	t.numGoroutines++
	if len(stack) == 0 {
		// t is a leaf for the stack that we just finished processing.
		t.numLeafGoroutines++
		return
	}
	child := t.findChild(stack[0].call.RemoteSrcPath, stack[0].call.Line)
	if child != nil {
		if len(stack[0].vars) > 0 {
			child.vars = append(child.vars, stack[0].vars)
		}
		child.addStack(stack[1:])
	} else {
		t.createPath(stack)
	}
}

// createPath adds children to t recursively such that the tree gets the path
// t -> stack[0] -> stack[1] -> ...
func (t *treeNode) createPath(stack []frame) {
	t.numGoroutines++
	if len(stack) == 0 {
		// The stack had t as a leaf function.
		t.numLeafGoroutines++
		return
	}
	call := &stack[0].call
	t.children = append(t.children, treeNode{
		function:          call.Func,
		file:              call.RemoteSrcPath,
		line:              call.Line,
		pcOffset:          call.PCOffset,
		path:              append(t.path, computeScopeID(call)),
		children:          nil,
		numLeafGoroutines: 0,
		vars:              [][]string{stack[0].vars},
	})
	t.children[len(t.children)-1].createPath(stack[1:])
}

func computeScopeID(call *pp.Call) weightedtree.ScopeID {
	return computeScopeIDInner(call.Func.Complete, call.RemoteSrcPath, uint32(call.Line))
}

func computeScopeIDInner(funcName string, file string, line uint32) weightedtree.ScopeID {
	hash := fnv.New64()
	hash.Write([]byte(funcName))
	hash.Write([]byte(file))
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, line)
	hash.Write(bs)
	return weightedtree.ScopeID(hash.Sum64())
}

// nodeBuilder abstracts the differences between a weightedtree.Tree and a
// weightedtree.Node, allowing either to be used to construct a tree.
type nodeBuilder interface {
	Node(selfMagnitude float64, properties ...tvutil.PropertyUpdate) *weightedtree.Node
}

// toWeightedTree uses the provided builder to transforms a SubtreeNode (and
// its children, recursively) into a weighted tree.
func toWeightedTree(node *weightedtree.SubtreeNode, builder nodeBuilder, colorSpace *color.Space) {
	t := node.TreeNode.(*treeNode)

	var varsProp []string
	var sb strings.Builder
	for _, frame := range t.vars {
		sb.Reset()
		for _, v := range frame {
			sb.WriteString(v)
			sb.WriteRune('\n')
		}
		varsProp = append(varsProp, sb.String())
	}

	n := builder.Node(float64(t.numLeafGoroutines),
		tvutil.StringProperty(nameKey, t.function.DirName+"."+t.function.Name),
		weightedtree.Path(t),
		tvutil.StringProperty(fullNameKey, t.function.Complete),
		tvutil.IntegerProperty(pcOffsetKey, t.pcOffset),
		tvutil.IntegerProperty(lineKey, int64(t.line)),
		tvutil.StringProperty(fileKey, t.file),
		tvutil.StringsProperty(varsKey, varsProp...),
		tvutil.StringProperty(
			detailsFormatKey,
			fmt.Sprintf("$(%s) - $(%s) $(%s):$(%s)",
				fullNameKey, varsKey, fileKey, lineKey,
			)),
		colorSpace.PrimaryColor(functionNameToColor(t.function.Complete)),
		label.Format(fmt.Sprintf("$(%s)", nameKey)))
	for _, c := range node.Children {
		toWeightedTree(c, n, colorSpace)
	}
}

// functionNameToColor takes in a function name (including the package) and
// returns a float between [0,1] signifying the color that should be used to
// fill boxes corresponding to this function. The result is supposed to be used
// to index into a color space.
// The same color will be returned for all functions in a package (if the
// package name can be identified). This matches pprof behavior.
func functionNameToColor(functionName string) float64 {
	// pkgRE extracts package name, It looks for the first "." or "::" that
	// occurs after the last "/". (Searching after the last / allows us to
	// correctly handle names that look like "some.url.com/foo.bar".)
	pkgRE := regexp.MustCompile(`^((.*/)?[\w\d_]+)(\.|::)([^/]*)$`)
	var pkg string
	m := pkgRE.FindStringSubmatch(functionName)
	if m == nil {
		pkg = functionName
	} else {
		pkg = m[1]
	}
	h := sha256.Sum256([]byte(pkg))
	hash := binary.LittleEndian.Uint32(h[:])
	return float64(hash) / math.MaxUint32
}

type frame struct {
	call pp.Call
	vars []string
}

// buildTree builds a trie out of the stack traces in snap.
func (ds *DataSource) buildTree(snap *pp.Snapshot, fois FramesOfInterest) *treeNode {
	root := &treeNode{
		function: pp.Func{
			Complete: "root",
			Name:     "root",
		},
		// The root doesn't have a path, as per weightedtree.TreeNode
		// convention.
		path: nil,
	}
	for _, s := range snap.Goroutines {
		// Join the stack trace with the variable data. Also invert the stack; we
		// want it ordered from top-level function to leaf function.
		l := len(s.Signature.Stack.Calls)
		myFois := fois[s.ID]
		stack := make([]frame, l)
		for i := range s.Signature.Stack.Calls {
			stack[l-i-1] = frame{
				call: s.Signature.Stack.Calls[i],
				vars: myFois[i],
			}
		}
		root.addStack(stack)
	}
	return root
}

// SupportedDataSeriesQueries implements the traceviz datasource interface. It
// returns the supported query names DataSeriesRequest.
func (ds *DataSource) SupportedDataSeriesQueries() []string {
	return []string{stacksRawQuery, stacksTreeQuery}
}

// HandleDataSeriesRequests handles the provided set of DataSeriesRequests, with
// the provided global filters.  It assembles its responses in the provided
// DataResponseBuilder.
func (ds *DataSource) HandleDataSeriesRequests(
	ctx context.Context,
	globalFilters map[string]*tvutil.V,
	drb *tvutil.DataResponseBuilder,
	reqs []*tvutil.DataSeriesRequest,
) error {
	snapshotIDVal, ok := globalFilters[snapshotIDKey]
	if !ok {
		return fmt.Errorf("missing required filter option '%s'", snapshotIDKey)
	}
	snapshotID, err := tvutil.ExpectIntegerValue(snapshotIDVal)
	if err != nil {
		return fmt.Errorf("required filter option '%s' must be an int", snapshotIDKey)
	}

	processSnapshot, err := ds.fetchCollection(ctx, int(snapshotID))
	if err != nil {
		log.Printf("Failed to fetch collection: %s", err)
		return err
	}
	log.Printf("Loaded snapshot: id: %d, processID: %q", snapshotID, processSnapshot.processID)

	var filter string
	filterVal, ok := globalFilters[filterKey]
	if ok {
		filter, err = tvutil.ExpectStringValue(filterVal)
		if err != nil {
			return fmt.Errorf("filter '%s' must be a string list", filter)
		}
	}
	pathPrefixVal, ok := globalFilters[pathPrefixKey]
	var pathPrefix []string
	if ok {
		pathPrefix, err = tvutil.ExpectStringsValue(pathPrefixVal)
		if err != nil {
			return fmt.Errorf("filter '%s' must be a string list", pathPrefix)
		}
	}
	log.Printf("!!! query (%d requests) filter: %s, path prefix: %s\n",
		len(reqs), filter, pathPrefix)

	snap := ds.filterStacks(processSnapshot.snapshot, filter)
	path := make([]weightedtree.ScopeID, len(pathPrefix))
	for i, p := range pathPrefix {
		sid, err := strconv.ParseUint(p, 10, 64)
		if err != nil {
			return err
		}
		path[i] = weightedtree.ScopeID(sid)
	}
	snap = ds.filterStacksByPrefix(snap, path)
	log.Printf("!!! goroutines after prefix filter: %d", len(snap.Goroutines))

	for _, req := range reqs {
		builder := drb.DataSeries(req)
		switch req.QueryName {
		case stacksTreeQuery:
			if err := ds.handleStacksTreeQuery(snap, processSnapshot.framesOfInterest, path, builder); err != nil {
				return err
			}
		case stacksRawQuery:
			ds.handleStacksRawQuery(snap, len(processSnapshot.snapshot.Goroutines), processSnapshot.framesOfInterest, builder)
		default:
			return fmt.Errorf("unsupported data query: %s", req.QueryName)
		}
	}
	return nil
}

// filterStacks returns a new Snapshot containing the goroutines in snap that
// contain at least a frame that matches filter.
func (ds *DataSource) filterStacks(snap *pp.Snapshot, filter string) *pp.Snapshot {
	if filter == "" {
		return snap
	}
	res := new(pp.Snapshot)
	*res = *snap // shallow copy
	res.Goroutines = nil
	for _, g := range snap.Goroutines {
		if ds.stackMatchesFilter(g, filter) {
			res.Goroutines = append(res.Goroutines, g)
		}
	}
	return res
}

// filterStacksByPrefix returns a new Snapshot containing the goroutines in snap
// that have the given prefix.
func (ds *DataSource) filterStacksByPrefix(snap *pp.Snapshot, prefix []weightedtree.ScopeID) *pp.Snapshot {
	if len(prefix) == 0 {
		return snap
	}
	res := new(pp.Snapshot)
	*res = *snap // shallow copy
	res.Goroutines = nil
	for _, g := range snap.Goroutines {
		if ds.stackMatchesPrefix(g, prefix) {
			res.Goroutines = append(res.Goroutines, g)
		}
	}
	return res
}

// handleStacksTreeQuery uses the provided builder to construct the response to
// the "stacks tree" query, turning a snapshot into a weighted tree. This
// function will filter the snapshot by the provided prefix path. Even if the
// snapshot is already filtered, the path should still be specified (if there is
// one), such that the prefix nodes in the resulting tree are correctly marked
// as such.
func (ds *DataSource) handleStacksTreeQuery(
	snap *pp.Snapshot,
	fois FramesOfInterest,
	path []weightedtree.ScopeID,
	builder tvutil.DataBuilder,
) error {
	tree := ds.buildTree(snap, fois)
	filtered, err := weightedtree.Walk(
		tree,
		compareByFunctionName, // within a level, sort alphabetically
		weightedtree.PathPrefix(path...),
	)
	if err != nil {
		return err
	}

	renderSettings := &weightedtree.RenderSettings{
		FrameHeightPx: 20,
	}
	wt := weightedtree.New(builder, renderSettings)
	// Include the color space.
	// This represents the color range yellow-ish (a bit towards orange) to
	// red-ish (a bit towards orange).
	// TODO(andrei): Try expressing this in the HSL color scheme in the hope
	// that a human can read it more easily.
	var colorSpace = color.NewSpace("yellow-to-red", "#E1A81E", "#E35A1C")
	wt.With(colorSpace.Define())
	toWeightedTree(filtered, wt, colorSpace)
	return nil
}

// Columns for the raw stacks.
var stackCol = table.Column(category.New("stack", "Backtrace", "The goroutine's stack backtrace"))

// Columns for the aggregated views.
var pkgCol = table.Column(category.New("package", "Package", "The name of the package that the function lives in."))
var fileLineCol = table.Column(category.New("fileLine", "file:line", "The source location."))
var funcCol = table.Column(category.New("function", "Function", "The name of the function."))
var pcOffsetCol = table.Column(category.New("pcoff", "PC offset", "instruction offset from function entry"))

func (ds *DataSource) handleStacksRawQuery(
	snap *pp.Snapshot, numTotalGoroutines int, fois FramesOfInterest, builder tvutil.DataBuilder,
) {
	renderSettings := &table.RenderSettings{
		RowHeightPx: 20,
		FontSizePx:  14,
	}
	builder.With(tvutil.IntegerProperty(numTotalGoroutinesKey, int64(numTotalGoroutines)))
	builder.With(tvutil.IntegerProperty(numFilteredGoroutinesKey, int64(numTotalGoroutines-len(snap.Goroutines))))
	agg := snap.Aggregate(pp.AnyValue)

	aggBuilder := builder.Child().With(tvutil.IntegerProperty(numBucketsKey, int64(len(agg.Buckets))))
	rawBuilder := builder.Child()

	for _, b := range agg.Buckets {
		tab := table.New(aggBuilder.Child(), renderSettings, pkgCol, fileLineCol, funcCol, pcOffsetCol).With(tvutil.IntegerProperty(numGoroutinesInBucketKey, int64(len(b.IDs))))
		for j := range b.Stack.Calls {
			c := &b.Stack.Calls[j]
			tab.Row(
				table.FormattedCell(pkgCol, c.Func.DirName),
				table.FormattedCell(fileLineCol, fmt.Sprintf("%s:%d", c.SrcName, c.Line)),
				table.FormattedCell(funcCol, c.Func.Name),
				table.FormattedCell(pcOffsetCol, fmt.Sprintf("0x%x", c.PCOffset)),
			)
		}
	}

	for _, g := range snap.Goroutines {
		tab := table.New(rawBuilder.Child(), renderSettings, stackCol).With(tvutil.IntegerProperty(goroutineIDKey, int64(g.ID)))
		// Render all the frames of interest for the goroutine, across all the frames.
		for _, foisByFrame := range fois[g.ID] {
			for _, foi := range foisByFrame {
				tab.Row(table.Cell(stackCol, tvutil.String(foi)))
			}
		}

		for j := range g.Stack.Calls {
			c := &g.Stack.Calls[j]
			tab.Row(table.FormattedCell(stackCol, c.Func.Complete))
		}
	}
}

// fetchCollection returns the specified processSnapshot from the LRU if it's
// present there.  If it isn't already in the LRU, it is fetched and added to
// the LRU before being returned.
func (ds *DataSource) fetchCollection(ctx context.Context, snapshotID int) (processSnapshot, error) {
	col, err := ds.fetcher.Fetch(ctx, snapshotID)
	if err != nil {
		return processSnapshot{}, err
	}
	return col, nil
}

func (ds *DataSource) stackMatchesFilter(g *pp.Goroutine, filter string) bool {
	for i := range g.Stack.Calls {
		if strings.Contains(g.Stack.Calls[i].Func.Complete, filter) {
			return true
		}
	}
	return false
}

func (ds *DataSource) stackMatchesPrefix(g *pp.Goroutine, prefix []weightedtree.ScopeID) bool {
	if len(g.Stack.Calls) < len(prefix) {
		return false
	}
	for i := range prefix {
		c := &g.Stack.Calls[len(g.Stack.Calls)-i-1]
		if computeScopeID(c) != prefix[i] {
			return false
		}
	}
	return true
}

// compareByFunctionName compares the function names of two nodes
// lexicographically, returning 1 if a < b and -1 if b < a. This corresponds to
// a descending sorting; this function is intended to be used with
// weightedtree.Walk(), which explores "higher" nodes first so, in order to get
// alphabetic sorting, we invert the regular comparison.
func compareByFunctionName(a, b weightedtree.TreeNode) (int, error) {
	aa := a.(*treeNode)
	bb := b.(*treeNode)
	res := -strings.Compare(aa.function.Complete, bb.function.Complete)
	return res, nil
}
