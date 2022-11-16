package qlng

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/spf13/cast"
)

type (
	ASTNode struct {
		pMeta *parserMeta

		Ref  string     `json:"ref,omitempty"`
		Args ASTNodeSet `json:"args,omitempty"`

		Symbol string      `json:"symbol,omitempty"`
		Value  *typedValue `json:"value,omitempty"`

		Raw string `json:"raw,omitempty"`
	}
	ASTNodeSet []*ASTNode

	typedValue struct {
		V expr.TypedValue
	}

	parserMeta struct {
		opDef *opDef
	}
)

func (n *ASTNode) String() string {
	// Leaf edge-cases
	switch {
	case n.Symbol != "":
		return n.Symbol
	case n.Value != nil:
		return "\"" + cast.ToString(n.Value.V.Get()) + "\""
	}

	// Process arguments for the op.
	args := make([]string, len(n.Args))
	for i, a := range n.Args {
		s := a.String()
		args[i] = s
	}

	// Default handlers
	return fmt.Sprintf("%s(%s)", n.Ref, strings.Join(args, ", "))
}

func MakeValueOf(t string, v interface{}) *typedValue {
	return &typedValue{
		V: expr.Must(qlTypeRegistry(t).Cast(v)),
	}
}

func WrapValue(v expr.TypedValue) *typedValue {
	return &typedValue{
		V: v,
	}
}

func (t *typedValue) UnmarshalJSON(in []byte) (err error) {
	var (
		aux = struct {
			Type  string      `json:"@type"`
			Value interface{} `json:"@value"`
		}{}
	)

	if len(in) == 0 {
		return nil
	}

	if err = json.Unmarshal(in, &aux); err != nil {
		return
	}

	if aux.Type == "" {
		return errors.New("invalid value definition: missing @type definition")
	}

	t.V, err = qlTypeRegistry(aux.Type).Cast(aux.Value)
	return
}

func (t *typedValue) MarshalJSON() ([]byte, error) {
	var (
		aux = struct {
			Type  string      `json:"@type"`
			Value interface{} `json:"@value"`
		}{}
	)

	if t.V == nil {
		return json.Marshal(aux)
	}

	aux.Type = t.V.Type()
	aux.Value = t.V.Get()

	switch aux.Type {
	case "ID", "Record", "User":
		v := aux.Value.(uint64)
		aux.Value = strconv.FormatUint(v, 10)
	}

	return json.Marshal(aux)
}

// Traverse traverses the AST down to leaf nodes.
//
// If fnc. returns false, the traversal of the current branch ends.
func (n *ASTNode) Traverse(f func(*ASTNode) (bool, *ASTNode, error)) (err error) {
	var ok bool
	var r *ASTNode
	if n == nil {
		return nil
	}

	ok, r, err = f(n)
	if err != nil {
		return err
	}
	*n = *r
	if !ok {
		return
	}

	for _, a := range n.Args {
		if err = a.Traverse(f); err != nil {
			return
		}
	}

	return
}

func (n ASTNode) Clone() *ASTNode {
	aa := n.Args

	if n.Args != nil {
		n.Args = make(ASTNodeSet, len(aa))
		for i, a := range aa {
			n.Args[i] = a.Clone()
		}
	}

	if n.Value != nil {
		n.Value = &typedValue{
			V: n.Value.V,
		}
	}

	return &n
}

func (n lNull) ToAST() (out *ASTNode) {
	return &ASTNode{
		Ref: "null",
	}
}

func (n lBoolean) ToAST() (out *ASTNode) {
	return &ASTNode{
		Value: MakeValueOf("Boolean", n.value),
	}
}

func (n lString) ToAST() (out *ASTNode) {
	return &ASTNode{
		Value: MakeValueOf("String", n.value),
	}
}

// @todo differentiate between floats and others
func (n lNumber) ToAST() (out *ASTNode) {
	if isFloaty(n.value) {
		return &ASTNode{
			Value: MakeValueOf("Float", n.value),
		}
	} else {
		return &ASTNode{
			Value: MakeValueOf("Integer", n.value),
		}
	}
}

