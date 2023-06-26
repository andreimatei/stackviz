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
	"io"
	"log"
	"math"
	"regexp"
	"stacksviz/ent"
	"stacksviz/ent/processsnapshot"
	"stacksviz/stacks"
	"strconv"
	"strings"
)

const (
	stacksTreeQuery = "stacks.tree"
	stacksRawQuery  = "stacks.raw"

	collectionIDKey          = "collection_id"
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

// ProcessSnapshot represents a single fetched log trace, along with any metadata it
// requires.
type ProcessSnapshot struct {
	processID        string
	Snapshot         *pp.Snapshot
	agg              *pp.Aggregated
	FramesOfInterest stacks.FOIS
}

// StacksFetcher describes types capable of fetching stack traces by collection
// name.
type StacksFetcher interface {
	// Fetch fetches the stacks specified by collectionName, returning a
	// LogTrace or an error if a failure is encountered.
	Fetch(ctx context.Context, collectionID int, snapshotID int) (ProcessSnapshot, error)
}

type stacksFetcherImpl struct {
	client *ent.Client
	// lru is a cache mapping from ProcessSnapshot ID to previously-loaded ProcessSnapshot.
	lru simplelru.LRUCache[int, ProcessSnapshot]
}

var _ StacksFetcher = &stacksFetcherImpl{}

// NewStacksFetcher creates a new StacksFetcher that will read collections from
// the specified directory.
func NewStacksFetcher(client *ent.Client) StacksFetcher {
	lru, err := lru.New[int, ProcessSnapshot](100)
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

func getSnapshotsForCollection(ctx context.Context, collectionID int, client *ent.Client) ([]*ent.ProcessSnapshot, error) {
	log.Printf("!!! getSnapshotsByCollection: id: %d", collectionID)
	c, err := client.Collection.Get(ctx, collectionID)
	if err != nil {
		return nil, err
	}
	return c.ProcessSnapshots(ctx)
}

func (f *stacksFetcherImpl) Fetch(ctx context.Context, collectionID int, snapshotID int) (ProcessSnapshot, error) {
	// Check the cache first.
	{
		snap, ok := f.lru.Get(snapshotID)
		if ok {
			return snap, nil
		}
	}

	snapRec, err := getSnapshot(ctx, snapshotID, f.client)
	if err != nil {
		return ProcessSnapshot{}, err
	}

	opts := pp.DefaultOpts()
	opts.ParsePC = true
	snap, _, err := pp.ScanSnapshot(strings.NewReader(snapRec.Snapshot), io.Discard, opts)
	if err != nil && err != io.EOF {
		return ProcessSnapshot{}, err
	}
	if snap == nil {
		return ProcessSnapshot{}, fmt.Errorf("failed to parse any stacks")
	}

	agg := snap.Aggregate(pp.AnyValue)

	var fois stacks.FramesOfInterest
	if snapRec.FramesOfInterest != "" {
		if err = json.Unmarshal([]byte(snapRec.FramesOfInterest), &fois); err != nil {
			log.Printf("!!! json: %s", snapRec.FramesOfInterest)
			return ProcessSnapshot{}, fmt.Errorf("failed to unmarshal frames of interest: %w", err)
		}
	}

	// Find links to other captured variables.
	allSnaps, err := getSnapshotsForCollection(ctx, collectionID, f.client)
	if err != nil {
		return ProcessSnapshot{}, err
	}
	processed := make(stacks.FOIS)
	for gid, m := range fois {
		prm := make(map[int]stacks.ProcessedFOI)
		for idx, vars := range m {
			var pf stacks.ProcessedFOI
			pf.Vars = make([]stacks.VarInfo, len(vars))
			for i, v := range vars {
				links := stacks.FindLinks(v, allSnaps)
				pf.Vars[i] = stacks.VarInfo{
					Val:   v,
					Links: links,
				}
			}
			prm[idx] = pf
		}
		processed[gid] = prm
	}

	// !!!
	//for i, foi := range snapRec.FramesOfInterest {
	//	if err = json.Unmarshal([]byte(foi), &fois[i]); err != nil {
	//		return processSnapshot{}, err
	//	}
	//}

	res := ProcessSnapshot{
		Snapshot:         snap,
		agg:              agg,
		FramesOfInterest: processed,
	}

	f.lru.Add(snapshotID, res)
	return res, nil
}

// nodeBuilder abstracts the differences between a weightedtree.Tree and a
// weightedtree.Node, allowing either to be used to construct a tree.
type nodeBuilder interface {
	Node(selfMagnitude float64, properties ...tvutil.PropertyUpdate) *weightedtree.Node
}

// toWeightedTree uses the provided builder to transforms a SubtreeNode (and
// its children, recursively) into a weighted tree.
func toWeightedTree(node *weightedtree.SubtreeNode, builder nodeBuilder, colorSpace *color.Space) {
	t := node.TreeNode.(*stacks.TreeNode)

	var varsProp []string
	var sb strings.Builder
	for _, frame := range t.Vars {
		sb.Reset()
		for _, v := range frame {
			sb.WriteString(v.Val)
			sb.WriteRune('\n')
		}
		varsProp = append(varsProp, sb.String())
	}

	n := builder.Node(float64(t.NumLeafGoroutines),
		tvutil.StringProperty(nameKey, t.Function.DirName+"."+t.Function.Name),
		weightedtree.Path(t),
		tvutil.StringProperty(fullNameKey, t.Function.Complete),
		tvutil.IntegerProperty(pcOffsetKey, t.PcOffset),
		tvutil.IntegerProperty(lineKey, int64(t.Line)),
		tvutil.StringProperty(fileKey, t.File),
		tvutil.StringsProperty(varsKey, varsProp...),
		tvutil.StringProperty(
			detailsFormatKey,
			fmt.Sprintf("$(%s) - $(%s) $(%s):$(%s)",
				fullNameKey, varsKey, fileKey, lineKey,
			)),
		colorSpace.PrimaryColor(functionNameToColor(t.Function.Complete)),
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
	log.Printf("!!! global filters: %v", globalFilters)
	snapshotIDVal, ok := globalFilters[snapshotIDKey]
	if !ok {
		return fmt.Errorf("missing required filter option '%s'", snapshotIDKey)
	}
	snapshotID, err := tvutil.ExpectIntegerValue(snapshotIDVal)
	if err != nil {
		return fmt.Errorf("required filter option '%s' must be an int", snapshotIDKey)
	}
	collectionIDVal, ok := globalFilters[collectionIDKey]
	if !ok {
		return fmt.Errorf("missing required filter option '%s'", collectionIDKey)
	}
	collectionID, err := tvutil.ExpectIntegerValue(collectionIDVal)
	if err != nil {
		return fmt.Errorf("required filter option '%s' must be an int", collectionIDKey)
	}

	processSnapshot, err := ds.fetchCollection(ctx, int(collectionID), int(snapshotID))
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

	snap := ds.filterStacks(processSnapshot.Snapshot, filter)
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
			if err := ds.handleStacksTreeQuery(snap, processSnapshot.FramesOfInterest, path, builder); err != nil {
				return err
			}
		case stacksRawQuery:
			ds.handleStacksRawQuery(snap, len(processSnapshot.Snapshot.Goroutines), processSnapshot.FramesOfInterest, builder)
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
	fois stacks.FOIS,
	path []weightedtree.ScopeID,
	builder tvutil.DataBuilder,
) error {
	tree := stacks.BuildTree(snap, fois)
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
	snap *pp.Snapshot, numTotalGoroutines int, fois stacks.FOIS, builder tvutil.DataBuilder,
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
		// !!!
		for _, foisByFrame := range fois[g.ID] {
			for _, v := range foisByFrame.Vars {
				var sb strings.Builder
				for i, l := range v.Links {
					if l.GoroutineID == g.ID {
						// This is the current goroutine; don't render a link to itself.
						continue
					}
					if i > 0 {
						sb.WriteString(", ")
					}
					sb.WriteString(fmt.Sprintf("[snap: %d, goroutine: %d (frame: %d)]", l.SnapshotID, l.GoroutineID, l.FrameIdx))
				}
				tab.Row(table.Cell(stackCol, tvutil.String(fmt.Sprintf("var: %s links: %s", v.Val, sb.String()))))
			}
		}

		for j := range g.Stack.Calls {
			c := &g.Stack.Calls[j]
			tab.Row(table.FormattedCell(stackCol, c.Func.Complete))
		}
	}
}

// fetchCollection returns the specified ProcessSnapshot from the LRU if it's
// present there.  If it isn't already in the LRU, it is fetched and added to
// the LRU before being returned.
func (ds *DataSource) fetchCollection(ctx context.Context, collectionID int, snapshotID int) (ProcessSnapshot, error) {
	col, err := ds.fetcher.Fetch(ctx, collectionID, snapshotID)
	if err != nil {
		return ProcessSnapshot{}, err
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
		if stacks.ComputeScopeID(c) != prefix[i] {
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
	aa := a.(*stacks.TreeNode)
	bb := b.(*stacks.TreeNode)
	res := -strings.Compare(aa.Function.Complete, bb.Function.Complete)
	return res, nil
}
