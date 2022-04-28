package report

import (
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/ql"
)

type (
	HandlerSig func(vv ...expr.TypedValue) expr.TypedValue

	argSet []*arg
	arg    struct {
		Required bool
		Type     string
	}
)

var (
	handlers = map[string]HandlerSig{
		// operators
		// - bool
		"and": andHandler,
		"or":  orHandler,

		// - comp.
		"eq": makeCmpHandler(0),
		"lt": makeCmpHandler(-1),
		"gt": makeCmpHandler(1),

		"is": existenceHandler,

		// generic stuff
		"null":  nullHandler,
		"nnull": notNullHandler,
	}
)

// eval evaluates the given AST over the provided frame row
//
// This is a simplified bool only implementation as nothing else is needed for now.
func (d *joinedDataset) eval(n *ql.ASTNode, row FrameRow, cc FrameColumnSet) bool {
	if v, ok := d.evalRec(n, true, row, cc).(*expr.Boolean); !ok {
		return false
	} else {
		return v.GetValue()
	}
}

func (d *joinedDataset) evalRec(n *ql.ASTNode, isRoot bool, row FrameRow, cc FrameColumnSet) expr.TypedValue {
	// Leaf edge-cases
	switch {
	case n.Symbol != "":
		return row[cc.Find(n.Symbol)]
	case n.Value != nil:
		return n.Value.V
	}

	// Process arguments for the op.
	args := make([]expr.TypedValue, len(n.Args))
	for i, a := range n.Args {
		s := d.evalRec(a, false, row, cc)
		args[i] = s
	}

	// Default handlers
	return handlers[n.Ref](args...)
}

func andHandler(aa ...expr.TypedValue) expr.TypedValue {
	for _, a := range aa {
		if v, ok := a.(*expr.Boolean); !ok || !v.GetValue() {
			return expr.Must(expr.NewBoolean(false))
		}
	}
	return expr.Must(expr.NewBoolean(true))
}

func orHandler(aa ...expr.TypedValue) expr.TypedValue {
	for _, a := range aa {
		if v, ok := a.(*expr.Boolean); ok && v.GetValue() {
			return expr.Must(expr.NewBoolean(true))
		}
	}
	return expr.Must(expr.NewBoolean(false))
}

func makeCmpHandler(val int) HandlerSig {
	return func(aa ...expr.TypedValue) expr.TypedValue {
		a := aa[0].(expr.Comparable)
		b := aa[1]

		c, _ := a.Compare(b)
		return expr.Must(expr.NewBoolean(c == val))
	}
}

func nullHandler(_ ...expr.TypedValue) expr.TypedValue {
	return nil
}

func notNullHandler(_ ...expr.TypedValue) expr.TypedValue {
	return expr.Must(expr.NewAny(nil))
}

func existenceHandler(aa ...expr.TypedValue) expr.TypedValue {
	a := isNil(aa[0])
	b := isNil(aa[1])

	return expr.Must(expr.NewBoolean(a == b))
}
