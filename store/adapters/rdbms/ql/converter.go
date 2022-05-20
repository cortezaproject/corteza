package ql

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/ql"
	"github.com/doug-martin/goqu/v9/exp"
)

type (
	refHandler func(*ql.ASTNode, ...exp.Expression) (exp.Expression, error)
	symHandler func(*ql.ASTNode) (exp.Expression, error)
	valHandler func(*ql.ASTNode) (exp.Expression, error)

	// converts query and AST from the parsed query into goqu expression
	// that can be used in the SQL
	converter struct {
		parser     *ql.Parser
		refHandler refHandler
		symHandler symHandler
		valHandler valHandler
	}

	op func(*converter)
)

// Initializes new converter
func Converter(oo ...op) *converter {
	c := &converter{
		parser:     ql.NewParser(),
		refHandler: DefaultRefHandler,
		symHandler: DefaultSymbolHandler,
		valHandler: DefaultValueHandler,
	}

	for _, o := range oo {
		o(c)
	}

	return c
}

// RefHandler sets custom ref handler
func RefHandler(h refHandler) op {
	return func(c *converter) {
		c.refHandler = h
	}
}

// SymHandler sets custom symbol handler
func SymHandler(h symHandler) op {
	return func(c *converter) {
		c.symHandler = h
	}
}

// ValHandler sets custom value handler
func ValHandler(h valHandler) op {
	return func(c *converter) {
		c.valHandler = h
	}
}

func (c *converter) Parse(q string) (exp.Expression, error) {
	n, err := c.parser.Parse(q)
	if err != nil {
		return nil, err
	}

	return c.Convert(n)
}

// Convert recursively walks through AST and converts the entire tree into expression
// that can be used in goqu SQL builder
func (c *converter) Convert(n *ql.ASTNode) (_ exp.Expression, err error) {
	switch {
	case n.Symbol != "":
		return c.symHandler(n)
	case n.Value != nil:
		return c.valHandler(n)
	}

	args := make([]exp.Expression, len(n.Args))
	for i, a := range n.Args {
		args[i], err = c.Convert(a)
		if err != nil {
			return
		}
	}

	return c.refHandler(n, args...)
}

// DefaultRefHandler converts ref from the AST node using ref2exp
func DefaultRefHandler(n *ql.ASTNode, args ...exp.Expression) (exp.Expression, error) {
	if ref2exp[n.Ref] == nil {
		return nil, fmt.Errorf("unknown ref %q", n.Ref)
	}

	return ref2exp[n.Ref].Handler(args...), nil
}

// DefaultSymbolHandler parses symbol from the AST node into an identifier
func DefaultSymbolHandler(n *ql.ASTNode) (exp.Expression, error) {
	return exp.ParseIdentifier(n.Symbol), nil
}

// DefaultValueHandler converts node into placeholder and a new attribute
func DefaultValueHandler(n *ql.ASTNode) (exp.Expression, error) {
	return exp.NewLiteralExpression("?", n.Value.V.Get()), nil
}
