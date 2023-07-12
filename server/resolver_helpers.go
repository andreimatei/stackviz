package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-delve/delve/service/debugger"
	pp "github.com/maruel/panicparse/v2/stack"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"stacksviz/ent"
	"stacksviz/graph"
	"stacksviz/stacks"
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

func (r *mutationResolver) getSnapshotFromDelveAgent(
	ctx context.Context, agentAddr string, spec *ent.CollectSpec,
) (agentrpc.Snapshot, error) {
	client, err := rpc.DialHTTP("tcp", agentAddr)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// Load the spec for the frame of interests.
	frames, err := spec.Frames(ctx)
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

	// Call the Delve agent asking for a snapshot with the specified frames of
	// interest.
	var res = &agentrpc.GetSnapshotOut{}
	err = client.Call("Agent.GetSnapshot", args, &res)
	if err != nil {
		log.Fatal("call to agent failed: ", err)
	}
	return res.Snapshot, nil
}

func reconcileFlightRecorder(ctx context.Context, agentAddr string, events []FlightRecorderEventSpecFull) error {
	client, err := rpc.DialHTTP("tcp", agentAddr)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var eventsIn []agentrpc.FlightRecorderEventSpec
	for _, e := range events {
		eventsIn = append(eventsIn, agentrpc.FlightRecorderEventSpec{
			Frame:   e.Frame,
			Expr:    e.Expr,
			KeyExpr: e.KeyExpr,
		})
	}

	log.Printf("calling Agent.ReconcileFlightRecorder with events: %#v", eventsIn)
	out := agentrpc.ReconcileFLightRecorderOut{}
	return client.Call("Agent.ReconcileFlightRecorder", &agentrpc.ReconcileFlightRecorderIn{
		Events: eventsIn,
	}, &out)
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

func stacksToString(s agentrpc.Snapshot) string {
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
			fi := dbClient.FrameSpec.Create().
				SetFrame(f.frame).
				SetParentCollection(cs).
				SetCollectExpressions(f.exprs).
				SetFlightRecorderEvents([]string{}).
				SaveX(ctx)
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

type target struct {
	serviceName, processName, URL string
}

func (r *Resolver) getTargets() []target {
	var svcName string
	for serviceName, _ := range r.conf.Targets {
		// TODO(andrei): deal with multiple services
		svcName = serviceName
		break
	}
	res := make([]target, 0, len(r.conf.Targets[svcName]))
	for processName, url := range r.conf.Targets[svcName] {
		res = append(res, target{
			serviceName: svcName,
			processName: processName,
			URL:         url,
		})
	}
	return res
}

func stackMatchesFilter(g *pp.Goroutine, filter string) bool {
	for i := range g.Stack.Calls {
		if strings.Contains(g.Stack.Calls[i].Func.Complete, filter) {
			return true
		}
	}
	return false
}

type Snapshot struct {
	stacks             *pp.Snapshot
	framesOfInterest   stacks.FOIS
	flightRecorderData map[string][]string
}

func (r *Resolver) loadSnapshot(
	ctx context.Context, collectionID int, snapshotID int,
) (Snapshot, error) {
	// Load the snapshot from the database.
	snapRec := r.dbClient.ProcessSnapshot.GetX(ctx, snapshotID)

	// Parse the stacks.
	opts := pp.DefaultOpts()
	opts.ParsePC = true
	snap, _, err := pp.ScanSnapshot(strings.NewReader(snapRec.Snapshot), io.Discard, opts)
	if err != nil && err != io.EOF {
		return Snapshot{}, err
	}
	if snap == nil {
		return Snapshot{}, fmt.Errorf("failed to parse any stacks")
	}

	var rawFois stacks.FramesOfInterest
	if snapRec.FramesOfInterest != "" {
		if err = json.Unmarshal([]byte(snapRec.FramesOfInterest), &rawFois); err != nil {
			return Snapshot{}, fmt.Errorf("failed to unmarshal frames of interest: %w", err)
		}
	} else {
		log.Printf("!!! snapshot %d does not have any frames of interest", snapshotID)
	}

	// Load all snapshots in this collection and find links between the captured variables.
	c := r.dbClient.Collection.GetX(ctx, collectionID)
	allSnaps, err := c.ProcessSnapshots(ctx)
	if err != nil {
		return Snapshot{}, err
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
					Expr:     v.Expr,
					Value:    v.Val,
					Links:    stacks.LinksExcludingSelf(varToLinks[v.Val], gid),
					FrameIdx: idx,
				}
			}
			frameIdxToFOI[idx] = pf
		}
		fois[gid] = frameIdxToFOI
	}

	return Snapshot{
		stacks:             snap,
		framesOfInterest:   fois,
		flightRecorderData: snapRec.FlightRecorderData,
	}, nil
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

type FlightRecorderEventSpec struct {
	Expr    string
	KeyExpr string
}

type FlightRecorderEventSpecFull struct {
	FlightRecorderEventSpec
	Frame string
}
