package datasource

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/google/traceviz/server/go/label"
	tvutil "github.com/google/traceviz/server/go/util"
	weightedtree "github.com/google/traceviz/server/go/weighted_tree"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/hashicorp/golang-lru/v2/simplelru"
	pp "github.com/maruel/panicparse/v2/stack"
	"hash/fnv"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	rawEntriesQuery = "stacks.raw_entries"
	stacksTreeQuery = "stacks.tree"

	collectionNameKey = "collection_name"
	pathPrefixKey     = "path_prefix"
	nameKey           = "name"
	pathKey           = "path"
)

type DataSource struct {
	fetcher StacksFetcher
}

func NewDataSource(fetcher StacksFetcher) *DataSource {
	return &DataSource{fetcher: fetcher}
}

// collection represents a single fetched log trace, along with any metadata it
// requires.
type collection struct {
	snapshot *pp.Snapshot
}

// StacksFetcher describes types capable of fetching stack traces by collection
// name.
type StacksFetcher interface {
	// Fetch fetches the stacks specified by collectionName, returning a
	// LogTrace or an error if a failure is encountered.
	Fetch(ctx context.Context, collectionName string) (collection, error)
}

type stacksFetcherImpl struct {
	// rootDir is the directory from which files containing stack traces are read.
	rootDir string
	// lru is a cache mapping from collection name to previously-loaded collection.
	lru simplelru.LRUCache[string, collection]
}

var _ StacksFetcher = &stacksFetcherImpl{}

func NewStacksFetcher(dir string) StacksFetcher {
	lru, err := lru.New[string, collection](100)
	if err != nil {
		panic(err)
	}
	return &stacksFetcherImpl{
		rootDir: dir,
		lru:     lru,
	}
}

func (f *stacksFetcherImpl) Fetch(ctx context.Context, collectionName string) (collection, error) {
	// Check the cache first.
	col, ok := f.lru.Get(collectionName)
	if ok {
		return col, nil
	}

	// Read the stacks from the file.
	file, err := os.Open(path.Join(f.rootDir, collectionName))
	if err != nil {
		return collection{}, err
	}
	defer file.Close()
	col, err = f.readStacks(file)
	if err != nil {
		return collection{}, err
	}

	f.lru.Add(collectionName, col)
	return col, nil
}

func (f *stacksFetcherImpl) readStacks(r io.Reader) (collection, error) {
	snap, _, err := pp.ScanSnapshot(r, io.Discard, pp.DefaultOpts())
	if err != nil && err != io.EOF {
		return collection{}, err
	}
	if snap == nil {
		return collection{}, fmt.Errorf("failed to parse any stacks")
	}
	fmt.Printf("!!! loaded %d stacks\n", len(snap.Goroutines))
	return collection{snapshot: snap}, nil
}

