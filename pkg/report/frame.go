package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/spf13/cast"
)

type (
	Frame struct {
		Name   string `json:"name"`
		Source string `json:"source"`
		Ref    string `json:"ref,omitempty"`

		RefValue  string `json:"refValue,omitempty"`
		RelColumn string `json:"relColumn,omitempty"`

		Columns FrameColumnSet `json:"columns"`
		Rows    FrameRowSet    `json:"rows"`

		Paging  *filter.Paging  `json:"paging"`
		Sorting *filter.Sorting `json:"sorting"`

		// params to help us perform things in place
		startIndex int
		size       int
		sliced     bool
	}

	FrameRowSet []FrameRow
	FrameRow    []expr.TypedValue

	frameCellCaster func(in interface{}) (expr.TypedValue, error)
	FrameColumnSet  []*FrameColumn
	FrameColumn     struct {
		Name   string          `json:"name"`
		Label  string          `json:"label"`
		Kind   string          `json:"kind"`
		Caster frameCellCaster `json:"-"`
	}

	FrameDefinitionSet []*FrameDefinition
	FrameDefinition    struct {
		Name    string         `json:"name"`
		Source  string         `json:"source"`
		Ref     string         `json:"ref"`
		Rows    *RowDefinition `json:"rows"`
		Columns FrameColumnSet `json:"columns"`

		Paging  *filter.Paging     `json:"paging"`
		Sorting filter.SortExprSet `json:"sorting"`
	}
)

const (
	columnWildcard = "*"
)

func MakeColumnOfKind(k string) *FrameColumn {
	return &FrameColumn{
		Kind: k,
		Caster: func(in interface{}) (expr.TypedValue, error) {
			switch k {
			case "Number":
				return expr.NewFloat(in)
			case "DateTime":
				return expr.NewDateTime(in)
			case "User",
				"Record":
				return expr.NewID(in)
			case "Checkbox":
				return expr.NewBoolean(in)
			default:
				return expr.NewString(in)
			}
		},
	}
}

func KindOf(v expr.TypedValue) string {
	// @todo ...
	if v == nil {
		return "String"
	}

	switch v.Type() {
	case "Integer",
		"UnsignedInteger",
		"Float":
		return "Number"
	case "DateTime":
		return "DateTime"
	case "ID":
		return "Ref"
	case "Boolean":
		return "Checkbox"
	default:
		return "String"
	}
}

func (b *CellDefinition) UnmarshalJSON(data []byte) (err error) {
	if b == nil {
		*b = *(&CellDefinition{})
	}

	aux := make(map[string]string)
	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}

	for op, val := range aux {
		b.Value = val
		b.Op = op
	}

	return nil
}

func (b CellDefinition) OpToCmp() string {
	switch b.Op {
	case "eq":
		return "="
	case "ne":
		return "!="
	case "lt":
		return "<"
	case "gt":
		return ">"
	case "le":
		return "<="
	case "ge":
		return ">="

	default:
		return "="
	}
}

// Slice in place
func (f *Frame) Slice(startIndex, size int) (a, b *Frame) {
	a = &Frame{
		Name:      f.Name,
		Source:    f.Source,
		Ref:       f.Ref,
		RefValue:  f.RefValue,
		RelColumn: f.RelColumn,

		Columns: f.Columns,
	}
	b = &Frame{
		Name:      f.Name,
		Source:    f.Source,
		Ref:       f.Ref,
		RefValue:  f.RefValue,
		RelColumn: f.RelColumn,

		Columns: f.Columns,
	}

	a.Rows = f.Rows[startIndex:size]
	b.Rows = f.Rows[size:]
	return a, b
}

// With guard element
func (f *Frame) WalkRowsG(cb func(i int, r FrameRow) error) (err error) {
	err = f.WalkRows(cb)
	if err != nil {
		return err
	}

	return cb(f.Size(), nil)
}

func (f *Frame) WalkRows(cb func(i int, r FrameRow) error) (err error) {
	for i := 0; i < f.Size(); i++ {
		if err = cb(i, f.Rows[i]); err != nil {
			return err
		}
	}

	return nil
}

func (f *Frame) WalkRowsR(cb func(i int, r FrameRow) error) (err error) {
	for i := f.Size(); i >= 0; i-- {
		if err = cb(i, f.Rows[i]); err != nil {
			return err
		}
	}

	return nil
}

func (f *Frame) PeekRow(i int) FrameRow {
	return f.Rows[i]
}

