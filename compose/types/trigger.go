package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/internal/permissions"
)

type (
	ActionSet []string
	Trigger   struct {
		ID          uint64    `json:"triggerID,string" db:"id"`
		NamespaceID uint64    `json:"namespaceID,string" db:"rel_namespace"`
		ModuleID    uint64    `json:"moduleID,string,omitempty" db:"rel_module"`
		Name        string    `json:"name" db:"name"`
		Actions     ActionSet `json:"actions" db:"actions"`
		Enabled     bool      `json:"enabled" db:"enabled"`
		Source      string    `json:"source" db:"source"`
		// Weight      int       `json:"weight" db:"weight"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
	}

	Script struct {
		Source   string `json:"source"`
		Language string `json:"language"`
		Critical bool   `json:"critical"`
		Async    bool   `json:"async"`
		Timeout  uint32 `json:"timeout"`
		RunAs    uint64 `json:"runAs,string"`
	}

	TriggerFilter struct {
		NamespaceID uint64 `json:"namespaceID,string"`
		Query       string `json:"query"`
		EnabledOnly bool   `json:"-"`
		Page        uint   `json:"page"`
		PerPage     uint   `json:"perPage"`
		// Sort    string `json:"sort"`
		Count    uint   `json:"count"`
		ModuleID uint64 `json:"moduleID,string"`
	}
)

func (t Trigger) IsCritical() bool {
	return true
}

func (t Trigger) IsAsync() bool {
	return false
}

func (t Trigger) GetRunnerID() uint64 {
	return 0
}

func (t Trigger) GetTimeout() uint32 {
	return 0
}

func (t Trigger) GetName() string {
	return fmt.Sprintf("%d %s", t.ID, t.Name)
}

func (t Trigger) GetSource() string {
	return t.Source
}

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

func (set ActionSet) Has(action ...string) bool {
	for _, a := range set {
		for _, i := range action {
			if i == a {
				return true
			}
		}
	}

	return false
}

func (set TriggerSet) WalkByAction(action string, fn func(t *Trigger) error) error {
	return set.Walk(func(t *Trigger) error {
		if !t.Actions.Has(action) {
			return nil
		}

		return fn(t)
	})
}

// Resource returns a system resource ID for this type
func (t Trigger) PermissionResource() permissions.Resource {
	return TriggerPermissionResource.AppendID(t.ID)
}
