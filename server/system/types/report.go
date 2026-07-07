package types

import (
	"bytes"
	"encoding/json"

	"github.com/crusttech/human/server/pkg/ast"
	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/ql"
	"github.com/spf13/cast"
)

type (
	ReportScenarioSet      []*ReportScenario
	ScenarioFilterMap      map[string]ReportFilterExpr
	ReportDataSourceSet    []*ReportDataSource
	ReportBlockSet         []*ReportBlock
	ReportStepSet          []*ReportStep
	ReportAggregateColumnSet []*ReportAggregateColumn

	ReportFilter struct {
		ReportID []string `json:"reportID"`

		Handle string `json:"handle"`
		Query  string `json:"query"`

		Deleted filter.State `json:"deleted"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		Check func(*Report) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}

	// ReportFilterExpr wraps ast.ASTNode for custom JSON unmarshal required by reporting.
	ReportFilterExpr struct {
		*ast.ASTNode
		Error string `json:"error,omitempty"`
	}
)

func (ss ReportDataSourceSet) ReportSteps() ReportStepSet {
	out := make(ReportStepSet, 0, 124)
	for _, s := range ss {
		out = append(out, s.Step)
	}
	return out
}

func (pp ReportBlockSet) ReportSteps() ReportStepSet {
	out := make(ReportStepSet, 0, 124)
	for _, p := range pp {
		out = append(out, p.Sources...)
	}
	return out
}

// Initial ReportBlock struct definition omitted string casting for the BlockID (sorry)
// so we need to handle that edge case when reading from DB.
// @todo consider dropping this in the next/one of the following releases
func (b *ReportBlock) UnmarshalJSON(data []byte) (err error) {
	type internalReportBlock ReportBlock
	i := struct {
		internalReportBlock
		BlockID interface{} `json:"blockID"`
	}{}

	if err = json.Unmarshal(data, &i); err != nil {
		return
	}

	bID, err := cast.ToUint64E(i.BlockID)
	if err != nil {
		return
	}

	*b = ReportBlock(i.internalReportBlock)
	b.BlockID = bID

	return nil
}

func (f *ReportFilterExpr) Node() *ast.ASTNode {
	if f == nil {
		return nil
	}
	return f.ASTNode
}

func (f *ReportFilterExpr) UnmarshalJSON(data []byte) (err error) {
	var aux interface{}
	if err = json.Unmarshal(data, &aux); err != nil {
		return
	}

	p := ql.NewParser()

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

	if bytes.Equal([]byte{'{', '}'}, data) {
		return
	}

	if err = json.Unmarshal(data, &f.ASTNode); err != nil {
		f.Error = err.Error()
		return nil
	}

	if f.ASTNode == nil {
		return nil
	}

	err = f.ASTNode.Traverse(func(n *ast.ASTNode) (bool, *ast.ASTNode, error) {
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
