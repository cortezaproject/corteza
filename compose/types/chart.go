package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/pkg/errors"
)

type (
	Chart struct {
		ID     uint64      `json:"chartID,string"`
		Handle string      `json:"handle"`
		Name   string      `json:"name"`
		Config ChartConfig `json:"config"`

		Labels map[string]string `json:"labels,omitempty"`

		NamespaceID uint64 `json:"namespaceID,string"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	ChartConfig struct {
		Reports     []*ChartConfigReport `json:"reports,omitempty"`
		ColorScheme string               `json:"colorScheme,omitempty"`
	}

	ChartConfigReport struct {
		Filter     string                   `json:"filter"`
		ModuleID   uint64                   `json:"moduleID,string,omitempty"`
		Metrics    []map[string]interface{} `json:"metrics,omitempty"`
		Dimensions []map[string]interface{} `json:"dimensions,omitempty"`
		YAxis      map[string]interface{}   `json:"yAxis,omitempty"`
		Renderer   struct {
			Version string `json:"version,omitempty" `
		} `json:"renderer,omitempty"`
	}

	ChartFilter struct {
		NamespaceID uint64   `json:"namespaceID,string"`
		ChartID     []uint64 `json:"chartID"`
		Handle      string   `json:"handle"`
		Name        string   `json:"name"`
		Query       string   `json:"query"`

		LabeledIDs []uint64          `json:"-"`
		Labels     map[string]string `json:"labels,omitempty"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Chart) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}
)

// Resource returns a system resource ID for this type
func (c Chart) RBACResource() rbac.Resource {
	return ChartRBACResource.AppendID(c.ID)
}

func (c Chart) DynamicRoles(userID uint64) []uint64 {
	return nil
}

// FindByHandle finds chart by it's handle
func (set ChartSet) FindByHandle(handle string) *Chart {
	for i := range set {
		if set[i].Handle == handle {
			return set[i]
		}
	}

	return nil
}

func (cc *ChartConfig) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*cc = ChartConfig{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, cc); err != nil {
			return errors.Wrapf(err, "cannot scan '%v' into ChartConfig", string(b))
		}
	}

	return nil
}

func (cc ChartConfig) Value() (driver.Value, error) {
	return json.Marshal(cc)
}
