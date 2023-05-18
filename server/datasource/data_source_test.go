package datasource

import (
	"context"
	"fmt"
	"github.com/google/traceviz/server/go/util"
	weightedtree "github.com/google/traceviz/server/go/weighted_tree"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/stretchr/testify/require"
	"path"
	"runtime"
	"testing"
)

func TestDataSource(t *testing.T) {
	ctx := context.Background()
	lru, err := lru.New[string, processSnapshot](100)
	require.NoError(t, err)
	_, filename, _, _ := runtime.Caller(0)
	src := DataSource{fetcher: &stacksFetcherImpl{
		rootDir: path.Dir(filename),
		lru:     lru,
	}}
	require.NoError(t, src.HandleDataSeriesRequests(ctx, map[string]*util.V{
		collectionNameKey: util.StringValue("cockroachdb_example_snapshot.txt"),
	}, nil, nil))
}

func TestBuildTree(t *testing.T) {
	ctx := context.Background()
	lru, err := lru.New[string, processSnapshot](100)
	require.NoError(t, err)
	_, filename, _, _ := runtime.Caller(0)
	src := DataSource{fetcher: &stacksFetcherImpl{
		rootDir: path.Dir(filename),
		lru:     lru,
	}}
	col, err := src.fetcher.Fetch(ctx, "example.txt")
	require.NoError(t, err)
	require.NotNil(t, col.snapshot)
	tree := src.buildTree(col.snapshot)
	require.NotNil(t, tree)
	fmt.Printf("!!! tree:\n")
	tree.prettyPrint()

	drb := util.NewDataResponseBuilder()
	req := &util.DataSeriesRequest{
		QueryName:  "test query",
		SeriesName: "my series",
		Options:    nil,
	}
	builder := drb.DataSeries(req)
	src.filterByPathPrefix(tree, []weightedtree.ScopeID{10723355014697041956}, builder)
	// !!! src.filterByPathPrefix(tree, []weightedtree.ScopeID{14477739511207775948, 17773099551792462126}, builder)

	data, err := drb.Data()
	require.NoError(t, err)
	fmt.Printf(data.PrettyPrint() + "\n")
}

// !!!
//func TestXXX(t *testing.T) {
//	buf := make([]byte, 1<<20)
//	n := runtime.Stack(buf, true)
//	buf = buf[:n]
//	fmt.Printf("%s\n", string(buf))
//
//	r := bytes.NewReader(buf)
//	snap, _, err := pp.ScanSnapshot(r, io.Discard, pp.DefaultOpts())
//	if err != nil && err != io.EOF {
//		fmt.Printf("!!! reading stacks - err: %s\n", err)
//	}
//	if snap == nil {
//		fmt.Printf("!!! nil snap\n")
//	}
//}
