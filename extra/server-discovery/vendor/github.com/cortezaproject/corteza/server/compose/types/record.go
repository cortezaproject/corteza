package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/modern-go/reflect2"
	"github.com/spf13/cast"
)

type (
	bulkPreProcess func(m, r *Record) (*Record, error)

	OperationType string

	RecordBulkOperation struct {
		Record      *Record
		RecordID    uint64
		NamespaceID uint64
		ModuleID    uint64

		LinkBy    string
		Operation OperationType
		ID        string
	}
	RecordBulkOperationResult struct {
		Record           *Record
		Error            error
		ValueError       *RecordValueErrorSet
		DuplicationError *RecordValueErrorSet
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

		Revision uint `json:"revision,omitempty"`

		module *Module

		Values RecordValueSet `json:"values,omitempty"`

		Meta map[string]any `json:"meta,omitempty"`

		NamespaceID uint64 `json:"namespaceID,string"`

		OwnedBy   uint64     `json:"ownedBy,string"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
	}

	RecordFilter struct {
		ModuleID    uint64 `json:"moduleID,string"`
		NamespaceID uint64 `json:"namespaceID,string"`
		Query       string `json:"query"`

		Meta map[string]any `json:"meta,omitempty"`

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

	SensitiveRecord struct {
		RecordID uint64
		Values   []map[string]any
	}

	SensitiveRecordSet struct {
		// Contextual metadata
		ConnectionID uint64
		Module       *Module
		Namespace    *Namespace

		Records []SensitiveRecord
	}
)

const (
	OperationTypeCreate   OperationType = "create"
	OperationTypeUpdate   OperationType = "update"
	OperationTypeDelete   OperationType = "delete"
	OperationTypePatch    OperationType = "patch"
	OperationTypeUndelete OperationType = "undelete"
)

func (f RecordFilter) ToConstraintedFilter(c map[string][]any) filter.Filter {
	return filter.Generic(
		// combine constraints with namespace and module
		filter.WithConstraints(c),
		filter.WithExpression(f.Query),
		filter.WithOrderBy(f.Sort),
		filter.WithLimit(f.Limit),
		filter.WithCursor(f.PageCursor),
		filter.WithMetaConstraints(f.Meta),
		filter.WithStateConstraint("deletedAt", f.Deleted),
	)
}

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

func (r *Record) getValue(name string, pos uint) (any, error) {
	if val := r.Values.Get(name, pos); val != nil {
		return val.Value, nil
	}

	return nil, nil
}

// CountValues returns how many values per field are there
func (r *Record) CountValues() (pos map[string]uint) {
	var (
		mod = r.GetModule()
	)

	pos = map[string]uint{
		"ID":          1,
		"moduleID":    1,
		"namespaceID": 1,
		"meta":        1,
		"revision":    1,
		"createdAt":   1,
		"createdBy":   1,
		"updatedAt":   1,
		"updatedBy":   1,
		"deletedAt":   1,
		"deletedBy":   1,
		"ownedBy":     1,
	}

	// if mod == nil {
	// 	// count record values
	// 	// only when module is known
	// 	return
	// }

	if mod != nil {
		for _, f := range mod.Fields {
			pos[f.Name] = 0
		}
	}

	for _, val := range r.Values {
		pos[val.Name]++
	}

	return
}

func (r *Record) setValue(name string, pos uint, value any) (err error) {
	if reflect2.IsNil(value) {
		r.Values, _ = r.Values.Filter(func(rv *RecordValue) (bool, error) {
			if rv.Name == name && rv.Place == pos {
				return false, nil
			}

			return true, nil
		})

		return
	}

	rv := &RecordValue{Name: name, Place: pos}

	if cv, ok := value.(*RecordValue); ok {
		rv = cv
	} else {
		var auxv string
		switch aux := value.(type) {
		case *time.Time:
			auxv = aux.Format(time.RFC3339)

		case time.Time:
			auxv = aux.Format(time.RFC3339)

		default:
			auxv, err = cast.ToStringE(aux)
		}
		if err != nil {
			return
		}

		rv.Value = auxv
	}

	// Try to utilize the module when possible
	// It can be omitted for some cases for easier test cases
	if r.module != nil {
		f := r.module.Fields.FindByName(name)
		if f != nil {
			switch f.Kind {
			case "Record", "User", "File":
				rv.Ref = cast.ToUint64(value)
			}
		}
	}

	r.Values = r.Values.Set(rv)

	return
}

func (r Record) Dict() map[string]interface{} {
	dict := map[string]interface{}{
		"ID":          r.ID,
		"recordID":    r.ID,
		"moduleID":    r.ModuleID,
		"revision":    r.Revision,
		"meta":        r.Meta,
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

// GetValuesByName filters values for records by names
func (set RecordSet) GetValuesByName(names ...string) (out RecordValueSet) {
	nameMap := make(map[string]bool)
	for _, n := range names {
		if len(n) > 0 {
			nameMap[n] = true
		}
	}

	err := set.Walk(func(rec *Record) error {
		_ = rec.Values.Walk(func(val *RecordValue) error {
			if val != nil && nameMap[val.Name] {
				val.RecordID = rec.ID
				out = append(out, val)
			}
			return nil
		})
		return nil
	})
	if err != nil {
		return
	}

	return
}
