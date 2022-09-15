package reporting

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/system/types"
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
		Rows    []FrameRow     `json:"rows"`

		Paging *filter.Paging          `json:"paging"`
		Sort   filter.SortExprSet      `json:"sort"`
		Filter *types.ReportFilterExpr `json:"filter"`
	}

	FrameDescription struct {
		Source  string         `json:"source"`
		Ref     string         `json:"ref,omitempty"`
		Columns FrameColumnSet `json:"columns"`
	}

	FrameRow []string

	FrameColumnSet []FrameColumn
	FrameColumn    struct {
		Name    string `json:"name"`
		Label   string `json:"label"`
		Kind    string `json:"kind"`
		Primary bool   `json:"primary"`
		Unique  bool   `json:"unique"`
		System  bool   `json:"system"`

		Multivalue          bool   `json:"multivalue"`
		MultivalueDelimiter string `json:"multivalueDelimiter"`
	}

	FrameDefinitionSet []*FrameDefinition
	FrameDefinition    struct {
		Name    string         `json:"name"`
		Source  string         `json:"source"`
		Ref     string         `json:"ref"`
		Columns FrameColumnSet `json:"columns"`

		Filter *types.ReportFilterExpr `json:"filter"`
		Paging *filter.Paging          `json:"paging"`
		Sort   filter.SortExprSet      `json:"sort"`
	}
)

// @todo nicer formatting and alignment
func (f Frame) String() string {
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

	for i, r := range f.Rows {
		out += fmt.Sprintf("%d| %s", i+1, strings.Join(r, ", "))
		out += "\n"
	}

	if f.Paging != nil {
		out += "\n"
		out += fmt.Sprintf("< %s; =%s; > %s", f.Paging.PrevPage.String(), f.Paging.PageCursor.String(), f.Paging.NextPage.String())
	}

	if len(f.Sort) > 0 {
		out += "\n"
		out += f.Sort.String()
	}

	out += "\n"
	out += fmt.Sprintf("ix %d; len %d", 0, len(f.Rows))

	return out
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

func (r FrameRow) String() string {
	return strings.Join(r, ", ")
}

// FilterBySource returns a set of definitions for the requested identifier
func (dd FrameDefinitionSet) FilterBySource(ident string) FrameDefinitionSet {
	out := make(FrameDefinitionSet, 0, len(dd))
	for _, d := range dd {
		if d.Source == ident {
			out = append(out, d)
		}
	}

	return out
}
