package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/revisions"
	"github.com/crusttech/human/server/pkg/sql"
	"reflect"
	"time"
)

type (
	Report struct {
		ID        uint64                           `json:"reportID,string"`
		TenantID  uint64                           `json:"tenantID,string,omitempty"`
		ProjectID uint64                           `json:"projectID,string,omitempty"`
		Handle    string                           `json:"handle"`
		Meta      *ReportMeta                      `json:"meta,omitempty"`
		Scenarios ReportScenarioSet                `json:"scenarios,omitempty"`
		Sources   ReportDataSourceSet              `json:"sources"`
		Blocks    ReportBlockSet                   `json:"blocks"`
		OwnedBy   uint64                           `json:"ownedBy"`
		CreatedAt time.Time                        `json:"createdAt"`
		UpdatedAt *time.Time                       `json:"updatedAt,omitempty"`
		DeletedAt *time.Time                       `json:"deletedAt,omitempty"`
		CreatedBy uint64                           `json:"createdBy"`
		UpdatedBy uint64                           `json:"updatedBy,omitempty"`
		DeletedBy uint64                           `json:"deletedBy,omitempty"`
		Labels    map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	}

	ReportMeta struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	ReportScenario struct {
		ScenarioID uint64            `json:"scenarioID,string,omitempty"`
		Label      string            `json:"label"`
		Filters    ScenarioFilterMap `json:"filters,omitempty"`
	}

	ReportDataSource struct {
		Meta interface{} `json:"meta,omitempty"`
		Step *ReportStep `json:"step"`
	}

	ReportBlock struct {
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

	ReportStep struct {
		Kind         string                 `json:"kind,omitempty"`
		Load         *ReportStepLoad        `json:"load,omitempty"`
		Join         *ReportStepJoin        `json:"join,omitempty"`
		Link         *ReportStepLink        `json:"link,omitempty"`
		Aggregate    *ReportStepAggregate   `json:"aggregate,omitempty"`
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

	ReportAggregateColumn struct {
		Name  string            `json:"name"`
		Label string            `json:"label"`
		Def   *ReportFilterExpr `json:"def"`
	}
)

func (r Report) Clone() *Report {
	dup := r
	if r.Meta != nil {
		dup.Meta = r.Meta.Clone()
	}

	if r.UpdatedAt != nil {
		v := *r.UpdatedAt
		dup.UpdatedAt = &v
	}

	if r.DeletedAt != nil {
		v := *r.DeletedAt
		dup.DeletedAt = &v
	}

	if r.Labels != nil {
		dup.Labels = make(map[string]labelTypes.LabelValue, len(r.Labels))
		for k, v := range r.Labels {
			dup.Labels[k] = v
		}
	}
	return &dup
}

func (r Report) Diff(cmp *Report) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &Report{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "reportID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.Handle != cmp.Handle {
		out = append(out, &revisions.Change{Key: "handle", Old: []any{cmp.Handle}, New: []any{r.Handle}})
	}

	if (r.Meta == nil) != (cmp.Meta == nil) {
		out = append(out, &revisions.Change{Key: "meta", Old: []any{cmp.Meta}, New: []any{r.Meta}})
	} else if r.Meta != nil {
		for _, c := range r.Meta.Diff(cmp.Meta) {
			c.Key = "meta." + c.Key
			out = append(out, c)
		}
	}

	if !reflect.DeepEqual(r.Scenarios, cmp.Scenarios) {
		out = append(out, &revisions.Change{Key: "scenarios", Old: []any{cmp.Scenarios}, New: []any{r.Scenarios}})
	}

	if !reflect.DeepEqual(r.Sources, cmp.Sources) {
		out = append(out, &revisions.Change{Key: "sources", Old: []any{cmp.Sources}, New: []any{r.Sources}})
	}

	if !reflect.DeepEqual(r.Blocks, cmp.Blocks) {
		out = append(out, &revisions.Change{Key: "blocks", Old: []any{cmp.Blocks}, New: []any{r.Blocks}})
	}

	if r.OwnedBy != cmp.OwnedBy {
		out = append(out, &revisions.Change{Key: "ownedBy", Old: []any{cmp.OwnedBy}, New: []any{r.OwnedBy}})
	}

	if !reflect.DeepEqual(r.CreatedAt, cmp.CreatedAt) {
		out = append(out, &revisions.Change{Key: "createdAt", Old: []any{cmp.CreatedAt}, New: []any{r.CreatedAt}})
	}

	if !reflect.DeepEqual(r.UpdatedAt, cmp.UpdatedAt) {
		out = append(out, &revisions.Change{Key: "updatedAt", Old: []any{cmp.UpdatedAt}, New: []any{r.UpdatedAt}})
	}

	if !reflect.DeepEqual(r.DeletedAt, cmp.DeletedAt) {
		out = append(out, &revisions.Change{Key: "deletedAt", Old: []any{cmp.DeletedAt}, New: []any{r.DeletedAt}})
	}

	if r.CreatedBy != cmp.CreatedBy {
		out = append(out, &revisions.Change{Key: "createdBy", Old: []any{cmp.CreatedBy}, New: []any{r.CreatedBy}})
	}

	if r.UpdatedBy != cmp.UpdatedBy {
		out = append(out, &revisions.Change{Key: "updatedBy", Old: []any{cmp.UpdatedBy}, New: []any{r.UpdatedBy}})
	}

	if r.DeletedBy != cmp.DeletedBy {
		out = append(out, &revisions.Change{Key: "deletedBy", Old: []any{cmp.DeletedBy}, New: []any{r.DeletedBy}})
	}

	if !reflect.DeepEqual(r.Labels, cmp.Labels) {
		out = append(out, &revisions.Change{Key: "labels", Old: []any{cmp.Labels}, New: []any{r.Labels}})
	}
	return out
}

func (r ReportMeta) Clone() *ReportMeta {
	dup := r
	return &dup
}

func (r ReportMeta) Diff(cmp *ReportMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ReportMeta{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	return out
}

func (r *ReportMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ReportMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ReportScenario) Clone() *ReportScenario {
	dup := r
	return &dup
}

func (r ReportScenario) Diff(cmp *ReportScenario) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ReportScenario{}
	}
	if r.ScenarioID != cmp.ScenarioID {
		out = append(out, &revisions.Change{Key: "scenarioID", Old: []any{cmp.ScenarioID}, New: []any{r.ScenarioID}})
	}

	if r.Label != cmp.Label {
		out = append(out, &revisions.Change{Key: "label", Old: []any{cmp.Label}, New: []any{r.Label}})
	}

	if !reflect.DeepEqual(r.Filters, cmp.Filters) {
		out = append(out, &revisions.Change{Key: "filters", Old: []any{cmp.Filters}, New: []any{r.Filters}})
	}

	return out
}

