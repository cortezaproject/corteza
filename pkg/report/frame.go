package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/qlng"
	"github.com/spf13/cast"
)

type (
	Frame struct {
		Name   string `json:"name"`
		Source string `json:"source"`
		Ref    string `json:"ref,omitempty"`

		// RefValue is the common value between the two datasources
		RefValue string `json:"refValue,omitempty"`
		// RelColumn is what column in the local ds was used
		RelColumn string `json:"relColumn,omitempty"`
		// RelSource is the ds that this frame is related to
		RelSource string `json:"relSource,omitempty"`

		Columns FrameColumnSet `json:"columns"`
		Rows    FrameRowSet    `json:"rows"`

		Paging *filter.Paging     `json:"paging"`
		Sort   filter.SortExprSet `json:"sort"`
		Filter *Filter            `json:"filter"`

		// params to help us perform things in place
		startIndex int
	}

	FrameRowSet []FrameRow
	FrameRow    []expr.TypedValue

	frameCellCaster func(in interface{}) (expr.TypedValue, error)
	FrameColumnSet  []*FrameColumn
	FrameColumn     struct {
		Name    string `json:"name"`
		Label   string `json:"label"`
		Kind    string `json:"kind"`
		Primary bool   `json:"primary"`
		Unique  bool   `json:"unique"`
		System  bool   `json:"system"`

		Multivalue          bool   `json:"multivalue"`
		MultivalueDelimiter string `json:"multivalueDelimiter"`

		Caster frameCellCaster `json:"-" yaml:"-"`
	}

	FrameDefinitionSet []*FrameDefinition
	FrameDefinition    struct {
		Name    string         `json:"name"`
		Source  string         `json:"source"`
		Ref     string         `json:"ref"`
		Columns FrameColumnSet `json:"columns"`

		Filter *Filter            `json:"filter"`
		Paging *filter.Paging     `json:"paging"`
		Sort   filter.SortExprSet `json:"sort"`
	}

	FrameDescriptionSet []*FrameDescription
	FrameDescription    struct {
		Source  string         `json:"source"`
		Ref     string         `json:"ref,omitempty"`
		Columns FrameColumnSet `json:"columns"`

		// @todo size and other shape related bits
	}

	// Filter is a qlng.ASTNode wrapper to get some unmarshal/marshal features
	Filter struct {
		*qlng.ASTNode
		Error string `json:"error,omitempty"`
	}
)

func (f *Filter) Clone() *Filter {
	if f == nil {
		return nil
	}
	if f.ASTNode == nil {
		return nil
	}

	return &Filter{
		ASTNode: f.ASTNode.Clone(),
		Error:   f.Error,
	}
}

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

