// Copyright 2023 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package stacks

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/andreimatei/delve-agent/agentrpc"
	weightedtree "github.com/google/traceviz/server/go/weighted_tree"
	pp "github.com/maruel/panicparse/v2/stack"
	"hash/fnv"
	"stacksviz/ent"
	"strconv"
	"strings"
)

// goroutineID to frame index to list of variables
type FramesOfInterest map[int]map[int][]agentrpc.CapturedExpr

// goroutineID to frame index to list of variables
type FOIS map[int]map[int]ProcessedFOI
type ProcessedFOI struct {
	Vars []VarInfo
}
type VarInfo struct {
	Expr  string
	Value string
	Links []Link
}
type Link struct {
	SnapshotID  int
	GoroutineID int
	FrameIdx    int
}
type Frame struct {
	call pp.Call
	vars []VarInfo
}

// BuildTree builds a trie out of the stack traces in snap.
//
// capturedData represents variables captured from different frames. It's OK for
// capturedData to reference goroutines not in snap; such data will be ignored.
func BuildTree(snap *pp.Snapshot, capturedData FOIS) *TreeNode {
	root := &TreeNode{
		Function: pp.Func{
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
		myFois := capturedData[s.ID]
		stack := make([]Frame, l)
		for i := range s.Signature.Stack.Calls {
			stack[l-i-1] = Frame{
				call: s.Signature.Stack.Calls[i],
				vars: myFois[i].Vars,
			}
		}
		root.addStack(stack)
	}
	return root
}

// !!!
//func FindLinks(v string, snaps []*ent.ProcessSnapshot) []Link {
//	var res []Link
//	for _, s := range snaps {
//		if s.FramesOfInterest == "" {
//			continue
//		}
//		var fois FramesOfInterest
//		if err := json.Unmarshal([]byte(s.FramesOfInterest), &fois); err != nil {
//			panic(err)
//		}
//		for gid, m := range fois {
//			for frameIdx, vars := range m {
//				for _, vv := range vars {
//					if vv == v {
//						res = append(res, Link{
//							SnapshotID:  s.ID,
//							GoroutineID: gid,
//							FrameIdx:    frameIdx,
//						})
//					}
//				}
//			}
//		}
//	}
//
//	return res
//}

func ComputeLinks(snaps []*ent.ProcessSnapshot) map[string][]Link {
	res := make(map[string][]Link)
	for _, s := range snaps {
		if s.FramesOfInterest == "" {
			continue
		}
		var fois FramesOfInterest
		if err := json.Unmarshal([]byte(s.FramesOfInterest), &fois); err != nil {
			panic(err)
		}

		for gid, frameIdxToVars := range fois {
			for frameIdx, vars := range frameIdxToVars {
				for _, v := range vars {
					res[v.Val] = append(res[v.Val],
						Link{
							SnapshotID:  s.ID,
							GoroutineID: gid,
							FrameIdx:    frameIdx,
						})
				}
			}
		}
	}
	return res
}

// TreeNode is a node in a trie of stack traces. Each node represents a
// function; its children are other functions called by the node's function in
// one or more stacks.
type TreeNode struct {
	Function pp.Func
	File     string
	Line     int
	PcOffset int64
	Vars     [][]VarInfo
	// path is the path from the root to this node, represented by hashes of
	// each ancestor's Function.
	path     []weightedtree.ScopeID
	children []TreeNode
	// NumLeafGoroutines counts how many goroutines have this node as their leaf
	// function. This results in the "self magnitude" of the node when rendered
	// as a flame graph - i.e. how much weight it needs to have in addition to
	// the sum of the children's weights.
	NumLeafGoroutines int
	NumGoroutines     int
}

// scopeID returns the identifier for this node.
func (t *TreeNode) scopeID() weightedtree.ScopeID {
	if len(t.path) > 0 {
		return t.path[len(t.path)-1]
	}
	return 0
}

var _ weightedtree.TreeNode = &TreeNode{}

// Path is part of the weightedtree.TreeNode interface.
func (t *TreeNode) Path() []weightedtree.ScopeID {
	return t.path
}

func (t *TreeNode) pathAsStrings() []string {
	path := make([]string, len(t.Path()))
	for i, p := range t.Path() {
		path[i] = strconv.FormatUint(uint64(p), 10)
	}
	return path
}

// Children is part of the weightedtree.TreeNode interface.
func (t *TreeNode) Children(ids ...weightedtree.ScopeID) ([]weightedtree.TreeNode, error) {
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

func (t *TreeNode) prettyPrint() {
	t.prettyPrintInner(0)
}

func (t *TreeNode) prettyPrintInner(indent int) {
	var sb strings.Builder
	for i := 0; i < indent; i++ {
		sb.WriteRune('\t')
	}
	fmt.Printf("%s(%d) %s (%s:%d) (%v)\n", sb.String(), t.NumLeafGoroutines, t.Function.Complete, t.File, t.Line, t.path)
	for i := range t.children {
		t.children[i].prettyPrintInner(indent + 1)
	}
}

// findChild finds the child of t for a call at File:line. If such a child
// doesn't exist, returns nil.
func (t *TreeNode) findChild(file string, line int) *TreeNode {
	for i := range t.children {
		c := &t.children[i]
		if c.File == file && c.Line == line {
			return c
		}
	}
	return nil
}

// addStack adds the stack to the tree rooted at t, creating new nodes for calls
// that don't yet exist.
func (t *TreeNode) addStack(stack []Frame) {
	t.NumGoroutines++
	if len(stack) == 0 {
		// t is a leaf for the stack that we just finished processing.
		t.NumLeafGoroutines++
		return
	}
	child := t.findChild(stack[0].call.RemoteSrcPath, stack[0].call.Line)
	if child != nil {
		if len(stack[0].vars) > 0 {
			child.Vars = append(child.Vars, stack[0].vars)
		}
		child.addStack(stack[1:])
	} else {
		t.createPath(stack)
	}
}

// createPath adds children to t recursively such that the tree gets the path
// t -> stack[0] -> stack[1] -> ...
func (t *TreeNode) createPath(stack []Frame) {
	t.NumGoroutines++
	if len(stack) == 0 {
		// The stack had t as a leaf function.
		t.NumLeafGoroutines++
		return
	}
	call := &stack[0].call
	var vars [][]VarInfo
	if len(stack[0].vars) > 0 {
		vars = [][]VarInfo{stack[0].vars}
	}
	t.children = append(t.children, TreeNode{
		Function:          call.Func,
		File:              call.RemoteSrcPath,
		Line:              call.Line,
		PcOffset:          call.PCOffset,
		path:              append(t.path, ComputeScopeID(call)),
		children:          nil,
		NumLeafGoroutines: 0,
		Vars:              vars,
	})
	t.children[len(t.children)-1].createPath(stack[1:])
}

var _ json.Marshaler = &TreeNode{}

func (t *TreeNode) MarshalJSON() ([]byte, error) {
	/*
		{
			"name": "...",
			"details": "...",
			"file": "...",
			"line": x,
			"pcoff": x,
			"value": x,
			"children": [
				...recurse...
			]
		}
	*/

	// !!!
	//var varsProp []string
	//{
	//	var sb strings.Builder
	//	for _, frame := range t.Vars {
	//		sb.Reset()
	//		for _, v := range frame {
	//			sb.WriteString(v.Value)
	//			sb.WriteRune('\n')
	//		}
	//		varsProp = append(varsProp, sb.String())
	//	}
	//}
	//if len(varsProp) != 0 {
	//	log.Printf("!!! found frame with vars: %s - (%d) %+v %q", t.Function.Complete, len(varsProp), varsProp, varsProp[0])
	//}

	var sb strings.Builder
	sb.WriteString("{\n")
	sb.WriteString("\t\"name\": ")
	sb.WriteString(fmt.Sprintf("%q", t.Function.DirName+"."+t.Function.Name))
	sb.WriteString(",\n\t\"file\": ")
	sb.WriteString(fmt.Sprintf("%q", t.File))
	sb.WriteString(",\n\t\"line\": ")
	sb.WriteString(strconv.Itoa(t.Line))
	sb.WriteString(",\n\t\"pcoff\": ")
	sb.WriteString(strconv.Itoa(int(t.PcOffset)))
	sb.WriteString(",\n\t\"details\": ")
	sb.WriteString(fmt.Sprintf("%q", t.Function.Complete))
	sb.WriteString(",\n\t\"vars\": ")
	// !!! varsJSON, err := json.Marshal(varsProp)
	varsJSON, err := json.Marshal(t.Vars)
	if err != nil {
		return nil, err
	}
	sb.WriteString(string(varsJSON))
	sb.WriteString(",\n\t\"value\": ")
	sb.WriteString(strconv.Itoa(t.NumGoroutines))
	if len(t.children) > 0 {
		sb.WriteString(",\n\t\"children\": [\n")
		for i := range t.children {
			c, err := t.children[i].MarshalJSON()
			if err != nil {
				panic(err)
			}
			sb.WriteString(string(c))
			if i < len(t.children)-1 {
				sb.WriteRune(',')
			}
			sb.WriteRune('\n')
		}
		sb.WriteString("]\n")
	} else {
		sb.WriteRune('\n')
	}
	sb.WriteRune('}')
	return []byte(sb.String()), nil
}

func (t *TreeNode) ToJSON() string {
	// !!!
	//xxx, err := t.MarshalJSON()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("!!! json: %s", xxx)

	s, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(s)
}

func ComputeScopeID(call *pp.Call) weightedtree.ScopeID {
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
