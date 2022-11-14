package types

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/ql"
	"github.com/cortezaproject/corteza-server/pkg/sql"
	"github.com/spf13/cast"
)

type (
	Report struct {
		ID     uint64      `json:"reportID,string"`
		Handle string      `json:"handle"`
		Meta   *ReportMeta `json:"meta,omitempty"`

		Scenarios ReportScenarioSet   `json:"scenarios,omitempty"`
		Sources   ReportDataSourceSet `json:"sources"`
		Blocks    ReportBlockSet      `json:"blocks"`

		Labels map[string]string `json:"labels,omitempty"`

		OwnedBy   uint64     `json:"ownedBy"`
		CreatedBy uint64     `json:"createdBy"`
		CreatedAt time.Time  `json:"createdAt"`
		UpdatedBy uint64     `json:"updatedBy,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	ReportMeta struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	ReportScenarioSet []*ReportScenario
	ScenarioFilterMap map[string]ReportFilterExpr
	ReportScenario    struct {
		ScenarioID uint64            `json:"scenarioID,string,omitempty"`
		Label      string            `json:"label"`
		Filters    ScenarioFilterMap `json:"filters,omitempty"`
	}

	ReportDataSourceSet []*ReportDataSource
	ReportDataSource    struct {
		Meta interface{} `json:"meta,omitempty"`
		Step *ReportStep `json:"step"`
	}

	ReportBlockSet []*ReportBlock
	ReportBlock    struct {
		BlockID     uint64                 `json:"blockID,string"`
		Title       string                 `json:"title"`
		Description string                 `json:"description"`
		Key         string                 `json:"key"`
		Kind        string                 `json:"kind"`
		Options     map[string]interface{} `json:"options,omitempty"`
		Elements    []interface{}          `json:"elements"`
		Sources     ReportStepSet          `json:"sources"`
		XYWH        [4]int                 `json:"xywh"`
		Layout      string                 `json:"layout"`
	}

	ReportStepSet []*ReportStep
	ReportStep    struct {
		Kind string `json:"kind,omitempty"`

		Load      *ReportStepLoad      `json:"load,omitempty"`
		Join      *ReportStepJoin      `json:"join,omitempty"`
		Link      *ReportStepLink      `json:"link,omitempty"`
		Aggregate *ReportStepAggregate `json:"aggregate,omitempty"`

		// @todo remove for the next set of patch/major releases.
		//       it exists just for the migration as we need to rename this one.
		Group_legacy *ReportLegacyStepGroup `json:"group,omitempty"`
	}

	ReportStepLoad struct {
		Name       string                 `json:"name"`
		Source     string                 `json:"source"`
		Definition map[string]interface{} `json:"definition"`
		Filter     *ReportFilterExpr      `json:"filter,omitempty"`
	}

	ReportStepJoin struct {
		Name          string            `json:"name"`
		LocalSource   string            `json:"localSource"`
		LocalColumn   string            `json:"localColumn"`
		ForeignSource string            `json:"foreignSource"`
		ForeignColumn string            `json:"foreignColumn"`
		Filter        *ReportFilterExpr `json:"filter,omitempty"`
	}

	ReportStepLink struct {
		Name          string            `json:"name"`
		LocalSource   string            `json:"localSource"`
		LocalColumn   string            `json:"localColumn"`
		ForeignSource string            `json:"foreignSource"`
		ForeignColumn string            `json:"foreignColumn"`
		Filter        *ReportFilterExpr `json:"filter,omitempty"`
	}

	ReportLegacyStepGroup struct {
		Name    string                   `json:"name"`
		Source  string                   `json:"source"`
		Keys    ReportAggregateColumnSet `json:"keys"`
		Columns ReportAggregateColumnSet `json:"columns"`
		Filter  *ReportFilterExpr        `json:"filter,omitempty"`
	}

	ReportStepAggregate struct {
		Name    string                   `json:"name"`
		Source  string                   `json:"source"`
		Keys    ReportAggregateColumnSet `json:"keys"`
		Columns ReportAggregateColumnSet `json:"columns"`
		Filter  *ReportFilterExpr        `json:"filter,omitempty"`
	}

	ReportAggregateColumnSet []*ReportAggregateColumn
	ReportAggregateColumn    struct {
		Name  string            `json:"name"`
		Label string            `json:"label"`
		Def   *ReportFilterExpr `json:"def"`
	}

	ReportFilter struct {
		ReportID []uint64 `json:"reportID"`

		Handle string `json:"handle"`
		Query  string `json:"query"`

		Deleted filter.State `json:"deleted"`

		LabeledIDs []uint64          `json:"-"`
		Labels     map[string]string `json:"labels,omitempty"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Report) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	// ReportFilterExpr is a wrapper for ql.ASTNode to implement custom JSON
	// unmarshal required by reporting.
	// @todo consider moving this to the ql package
	ReportFilterExpr struct {
		*ql.ASTNode
		Error string `json:"error,omitempty"`
	}
)

// ReportSteps returns a ReportStepSet collected from the ReportDataSourceSet
func (ss ReportDataSourceSet) ReportSteps() ReportStepSet {
	out := make(ReportStepSet, 0, 124)

	for _, s := range ss {
		out = append(out, s.Step)
	}

	return out
}

// ReportSteps returns a ReportStepSet collected from the ReportBlockSet
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

// Store stuff

func (vv *ReportMeta) Scan(src any) error           { return sql.ParseJSON(src, vv) }
func (vv *ReportMeta) Value() (driver.Value, error) { return json.Marshal(vv) }

func (vv *ReportBlockSet) Scan(src any) error          { return sql.ParseJSON(src, vv) }
func (vv ReportBlockSet) Value() (driver.Value, error) { return json.Marshal(vv) }

// Scan on ReportDataSourceSet gracefully handles conversion from NULL
func (vv *ReportDataSourceSet) Scan(src any) error          { return sql.ParseJSON(src, vv) }
func (vv ReportDataSourceSet) Value() (driver.Value, error) { return json.Marshal(vv) }

func (vv *ReportScenarioSet) Scan(src any) error          { return sql.ParseJSON(src, vv) }
func (vv ReportScenarioSet) Value() (driver.Value, error) { return json.Marshal(vv) }

// Node is a helper for accessing the wrapped QL node to omit nil checks
func (f *ReportFilterExpr) Node() *ql.ASTNode {
	if f == nil {
		return nil
	}
	return f.ASTNode
}

// UnmarshalJSON parses the wrap into a proper QL node and an optional error
//
// The function can work over JSON strings (where FE provides a QL node) or
// raw expression strings (where FE sends over the easeier stringified expression).
func (f *ReportFilterExpr) UnmarshalJSON(data []byte) (err error) {
	var aux interface{}
	if err = json.Unmarshal(data, &aux); err != nil {
		return
	}

	p := ql.NewParser()

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
	err = f.ASTNode.Traverse(func(n *ql.ASTNode) (bool, *ql.ASTNode, error) {
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
