package types

import (
	"context"
	"fmt"

	nx "github.com/cortezaproject/corteza/server/pkg/automation_exec/types"
	"github.com/cortezaproject/corteza/server/pkg/ast"
	"github.com/cortezaproject/corteza/server/pkg/expr"
)

// inclusiveGateway evaluates all conditions; activates every matching child as a
// parallel branch. A nil condition matches unconditionally (else path).
// If nothing matches, falls back to the last child.
type inclusiveGateway struct {
	conditions []*ast.ASTNode
}

// Select implements runtime.GatewayHandler for inclusive gateways.
func (g *inclusiveGateway) Select(_ context.Context, scope map[string]*expr.Vars) ([]int, error) {
	var matched []int
	for i, cond := range g.conditions {
		if cond == nil {
			matched = append(matched, i)
			continue
		}
		ok, err := ast.EvalBool(cond, scope)
		if err != nil {
			return nil, fmt.Errorf("inclusive gateway condition[%d]: %w", i, err)
		}
		if ok {
			matched = append(matched, i)
		}
	}
	if len(matched) == 0 && len(g.conditions) > 0 {
		// Fallback: last child (else)
		return []int{len(g.conditions) - 1}, nil
	}
	return matched, nil
}

// SetConditions updates the conditions slice after handler creation.
func (g *inclusiveGateway) SetConditions(conds []*ast.ASTNode) { g.conditions = conds }

// ExecN implements execTypes.StepHandler.
func (g *inclusiveGateway) ExecN(_ context.Context, _ *nx.ExecRequest) (nx.ExecResponse, error) {
	out, err := expr.NewVars(nil)
	return out, err
}

// InclusiveGatewayStep returns a StepHandler for an inclusive (OR) gateway.
// conditions must be in child order; nil = else/default path.
func InclusiveGatewayStep(conditions []*ast.ASTNode) nx.StepHandler {
	return &inclusiveGateway{conditions: conditions}
}
