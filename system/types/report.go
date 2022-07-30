package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/cortezaproject/corteza-server/pkg/sql"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/report"
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

		// Report labels
		Labels map[string]string `json:"labels,omitempty"`

		OwnedBy   uint64     `json:"ownedBy"`
		CreatedBy uint64     `json:"createdBy"`
		CreatedAt time.Time  `json:"createdAt"`
		UpdatedBy uint64     `json:"updatedBy,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	ReportScenarioSet []*ReportScenario
	ScenarioFilterMap map[string]*report.Filter
	ReportScenario    struct {
		ScenarioID uint64            `json:"scenarioID,string,omitempty"`
		Label      string            `json:"label"`
		Filters    ScenarioFilterMap `json:"filters,omitempty"`
	}

	ReportDataSource struct {
		Meta interface{}            `json:"meta,omitempty"`
		Step *report.StepDefinition `json:"step"`
	}
	ReportDataSourceSet []*ReportDataSource

	ReportMeta struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	ReportBlock struct {
		BlockID     uint64                   `json:"blockID,string"`
		Title       string                   `json:"title"`
		Description string                   `json:"description"`
		Key         string                   `json:"key"`
		Kind        string                   `json:"kind"`
		Options     map[string]interface{}   `json:"options,omitempty"`
		Elements    []interface{}            `json:"elements"`
		Sources     report.StepDefinitionSet `json:"sources"`
		XYWH        [4]int                   `json:"xywh"`
		Layout      string                   `json:"layout"`
	}
	ReportBlockSet []*ReportBlock

	ReportFilter struct {
		ReportID []uint64 `json:"reportID"`

		Handle string `json:"handle"`

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
)

// Initial ReportBlock struct definition omitted string casting for the BlockID (sorry)
// so we need to handle that edge case when reading from DB.
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

func (ss ReportDataSourceSet) ModelSteps() report.StepDefinitionSet {
	out := make(report.StepDefinitionSet, 0, 124)

	for _, s := range ss {
		out = append(out, s.Step)
	}

	return out
}

func (pp ReportBlockSet) ModelSteps() report.StepDefinitionSet {
	out := make(report.StepDefinitionSet, 0, 124)

	for _, p := range pp {
		out = append(out, p.Sources...)
	}

	return out
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

// func (r *Report) decodeTranslations(tt locale.ResourceTranslationIndex) {
// 	var aux *locale.ResourceTranslation

// 	for i, p := range r.Blocks {
// 		blockID := locale.ContentID(p.BlockID, i)
// 		rpl := strings.NewReplacer(
// 			"{{blockID}}", strconv.FormatUint(blockID, 10),
// 		)

// 		// - generic page block stuff
// 		if aux = tt.FindByKey(rpl.Replace(LocaleKeyReportBlockTitle.Path)); aux != nil {
// 			p.Title = aux.Msg
// 		}
// 		if aux = tt.FindByKey(rpl.Replace(LocaleKeyReportBlockDescription.Path)); aux != nil {
// 			p.Description = aux.Msg
// 		}
// 	}
// }

// func (r *Report) encodeTranslations() (out locale.ResourceTranslationSet) {
// 	out = make(locale.ResourceTranslationSet, 0, 3)

// 	// Page blocks
// 	for i, block := range r.Blocks {
// 		blockID := locale.ContentID(block.BlockID, i)
// 		rpl := strings.NewReplacer(
// 			"{{blockID}}", strconv.FormatUint(blockID, 10),
// 		)

// 		// - generic page block stuff
// 		out = append(out, &locale.ResourceTranslation{
// 			Resource: r.ResourceTranslation(),
// 			Key:      rpl.Replace(rpl.Replace(LocaleKeyReportBlockTitle.Path)),
// 			Msg:      block.Title,
// 		})

// 		out = append(out, &locale.ResourceTranslation{
// 			Resource: r.ResourceTranslation(),
// 			Key:      rpl.Replace(rpl.Replace(LocaleKeyReportBlockDescription.Path)),
// 			Msg:      block.Description,
// 		})
// 	}

// 	return
// }
