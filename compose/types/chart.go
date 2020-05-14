package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	Chart struct {
		ID     uint64      `json:"chartID,string" db:"id"`
		Handle string      `json:"handle" db:"handle"`
		Name   string      `json:"name" db:"name"`
		Config ChartConfig `json:"config" db:"config"`

		NamespaceID uint64 `json:"namespaceID,string" db:"rel_namespace,string"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
	}

	ChartConfig struct {
		Reports     []*ChartConfigReport `json:"reports,omitempty" yaml:",omitempty"`
		ColorScheme string               `json:"colorScheme,omitempty" yaml:",omitempty"`
	}

	ChartConfigReport struct {
		Filter     string                   `json:"filter"                    yaml:",omitempty"`
		ModuleID   uint64                   `json:"moduleID,string,omitempty" yaml:"moduleID,omitempty"`
		Metrics    []map[string]interface{} `json:"metrics,omitempty"         yaml:",omitempty"`
		Dimensions []map[string]interface{} `json:"dimensions,omitempty"      yaml:",omitempty"`
		YAxis      map[string]interface{}   `json:"yAxis,omitempty"           yaml:",omitempty"`
		Renderer   struct {
			Version string `json:"version,omitempty"  yaml:",omitempty"`
		} `json:"renderer,omitempty" yaml:",omitempty"`
	}

	ChartFilter struct {
		NamespaceID uint64 `json:"namespaceID,string"`
		Handle      string `json:"handle"`
		Query       string `json:"query"`

		Sort string `json:"sort"`

		// Standard paging fields & helpers
		rh.PageFilter

		// Resource permission check filter
		IsReadable *permissions.ResourceFilter `json:"-"`
	}
)

// Resource returns a system resource ID for this type
func (c Chart) PermissionResource() permissions.Resource {
	return ChartPermissionResource.AppendID(c.ID)
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
			return errors.Wrapf(err, "Can not scan '%v' into ChartConfig", string(b))
		}
	}

	return nil
}

func (cc ChartConfig) Value() (driver.Value, error) {
	return json.Marshal(cc)
}
