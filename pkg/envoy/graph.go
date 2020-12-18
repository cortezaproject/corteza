package envoy

import (
	"context"
	"errors"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	nodeMap map[*node]bool

	graph struct {
		nn       nodeSet
		resIndex map[resource.Interface]*node

		// Config flags
		inverted bool

		processed   nodeMap
		conflicting nodeMap
	}

	ResourceState struct {
		Res             resource.Interface
		Conflicting     bool
		DepResources    []resource.Interface
		ParentResources []resource.Interface
	}
)

func newGraph() *graph {
	return &graph{
		nn:          make(nodeSet, 0, 100),
		resIndex:    make(map[resource.Interface]*node),
		inverted:    false,
		processed:   make(nodeMap),
		conflicting: make(nodeMap),
	}
}

func (g *graph) addNode(nn ...*node) {
	for _, n := range nn {
		g.resIndex[n.res] = n
	}
	g.nn = g.nn.add(nn...)
}

func (g *graph) removeNode(nn ...*node) {
	g.nn = g.nn.remove(nn...)
	for _, n := range nn {
		g.removeChild(n)
		g.removeParent(n)
	}
}

func (g *graph) invert() {
	g.inverted = !g.inverted
}

func (g *graph) childNodes(n *node) nodeSet {
	if g.inverted {
		return n.pp
	}

	return n.cc
}

func (g *graph) parentNodes(n *node) nodeSet {
	if g.inverted {
		return n.cc
	}

	return n.pp
}

func (g *graph) nodes() nodeSet {
	return g.nn
}

func (g *graph) markProcessed(n *node) {
	g.processed = g.processed.add(n)
}
func (g *graph) markConflicting(n *node) {
	g.conflicting = g.conflicting.add(n)
}

func (g *graph) NextInverted(ctx context.Context) (s *ResourceState, err error) {
	g.inverted = true
	defer func() {
		g.inverted = false
	}()

	return g.Next(ctx)
}

func (g *graph) Next(ctx context.Context) (s *ResourceState, err error) {
	upNN := g.removeProcessed(g.nodes())

	// We are done here
	if len(upNN) <= 0 {
		return nil, nil
	}

	var nx *node

	for _, upN := range upNN {
		// We should not take into account conflicted parent nodes,
		// as they already resolved the conflict.
		if len(g.removeProcessed(g.removeConflicting(g.parentNodes(upN)))) <= 0 {
			nx = upN
			break
		}
	}

	if nx != nil {
		// Prepare the required context for the processing.
		// Perform some basic cleanup.
		es := g.prepExecState(nx, false)
		g.markProcessed(nx)

		return es, nil
	}

	// There are only conflicting nodes
	// Try to get a cycle node to resolve the conflict
	nx = g.findCycleNode(g.removeConflicting(upNN))

	if nx == nil {
		// This is basically impossible, unless I've messed up the algorithm
		return nil, errors.New("could not determine non-conflicting node")
	}

	// Prepare the required context for the processing.
	es := g.prepExecState(nx, true)
	g.markConflicting(nx)

	return es, nil
}

// findCycleNode returns the first graph node that caused a cycle
//
// General outline:
//  * DFS from a start node(s)
//  * if a child node is already in path, return that node
//  * else return nil and cleanup the path until the first node with
//    unprocessed child nodes
//
// @note we could complicate this further by doing cycle enumeration algorithms.
//       I might do it when no one is watching :)
func (g *graph) findCycleNode(nn nodeSet) *node {
	path := make(nodeMap)
	processed := make(nodeMap)

	for _, n := range nn {
		if !processed.has(n) {
			m := g.traverse(n, path, processed)
			if m != nil {
				return m
			}
		}
	}

	return nil
}

func (g *graph) traverse(n *node, path, processed nodeMap) *node {
	processed.add(n)

	// Found ourselves a cycle
	if path.has(n) {
		return n
	}
	cnn := g.removeProcessed(g.removeConflicting(g.childNodes(n)))
	// Nothing else to look at
	if len(cnn) == 0 {
		return nil
	}

	path.add(n)

	for _, c := range cnn {
		m := g.traverse(c, path, processed)
		if m != nil {
			return m
		}
	}

	path.remove(n)
	return nil
}

// util

func (g *graph) removeProcessed(nn nodeSet) nodeSet {
	mm := make(nodeSet, 0, len(nn))

	for _, n := range nn {
		if !g.processed.has(n) {
			mm = mm.add(n)
		}
	}

	return mm
}

func (g *graph) removeConflicting(nn nodeSet) nodeSet {
	mm := make(nodeSet, 0, len(nn))

	for _, n := range nn {
		if !g.conflicting.has(n) {
			mm = mm.add(n)
		}
	}

	return mm
}

func (g *graph) addChild(n *node, mm ...*node) {
	if g.inverted {
		n.pp = n.pp.add(mm...)
	} else {
		n.cc = n.cc.add(mm...)
	}
}

func (g *graph) removeChild(n *node, mm ...*node) {
	if len(mm) <= 0 {
		if g.inverted {
			n.pp = make(nodeSet, 0, 10)
		} else {
			n.cc = make(nodeSet, 0, 10)
		}
	} else {
		if g.inverted {
			n.pp = n.pp.remove(mm...)
		} else {
			n.cc = n.cc.remove(mm...)
		}
	}
}

func (g *graph) addParent(n *node, mm ...*node) {
	if g.inverted {
		n.cc = n.cc.add(mm...)
	} else {
		n.pp = n.pp.add(mm...)
	}
}

func (g *graph) removeParent(n *node, mm ...*node) {
	if len(mm) <= 0 {
		if g.inverted {
			n.cc = make(nodeSet, 0, 10)
		} else {
			n.pp = make(nodeSet, 0, 10)
		}
	} else {
		if g.inverted {
			n.cc = n.cc.remove(mm...)
		} else {
			n.pp = n.pp.remove(mm...)
		}
	}
}

func (g *graph) nodeResource(nn ...*node) []resource.Interface {
	rr := make([]resource.Interface, 0, len(nn))
	for _, n := range nn {
		rr = append(rr, n.res)
	}

	return rr
}

func (g *graph) prepExecState(n *node, conflicting bool) *ResourceState {
	return &ResourceState{
		Res:             n.res,
		DepResources:    g.nodeResource(g.childNodes(n)...),
		ParentResources: g.nodeResource(g.parentNodes(n)...),
		Conflicting:     conflicting,
	}
}

func (nm nodeMap) add(n *node) nodeMap {
	nm[n] = true
	return nm
}
func (nm nodeMap) has(n *node) bool {
	return nm[n]
}
func (nm nodeMap) remove(n *node) nodeMap {
	delete(nm, n)
	return nm
}
