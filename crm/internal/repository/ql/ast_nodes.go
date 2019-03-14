package ql

import (
	"fmt"

	"gopkg.in/Masterminds/squirrel.v1"
)

// SelectStatement represents a SQL SELECT statement.
type (
	ASTNode interface {
		fmt.Stringer
		squirrel.Sqlizer
	}

	ASTSet   []ASTNode // Stream of comma delimited nodes
	ASTNodes []ASTNode // Stream of un-delimited nodes

	Columns []Column

	String struct {
		Value string
		Args  []interface{}
	}

	Number struct {
		Value string
	}

	Operator struct {
		Kind string
	}

	Ident struct {
		Value string
		Args  []interface{}
	}

	Keyword struct {
		Keyword string
	}

	Interval struct {
		Value string
		Unit  string
	}

	Column struct {
		Expr  ASTNodes
		Alias string
	}
	Function struct {
		Name      string
		Arguments ASTSet
	}
)

func (n String) String() string { return fmt.Sprintf("%q", n.Value) }

func (n Number) String() string { return n.Value }

func (n Operator) String() string { return n.Kind }

func (n Keyword) String() string { return n.Keyword }

func (n Interval) String() string { return fmt.Sprintf("INTERVAL %s %s", n.Value, n.Unit) }

func (n Function) String() string { return fmt.Sprintf("%s(%s)", n.Name, n.Arguments) }

func (n Ident) String() string { return n.Value }

func (n Column) String() (out string) {
	out = n.Expr.String()
	if n.Alias != "" {
		out = out + " AS " + n.Alias
	}

	return
}

func (nn ASTNodes) String() (out string) {
	for _, n := range nn {
		out = out + n.String()
	}

	return
}

func (nn ASTSet) String() (out string) {
	for i, n := range nn {
		if i > 0 {
			out = out + ", "
		}
		out = out + n.String()
	}

	return
}

func (nn Columns) String() (out string) {
	for i, n := range nn {
		if i > 0 {
			out = out + ", "
		}
		out = out + n.String()
	}

	return
}

func (nn Columns) Strings() (out []string) {
	out = make([]string, len(nn))
	for i, n := range nn {
		out[i] = n.String()
	}

	return
}
