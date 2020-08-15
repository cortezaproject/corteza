package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	bulkPreProcess func(m, r *Record) (*Record, error)

	OperationType string

	RecordBulkOperation struct {
		Record    *Record
		LinkBy    string
		Operation OperationType
		ID        string
	}

	RecordBulk struct {
		RefField string    `json:"refField,omitempty"`
		IDPrefix string    `json:"idPrefix,omitempty"`
		Set      RecordSet `json:"set,omitempty"`
	}

	RecordBulkSet []*RecordBulk

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
		Query       string `json:"query"`
		Sort        string `json:"sort"`

		// Standard paging fields & helpers
		rh.PageFilter

		Deleted rh.FilterState `json:"deleted"`
	}
)

const (
	OperationTypeCreate OperationType = "create"
	OperationTypeUpdate OperationType = "update"
	OperationTypeDelete OperationType = "delete"
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

func (r Record) DynamicRoles(userID uint64) []uint64 {
	return permissions.DynamicRoles(
		userID,
		r.OwnedBy, permissions.OwnersDynamicRoleID,
		r.CreatedBy, permissions.CreatorsDynamicRoleID,
		r.UpdatedBy, permissions.UpdatersDynamicRoleID,
		r.DeletedBy, permissions.DeletersDynamicRoleID,
	)
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

// ToBulkOperations converts BulkRecordSet to a list of BulkRecordOperations
func (set RecordBulkSet) ToBulkOperations(dftModule uint64, dftNamespace uint64) (oo []*RecordBulkOperation, err error) {
	for _, br := range set {
		// can't use for loop's index, since some records can already have an ID
		i := 0
		for _, rr := range br.Set {
			// No use in allowing cross-namespace record creation.
			rr.NamespaceID = dftNamespace

			// default module
			if rr.ModuleID == 0 {
				rr.ModuleID = dftModule
			}
			b := &RecordBulkOperation{
				Record:    rr,
				Operation: OperationTypeUpdate,
				LinkBy:    br.RefField,
			}

			if rr.ID == 0 {
				b.ID = fmt.Sprintf("%s:%d", br.IDPrefix, i)
				i++
			} else {
				b.ID = strconv.FormatUint(rr.ID, 10)
			}

			// If no RecordID is defined, we should create it
			if rr.ID == 0 {
				b.Operation = OperationTypeCreate
			}

			if rr.DeletedAt != nil {
				b.Operation = OperationTypeDelete
			}

			oo = append(oo, b)
		}
	}

	return
}
