package types

import (
	"context"
	"fmt"

	nx "github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/ast"
	"github.com/cortezaproject/corteza/server/pkg/expr"
)

// exclusiveGateway evaluates conditions in order; returns the index of the first
// matching condition. A nil condition matches unconditionally (else path).
// If no condition matches, falls back to the last child.
type exclusiveGateway struct {
	conditions []*ast.ASTNode
}

// Select implements runtime.GatewayHandler for exclusive gateways.
func (g *exclusiveGateway) Select(_ context.Context, scope map[string]*expr.Vars) ([]int, error) {
	for i, cond := range g.conditions {
		if cond == nil {
			return []int{i}, nil
		}
		ok, err := ast.EvalBool(cond, scope)
		if err != nil {
			return nil, fmt.Errorf("exclusive gateway condition[%d]: %w", i, err)
		}
		if ok {
			return []int{i}, nil
		}
	}
	// Fallback: last child (else)
	if len(g.conditions) > 0 {
		return []int{len(g.conditions) - 1}, nil
	}
	return nil, nil
}

// SetConditions updates the conditions slice after handler creation.
func (g *exclusiveGateway) SetConditions(conds []*ast.ASTNode) { g.conditions = conds }

// ExecN implements execTypes.StepHandler.
func (g *exclusiveGateway) ExecN(_ context.Context, _ *nx.ExecRequest) (nx.ExecResponse, error) {
	out, err := expr.NewVars(nil)
	return out, err
}

// ExclusiveGatewayStep returns a StepHandler for an exclusive (XOR) gateway.
// conditions must be in child order; nil = else/default path.
func ExclusiveGatewayStep(conditions []*ast.ASTNode) nx.StepHandler {
	return &exclusiveGateway{conditions: conditions}
}
