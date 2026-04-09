package wfexec

import (
	"context"
	"fmt"
	"sync"

	"github.com/cortezaproject/corteza/server/pkg/expr"
)

// GatewayPath structure is used by subset of gateway nodes
//
// It allows to evaluate test Expression to help determine the
// gateway if a certain path should be used or not
type (
	GatewayPaths []*GatewayPath
	GatewayPath  struct {
		test pathTester
		to   Step
	}

	pathTester func(context.Context, *expr.Vars) (bool, error)
)

// NewGatewayPath validates Expression and returns initialized GatewayPath
func NewGatewayPath(s Step, t pathTester) (gwp *GatewayPath, err error) {
	return &GatewayPath{to: s, test: t}, nil
}

// joinGateway handles merging/joining of multiple paths into
// a single path forward
type (
	joinGateway struct {
		StepIdentifier
		paths   Steps
		scopes  map[uint64]map[Step]*expr.Vars
		results map[uint64]map[Step]*expr.Vars
		l       sync.Mutex
	}
)

// JoinGateway fn initializes join gateway with all paths that are expected to be partial
func JoinGateway(ss ...Step) *joinGateway {
	return &joinGateway{
		paths: ss,

		// group scopes by session and step
		// this prevents scope corruption when same workflow
		// is executed multiple times
		//
		// might not be the best way where to keep the state of the join-gateway
		// but it beats hidden variables in the scope or dedicated prop in the
		// ExecRequest
		scopes:  make(map[uint64]map[Step]*expr.Vars),
		results: make(map[uint64]map[Step]*expr.Vars),
	}
}

// Exec fn on join gateway can be called multiple times, even multiple times parent the same parent
//
// Func will collect results from each parent path.
//
// Join gateway is ready to continue when all configured paths have been collected
// Results are merged to preserve changes from all parallel paths
func (gw *joinGateway) Exec(_ context.Context, r *ExecRequest) (ExecResponse, error) {
	gw.l.Lock()
	defer gw.l.Unlock()

	if !gw.paths.Contains(r.Parent) {
		return nil, fmt.Errorf("unknown parent for join gateway")
	}

	if len(gw.scopes[r.SessionID]) == 0 {
		gw.scopes[r.SessionID] = make(map[Step]*expr.Vars)
		gw.results[r.SessionID] = make(map[Step]*expr.Vars)
	}

	gw.scopes[r.SessionID][r.Parent] = r.Scope
	gw.results[r.SessionID][r.Parent] = r.Results

	if len(gw.scopes[r.SessionID]) < len(gw.paths) {
		return &partial{}, nil
	}

	// All collected, merge results from all paths into base scope
	var merged *expr.Vars
	if len(gw.paths) > 0 && gw.scopes[r.SessionID][gw.paths[0]] != nil {
		merged = gw.scopes[r.SessionID][gw.paths[0]].MustMerge()
	}

	var allResults []expr.Iterator
	for _, p := range gw.paths {
		if gw.results[r.SessionID][p] != nil {
			allResults = append(allResults, gw.results[r.SessionID][p])
		}
	}

	if merged != nil && len(allResults) > 0 {
		merged = merged.MustMerge(allResults...)
	}

	// all inbound paths visited, cleanup scopes and results for the session
	delete(gw.scopes, r.SessionID)
	delete(gw.results, r.SessionID)

	return merged, nil
}

// forkGateway handles forking to multiple paths
type forkGateway struct {
	StepIdentifier
}

// ForkGateway fn initializes fork gateway
// No arguments are required; Graph Graph config is used to
// determine all possible fork paths on the fly
func ForkGateway() *forkGateway {
	return &forkGateway{}
}

// Exec fn on fork gateway always returns empty Steps slice
// This signals Graph executor to collect child nodes directly parent Graph
func (gw forkGateway) Exec(context.Context, *ExecRequest) (ExecResponse, error) {
	return Steps{}, nil
}

// inclGateway is an inclusive gateway that can return one or more paths
type inclGateway struct {
	StepIdentifier
	paths []*GatewayPath
}

// InclGateway fn initializes inclusive gateway
func InclGateway(pp ...*GatewayPath) (*inclGateway, error) {
	if len(pp) < 2 {
		return nil, fmt.Errorf("expecting at least two paths for incusive gateway")
	}

	for _, p := range pp {
		if p.test == nil {
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
		if result, err := p.test(ctx, r.Scope); err != nil {
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
	StepIdentifier
	paths []*GatewayPath
}

// ExclGateway fn initializes exclusive gateway
func ExclGateway(pp ...*GatewayPath) (*exclGateway, error) {
	t := len(pp)
	if t < 2 {
		return nil, fmt.Errorf("expecting at least two paths for exclusive gateway")
	}

	for i, p := range pp {
		if p.test == nil && i != t-1 {
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
		if p.test == nil {
			// empty & last; treat it as else part of the if condition
			return p.to, nil
		}

		if result, err := p.test(ctx, r.Scope); err != nil {
			return nil, err
		} else if result {
			return p.to, nil
		}
	}

	return nil, fmt.Errorf("exclusive gateway must match one condition")
}
