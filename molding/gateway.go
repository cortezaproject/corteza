package molding

import (
	"context"
	"fmt"
	"github.com/PaesslerAG/gval"
	"sync"
)

type (
	gatewayPath struct {
		expr string
		eval gval.Evaluable
		to   Node
	}

	joinGateway struct {
		nodeRef string

		// all parent nodes we'll wait to be executed before going to next node
		paths []Iterator
		next  Node

		l sync.Mutex

		scope map[Node]Variables
	}

	forkGateway struct {
		nodeRef string
		paths   Nodes
	}

	inclGateway struct {
		nodeRef string
		paths   []*gatewayPath
	}

	exclGateway struct {
		nodeRef string
		paths   []*gatewayPath
	}
)

var (
	_ Joiner = &joinGateway{}
	_ Tester = &forkGateway{}
	_ Tester = &inclGateway{}
	_ Tester = &exclGateway{}
)

func NewGatewayCondition(expr string, to Node) *gatewayPath {
	return &gatewayPath{expr: expr, to: to}
}

func NewGatewayNoCondition(to Node) *gatewayPath {
	return &gatewayPath{to: to}
}

func initGatewayPaths(paths ...*gatewayPath) ([]*gatewayPath, error) {
	var (
		err error
	)

	for _, p := range paths {
		if p.expr == "" {
			continue
		}

		if p.eval, err = gval.Full().NewEvaluable(p.expr); err != nil {
			return nil, fmt.Errorf("can not parse %s: %w", p.expr, err)
		}
	}

	return paths, err
}

func NewJoinGateway(paths ...Iterator) (*joinGateway, error) {
	join := &joinGateway{
		paths: paths,
		scope: make(map[Node]Variables),
	}
	for _, p := range paths {
		p.SetNext(join)
	}
	return join, nil
}

func (gw *joinGateway) NodeRef() string { return gw.nodeRef }
func (gw *joinGateway) Next() Node      { return gw.next }
func (gw *joinGateway) SetNext(n Node)  { gw.next = n }
func (gw *joinGateway) Paths() Nodes {
	var pp = make(Nodes, 0, len(gw.paths))
	for _, p := range gw.paths {
		pp = append(pp, p)
	}

	return pp
}

func (gw *joinGateway) Join(p Node, scope Variables) (Node, Variables, error) {
	gw.l.Lock()
	defer gw.l.Unlock()

	// Allow scope overriding (in case when
	// parent is executed again
	//
	// This covers scenario where we route workflow back to one
	// of the nodes that is then joined
	gw.scope[p] = scope

	if len(gw.scope) < len(gw.paths) {
		// Not all collected
		return nil, nil, nil
	}

	// All collected, merge scope from all paths in the
	// defined order
	var out = Variables{}
	for _, p := range gw.paths {
		out = out.Merge(gw.scope[p])
	}

	return gw.next, out, nil
}

func NewForkGateway(paths ...Node) (*forkGateway, error) { return &forkGateway{paths: paths}, nil }

func (gw forkGateway) NodeRef() string                                    { return gw.nodeRef }
func (gw forkGateway) Paths() Nodes                                       { return gw.paths }
func (gw forkGateway) Test(_ context.Context, _ Variables) (Nodes, error) { return gw.paths, nil }

// multiple matches
func NewInclGateway(paths ...*gatewayPath) (*inclGateway, error) {
	var err error
	paths, err = initGatewayPaths(paths...)
	return &inclGateway{paths: paths}, err
}

func (gw inclGateway) NodeRef() string { return gw.nodeRef }
func (gw inclGateway) Paths() Nodes {
	var paths Nodes
	for _, p := range gw.paths {
		paths = append(paths, p.to)
	}
	return paths
}

// Test returns nodes from all paths that have a matching condition (or no condition at all)
func (gw inclGateway) Test(ctx context.Context, scope Variables) (to Nodes, err error) {
	for _, p := range gw.paths {
		if result, err := p.eval.EvalBool(ctx, scope); err != nil {
			return nil, err
		} else if result {
			to = append(to, p.to)
		}
	}

	return
}

// single match
func NewExclGateway(paths ...*gatewayPath) (*exclGateway, error) {
	var err error
	paths, err = initGatewayPaths(paths...)
	return &exclGateway{paths: paths}, err
}

func (gw exclGateway) NodeRef() string { return gw.nodeRef }
func (gw exclGateway) Paths() Nodes {
	var paths Nodes
	for _, p := range gw.paths {
		paths = append(paths, p.to)
	}
	return paths
}

// Test returns first path with matching condition
func (gw exclGateway) Test(ctx context.Context, scope Variables) (to Nodes, err error) {
	for i, p := range gw.paths {
		if len(p.expr) == 0 && i == len(gw.paths)-1 {
			// empty & last; treat it as else
			return Nodes{p.to}, nil
		}

		if result, err := p.eval.EvalBool(ctx, scope); err != nil {
			return nil, err
		} else if result {
			return Nodes{p.to}, nil
		}
	}

	return nil, fmt.Errorf("could not match any of conditions")
}
