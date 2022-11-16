package report

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/cortezaproject/corteza/server/pkg/qlng"
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
		Keys    []*GroupColumn `json:"keys"`
		Columns []*GroupColumn `json:"columns"`
		Filter  *Filter        `json:"filter,omitempty"`
	}

	GroupStepDefinition struct {
		Name   string `json:"name"`
		Source string `json:"source"`
		GroupDefinition
	}

	GroupColumn struct {
		Name  string  `json:"name"`
		Label string  `json:"label"`
		Def   *colDef `json:"def"`
	}

	// Wrapper for easier marshal/unmarshal
	colDef struct {
		*qlng.ASTNode
	}
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
	case len(j.def.Keys) == 0:
		return errors.New(pfx + "no group defined")
	}

	// columns...
	for i, g := range j.def.Keys {
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

func (def *colDef) UnmarshalJSON(data []byte) (err error) {
	var aux interface{}
	if err = json.Unmarshal(data, &aux); err != nil {
		return
	}

	p := qlng.NewParser()

	// String expr. needs to be parsed to the AST
	switch v := aux.(type) {
	case string:
		if def.ASTNode, err = p.Parse(v); err != nil {
			return
		}
		def.ASTNode.Raw = v
		return
	}

	// special case for empty JSON
	if bytes.Equal([]byte{'{', '}'}, data) {
		return
	}

	// non-string is considered an AST and we parse that
	if err = json.Unmarshal(data, &def.ASTNode); err != nil {
		return
	}

	// traverse the AST to parse any raw exprs.
	if def.ASTNode == nil {
		return
	}

	// A raw expression takes priority and replaces the original AST sub-tree
	return def.ASTNode.Traverse(func(n *qlng.ASTNode) (bool, *qlng.ASTNode, error) {
		if n.Raw == "" {
			return true, n, nil
		}

		aux, err := p.Parse(n.Raw)
		aux.Raw = n.Raw
		if err != nil {
			return false, n, err
		}

		return false, aux, nil
	})
}

// // // //

// @todo manual group step implementation for Datasources that don't provide it