func (r *ReportScenario) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ReportScenario) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ReportDataSource) Clone() *ReportDataSource {
	dup := r
	if r.Step != nil {
		dup.Step = r.Step.Clone()
	}

	return &dup
}

func (r ReportDataSource) Diff(cmp *ReportDataSource) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ReportDataSource{}
	}
	if r.Meta != cmp.Meta {
		out = append(out, &revisions.Change{Key: "meta", Old: []any{cmp.Meta}, New: []any{r.Meta}})
	}

	if (r.Step == nil) != (cmp.Step == nil) {
		out = append(out, &revisions.Change{Key: "step", Old: []any{cmp.Step}, New: []any{r.Step}})
	} else if r.Step != nil {
		for _, c := range r.Step.Diff(cmp.Step) {
			c.Key = "step." + c.Key
			out = append(out, c)
		}
	}

	return out
}

func (r *ReportDataSource) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ReportDataSource) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ReportBlock) Clone() *ReportBlock {
	dup := r
	if r.Options != nil {
		dup.Options = make(map[string]interface{}, len(r.Options))
		for k, v := range r.Options {
			dup.Options[k] = v
		}
	}

	if r.Elements != nil {
		dup.Elements = make([]interface{}, len(r.Elements))
		copy(dup.Elements, r.Elements)
	}

	return &dup
}

