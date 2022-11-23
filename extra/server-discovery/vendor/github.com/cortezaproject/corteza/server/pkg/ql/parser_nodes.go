package ql

import (
	"fmt"
	"strings"
)

// SelectStatement represents a SQL SELECT statement.
type (
	parserNode interface {
		fmt.Stringer

		Validate() error
		ToAST() *ASTNode
	}

	parserNodeSet []parserNode // Stream of comma delimited nodes
	parserNodes   []parserNode // Stream of space delimited nodes

	lNull    struct{}
	lBoolean struct{ value bool }

	lString struct {
		value string
		args  []interface{}
	}

	lNumber struct {
		value string
	}

	operator struct {
		kind string
	}

	Ident struct {
		Value string
		args  []interface{}
	}

	keyword struct {
		keyword string
	}

	interval struct {
		value string
		unit  string
	}

	function struct {
		name      string
		arguments parserNodeSet
	}

	opDef struct {
		name   string
		weight int
	}
)

var (
	ops = map[string]opDef{
		// generic comparison
		`=`:      {name: `eq`, weight: 40},
		`==`:     {name: `eq`, weight: 40},
		`===`:    {name: `eq`, weight: 40},
		`!=`:     {name: `ne`, weight: 40},
		`!==`:    {name: `ne`, weight: 40},
		`<>`:     {name: `ne`, weight: 40},
		`<`:      {name: `lt`, weight: 30},
		`<=`:     {name: `le`, weight: 30},
		`>`:      {name: `gt`, weight: 30},
		`>=`:     {name: `ge`, weight: 30},
		`IS`:     {name: `is`, weight: 40},
		`IS NOT`: {name: `nis`, weight: 40},
		`IN`:     {name: `in`, weight: 40},
		`NOT IN`: {name: `nin`, weight: 40},

		// conjunction
		`AND`: {name: `and`, weight: 50},
		`&&`:  {name: `and`, weight: 50},
		`OR`:  {name: `or`, weight: 60},
		`||`:  {name: `or`, weight: 60},
		`XOR`: {name: `xor`, weight: 60},

		// math
		`+`: {name: `add`, weight: 20},
		`-`: {name: `sub`, weight: 20},
		`*`: {name: `mult`, weight: 10},
		`/`: {name: `div`, weight: 10},

		// modifiers
		`!`: {name: `not`, weight: 0},

		// str comp.
		`LIKE`:     {name: `like`, weight: 40},
		`NOT LIKE`: {name: `nlike`, weight: 40},
	}
)

func isUnary(s string) bool {
	return s == "!" || s == "not"
}

func getOp(op string) *opDef {
	o, ok := ops[strings.ToUpper(op)]
	if !ok {
		return nil
	}
	return &o
}

func (n lNull) Validate() (err error) { return }
func (n lNull) String() string        { return "NULL" }

func (n lBoolean) Validate() (err error) { return }
func (n lBoolean) String() string {
	if n.value {
		return "TRUE"
	} else {
		return "FALSE"
	}
}

func (n lString) Validate() (err error) { return }
func (n lString) String() string        { return fmt.Sprintf("%q", n.value) }

func (n lNumber) Validate() (err error) { return }
func (n lNumber) String() string        { return n.value }

func (n operator) Validate() (err error) {
	if getOp(n.kind) == nil {
		return fmt.Errorf("unknown operator '%s'", n.kind)
	}
	return
}
func (n operator) String() string { return n.kind }

func (n keyword) Validate() (err error) { return }
func (n keyword) String() string        { return n.keyword }

func (n interval) Validate() (err error) { return }
func (n interval) String() string        { return fmt.Sprintf("INTERVAL %s %s", n.value, n.unit) }

func (n function) Validate() (err error) { return }
func (n function) String() string        { return fmt.Sprintf("%s(%s)", n.name, n.arguments) }

func (n Ident) Validate() (err error) { return }
func (n Ident) String() string        { return n.Value }

func (nn parserNodes) Validate() (err error) {
	if err = validate(nn); err != nil {
		return
	}

	l := len(nn)
	if l == 0 {
		return fmt.Errorf("empty set")
	}

	if op, ok := nn[0].(operator); ok && !isUnary(op.kind) {
		return fmt.Errorf("malformed expression, unexpected operator '%s' at first node", op)
	}

	if l > 1 {
		if op, ok := nn[l-1].(operator); ok {
			return fmt.Errorf("malformed expression, unexpected operator '%s' at last node", op)
		}
	}

	return
}

func (nn parserNodes) String() (out string) {
	for i, n := range nn {
		if i > 0 {
			out = out + "  "
		}

		out = out + n.String()
	}

	return
}

func (nn parserNodeSet) Validate() (err error) {
	return validate(nn)
}

func (nn parserNodeSet) String() (out string) {
	for i, n := range nn {
		if i > 0 {
			out = out + ", "
		}
		out = out + n.String()
	}

	return
}

func validate(nn []parserNode) (err error) {
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