func (f *Frame) PeekRowSafe(i int) FrameRow {
	if i >= f.Size() {
		return nil
	}

	return f.Rows[i]
}

func (f *Frame) Size() int {
	return len(f.Rows)
}

func (f *Frame) FirstRow() FrameRow {
	if f.Size() == 0 {
		return nil
	}

	return f.Rows[0]
}

func (f *Frame) LastRow() FrameRow {
	if f.Size() == 0 {
		return nil
	}

	return f.Rows[f.Size()-1]
}

// @todo nicer formatting and alignment
func (f *Frame) String() string {
	if f == nil {
		return "<NIL>"
	}
	out := fmt.Sprintf("n: %10s; src: %10s\n", f.Name, f.Source)

	if f.Ref != "" {
		out += fmt.Sprintf("ref: %10s; col: %10s\n; key: %10s\n", f.Ref, f.RelColumn, f.RefValue)
	}

	for _, c := range f.Columns {
		out += fmt.Sprintf("%s<%s>, ", c.Name, c.Kind)
	}
	out = strings.TrimRight(out, " ,")
	out += "\n"

	f.WalkRows(func(i int, r FrameRow) error {
		out += fmt.Sprintf("%d| ", i+1)
		for _, c := range r {
			if c == nil {
				out += "<N/A>, "
			} else {
				v := cast.ToString(c.Get())
				out += fmt.Sprintf("%s, ", v)
			}
		}
		out = strings.TrimRight(out, ", ") + "\n"
		return nil
	})

	if f.Paging != nil {
		out += "\n"
		out += fmt.Sprintf("< %s; =%s; > %s", f.Paging.PrevPage.String(), f.Paging.PageCursor.String(), f.Paging.NextPage.String())
	}

	if f.Sorting != nil {
		out += "\n"
		out += f.Sorting.Sort.String()
	}

	out += "\n"
	out += fmt.Sprintf("ix %d; len %d", f.startIndex, f.Size())

	return out
}

func (cc FrameColumnSet) Find(name string) int {
	for i, c := range cc {
		if c.Name == name {
			return i
		}
	}

	return -1
}

func (cc FrameColumnSet) String() string {
	out := ""
	for _, c := range cc {
		out += fmt.Sprintf("%s<%s>, ", c.Name, c.Kind)
	}
	return strings.TrimRight(out, " ,")
}

// Receivers to conform to rdbms field matcher
func (c *FrameColumn) IsBoolean() bool {
	return c.Kind == "Bool"
}

func (c *FrameColumn) IsNumeric() bool {
	return c.Kind == "Number"
}

func (c *FrameColumn) IsDateTime() bool {
	return c.Kind == "DateTime"
}

func (c *FrameColumn) IsRef() bool {
	// @todo not quite right
	return c.Kind == "Record"
}

func (r FrameRow) ToVars(cc FrameColumnSet) (vv *expr.Vars, err error) {
	vv, _ = expr.NewVars(nil)

	// The row
	for i, c := range r {
		if c == nil {
			err = vv.AssignFieldValue(cc[i].Name, nil)
			if err != nil {
				return nil, err
			}
		} else {
			err := vv.AssignFieldValue(cc[i].Name, c)
			if err != nil {
				return nil, err
			}
		}
	}
	return
}

func (dd FrameDefinitionSet) Find(name string) *FrameDefinition {
	for _, d := range dd {
		if d.Name == name {
			return d
		}
	}

	return nil
}

func (dd FrameDefinitionSet) FindBySourceRef(source, ref string) *FrameDefinition {
	for _, d := range dd {
		if d.Source == source && d.Ref == ref {
			return d
		}
	}

	return nil
}

func (r FrameRow) MarshalJSON() (out []byte, err error) {
	aux := make([]string, len(r))
	var s string

	for i, c := range r {
		if c == nil {
			continue
		}

		s, err = cast.ToStringE(c.Get())
		if err != nil {
			return nil, err
		}

		aux[i] = s
	}

	return json.Marshal(aux)
}

func (r FrameRow) String() string {
	out := ""
	var s string
	var err error
	for _, c := range r {
		s, err = cast.ToStringE(c.Get())
		if err != nil {
			out = fmt.Sprintf("%s, [STRING CAST ERROR]%s", out, err.Error())
		} else {
			out = fmt.Sprintf("%s, %s", out, s)
		}
	}

	return strings.Trim(out, ", ")
}
