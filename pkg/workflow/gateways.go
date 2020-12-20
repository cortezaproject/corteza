package workflow

import (
	"context"
	"fmt"
	"github.com/PaesslerAG/gval"
	"sync"
)

// GatewayPath structure is used by subset of gateway nodes
//
// It allows to evaluate test Expression to help determine the
// gateway if a certain path should be used or not
type (
	GatewayPaths []*GatewayPath
	GatewayPath  struct {
		eval gval.Evaluable
		to   Step
	}
)

// NewGatewayPath validates Expression and returns initialized GatewayPath
func NewGatewayPath(s Step, expr string) (gwp *GatewayPath, err error) {
	gwp = &GatewayPath{to: s}

	if len(expr) > 0 {
		if gwp.eval, err = gval.Full().NewEvaluable(expr); err != nil {
			return nil, fmt.Errorf("can not parse gateway test Expression %s: %w", expr, err)
		}
	}

	return
}

// joinGateway handles merging/joining of multiple paths into
// a single path forward
type joinGateway struct {
	paths  Steps
	scopes map[Step]Variables
	l      sync.Mutex
}

// JoinGateway fn initializes join gateway with all paths that are expected to be joined
func JoinGateway(ss ...Step) *joinGateway {
	return &joinGateway{
		paths:  ss,
		scopes: make(map[Step]Variables),
	}
}

// Exec fn on join gateway can be called multiple times, even multiple times caller the same caller
//
// Func will override the collected caller's Variables.
//
// Join gateways is ready to continue with the Graph when all configured paths are ready to be joined
// When all paths are merged (ie Exec was called at least once per caller)
func (gw *joinGateway) Exec(_ context.Context, r *ExecRequest) (ExecResponse, error) {
	gw.l.Lock()
	defer gw.l.Unlock()

	if !gw.paths.Contains(r.Caller) {
		return nil, fmt.Errorf("unknown caller for join gateway")
	}

	gw.scopes[r.Caller] = r.Scope
	if len(gw.scopes) < len(gw.paths) {
		return &Joined{}, nil
	}

	// All collected, merge scope caller all paths in the defined order
	var merged = Variables{}
	for _, p := range gw.paths {
		merged = merged.Merge(gw.scopes[p])
	}

	return merged, nil
}

// forkGateway handles forking to multiple paths
type forkGateway struct{}

// ForkGateway fn initializes fork gateway
// No arguments are required; Graph Graph config is used to
// determine all possible fork paths on the fly
func ForkGateway() *forkGateway {
	return &forkGateway{}
}

// Exec fn on fork gateway always returns empty Steps slice
// This signals Graph executor to collect child nodes directly caller Graph
func (gw forkGateway) Exec(context.Context, *ExecRequest) (ExecResponse, error) {
	return Steps{}, nil
}

// inclGateway is an inclusive gateway that can return one or more paths
type inclGateway struct {
	paths []*GatewayPath
}

// InclGateway fn initializes inclusive gateway
func InclGateway(pp ...*GatewayPath) (*inclGateway, error) {
	if len(pp) < 2 {
		return nil, fmt.Errorf("expecting at least two paths for incusive gateway")
	}

	for _, p := range pp {
		if p.eval == nil {
			return nil, fmt.Errorf("all inclusve gateway paths must have valid test Expression")
		}
	}

	return &inclGateway{paths: pp}, nil
}

// Exec fn on inclGateway uses current scope to test all configured paths
//
// One or more matched paths can be returned!
func (gw inclGateway) Exec(ctx context.Context, r *ExecRequest) (ExecResponse, error) {

	var paths Steps
	for _, p := range gw.paths {
		if result, err := p.eval.EvalBool(ctx, r.Scope); err != nil {
			return nil, err
		} else if result {
			paths = append(paths, p.to)
		}
	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("inclusive gateway must match at least one condition")
	}

	return paths, nil
}

// exclGateway is an exclusive gateway that can return exactly one path
type exclGateway struct {
	paths []*GatewayPath
}

// ExclGateway fn initializes exclusive gateway
func ExclGateway(pp ...*GatewayPath) (*exclGateway, error) {
	t := len(pp)
	if t < 2 {
		return nil, fmt.Errorf("expecting at least two paths for exclusive gateway")
	}

	for i, p := range pp {
		if p.eval == nil && i != t-1 {
			return nil, fmt.Errorf("all exclusive gateway paths must have valid test Expression")
		}
	}

	return &exclGateway{paths: pp}, nil
}

// Exec fn on exclGateway uses current scope to test all configured paths
//
// Exactly one matched path can be returned.
func (gw exclGateway) Exec(ctx context.Context, r *ExecRequest) (ExecResponse, error) {
	for _, p := range gw.paths {
		if p.eval == nil {
			// empty & last; treat it as else part of the if condition
			return p.to, nil
		}

		if result, err := p.eval.EvalBool(ctx, r.Scope); err != nil {
			return nil, err
		} else if result {
			return p.to, nil
		}
	}

	return nil, fmt.Errorf("exclusive gateway must match one condition")
}