func (n operator) ToAST() (out *ASTNode) {
	op := getOp(n.kind)

	return &ASTNode{
		Ref:   op.name,
		Args:  make(ASTNodeSet, 0, 2),
		pMeta: &parserMeta{opDef: op},
	}
}

func (n Ident) ToAST() (out *ASTNode) {
	return &ASTNode{
		Symbol: n.Value,
	}
}

func (n keyword) ToAST() (out *ASTNode) {
	return &ASTNode{
		Ref: n.keyword,
	}
}

func (n interval) ToAST() (out *ASTNode) {
	return &ASTNode{
		Ref: "interval",
		Args: ASTNodeSet{
			{Symbol: n.unit},
			{Value: MakeValueOf("Number", n.value)},
		},
	}
}

func (n function) ToAST() (out *ASTNode) {
	auxA := n.arguments.ToAST()

	return &ASTNode{
		Ref:  n.name,
		Args: auxA.Args,
	}
}

func (nn parserNodeSet) ToAST() (out *ASTNode) {
	auxArgs := make(ASTNodeSet, 0, len(nn))

	for _, n := range nn {
		auxArgs = append(auxArgs, n.ToAST())
	}

	return &ASTNode{
		Ref:  "group",
		Args: auxArgs,
	}
}

func (nn parserNodes) ToAST() (out *ASTNode) {
	// Prep
	auxArgs := make(ASTNodeSet, 0, len(nn))

	// Convert the entire level to AST nodes
	for _, n := range nn {
		auxArgs = append(auxArgs, n.ToAST())
	}

	// In the current level, have the operators consume operands
	// based on their defined weight
	//
	// - find the highest prio. op and have it consume what it wants
	// - repeat until all ops are satisfied
	// -- post optimizations?
	for {
		var bestOp *opDef
		bestOpIx := -1

		// We're done when it is reduced to 1
		if len(auxArgs) <= 1 {
			break
		}

		for _i, _a := range auxArgs {
			i := _i
			a := _a

			// use this as a delimiter for now
			if a.pMeta == nil || a.pMeta.opDef == nil {
				continue
			}

			if bestOp == nil {
				bestOp = a.pMeta.opDef
				bestOpIx = i
				continue
			}

			if a.pMeta.opDef.weight < bestOp.weight {
				bestOp = a.pMeta.opDef
				bestOpIx = i
				continue
			}
		}

		// Have the op consume what it needs.
		arg := auxArgs[bestOpIx]
		if !isUnary(arg.Ref) {
			arg.Args = append(arg.Args, auxArgs[bestOpIx-1], auxArgs[bestOpIx+1])
			// this is not needed anymore so we can remove it
			arg.pMeta = nil

			// Remove the consumed bits and replace it with the new bit
			aux := auxArgs[0 : bestOpIx-1]
			aux = append(aux, arg)
			// +1 for right side, +1 because the left index is inclusive
			aux = append(aux, auxArgs[bestOpIx+2:]...)
			auxArgs = aux
		} else {
			arg.Args = append(arg.Args, auxArgs[bestOpIx+1])
			// this is not needed anymore so we can remove it
			arg.pMeta = nil

			// Remove the consumed bits and replace it with the new bit
			aux := auxArgs[0:bestOpIx]
			aux = append(aux, arg)
			// +1 because the left index is inclusive
			aux = append(aux, auxArgs[bestOpIx+2:]...)
			auxArgs = aux
		}
	}

	if len(auxArgs) > 1 {
		return &ASTNode{
			Ref:  "group",
			Args: auxArgs,
		}
	}

	return auxArgs[0]
}

// A simplified type registry for the types that QL needs to understand
func qlTypeRegistry(ref string) expr.Type {
	switch ref {
	case "ID", "Record", "User":
		return &expr.ID{}
	case "Boolean", "Bool":
		return &expr.Boolean{}
	case "Integer":
		return &expr.Integer{}
	case "UnsignedInteger":
		return &expr.UnsignedInteger{}
	case "Float", "Number":
		return &expr.Float{}
	case "String", "Select":
		return &expr.String{}
	case "DateTime":
		return &expr.DateTime{}
	}

	return nil
}
