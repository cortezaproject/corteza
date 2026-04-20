package types

import "github.com/crusttech/human/server/pkg/ast"

// GatewayConditionsSetter is implemented by both gateway handler types.
// The converter uses this to populate the conditions slice during Phase 3.
type GatewayConditionsSetter interface {
	SetConditions([]*ast.ASTNode)
}
