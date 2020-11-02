package envoy

import (
	"context"
	"errors"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	graph struct {
		nn nodeSet

		// Config flags
		inverted bool
		dry      bool

		processed   nodeSet
		conflicting nodeSet
	}

	ExecState struct {
		Res             resource.Interface
		MissingDeps     resource.RefSet
		Conflicting     bool
		DepResources    []resource.Interface
		ParentResources []resource.Interface
	}
)

func newGraph() *graph {
	return &graph{
		nn:          make(nodeSet, 0, 100),
		inverted:    false,
		processed:   make(nodeSet, 0, 100),
		conflicting: make(nodeSet, 0, 100),
	}
}

func (g *graph) addNode(nn ...*node) {
	g.nn = g.nn.add(nn...)
}

func (g *graph) removeNode(nn ...*node) {
	g.nn = g.nn.remove(nn...)
}

func (g *graph) invert() {
	g.inverted = !g.inverted
}

func (g *graph) DryRun() {
	g.dry = true
	g.purge(true)
}

func (g *graph) ProdRun() {
	g.dry = false
	g.purge(true)
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

func (g *graph) Next(ctx context.Context) (s *ExecState, err error) {
	upNN := g.removeProcessed(g.nodes())

	// We are done here
	if len(upNN) <= 0 {
		g.purge(g.dry)
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
		pp := g.parentNodes(nx)
		g.removeParent(nx)

		for _, p := range pp {
			g.removeChild(p, nx)
			if len(g.childNodes(p)) <= 0 && g.processed.has(p) {
				g.removeNode(p)
				g.processed = g.processed.remove(p)
				g.conflicting = g.conflicting.remove(p)
			}
		}

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

func (g *graph) purge(partial bool) {
	g.conflicting = nil
	g.processed = nil

	if !partial {
		g.nn = nil
	}
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

func (g *graph) prepExecState(n *node, conflicting bool) *ExecState {
	return &ExecState{
		Res:             n.res,
		DepResources:    g.nodeResource(g.childNodes(n)...),
		ParentResources: g.nodeResource(g.parentNodes(n)...),
		Conflicting:     conflicting,
	}
}
