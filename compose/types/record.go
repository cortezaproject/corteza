package types

import (
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	// Record is a stored row in the `record` table
	Record struct {
		ID       uint64 `json:"recordID,string" db:"id"`
		ModuleID uint64 `json:"moduleID,string" db:"module_id"`

		Values RecordValueSet `json:"values,omitempty" db:"-"`

		NamespaceID uint64 `json:"namespaceID,string" db:"rel_namespace"`

		OwnedBy   uint64     `db:"owned_by"   json:"ownedBy,string"`
		CreatedAt time.Time  `db:"created_at" json:"createdAt,omitempty"`
		CreatedBy uint64     `db:"created_by" json:"createdBy,string" `
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		UpdatedBy uint64     `db:"updated_by" json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
		DeletedBy uint64     `db:"deleted_by" json:"deletedBy,string,omitempty" `
	}

	RecordFilter struct {
		ModuleID    uint64 `json:"moduleID,string"`
		NamespaceID uint64 `json:"namespaceID,string"`
		Filter      string `json:"query"`
		Sort        string `json:"sort"`

		// Standard paging fields & helpers
		rh.PageFilter
	}
)

// UserIDs returns a slice of user IDs from all items in the set
func (set RecordSet) UserIDs() (IDs []uint64) {
	IDs = make([]uint64, 0)

loop:
	for i := range set {
		for _, id := range IDs {
			if id == set[i].OwnedBy {
				continue loop
			}
		}

		IDs = append(IDs, set[i].OwnedBy)
	}

	return
}

// Resource returns a system resource ID for this type
func (r Record) PermissionResource() permissions.Resource {
	return ModulePermissionResource.AppendID(r.ModuleID)
}

// UnmarshalJSON for custom record deserialization
//
// Due to https://github.com/golang/go/issues/21092, we should manually reset the given record value set.
// If this is skipped there is a chance of data corruption; ie. wrong value is removed/edited
func (r *Record) UnmarshalJSON(data []byte) error {
	// Reset value set
	r.Values = nil

	// Deserialize to r (*Record) via auxRecord auxiliary record type alias
	//
	// This prevents inf. loop where json.Unmarshal directly on Record type
	// calls this function
	type auxRecord Record
	return json.Unmarshal(data, &struct{ *auxRecord }{auxRecord: (*auxRecord)(r)})
}
