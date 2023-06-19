package server

import (
	"context"
	"github.com/go-delve/delve/service/debugger"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"stacksviz/ent"
	"strings"

	"github.com/andreimatei/delve-agent/agentrpc"
	"github.com/kr/pretty"
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

	frames, err := spec.Frames(ctx)
	if err != nil {
		return agentrpc.Snapshot{}, err
	}
	framesSpec := make(map[string][]string)
	for _, f := range frames {
		for _, e := range f.Exprs {
			framesSpec[f.Frame] = append(framesSpec[f.Frame], e)
		}
	}
	args := &agentrpc.GetSnapshotIn{FramesSpec: framesSpec}

	var res = &agentrpc.GetSnapshotOut{}
	err = client.Call("Agent.GetSnapshot", args, &res)
	if err != nil {
		log.Fatal("call to agent failed: ", err)
	}
	pretty.Print(res) // !!!

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
