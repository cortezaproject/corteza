package ql

import (
	"github.com/crusttech/human/server/pkg/ast"
	"github.com/crusttech/human/server/pkg/expr"
)

func (n lNull) ToAST() *ast.ASTNode {
	return &ast.ASTNode{
		Ref: "null",
	}
}

func (n lBoolean) ToAST() *ast.ASTNode {
	return &ast.ASTNode{
		Value: ast.MakeValueOf("Boolean", n.value),
	}
}

func (n lString) ToAST() *ast.ASTNode {
	return &ast.ASTNode{
		Value: ast.MakeValueOf("String", n.value),
	}
}

// @todo differentiate between floats and others
func (n lNumber) ToAST() *ast.ASTNode {
	if isFloaty(n.value) {
		return &ast.ASTNode{
			Value: ast.MakeValueOf("Float", n.value),
		}
	}
	return &ast.ASTNode{
		Value: ast.MakeValueOf("Integer", n.value),
	}
}

func (n operator) ToAST() *ast.ASTNode {
	op := getOp(n.kind)

	node := &ast.ASTNode{
		Ref:  op.name,
		Args: make(ast.ASTNodeSet, 0, 2),
	}
	// Store opDef in Meta for use by parserNodes.ToAST
	node.Meta = map[string]any{"opDef": op}
	return node
}

func (n Ident) ToAST() *ast.ASTNode {
	return &ast.ASTNode{
		Symbol: n.Value,
	}
}

func (n keyword) ToAST() *ast.ASTNode {
	return &ast.ASTNode{
		Ref: n.keyword,
	}
}

func (n interval) ToAST() *ast.ASTNode {
	return &ast.ASTNode{
		Ref: "interval",
		Args: ast.ASTNodeSet{
			// @todo consider introducing a new concept on the ASTNode to cover
			//       system defined symbols/keywords. This is a temporary solutions
			//       to fix reporting interval definitions; it should eventually be reworked
			{Value: ast.MakeValueOf("String", n.unit)},
			{Value: ast.MakeValueOf("Number", n.value)},
		},
	}
}

func (n function) ToAST() *ast.ASTNode {
	auxA := n.arguments.ToAST()

	return &ast.ASTNode{
		Ref:  n.name,
		Args: auxA.Args,
	}
}

func (nn parserNodeSet) ToAST() *ast.ASTNode {
	auxArgs := make(ast.ASTNodeSet, 0, len(nn))

	for _, n := range nn {
		auxArgs = append(auxArgs, n.ToAST())
	}

	return &ast.ASTNode{
		Ref:  "group",
		Args: auxArgs,
	}
}

func (nn parserNodes) ToAST() *ast.ASTNode {
	auxArgs := make(ast.ASTNodeSet, 0, len(nn))

	for _, n := range nn {
		auxArgs = append(auxArgs, n.ToAST())
	}

	for {
		var bestOp *opDef
		bestOpIx := -1

		if len(auxArgs) <= 1 {
			break
		}

		for _i, _a := range auxArgs {
			i := _i
			a := _a

			op, hasOp := opDefFromMeta(a)
			if !hasOp {
				continue
			}

			if bestOp == nil {
				bestOp = op
				bestOpIx = i
				continue
			}

			if op.weight < bestOp.weight {
				bestOp = op
				bestOpIx = i
				continue
			}
		}

		if bestOpIx < 0 {
			// Nothing left to reduce: the remaining nodes cannot be folded into
			// a single expression. parserNodes.Validate rejects the token
			// sequences that get here, so this is a guard against a reduction
			// bug rather than bad input — fall through to the group below,
			// because indexing with -1 would panic on a user-supplied query.
			break
		}

		arg := auxArgs[bestOpIx]
		if !isUnary(arg.Ref) {
			skip := 2
			arg.Args = append(arg.Args, auxArgs[bestOpIx-1], auxArgs[bestOpIx+1])
			if len(auxArgs) > bestOpIx+2 {
				// Only absorb a third argument when what follows is an operand.
				// Testing Ref against isOperator missed every comparison
				// operator (it knows only and/or/xor), so "a-b = 'c'" swallowed
				// its own '=' and left an unreducible pair behind.
				if _, isOp := opDefFromMeta(auxArgs[bestOpIx+2]); !isOp {
					skip = 3
					arg.Args = append(arg.Args, auxArgs[bestOpIx+2])
				}
			}

			arg.Meta = nil

			aux := auxArgs[0 : bestOpIx-1]
			aux = append(aux, arg)
			aux = append(aux, auxArgs[bestOpIx+skip:]...)
			auxArgs = aux
		} else {
			arg.Args = append(arg.Args, auxArgs[bestOpIx+1])
			arg.Meta = nil

			aux := auxArgs[0:bestOpIx]
			aux = append(aux, arg)
			aux = append(aux, auxArgs[bestOpIx+2:]...)
			auxArgs = aux
		}
	}

	if len(auxArgs) > 1 {
		return &ast.ASTNode{
			Ref:  "group",
			Args: auxArgs,
		}
	}

	return auxArgs[0]
}

// opDefFromMeta reads the opDef stored in a node's Meta map (set by operator.ToAST).
func opDefFromMeta(n *ast.ASTNode) (*opDef, bool) {
	if n.Meta == nil {
		return nil, false
	}
	op, ok := n.Meta["opDef"].(*opDef)
	return op, ok
}

// qlTypeRegistry delegates to ast.TypeRegistry.
func qlTypeRegistry(ref string) expr.Type {
	return ast.TypeRegistry(ref)
}
