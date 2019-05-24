package types

import (
	"database/sql/driver"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/internal/permissions"
)

type (
	ActionSet []string
	Trigger   struct {
		ID       uint64    `json:"triggerID,string" db:"id"`
		ModuleID uint64    `json:"moduleID,string,omitempty" db:"rel_module"`
		Name     string    `json:"name" db:"name"`
		Actions  ActionSet `json:"actions" db:"actions"`
		Enabled  bool      `json:"enabled" db:"enabled"`
		Source   string    `json:"source" db:"source"`

		NamespaceID uint64 `json:"namespaceID,string" db:"rel_namespace"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
	}

	TriggerFilter struct {
		NamespaceID uint64 `json:"namespaceID,string"`
		Query       string `json:"query"`
		Page        uint   `json:"page"`
		PerPage     uint   `json:"perPage"`
		// Sort    string `json:"sort"`
		Count    uint   `json:"count"`
		ModuleID uint64 `json:"moduleID,string"`
	}
)

func (set *ActionSet) Scan(src interface{}) error {
	if ser, ok := src.([]uint8); ok {
		var tmp = make([]string, 0)
		for _, a := range strings.Split(string(ser), ",") {
			if a = strings.TrimSpace(a); len(a) > 0 {
				tmp = append(tmp, a)
			}
		}

		*set = ActionSet(tmp)
	}
	return nil
}

func (set ActionSet) Value() (driver.Value, error) {
	return strings.Trim(strings.Join(set, ","), " ,"), nil
}

// Resource returns a system resource ID for this type
func (t Trigger) PermissionResource() permissions.Resource {
	return TriggerPermissionResource.AppendID(t.ID)
}
