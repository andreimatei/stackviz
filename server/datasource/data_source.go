package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/hashicorp/golang-lru/v2/simplelru"
	pp "github.com/maruel/panicparse/v2/stack"
	"io"
	"log"
	"stacksviz/ent"
	"stacksviz/ent/processsnapshot"
	"stacksviz/graph"
	"stacksviz/stacks"
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
	Agg              *pp.Aggregated
	FramesOfInterest stacks.FOIS
}

// StacksFetcher describes types capable of fetching stack traces by collection
// name.
type StacksFetcher interface {
	// Fetch fetches the stacks specified by collectionName, returning a
	// LogTrace or an error if a failure is encountered.
	Fetch(ctx context.Context, collectionID int, snapshotID int) (*pp.Snapshot, stacks.FOIS, error)
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

func (f *stacksFetcherImpl) Fetch(ctx context.Context, collectionID int, snapshotID int) (*pp.Snapshot, stacks.FOIS, error) {
	//// Check the cache first.
	//{
	//	snap, ok := f.lru.Get(snapshotID)
	//	if ok {
	//		return snap, nil
	//	}
	//}

	snapRec, err := getSnapshot(ctx, snapshotID, f.client)
	if err != nil {
		return nil, nil, err
	}

	opts := pp.DefaultOpts()
	opts.ParsePC = true
	snap, _, err := pp.ScanSnapshot(strings.NewReader(snapRec.Snapshot), io.Discard, opts)
	if err != nil && err != io.EOF {
		return nil, nil, err
	}
	if snap == nil {
		return nil, nil, fmt.Errorf("failed to parse any stacks")
	}

	var rawFois stacks.FramesOfInterest
	if snapRec.FramesOfInterest != "" {
		if err = json.Unmarshal([]byte(snapRec.FramesOfInterest), &rawFois); err != nil {
			log.Printf("!!! failed to unmarshal: %s", snapRec.FramesOfInterest)
			return nil, nil, fmt.Errorf("failed to unmarshal frames of interest: %w", err)
		}
	} else {
		log.Printf("!!! snapshot %d does not have any frames of interest", snapshotID)
	}

	// Find links to other captured variables.
	allSnaps, err := getSnapshotsForCollection(ctx, collectionID, f.client)
	if err != nil {
		return nil, nil, err
	}

	varToLinks := stacks.ComputeLinks(allSnaps)
	fois := make(stacks.FOIS)
	for gid, frameIdxToVars := range rawFois {
		frameIdxToFOI := make(map[int]stacks.ProcessedFOI)
		for idx, vars := range frameIdxToVars {
			var pf stacks.ProcessedFOI
			pf.Vars = make([]graph.CollectedVar, len(vars))
			for i, v := range vars {
				pf.Vars[i] = graph.CollectedVar{
					Expr:  v.Expr,
					Value: v.Val,
					Links: stacks.LinksExcludingSelf(varToLinks[v.Val], gid),
				}
			}
			frameIdxToFOI[idx] = pf
		}
		fois[gid] = frameIdxToFOI
	}

	// !!! I've removed the caching for debugging.
	// !!! f.lru.Add(snapshotID, res)
	return snap, fois, nil
}

//// filterStacks returns a new Snapshot containing the goroutines in snap that
//// contain at least a frame that matches filter.
//func (ds *DataSource) filterStacks(snap *pp.Snapshot, filter string) *pp.Snapshot {
//	if filter == "" {
//		return snap
//	}
//	res := new(pp.Snapshot)
//	*res = *snap // shallow copy
//	res.Goroutines = nil
//	for _, g := range snap.Goroutines {
//		if ds.stackMatchesFilter(g, filter) {
//			res.Goroutines = append(res.Goroutines, g)
//		}
//	}
//	return res
//}
//
//// filterStacksByPrefix returns a new Snapshot containing the goroutines in snap
//// that have the given prefix.
//func (ds *DataSource) filterStacksByPrefix(snap *pp.Snapshot, prefix []weightedtree.ScopeID) *pp.Snapshot {
//	if len(prefix) == 0 {
//		return snap
//	}
//	res := new(pp.Snapshot)
//	*res = *snap // shallow copy
//	res.Goroutines = nil
//	for _, g := range snap.Goroutines {
//		if ds.stackMatchesPrefix(g, prefix) {
//			res.Goroutines = append(res.Goroutines, g)
//		}
//	}
//	return res
//}
//
//func (ds *DataSource) stackMatchesFilter(g *pp.Goroutine, filter string) bool {
//	for i := range g.Stack.Calls {
//		if strings.Contains(g.Stack.Calls[i].Func.Complete, filter) {
//			return true
//		}
//	}
//	return false
//}
//
//func (ds *DataSource) stackMatchesPrefix(g *pp.Goroutine, prefix []weightedtree.ScopeID) bool {
//	if len(g.Stack.Calls) < len(prefix) {
//		return false
//	}
//	for i := range prefix {
//		c := &g.Stack.Calls[len(g.Stack.Calls)-i-1]
//		if stacks.ComputeScopeID(c) != prefix[i] {
//			return false
//		}
//	}
//	return true
//}
//
//// compareByFunctionName compares the function names of two nodes
//// lexicographically, returning 1 if a < b and -1 if b < a. This corresponds to
//// a descending sorting; this function is intended to be used with
//// weightedtree.Walk(), which explores "higher" nodes first so, in order to get
//// alphabetic sorting, we invert the regular comparison.
//func compareByFunctionName(a, b weightedtree.TreeNode) (int, error) {
//	aa := a.(*stacks.TreeNode)
//	bb := b.(*stacks.TreeNode)
//	res := -strings.Compare(aa.Function.Complete, bb.Function.Complete)
//	return res, nil
//}
