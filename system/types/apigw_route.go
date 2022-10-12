package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/cortezaproject/corteza-server/pkg/sql"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	ApigwRoute struct {
		ID       uint64         `json:"routeID,string"`
		Endpoint string         `json:"endpoint"`
		Method   string         `json:"method"`
		Enabled  bool           `json:"enabled"`
		Group    uint64         `json:"group,string"`
		Meta     ApigwRouteMeta `json:"meta"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty" `
	}

	ApigwRouteMeta struct {
		Debug bool `json:"debug"`
		Async bool `json:"async"`
	}

	ApigwRouteFilter struct {
		Route    string `json:"route"`
		Endpoint string `json:"endpoint"`
		Method   string `json:"method"`

		Deleted  filter.State `json:"deleted"`
		Disabled filter.State `json:"disabled"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*ApigwRoute) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)

func (cc *ApigwRouteMeta) Scan(src any) error          { return sql.ParseJSON(src, cc) }
func (cc ApigwRouteMeta) Value() (driver.Value, error) { return json.Marshal(cc) }
