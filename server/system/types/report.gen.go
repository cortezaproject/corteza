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
	"github.com/crusttech/human/server/pkg/sql"
	"time"
)

type Report struct {
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

func (r Report) Clone() *Report {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}

func (m *ReportMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ReportMeta) Value() (driver.Value, error) { return json.Marshal(m) }

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
