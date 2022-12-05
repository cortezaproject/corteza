package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/sql"

	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	ApigwFilterParams map[string]interface{}

	ApigwFilter struct {
		ID      uint64            `json:"filterID,string"`
		Route   uint64            `json:"routeID,string"`
		Weight  uint64            `json:"weight,string"`
		Ref     string            `json:"ref,omitempty"`
		Kind    string            `json:"kind,omitempty"`
		Enabled bool              `json:"enabled,omitempty"`
		Params  ApigwFilterParams `json:"params"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty" `
	}

	ApigwFilterFilter struct {
		RouteID uint64 `json:"routeID,string"`

		Deleted  filter.State `json:"deleted"`
		Disabled filter.State `json:"disabled"`

		Kind string `json:"kind"`
		Ref  string `json:"ref"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*ApigwFilter) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)

func (vv *ApigwFilterParams) Scan(src any) error          { return sql.ParseJSON(src, vv) }
func (vv ApigwFilterParams) Value() (driver.Value, error) { return json.Marshal(vv) }
