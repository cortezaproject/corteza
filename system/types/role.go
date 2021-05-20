package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	Role struct {
		ID     uint64 `json:"roleID,string"`
		Name   string `json:"name"`
		Handle string `json:"handle"`

		Meta   *RoleMeta         `json:"meta"`
		Labels map[string]string `json:"labels,omitempty"`

		CreatedAt  time.Time  `json:"createdAt,omitempty"`
		UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
		ArchivedAt *time.Time `json:"archivedAt,omitempty"`
		DeletedAt  *time.Time `json:"deletedAt,omitempty"`
	}

	RoleMeta struct {
		Description string       `json:"description,omitempty"`
		Context     *RoleContext `json:"context,omitempty"`
	}

	RoleContext struct {
		Resource []string `json:"resource,omitempty"`
		Expr     string   `json:"expr,omitempty"`
	}

	RoleFilter struct {
		RoleID   []uint64 `json:"roleID"`
		MemberID uint64   `json:"memberID"`

		Query string `json:"query"`

		Handle string `json:"handle"`
		Name   string `json:"name"`

		Deleted  filter.State `json:"deleted"`
		Archived filter.State `json:"archived"`

		LabeledIDs []uint64          `json:"-"`
		Labels     map[string]string `json:"labels,omitempty"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Role) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	RoleMetrics struct {
		Total         uint   `json:"total"`
		Valid         uint   `json:"valid"`
		Deleted       uint   `json:"deleted"`
		Archived      uint   `json:"archived"`
		DailyCreated  []uint `json:"dailyCreated"`
		DailyDeleted  []uint `json:"dailyDeleted"`
		DailyUpdated  []uint `json:"dailyUpdated"`
		DailyArchived []uint `json:"dailyArchived"`
	}
)

func (r *Role) DynamicRoles(userID uint64) []uint64 {
	return nil
}

// FindByHandle finds role by it's handle
func (set RoleSet) FindByHandle(handle string) *Role {
	for i := range set {
		if set[i].Handle == handle {
			return set[i]
		}
	}

	return nil
}

func (vv *RoleMeta) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*vv = RoleMeta{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, vv); err != nil {
			return fmt.Errorf("cannot scan '%v' into RoleMeta: %w", string(b), err)
		}
	}

	return nil
}

// Scan on RoleMeta gracefully handles conversion from NULL
func (vv *RoleMeta) Value() (driver.Value, error) {
	if vv == nil {
		return []byte("null"), nil
	}

	return json.Marshal(vv)
}