func (r ReportBlock) Diff(cmp *ReportBlock) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ReportBlock{}
	}
	if r.BlockID != cmp.BlockID {
		out = append(out, &revisions.Change{Key: "blockID", Old: []any{cmp.BlockID}, New: []any{r.BlockID}})
	}

	if r.Title != cmp.Title {
		out = append(out, &revisions.Change{Key: "title", Old: []any{cmp.Title}, New: []any{r.Title}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if r.Key != cmp.Key {
		out = append(out, &revisions.Change{Key: "key", Old: []any{cmp.Key}, New: []any{r.Key}})
	}

	if r.Kind != cmp.Kind {
		out = append(out, &revisions.Change{Key: "kind", Old: []any{cmp.Kind}, New: []any{r.Kind}})
	}

	if !reflect.DeepEqual(r.Options, cmp.Options) {
		out = append(out, &revisions.Change{Key: "options", Old: []any{cmp.Options}, New: []any{r.Options}})
	}

	if !reflect.DeepEqual(r.Elements, cmp.Elements) {
		out = append(out, &revisions.Change{Key: "elements", Old: []any{cmp.Elements}, New: []any{r.Elements}})
	}

	if !reflect.DeepEqual(r.Sources, cmp.Sources) {
		out = append(out, &revisions.Change{Key: "sources", Old: []any{cmp.Sources}, New: []any{r.Sources}})
	}

	if r.XYWH != cmp.XYWH {
		out = append(out, &revisions.Change{Key: "xywh", Old: []any{cmp.XYWH}, New: []any{r.XYWH}})
	}

	if r.Layout != cmp.Layout {
		out = append(out, &revisions.Change{Key: "layout", Old: []any{cmp.Layout}, New: []any{r.Layout}})
	}

	return out
}

func (r *ReportBlock) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ReportBlock) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ReportStep) Clone() *ReportStep {
	dup := r
	if r.Load != nil {
		dup.Load = r.Load.Clone()
	}

	if r.Join != nil {
		dup.Join = r.Join.Clone()
	}

	if r.Link != nil {
		dup.Link = r.Link.Clone()
	}

	if r.Aggregate != nil {
		dup.Aggregate = r.Aggregate.Clone()
	}

	if r.Group_legacy != nil {
		dup.Group_legacy = r.Group_legacy.Clone()
	}

	return &dup
}

func (r ReportStep) Diff(cmp *ReportStep) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ReportStep{}
	}
	if r.Kind != cmp.Kind {
		out = append(out, &revisions.Change{Key: "kind", Old: []any{cmp.Kind}, New: []any{r.Kind}})
	}

	if (r.Load == nil) != (cmp.Load == nil) {
		out = append(out, &revisions.Change{Key: "load", Old: []any{cmp.Load}, New: []any{r.Load}})
	} else if r.Load != nil {
		for _, c := range r.Load.Diff(cmp.Load) {
			c.Key = "load." + c.Key
			out = append(out, c)
		}
	}

	if (r.Join == nil) != (cmp.Join == nil) {
		out = append(out, &revisions.Change{Key: "join", Old: []any{cmp.Join}, New: []any{r.Join}})
	} else if r.Join != nil {
		for _, c := range r.Join.Diff(cmp.Join) {
			c.Key = "join." + c.Key
			out = append(out, c)
		}
	}

	if (r.Link == nil) != (cmp.Link == nil) {
		out = append(out, &revisions.Change{Key: "link", Old: []any{cmp.Link}, New: []any{r.Link}})
	} else if r.Link != nil {
		for _, c := range r.Link.Diff(cmp.Link) {
			c.Key = "link." + c.Key
			out = append(out, c)
		}
	}

	if (r.Aggregate == nil) != (cmp.Aggregate == nil) {
		out = append(out, &revisions.Change{Key: "aggregate", Old: []any{cmp.Aggregate}, New: []any{r.Aggregate}})
	} else if r.Aggregate != nil {
		for _, c := range r.Aggregate.Diff(cmp.Aggregate) {
			c.Key = "aggregate." + c.Key
			out = append(out, c)
		}
	}

	if (r.Group_legacy == nil) != (cmp.Group_legacy == nil) {
		out = append(out, &revisions.Change{Key: "group", Old: []any{cmp.Group_legacy}, New: []any{r.Group_legacy}})
	} else if r.Group_legacy != nil {
		for _, c := range r.Group_legacy.Diff(cmp.Group_legacy) {
			c.Key = "group." + c.Key
			out = append(out, c)
		}
	}

	return out
}

