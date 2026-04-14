package ast

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/spf13/cast"
)

type (
	ASTNode struct {
		// Meta holds additional node data such as scope references
		Meta map[string]any `json:"meta,omitempty"`

		Ref  string     `json:"ref,omitempty"`
		Args ASTNodeSet `json:"args,omitempty"`

		Symbol string      `json:"symbol,omitempty"`
		Value  *TypedValue `json:"value,omitempty"`

		Raw string `json:"raw,omitempty"`
	}
	ASTNodeSet []*ASTNode

	TypedValue struct {
		V expr.TypedValue
	}
)

func (n *ASTNode) String() string {
	if n == nil {
		return "<nil>"
	}

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

func MakeValueOf(t string, v interface{}) *TypedValue {
	return &TypedValue{
		V: expr.Must(TypeRegistry(t).Cast(v)),
	}
}

func WrapValue(v expr.TypedValue) *TypedValue {
	return &TypedValue{
		V: v,
	}
}

func (t *TypedValue) UnmarshalJSON(in []byte) (err error) {
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
		return fmt.Errorf("invalid value definition: missing @type definition")
	}

	t.V, err = TypeRegistry(aux.Type).Cast(aux.Value)
	return
}

func (t *TypedValue) MarshalJSON() ([]byte, error) {
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

func (n *ASTNode) CollectSymbols() (out []string) {
	out = make([]string, 0, 8)

	n.Traverse(func(a *ASTNode) (bool, *ASTNode, error) {
		if a.Symbol != "" {
			out = append(out, a.Symbol)
		}

		return true, a, nil
	})

	return
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
		n.Value = &TypedValue{
			V: n.Value.V,
		}
	}

	return &n
}

// TypeRegistry is a simplified type registry for the types that AST nodes need to understand.
func TypeRegistry(ref string) expr.Type {
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
