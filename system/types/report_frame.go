package types

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/ql"
)

type (
	ReportFrame struct {
		Name   string `json:"name"`
		Source string `json:"source"`
		Ref    string `json:"ref,omitempty"`

		// RefValue is the common value between the two datasources
		RefValue string `json:"refValue,omitempty"`
		// RelColumn is what column in the local ds was used
		RelColumn string `json:"relColumn,omitempty"`
		// RelSource is the ds that this frame is related to
		RelSource string `json:"relSource,omitempty"`

		Columns ReportFrameColumnSet `json:"columns"`
		Rows    []ReportFrameRow     `json:"rows"`

		Paging *filter.Paging     `json:"paging"`
		Sort   filter.SortExprSet `json:"sort"`
		Filter *ql.ASTNode        `json:"filter"`
	}

	FrameDescriptionSet    []ReportFrameDescription
	ReportFrameDescription struct {
		Source  string               `json:"source"`
		Ref     string               `json:"ref,omitempty"`
		Columns ReportFrameColumnSet `json:"columns"`
	}

	ReportFrameRow []string

	ReportFrameColumnSet []ReportFrameColumn
	ReportFrameColumn    struct {
		Name    string `json:"name"`
		Label   string `json:"label"`
		Kind    string `json:"kind"`
		Primary bool   `json:"primary"`
		Unique  bool   `json:"unique"`
		System  bool   `json:"system"`
	}

	ReportFrameDefinitionSet []*ReportFrameDefinition
	ReportFrameDefinition    struct {
		Name    string               `json:"name"`
		Source  string               `json:"source"`
		Ref     string               `json:"ref"`
		Columns ReportFrameColumnSet `json:"columns"`

		Filter *ql.ASTNode        `json:"filter"`
		Paging *filter.Paging     `json:"paging"`
		Sort   filter.SortExprSet `json:"sort"`
	}
)

// @todo nicer formatting and alignment
func (f ReportFrame) String() string {
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

func (cc ReportFrameColumnSet) String() string {
	out := ""
	for _, c := range cc {
		out += fmt.Sprintf("%s<%s>, ", c.Name, c.Kind)
	}
	return strings.TrimRight(out, " ,")
}

// OmitSys returns the columns that are not system defined
func (cc ReportFrameColumnSet) OmitSys() ReportFrameColumnSet {
	out := make(ReportFrameColumnSet, 0, len(cc))
	for _, c := range cc {
		if !c.System {
			out = append(out, c)
		}
	}

	return out
}

func (r ReportFrameRow) String() string {
	return strings.Join(r, ", ")
}
