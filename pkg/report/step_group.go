package report

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type (
	stepGroup struct {
		def *GroupStepDefinition
	}

	groupedDataset struct {
		def *GroupStepDefinition
		ds  Datasource
	}

	GroupDefinition struct {
		Groups  []*GroupKey    `json:"groups"`
		Columns []GroupColumn  `json:"columns"`
		Rows    *RowDefinition `json:"rows,omitempty"`
	}

	GroupStepDefinition struct {
		Name   string `json:"name"`
		Source string `json:"source"`
		GroupDefinition
	}

	GroupKey struct {
		// Name defines the alias for the new column
		Name string `json:"name"`
		// Expr defines the expression to transform the column
		Expr string `json:"expr"`

		// @todo imply from context
		Kind string `json:"kind"`
	}

	// Group columns define what columns we wish to produce and what operations
	// we should perform over them.
	//
	// alias: operation: args; for example -- { "total": { "sum": "cost" } }
	GroupColumn     map[string]AggregateColumn
	AggregateColumn map[string]string
)

var (
	simpleExprMatcher = regexp.MustCompile("^\\*|\\w+$")
)

const (
	stepGroupMaxFramers    = 6
	stepGroupMaxFinalizers = 2
)

func (j *stepGroup) Run(ctx context.Context, dd ...Datasource) (Datasource, error) {
	if len(dd) == 0 {
		return nil, fmt.Errorf("unknown group dimension: %s", j.def.Source)
	}

	return nil, nil
	// @todo
	// return &groupedDataset{
	// 	def: j.def,
	// 	ds:  dd[0],
	// }, nil
}

func (j *stepGroup) Validate() error {
	pfx := "invalid group step: "

	// base things...
	switch {
	case j.def.Name == "":
		return errors.New(pfx + "dimension name not defined")

	case j.def.Source == "":
		return errors.New(pfx + "groupping dimension not defined")
	case len(j.def.Groups) == 0:
		return errors.New(pfx + "no group defined")
	}

	// columns...
	for i, g := range j.def.Groups {
		if g.Name == "" {
			return fmt.Errorf("%sgroup key alias missing for group: %d", pfx, i)
		}
	}

	return nil
}

func (d *stepGroup) Name() string {
	return d.def.Name
}

func (d *stepGroup) Source() []string {
	return []string{d.def.Source}
}

func (d *stepGroup) Def() *StepDefinition {
	return &StepDefinition{Group: d.def}
}

func (c AggregateColumn) GetOp() string {
	for k := range c {
		return strings.ToLower(k)
	}
	return ""
}

// // // //

// @todo manual group step implementation for Datasources that don't provide it