func (r *ReportStep) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ReportStep) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ReportStepLoad) Clone() *ReportStepLoad {
	dup := r
	if r.Definition != nil {
		dup.Definition = make(map[string]interface{}, len(r.Definition))
		for k, v := range r.Definition {
			dup.Definition[k] = v
		}
	}

	if r.Filter != nil {
		v := *r.Filter
		dup.Filter = &v
	}

	return &dup
}

func (r ReportStepLoad) Diff(cmp *ReportStepLoad) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ReportStepLoad{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Source != cmp.Source {
		out = append(out, &revisions.Change{Key: "source", Old: []any{cmp.Source}, New: []any{r.Source}})
	}

	if !reflect.DeepEqual(r.Definition, cmp.Definition) {
		out = append(out, &revisions.Change{Key: "definition", Old: []any{cmp.Definition}, New: []any{r.Definition}})
	}

	if !reflect.DeepEqual(r.Filter, cmp.Filter) {
		out = append(out, &revisions.Change{Key: "filter", Old: []any{cmp.Filter}, New: []any{r.Filter}})
	}

	return out
}

func (r *ReportStepLoad) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ReportStepLoad) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ReportStepJoin) Clone() *ReportStepJoin {
	dup := r
	if r.Filter != nil {
		v := *r.Filter
		dup.Filter = &v
	}

	return &dup
}

func (r ReportStepJoin) Diff(cmp *ReportStepJoin) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ReportStepJoin{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.LocalSource != cmp.LocalSource {
		out = append(out, &revisions.Change{Key: "localSource", Old: []any{cmp.LocalSource}, New: []any{r.LocalSource}})
	}

	if r.LocalColumn != cmp.LocalColumn {
		out = append(out, &revisions.Change{Key: "localColumn", Old: []any{cmp.LocalColumn}, New: []any{r.LocalColumn}})
	}

	if r.ForeignSource != cmp.ForeignSource {
		out = append(out, &revisions.Change{Key: "foreignSource", Old: []any{cmp.ForeignSource}, New: []any{r.ForeignSource}})
	}

	if r.ForeignColumn != cmp.ForeignColumn {
		out = append(out, &revisions.Change{Key: "foreignColumn", Old: []any{cmp.ForeignColumn}, New: []any{r.ForeignColumn}})
	}

	if !reflect.DeepEqual(r.Filter, cmp.Filter) {
		out = append(out, &revisions.Change{Key: "filter", Old: []any{cmp.Filter}, New: []any{r.Filter}})
	}

	return out
}

func (r *ReportStepJoin) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ReportStepJoin) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ReportStepLink) Clone() *ReportStepLink {
	dup := r
	if r.Filter != nil {
		v := *r.Filter
		dup.Filter = &v
	}

	return &dup
}

func (r ReportStepLink) Diff(cmp *ReportStepLink) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ReportStepLink{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.LocalSource != cmp.LocalSource {
		out = append(out, &revisions.Change{Key: "localSource", Old: []any{cmp.LocalSource}, New: []any{r.LocalSource}})
	}

	if r.LocalColumn != cmp.LocalColumn {
		out = append(out, &revisions.Change{Key: "localColumn", Old: []any{cmp.LocalColumn}, New: []any{r.LocalColumn}})
	}

	if r.ForeignSource != cmp.ForeignSource {
		out = append(out, &revisions.Change{Key: "foreignSource", Old: []any{cmp.ForeignSource}, New: []any{r.ForeignSource}})
	}

	if r.ForeignColumn != cmp.ForeignColumn {
		out = append(out, &revisions.Change{Key: "foreignColumn", Old: []any{cmp.ForeignColumn}, New: []any{r.ForeignColumn}})
	}

	if !reflect.DeepEqual(r.Filter, cmp.Filter) {
		out = append(out, &revisions.Change{Key: "filter", Old: []any{cmp.Filter}, New: []any{r.Filter}})
	}

	return out
}