func (f *Filter) UnmarshalJSON(data []byte) (err error) {
	var aux interface{}
	if err = json.Unmarshal(data, &aux); err != nil {
		return
	}

	p := qlng.NewParser()

	// String expr. needs to be parsed to the AST
	switch v := aux.(type) {
	case string:
		if v == "" {
			return
		}

		f.ASTNode, err = p.Parse(v)
		f.ASTNode.Raw = v
		if err != nil {
			f.Error = err.Error()
		}
		return nil
	}

	// special case for empty JSON
	if bytes.Equal([]byte{'{', '}'}, data) {
		return
	}

	// non-string is considered an AST and we parse that
	if err = json.Unmarshal(data, &f.ASTNode); err != nil {
		f.Error = err.Error()
		return nil
	}

	// traverse the AST to parse any raw exprs.
	if f.ASTNode == nil {
		return nil
	}

	// A raw expression takes priority and replaces the original AST sub-tree
	err = f.ASTNode.Traverse(func(n *qlng.ASTNode) (bool, *qlng.ASTNode, error) {
		if n.Raw == "" {
			return true, n, nil
		}

		aux, err := p.Parse(n.Raw)
		if err != nil {
			return false, n, err
		}
		aux.Raw = n.Raw

		return false, aux, nil
	})

	if err != nil {
		f.Error = err.Error()
	} else {
		f.Error = ""
	}

	return nil
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

func (f *Frame) LastLastRow() FrameRow {
	if f.Size() == 0 {
		return nil
	}

	return f.Rows[f.Size()-2]
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

	if f.RelSource != "" {
		out += fmt.Sprintf("Rel source: %10s;\n", f.RelSource)
	}

	for _, c := range f.Columns {
		out += fmt.Sprintf("%s<%s>, ", c.Name, c.Kind)
	}
	out = strings.TrimRight(out, " ,")
	out += "\n"

	f.WalkRows(func(i int, r FrameRow) error {
		out += fmt.Sprintf("%d| ", i+1)
		out += fmt.Sprintf("%s, ", r.String())
		out = strings.TrimRight(out, ", ") + "\n"
		return nil
	})

	if f.Paging != nil {
		out += "\n"
		out += fmt.Sprintf("< %s; =%s; > %s", f.Paging.PrevPage.String(), f.Paging.PageCursor.String(), f.Paging.NextPage.String())
	}

	if len(f.Sort) > 0 {
		out += "\n"
		out += f.Sort.String()
	}

	out += "\n"
	out += fmt.Sprintf("ix %d; len %d", f.startIndex, f.Size())

	return out
}

func (f *Frame) CollectCursorValues(r FrameRow, cc ...*filter.SortExpr) *filter.PagingCursor {
	cursor := &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

	for _, c := range cc {
		var v interface{}
		if r[f.Columns.Find(c.Column)] != nil {
			v = r[f.Columns.Find(c.Column)].Get()
		}
		// the check for existence should be performed way in advanced so we won't bother here
		cursor.Set(c.Column, v, c.Descending)
	}

	return cursor
}

func (f *Frame) CloneMeta() *Frame {
	return &Frame{
		Name:      f.Name,
		Source:    f.Source,
		Ref:       f.Ref,
		RefValue:  f.RefValue,
		RelColumn: f.RelColumn,
		RelSource: f.RelSource,
		Columns:   f.Columns.Clone(),
		Paging:    f.Paging.Clone(),
		Sort:      f.Sort.Clone(),
		Filter:    f.Filter.Clone(),
	}
}

func (cc FrameColumnSet) Clone() (out FrameColumnSet) {
	out = make(FrameColumnSet, len(cc))
	for i, c := range cc {
		out[i] = &FrameColumn{
			Name:   c.Name,
			Label:  c.Label,
			Kind:   c.Kind,
			Caster: c.Caster,

			Primary: c.Primary,
			Unique:  c.Unique,
			System:  c.System,
		}
	}

	return
}

func (cc FrameColumnSet) Find(name string) int {
	for i, c := range cc {
		if c.Name == name {
			return i
		}
	}

	return -1
}

func (cc FrameColumnSet) FirstPrimary() *FrameColumn {
	for _, c := range cc {
		if c.Primary {
			return c
		}
	}
	return nil
}

func (cc FrameColumnSet) String() string {
	out := ""
	for _, c := range cc {
		out += fmt.Sprintf("%s<%s>, ", c.Name, c.Kind)
	}
	return strings.TrimRight(out, " ,")
}

// OmitSys returns the columns that are not system defined
func (cc FrameColumnSet) OmitSys() FrameColumnSet {
	out := make(FrameColumnSet, 0, len(cc))
	for _, c := range cc {
		if !c.System {
			out = append(out, c)
		}
	}

	return out
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

func (f *FrameDefinition) Clone() (out *FrameDefinition) {
	return &FrameDefinition{
		Name:   f.Name,
		Source: f.Source,
		Ref:    f.Ref,

		Columns: f.Columns.Clone(),
		Paging:  f.Paging.Clone(),
		Sort:    f.Sort.Clone(),
		Filter:  f.Filter.Clone(),
	}
}

func (d FrameDefinition) Validate() error {
	if d.Filter != nil && d.Filter.Error != "" {
		return errors.InvalidData(d.Filter.Error)
	}
	return nil
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

func (row FrameRow) String() string {
	out := ""
	var s string

	format := func(c expr.TypedValue) (out string) {
		out, err := cast.ToStringE(c.Get())
		if err != nil {
			return fmt.Sprintf("[STRING CAST ERROR]%s", err.Error())
		}
		return out
	}

	for _, cell := range row {
		if cell == nil {
			out += "<N/A>, "
		} else {
			switch c := cell.(type) {
			case *expr.Array:
				aux := make([]string, 0, 3)
				for _, c := range c.GetValue() {
					aux = append(aux, format(c))
				}

				out = fmt.Sprintf("%s, %s", out, fmt.Sprintf("[%s]", strings.Join(aux, ", ")))
			default:
				s = format(cell)
				out = fmt.Sprintf("%s, %s", out, s)
			}
		}
	}

	return strings.Trim(out, ", ")
}

func (a FrameRow) Compare(b FrameRow, cols ...int) int {
	for _, col := range cols {
		if cmp, _ := a[col].(expr.Comparable).Compare(b[col]); cmp != 0 {
			return cmp
		}
	}

	return 0
}

// AppendCell adds an extra value to the cell at index i
func (a FrameRow) AppendCell(i int, v expr.TypedValue) (err error) {
	out := make([]expr.TypedValue, 0, 10)

	handle := func(v expr.TypedValue) {
		switch c := v.(type) {
		case *expr.Array:
			out = append(out, c.GetValue()...)
		default:
			out = append(out, v)
		}
	}

	handle(a[i])
	handle(v)

	a[i], err = expr.NewArray(out)
	return err
}

func (dd FrameDescriptionSet) FilterBySource(source string) FrameDescriptionSet {
	out := make(FrameDescriptionSet, 0, len(dd))

	for _, d := range dd {
		if d.Source == source {
			out = append(out, d)
		}
	}

	return out
}

func (dd FrameDescriptionSet) FilterByRef(ref string) FrameDescriptionSet {
	out := make(FrameDescriptionSet, 0, len(dd))

	for _, d := range dd {
		if d.Ref == ref {
			out = append(out, d)
		}
	}

	return out
}
