package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/sql"
)

type (
	ProjectGroup struct {
		ID        uint64           `json:"projectGroupID,string"`
		TenantID  uint64           `json:"tenantID,string,omitempty"`
		ProjectID uint64           `json:"projectID,string"`
		Handle    string           `json:"handle"`
		Meta      ProjectGroupMeta `json:"meta"`
		CreatedAt time.Time        `json:"createdAt,omitempty"`
		UpdatedAt *time.Time       `json:"updatedAt,omitempty"`
		DeletedAt *time.Time       `json:"deletedAt,omitempty"`
	}

	ProjectGroupMeta struct {
		Short       string `json:"short"`
		Description string `json:"description,omitempty"`
	}

	ProjectGroupFilter struct {
		ProjectGroupID []uint64     `json:"projectGroupID"`
		TenantID       uint64       `json:"tenantID,string,omitempty"`
		ProjectID      uint64       `json:"projectID,string,omitempty"`
		Handle         string       `json:"handle,omitempty"`
		Query          string       `json:"query,omitempty"`
		Deleted        filter.State `json:"deleted"`

		LabeledIDs []uint64 `json:"-"`

		Check func(*ProjectGroup) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}

	ProjectGroupEntry struct {
		ID             uint64    `json:"-"`
		ProjectGroupID uint64    `json:"projectGroupID,string"`
		ResourceRef    string    `json:"resourceRef"`
		CreatedAt      time.Time `json:"createdAt,omitempty"`
	}

	ProjectGroupEntryFilter struct {
		ProjectGroupID uint64 `json:"projectGroupID,string,omitempty"`
		ResourceRef    string `json:"resourceRef,omitempty"`
		Limit          uint   `json:"-"`
	}
)

func (m *ProjectGroupMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ProjectGroupMeta) Value() (driver.Value, error) { return json.Marshal(m) }