// treeNode is a node in a trie of stack traces. Each node represents a
// function; its children are other functions called by the node's function in
// one or more stacks.
type treeNode struct {
	// functionName identifies that function represented by this node.
	functionName string
	file         string
	line         int
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
	fmt.Printf("%s(%d) %s (%s:%d) (%v)\n", sb.String(), t.numLeafGoroutines, t.functionName, t.file, t.line, t.path)
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
func (t *treeNode) addStack(stack []pp.Call) {
	t.numGoroutines++
	if len(stack) == 0 {
		// t is a leaf for the stack that we just finished processing.
		t.numLeafGoroutines++
		return
	}
	child := t.findChild(stack[0].RemoteSrcPath, stack[0].Line)
	if child != nil {
		child.addStack(stack[1:])
	} else {
		t.createPath(stack)
	}
}

// createPath adds children to t recursively such that the tree gets the path
// t -> stack[0] -> stack[1] -> ...
func (t *treeNode) createPath(stack []pp.Call) {
	t.numGoroutines++
	if len(stack) == 0 {
		// The stack had t as a leaf function.
		t.numLeafGoroutines++
		return
	}
	call := stack[0]
	t.children = append(t.children, treeNode{
		functionName:      call.Func.Complete,
		file:              call.RemoteSrcPath,
		line:              call.Line,
		path:              append(t.path, computeScopeID(call.Func.Complete, call.RemoteSrcPath, uint32(call.Line))),
		children:          nil,
		numLeafGoroutines: 0,
	})
	t.children[len(t.children)-1].createPath(stack[1:])
}

func computeScopeID(funcName string, file string, line uint32) weightedtree.ScopeID {
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

// toWeightedTree uses the provided dataBuilder to transforms a treeNode (and
// its children, recursively) into a weighted tree.
func (t *treeNode) toWeightedTree(dataBuilder tvutil.DataBuilder) {
	renderSettings := &weightedtree.RenderSettings{
		FrameHeightPx: 20,
	}
	t.toWeightedTreeInner(weightedtree.New(dataBuilder, renderSettings))
}

// toWeightedTreeInner creates a node in wt corresponding to t, and recurses in
// t's children.
func (t *treeNode) toWeightedTreeInner(builder nodeBuilder) {
	node := builder.Node(float64(t.numLeafGoroutines),
		tvutil.StringProperty(nameKey, t.functionName),
		tvutil.StringsProperty(pathKey, t.pathAsStrings()...),
		// !!! use StringsProperty instead?
		// tvutil.StringProperty(pathKey, pathBuilder.String()),
		label.Format(t.functionName),
	)
	for i := range t.children {
		t.children[i].toWeightedTreeInner(node)
	}
}

// toWeightedTree uses the provided dataBuilder to transforms a treeNode (and
// its children, recursively) into a weighted tree.
func toWeightedTree(view *weightedtree.SubtreeNode, dataBuilder tvutil.DataBuilder) {
	renderSettings := &weightedtree.RenderSettings{
		FrameHeightPx: 20,
	}
	wt := weightedtree.New(dataBuilder, renderSettings)
	t := view.TreeNode.(*treeNode)
	root := wt.Node(float64(t.numLeafGoroutines),
		tvutil.StringProperty(nameKey, "root"),
		tvutil.StringsProperty(pathKey, t.pathAsStrings()...),
		label.Format(t.functionName))
	toWeightedTreeInner(root, view)
}

func toWeightedTreeInner(builderNode *weightedtree.Node, view *weightedtree.SubtreeNode) {
	for _, c := range view.Children {
		t := c.TreeNode.(*treeNode)
		n := builderNode.Node(float64(t.numLeafGoroutines),
			tvutil.StringProperty(nameKey, t.functionName),
			tvutil.StringsProperty(pathKey, t.pathAsStrings()...),
			label.Format(t.functionName))
		toWeightedTreeInner(n, c)
	}
}

// buildTree builds a trie out of the stack traces in snap.
func (ds *DataSource) buildTree(snap *pp.Snapshot) *treeNode {
	root := &treeNode{
		functionName: "root",
		// The root doesn't have a path, as per weightedtree.TreeNode
		// convention.
		path: nil,
	}
	for _, s := range snap.Goroutines {
		// Invert the stack; we want it ordered from top-level function to leaf
		// function.
		l := len(s.Signature.Stack.Calls)
		stack := make([]pp.Call, l)
		for i := range s.Signature.Stack.Calls {
			stack[l-i-1] = s.Signature.Stack.Calls[i]
		}
		root.addStack(stack)
	}
	return root
}

// SupportedDataSeriesQueries implements the traceviz datasource interface. It
// returns the supported query names DataSeriesRequest.
func (ds *DataSource) SupportedDataSeriesQueries() []string {
	return []string{rawEntriesQuery, stacksTreeQuery}
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
	// Pull the collection name from the global filters.
	collectionNameVal, ok := globalFilters[collectionNameKey]
	if !ok {
		return fmt.Errorf("missing required filter option '%s'", collectionNameKey)
	}
	collectionName, err := tvutil.ExpectStringValue(collectionNameVal)
	if err != nil {
		return fmt.Errorf("required filter option '%s' must be a string", collectionNameKey)
	}

	pathPrefixVal, ok := globalFilters[pathPrefixKey]
	var pathPrefix []string
	if ok {
		pathPrefix, err = tvutil.ExpectStringsValue(pathPrefixVal)
		if err != nil {
			return fmt.Errorf("required filter option '%s' must be a string list", pathPrefix)
		}
		fmt.Printf("!!! path prefix: %s\n", pathPrefix)
	}

	// Fetch the collection, from the cache if it's there.
	col, err := ds.fetchCollection(ctx, collectionName)
	if err != nil {
		log.Printf("Failed to fetch collection: %s", err)
		return err
	}
	log.Printf("Loaded collection %s", collectionName)
	// !!!
	//// Build the queryFilters, just once, for all DataSeriesRequests.
	//qf, err := filterFromGlobalFilters(coll.lt, globalFilters)
	//if err != nil {
	//	return err
	//}

	for _, req := range reqs {
		series := drb.DataSeries(req)
		var err error
		switch req.QueryName {
		case rawEntriesQuery:
			err = ds.handleRawEntriesQuery(col, nil /* !!! filters */, series, req.Options)
		// !!!
		case stacksTreeQuery:
			builder := drb.DataSeries(req)
			tree := ds.buildTree(col.snapshot)
			if pathPrefix != nil {
				path := make([]weightedtree.ScopeID, len(pathPrefix))
				for i, p := range pathPrefix {
					sid, err := strconv.ParseUint(p, 10, 64)
					if err != nil {
						return err
					}
					path[i] = weightedtree.ScopeID(sid)
				}
				ds.filterByPathPrefix(tree, path, builder)
			} else {
				tree.toWeightedTree(builder)
			}
			// !!! err = handleStacksTreeQuery(coll, qf, series, req.Options)
			return nil
		default:
			err = fmt.Errorf("unsupported data query")
		}
		if err != nil {
			return fmt.Errorf("error handling data query %s: %w", req.QueryName, err)
		}
	}
	return nil
}

// fetchCollection returns the specified collection from the LRU if it's
// present there.  If it isn't already in the LRU, it is fetched and added to
// the LRU before being returned.
func (ds *DataSource) fetchCollection(ctx context.Context, collectionName string) (collection, error) {
	col, err := ds.fetcher.Fetch(ctx, collectionName)
	if err != nil {
		return collection{}, err
	}
	return col, nil
}

func compareByNumGoroutines(a, b weightedtree.TreeNode) (int, error) {
	aa := a.(*treeNode)
	bb := b.(*treeNode)
	if aa.numGoroutines < bb.numGoroutines {
		return -1, nil
	}
	if aa.numGoroutines == bb.numGoroutines {
		return 0, nil
	}
	return 1, nil
}

func (ds *DataSource) filterByPathPrefix(tree *treeNode, path []weightedtree.ScopeID, builder tvutil.DataBuilder) {
	filtered, err := weightedtree.Walk(tree, compareByNumGoroutines, weightedtree.PathPrefix(path...))
	if err != nil {
		panic(err)
	}
	toWeightedTree(filtered, builder)
}

func (ds *DataSource) handleRawEntriesQuery(col collection, qf interface{}, series tvutil.DataBuilder, options map[string]*tvutil.V) error {
	return nil
}
