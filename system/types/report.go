package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/pkg/report"
)

type (
	Report struct {
		ID     uint64      `json:"reportID,string"`
		Handle string      `json:"handle"`
		Meta   *ReportMeta `json:"meta,omitempty"`

		Sources     ReportDataSourceSet `json:"sources"`
		Projections ReportProjectionSet `json:"projections"`

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

	ReportDataSource struct {
		Meta interface{}            `json:"meta,omitempty"`
		Step *report.StepDefinition `json:"step"`
	}
	ReportDataSourceSet []*ReportDataSource

	ReportMeta struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	ReportProjection struct {
		ProjectionID uint64                   `json:"projectionID"`
		Title        string                   `json:"title"`
		Description  string                   `json:"description"`
		Key          string                   `json:"key"`
		Kind         string                   `json:"kind"`
		Options      map[string]interface{}   `json:"options,omitempty"`
		Elements     []interface{}            `json:"elements"`
		Sources      report.StepDefinitionSet `json:"sources"`
		XYWH         [4]int                   `json:"xywh"`
		Layout       string                   `json:"layout"`
	}
	ReportProjectionSet []*ReportProjection

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

func (ss ReportDataSourceSet) ModelSteps() report.StepDefinitionSet {
	out := make(report.StepDefinitionSet, 0, 124)

	for _, s := range ss {
		out = append(out, s.Step)
	}

	return out
}

func (pp ReportProjectionSet) ModelSteps() report.StepDefinitionSet {
	out := make(report.StepDefinitionSet, 0, 124)

	for _, p := range pp {
		out = append(out, p.Sources...)
	}

	return out
}

// Store stuff

func (vv *ReportMeta) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*vv = ReportMeta{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, vv); err != nil {
			return fmt.Errorf("cannot scan '%v' into ReportMeta: %w", string(b), err)
		}
	}

	return nil
}

// Scan on ReportMeta gracefully handles conversion from NULL
func (vv *ReportMeta) Value() (driver.Value, error) {
	if vv == nil {
		return []byte("null"), nil
	}

	return json.Marshal(vv)
}

// Scan on ReportProjectionSet gracefully handles conversion from NULL
func (vv ReportProjectionSet) Value() (driver.Value, error) {
	return json.Marshal(vv)
}

func (vv *ReportProjectionSet) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*vv = ReportProjectionSet{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, vv); err != nil {
			return fmt.Errorf("cannot scan '%v' into ReportProjectionSet: %w", string(b), err)
		}
	}

	return nil
}

// Scan on ReportDataSourceSet gracefully handles conversion from NULL
func (vv ReportDataSourceSet) Value() (driver.Value, error) {
	return json.Marshal(vv)
}

func (vv *ReportDataSourceSet) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*vv = ReportDataSourceSet{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, vv); err != nil {
			return fmt.Errorf("cannot scan '%v' into ReportDataSourceSet: %w", string(b), err)
		}
	}

	return nil
}

func (r *Report) decodeTranslations(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation

	for i, p := range r.Projections {
		projectionID := locale.ContentID(p.ProjectionID, i)
		rpl := strings.NewReplacer(
			"{{projectionID}}", strconv.FormatUint(projectionID, 10),
		)

		// - generic page block stuff
		if aux = tt.FindByKey(rpl.Replace(LocaleKeyReportProjectionTitle.Path)); aux != nil {
			p.Title = aux.Msg
		}
		if aux = tt.FindByKey(rpl.Replace(LocaleKeyReportProjectionDescription.Path)); aux != nil {
			p.Description = aux.Msg
		}
	}
}

func (r *Report) encodeTranslations() (out locale.ResourceTranslationSet) {
	out = make(locale.ResourceTranslationSet, 0, 3)

	// Page blocks
	for i, projection := range r.Projections {
		projectionID := locale.ContentID(projection.ProjectionID, i)
		rpl := strings.NewReplacer(
			"{{projectionID}}", strconv.FormatUint(projectionID, 10),
		)

		// - generic page block stuff
		out = append(out, &locale.ResourceTranslation{
			Resource: r.ResourceTranslation(),
			Key:      rpl.Replace(rpl.Replace(LocaleKeyReportProjectionTitle.Path)),
			Msg:      projection.Title,
		})

		out = append(out, &locale.ResourceTranslation{
			Resource: r.ResourceTranslation(),
			Key:      rpl.Replace(rpl.Replace(LocaleKeyReportProjectionDescription.Path)),
			Msg:      projection.Description,
		})
	}

	return
}
