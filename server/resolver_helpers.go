package server

import (
	"context"
	"github.com/go-delve/delve/service/debugger"
	pp "github.com/maruel/panicparse/v2/stack"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"stacksviz/ent"
	"stacksviz/graph"
	"strings"

	"github.com/andreimatei/delve-agent/agentrpc"
)

func (r *mutationResolver) getSnapshotFromPprof(targetURL string) (string, error) {
	resp, err := http.Get(targetURL)
	// TODO(andrei): try the other processes instead of bailing out
	if err != nil {
		return "", err
	}
	body := resp.Body
	defer body.Close()

	snap, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(snap), nil
}

func (r *mutationResolver) getSnapshotFromDelveAgent(ctx context.Context, agentAddr string, spec *ent.CollectSpec) (agentrpc.Snapshot, error) {
	client, err := rpc.DialHTTP("tcp", agentAddr)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	frames, err := spec.QueryFrames().All(ctx)
	//frames, err := spec.Frames(ctx)
	if err != nil {
		return agentrpc.Snapshot{}, err
	}
	framesSpec := make(map[string][]string)
	for _, f := range frames {
		for _, e := range f.CollectExpressions {
			framesSpec[f.Frame] = append(framesSpec[f.Frame], e)
		}
	}
	args := &agentrpc.GetSnapshotIn{FramesSpec: framesSpec}

	var res = &agentrpc.GetSnapshotOut{}
	err = client.Call("Agent.GetSnapshot", args, &res)
	if err != nil {
		log.Fatal("call to agent failed: ", err)
	}
	//pretty.Print(res) // !!!

	//var sb strings.Builder
	//for _, stack := range res.Snapshot.Stacks {
	//	sb.WriteString(stack)
	//	sb.WriteRune('\n')
	//}
	//
	//return sb.String(), nil
	return res.Snapshot, nil
}

func (r *queryResolver) getAvailableVarsFromDelveAgent(agentAddr string, fn string, pcOff int64) ([]debugger.VarInfo, []debugger.TypeInfo, error) {
	log.Printf("!!! getting available vars for %s:0x%x", fn, pcOff)
	client, err := rpc.DialHTTP("tcp", agentAddr)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	args := &agentrpc.ListVarsIn{Func: fn, PCOff: pcOff}
	var res = &agentrpc.ListVarsOut{}
	err = client.Call("Agent.ListVars", args, &res)
	if err != nil {
		log.Fatal("call to agent failed: ", err)
	}
	return res.Vars, res.Types, nil
}

func snapToString(s agentrpc.Snapshot) string {
	var sb strings.Builder
	for _, stack := range s.Stacks {
		sb.WriteString(stack)
		sb.WriteRune('\n')
	}
	return sb.String()
}

func getOrCreateCollectSpec(ctx context.Context, dbClient *ent.Client) *ent.CollectSpec {
	cis := dbClient.CollectSpec.Query().AllX(ctx)
	if len(cis) > 1 {
		log.Fatalf("expected at most one CollectSpec, got: %d", len(cis))
	}
	// If there isn't a CollectSpec already, create one.
	if len(cis) == 0 {
		cs := dbClient.CollectSpec.Create().SaveX(ctx)

		for _, f := range []struct {
			frame string
			exprs []string
		}{
			{
				frame: "google.golang.org/grpc.(*csAttempt).recvMsg",
				exprs: []string{"a.s.id"},
			},
			{
				frame: "google.golang.org/grpc.(*Server).processUnaryRPC",
				exprs: []string{"stream.id"},
			},
		} {
			fi := dbClient.FrameSpec.Create().SetFrame(f.frame).SetCollectSpecRef(cs).SetCollectExpressions(f.exprs).SetFlightRecorderEvents([]string{}).SaveX(ctx)
			_ = fi
			// !!! cs = cs.Update().AddFrames(fi).SaveX(ctx)
		}
		return cs
	}

	return cis[0]
}

func (r *queryResolver) getTypeInfoFromDelveAgent(agentAddr string, typeName string) ([]graph.FieldInfo, error) {
	client, err := rpc.DialHTTP("tcp", agentAddr)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	args := &agentrpc.GetTypeInfoIn{Name: typeName}
	var res = &agentrpc.GetTypeInfoOut{}
	err = client.Call("Agent.GetTypeInfo", args, &res)
	if err != nil {
		log.Fatal("call to agent failed: ", err)
	}

	fields := make([]graph.FieldInfo, len(res.Fields))
	for i, f := range res.Fields {
		fields[i] = graph.FieldInfo{
			Name:     f.Name,
			Type:     f.TypeName,
			Embedded: f.Embedded,
		}
	}
	return fields, nil
}

func stackMatchesFilter(g *pp.Goroutine, filter string) bool {
	for i := range g.Stack.Calls {
		if strings.Contains(g.Stack.Calls[i].Func.Complete, filter) {
			return true
		}
	}
	return false
}

// filterStacks returns a new Snapshot containing the goroutines in snap that
// contain at least a frame that matches filter.
func filterStacks(snap *pp.Snapshot, gid *int, filter *string) *pp.Snapshot {
	res := new(pp.Snapshot)
	*res = *snap // shallow copy

	if gid != nil {
		origGs := res.Goroutines
		res.Goroutines = nil
		for _, g := range origGs {
			if g.ID == *gid {
				res.Goroutines = append(res.Goroutines, g)
				break
			}
		}
	}

	if filter != nil && *filter != "" {
		origGs := res.Goroutines
		res.Goroutines = nil
		for _, g := range origGs {
			if stackMatchesFilter(g, *filter) {
				res.Goroutines = append(res.Goroutines, g)
			}
		}
	}

	return res
}

func flatten[T any](src []*T) []T {
	res := make([]T, len(src))
	for i := range src {
		res[i] = *src[i]
	}
	return res
}
