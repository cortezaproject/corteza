package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	ApigwFilterParams map[string]interface{}

	ApigwFilter struct {
		ID     uint64            `json:"filterID,string"`
		Route  uint64            `json:"routeID,string"`
		Weight uint64            `json:"weight,string"`
		Ref    string            `json:"ref,omitempty"`
		Kind   string            `json:"kind,omitempty"`
		Params ApigwFilterParams `json:"params"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty" `
	}

	ApigwFilterFilter struct {
		RouteID uint64 `json:"routeID,string"`
		Enabled bool   `json:"enabled"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*ApigwFilter) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)

func (vv *ApigwFilterParams) Scan(value interface{}) (err error) {
	if err := json.Unmarshal(value.([]byte), vv); err != nil {
		return fmt.Errorf("cannot scan '%v' into FuncParams", value)
	}

	return
}

func (vv ApigwFilterParams) Value() (driver.Value, error) {
	return json.Marshal(vv)
}