func (r *ReportStepLink) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ReportStepLink) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ReportLegacyStepGroup) Clone() *ReportLegacyStepGroup {
	dup := r
	if r.Filter != nil {
		v := *r.Filter
		dup.Filter = &v
	}

	return &dup
}

func (r ReportLegacyStepGroup) Diff(cmp *ReportLegacyStepGroup) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ReportLegacyStepGroup{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Source != cmp.Source {
		out = append(out, &revisions.Change{Key: "source", Old: []any{cmp.Source}, New: []any{r.Source}})
	}

	if !reflect.DeepEqual(r.Keys, cmp.Keys) {
		out = append(out, &revisions.Change{Key: "keys", Old: []any{cmp.Keys}, New: []any{r.Keys}})
	}

	if !reflect.DeepEqual(r.Columns, cmp.Columns) {
		out = append(out, &revisions.Change{Key: "columns", Old: []any{cmp.Columns}, New: []any{r.Columns}})
	}

	if !reflect.DeepEqual(r.Filter, cmp.Filter) {
		out = append(out, &revisions.Change{Key: "filter", Old: []any{cmp.Filter}, New: []any{r.Filter}})
	}

	return out
}

func (r *ReportLegacyStepGroup) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ReportLegacyStepGroup) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ReportStepAggregate) Clone() *ReportStepAggregate {
	dup := r
	if r.Filter != nil {
		v := *r.Filter
		dup.Filter = &v
	}

	return &dup
}

func (r ReportStepAggregate) Diff(cmp *ReportStepAggregate) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ReportStepAggregate{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Source != cmp.Source {
		out = append(out, &revisions.Change{Key: "source", Old: []any{cmp.Source}, New: []any{r.Source}})
	}

	if !reflect.DeepEqual(r.Keys, cmp.Keys) {
		out = append(out, &revisions.Change{Key: "keys", Old: []any{cmp.Keys}, New: []any{r.Keys}})
	}

	if !reflect.DeepEqual(r.Columns, cmp.Columns) {
		out = append(out, &revisions.Change{Key: "columns", Old: []any{cmp.Columns}, New: []any{r.Columns}})
	}

	if !reflect.DeepEqual(r.Filter, cmp.Filter) {
		out = append(out, &revisions.Change{Key: "filter", Old: []any{cmp.Filter}, New: []any{r.Filter}})
	}

	return out
}

func (r *ReportStepAggregate) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ReportStepAggregate) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ReportAggregateColumn) Clone() *ReportAggregateColumn {
	dup := r
	if r.Def != nil {
		v := *r.Def
		dup.Def = &v
	}

	return &dup
}

func (r ReportAggregateColumn) Diff(cmp *ReportAggregateColumn) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ReportAggregateColumn{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Label != cmp.Label {
		out = append(out, &revisions.Change{Key: "label", Old: []any{cmp.Label}, New: []any{r.Label}})
	}

	if !reflect.DeepEqual(r.Def, cmp.Def) {
		out = append(out, &revisions.Change{Key: "def", Old: []any{cmp.Def}, New: []any{r.Def}})
	}

	return out
}

func (r *ReportAggregateColumn) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ReportAggregateColumn) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseReportMeta(ss []string) (p *ReportMeta, err error) {
	p = &ReportMeta{}
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), p)
}

func (m *ReportScenarioSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ReportScenarioSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseReportScenarioSet(ss []string) (p ReportScenarioSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *ReportDataSourceSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ReportDataSourceSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseReportDataSourceSet(ss []string) (p ReportDataSourceSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *ReportBlockSet) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ReportBlockSet) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseReportBlockSet(ss []string) (p ReportBlockSet, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
