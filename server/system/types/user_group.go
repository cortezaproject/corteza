package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/sql"

	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	UserGroup struct {
		ID     uint64 `json:"userGroupID,string"`
		Handle string `json:"handle"`

		SelfID uint64 `json:"selfID,string"`

		Meta   *UserGroupMeta    `json:"meta"`
		Labels map[string]string `json:"labels,omitempty"`

		CreatedAt  time.Time  `json:"createdAt,omitempty"`
		UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
		ArchivedAt *time.Time `json:"archivedAt,omitempty"`
		DeletedAt  *time.Time `json:"deletedAt,omitempty"`
	}

	UserGroupMeta struct {
		Description string `json:"description"`
		Short       string `json:"short"`
	}

	UserGroupFilter struct {
		UserGroupID []string `json:"userGroupID"`
		MemberID    uint64   `json:"memberID,string"`

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
		Check func(*UserGroup) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}
)

func (r *UserGroup) Clone() *UserGroup {
	if r == nil {
		return nil
	}

	return &UserGroup{
		ID:         r.ID,
		Handle:     r.Handle,
		Meta:       r.Meta,
		Labels:     r.Labels,
		ArchivedAt: r.ArchivedAt,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
		DeletedAt:  r.DeletedAt,
	}
}

// FindByHandle finds userGroup by it's handle
func (set UserGroupSet) FindByHandle(handle string) *UserGroup {
	for i := range set {
		if set[i].Handle == handle {
			return set[i]
		}
	}

	return nil
}

func (vv *UserGroupMeta) Scan(src any) error           { return sql.ParseJSON(src, vv) }
func (vv *UserGroupMeta) Value() (driver.Value, error) { return json.Marshal(vv) }
