package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/cortezaproject/corteza/server/pkg/sql"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	DalSensitivityLevelMeta struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	DalSensitivityLevel struct {
		ID     uint64 `json:"sensitivityLevelID,string"`
		Handle string `json:"handle"`
		Level  int    `json:"level"`

		Meta DalSensitivityLevelMeta `json:"meta"`

		Labels map[string]string `json:"labels,omitempty"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty" `
	}

	DalSensitivityLevelFilter struct {
		SensitivityLevelID []uint64 `json:"sensitivityLevelID,string"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*DalSensitivityLevel) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Paging
		filter.Sorting
	}
)

func ParseDalSensitivityLevelMeta(ss []string) (m DalSensitivityLevelMeta, err error) {
	if len(ss) == 0 {
		return
	}

	err = json.Unmarshal([]byte(ss[0]), &m)
	return
}

func (nm *DalSensitivityLevelMeta) Scan(src any) error          { return sql.ParseJSON(src, nm) }
func (nm DalSensitivityLevelMeta) Value() (driver.Value, error) { return json.Marshal(nm) }

func (ss DalSensitivityLevelSet) Len() int { return len(ss) }
func (ss DalSensitivityLevelSet) Less(i, j int) bool {
	return ss[i].Level < ss[j].Level
}
func (ss DalSensitivityLevelSet) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}
