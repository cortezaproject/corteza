package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
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
		Name   string `json:"name,omitempty"`
		Groups []*ReportStepGroup
	}
	ReportDataSourceSet []*ReportDataSource

	ReportStepGroup struct {
		Name  string                   `json:"name,omitempty"`
		Steps report.StepDefinitionSet `json:"steps"`
	}

	ReportMeta struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	ReportProjection struct {
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

// @todo make better
func (ss ReportDataSourceSet) ModelSteps() report.StepDefinitionSet {
	out := make(report.StepDefinitionSet, 0, 124)

	for _, s := range ss {
		for _, g := range s.Groups {
			out = append(out, g.Steps...)
		}
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
