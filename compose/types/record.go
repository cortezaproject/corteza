package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"

	"github.com/cortezaproject/corteza-server/pkg/rbac"
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
		ID       uint64 `json:"recordID,string"`
		ModuleID uint64 `json:"moduleID,string"`

		module *Module

		Values RecordValueSet `json:"values,omitempty"`

		Labels map[string]string `json:"labels,omitempty"`

		NamespaceID uint64 `json:"namespaceID,string"`

		OwnedBy   uint64     `json:"ownedBy,string"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty" `
	}

	RecordFilter struct {
		ModuleID    uint64 `json:"moduleID,string"`
		NamespaceID uint64 `json:"namespaceID,string"`
		Query       string `json:"query"`

		LabeledIDs []uint64          `json:"-"`
		Labels     map[string]string `json:"labels,omitempty"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Record) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
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

// Sets/updates module ptr
//
// Only if not previously set and if matches record specs
func (r *Record) SetModule(m *Module) {
	if (r.module == nil || r.module.ID == m.ID) && r.ModuleID == m.ID {
		r.module = m
	}
}

func (r *Record) GetModule() *Module {
	return r.module
}

func (r Record) Clone() *Record {
	c := &r
	c.Values = r.Values.Clone()
	return c
}

// Resource returns a system resource ID for this type
func (r Record) RBACResource() rbac.Resource {
	return ModuleRBACResource.AppendID(r.ModuleID)
}

func (r Record) DynamicRoles(userID uint64) []uint64 {
	return rbac.DynamicRoles(
		userID,
		r.OwnedBy, rbac.OwnersDynamicRoleID,
		r.CreatedBy, rbac.CreatorsDynamicRoleID,
		r.UpdatedBy, rbac.UpdatersDynamicRoleID,
		r.DeletedBy, rbac.DeletersDynamicRoleID,
	)
}

func (r Record) Dict() map[string]interface{} {
	dict := map[string]interface{}{
		"ID":          r.ID,
		"moduleID":    r.ModuleID,
		"labels":      r.Labels,
		"namespaceID": r.NamespaceID,
		"ownedBy":     r.OwnedBy,
		"createdAt":   r.CreatedAt,
		"createdBy":   r.CreatedBy,
		"updatedAt":   r.UpdatedAt,
		"updatedBy":   r.UpdatedBy,
		"deletedAt":   r.DeletedAt,
		"deletedBy":   r.DeletedBy,
	}

	if r.module != nil {
		dict["values"] = r.Values.Dict(r.module.Fields)
	}

	return dict
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
