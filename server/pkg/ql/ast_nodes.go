package ql

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
)

// SelectStatement represents a SQL SELECT statement.
type (
	ASTNode interface {
		fmt.Stringer
		squirrel.Sqlizer

		Validate() error
	}

	replacer func(string) string

	ASTSet   []ASTNode // Stream of comma delimited nodes
	ASTNodes []ASTNode // Stream of space delimited nodes

	Columns []Column

	LNull    struct{}
	LBoolean struct{ Value bool }

	LString struct {
		Value string
		Args  []interface{}
	}

	LNumber struct {
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

	NodeF struct {
		Expr      string
		Arguments []ASTNode
		replacer  replacer
	}
)

var (
	operators = map[string]bool{
		`=`:        true,
		`!=`:       true,
		`<`:        true,
		`>`:        true,
		`>=`:       true,
		`<=`:       true,
		`<>`:       true,
		`+`:        true,
		`-`:        true,
		`*`:        true,
		`/`:        true,
		`AND`:      true,
		`OR`:       true,
		`XOR`:      true,
		`NOT LIKE`: true,
		`LIKE`:     true,
		`IS`:       true,
		`IS NOT`:   true,
	}
)

func (n LNull) Validate() (err error) { return }
func (n LNull) String() string        { return "NULL" }

func (n LBoolean) Validate() (err error) { return }
func (n LBoolean) String() string {
	if n.Value {
		return "TRUE"
	} else {
		return "FALSE"
	}
}

func (n LString) Validate() (err error) { return }
func (n LString) String() string        { return fmt.Sprintf("%q", n.Value) }

func (n LNumber) Validate() (err error) { return }
func (n LNumber) String() string        { return n.Value }

func (n Operator) Validate() (err error) {
	if !operators[strings.ToUpper(n.Kind)] {
		return fmt.Errorf("unknown operator '%s'", n.Kind)
	}
	return
}
func (n Operator) String() string { return n.Kind }

func (n Keyword) Validate() (err error) { return }
func (n Keyword) String() string        { return n.Keyword }

func (n Interval) Validate() (err error) { return }
func (n Interval) String() string        { return fmt.Sprintf("INTERVAL %s %s", n.Value, n.Unit) }

func (n Function) Validate() (err error) { return }
func (n Function) String() string        { return fmt.Sprintf("%s(%s)", n.Name, n.Arguments) }

func (n Ident) Validate() (err error) { return }
func (n Ident) String() string        { return n.Value }

func (n Column) Validate() (err error) {
	return n.Expr.Validate()
}

func (n Column) String() (out string) {
	out = n.Expr.String()
	if n.Alias != "" {
		out = out + " AS " + n.Alias
	}

	return
}

func (nn ASTNodes) Validate() (err error) {
	if err = validate(nn); err != nil {
		return
	}

	l := len(nn)
	if l == 0 {
		return fmt.Errorf("empty set")
	}

	if op, ok := nn[0].(Operator); ok {
		return fmt.Errorf("malformed expression, unexpected operator '%s' at first node", op)
	}

	if l > 1 {
		if op, ok := nn[l-1].(Operator); ok {
			return fmt.Errorf("malformed expression, unexpected operator '%s' at last node", op)
		}
	}

	return
}

func (nn ASTNodes) String() (out string) {
	for i, n := range nn {
		if i > 0 {
			out = out + "  "
		}

		out = out + n.String()
	}

	return
}

func (nn ASTSet) Validate() (err error) {
	return validate(nn)
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

func (nn Columns) Validate() (err error) {
	for _, n := range nn {
		if err = n.Validate(); err != nil {
			return
		}
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

// MakeReplacedFormattedNode also accepts the replacer to apply to the arguments
func MakeReplacedFormattedNode(expr string, r replacer, nn ...ASTNode) *NodeF {
	return &NodeF{Expr: expr, Arguments: nn, replacer: r}
}

func MakeFormattedNode(expr string, nn ...ASTNode) *NodeF {
	return &NodeF{Expr: expr, Arguments: nn}
}

func (n NodeF) Validate() (err error) { return }
func (n NodeF) String() string        { return n.Expr }

func validate(nn []ASTNode) (err error) {
	if len(nn) == 0 {
		return fmt.Errorf("empty node set")
	}

	for _, n := range nn {
		if err = n.Validate(); err != nil {
			return
		}
	}

	return
}
