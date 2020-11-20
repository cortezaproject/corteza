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
		dry      bool

		processed   nodeMap
		conflicting nodeMap
	}

	ResourceState struct {
		Res             resource.Interface
		MissingDeps     resource.RefSet
		Conflicting     bool
		DepResources    []resource.Interface
		ParentResources []resource.Interface
	}

	// Provides a filter for graph iterator
	iterFilter struct {
		resourceType string
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

func (g *graph) Relink() {
	for _, n := range g.nn {
		n.cc = make(nodeSet, 0, len(n.cc))
		n.pp = make(nodeSet, 0, len(n.pp))
	}

	for res := range g.resIndex {
		n := g.resIndex[res]
		if n == nil {
			return
		}

		for _, ref := range res.Refs() {
			// else find the node and link to it (if we can)
			m := g.nn.findByRef(ref)
			if m != nil {
				g.addChild(n, m)
				g.addParent(m, n)
			}
		}
	}
}

func (g *graph) NextInverted(ctx context.Context) (s *ResourceState, err error) {
	g.inverted = true
	defer func() {
		g.inverted = false
	}()

	return g.next(ctx, nil)
}

func (g *graph) Next(ctx context.Context) (s *ResourceState, err error) {
	return g.next(ctx, nil)
}

func (g *graph) NextOf(ctx context.Context, resTrype string) (s *ResourceState, err error) {
	return g.next(ctx, &iterFilter{resourceType: resTrype})
}

func (g *graph) next(ctx context.Context, f *iterFilter) (s *ResourceState, err error) {
	upNN := g.removeProcessed(g.nodes())

	if f != nil {
		upC := make(nodeSet, 0, len(upNN))
		if f.resourceType != "" {
			for _, n := range upNN {
				if n.res.ResourceType() == f.resourceType {
					upC = append(upC, n)
				}
			}
		}
		upNN = upC
	}

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

	// There are only conflicting nodes...
	// Determine a conflicting node that should be partially resolved
	for _, upN := range upNN {
		if !g.conflicting.has(upN) {
			nx = upN
			break
		}
	}

	if nx != nil {
		// Prepare the required context for the processing.
		es := g.prepExecState(nx, true)

		g.markConflicting(nx)

		return es, nil
	} else {
		return nil, errors.New("Uhoh, couldn't determine node; @todo error")
	}
}

func (g *graph) reset() {
	g.conflicting = make(nodeMap)
	g.processed = make(nodeMap)
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
	// Dryrun shouldn't remove any nodes
	if g.dry {
		return
	}

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
	// Dryrun shouldn't remove any nodes
	if g.dry {
		return
	}

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
