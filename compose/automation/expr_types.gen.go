package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// compose/automation/expr_types.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
	. "github.com/cortezaproject/corteza-server/pkg/expr"
)

var _ = context.Background
var _ = fmt.Errorf

// ComposeModule is an expression type, wrapper for *types.Module type
type ComposeModule struct{ value *types.Module }

// NewComposeModule creates new instance of ComposeModule expression type
func NewComposeModule(val interface{}) (*ComposeModule, error) {
	if c, err := CastToComposeModule(val); err != nil {
		return nil, fmt.Errorf("unable to create ComposeModule: %w", err)
	} else {
		return &ComposeModule{value: c}, nil
	}
}

// Return underlying value on ComposeModule
func (t ComposeModule) Get() interface{} { return t.value }

// Return underlying value on ComposeModule
func (t ComposeModule) GetValue() *types.Module { return t.value }

// Return type name
func (ComposeModule) Type() string { return "ComposeModule" }

// Convert value to *types.Module
func (ComposeModule) Cast(val interface{}) (TypedValue, error) {
	return NewComposeModule(val)
}

// Assign new value to ComposeModule
//
// value is first passed through CastToComposeModule
func (t *ComposeModule) Assign(val interface{}) error {
	if c, err := CastToComposeModule(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

func (t *ComposeModule) AssignFieldValue(key string, val interface{}) error {
	return assignToComposeModule(t.value, key, val)
}

// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access ComposeModule's underlying value (*types.Module)
// and it's fields
//
func (t ComposeModule) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	return composeModuleGValSelector(t.value, k)
}

// Select is field accessor for *types.Module
//
// Similar to SelectGVal but returns typed values
func (t ComposeModule) Select(k string) (TypedValue, error) {
	return composeModuleTypedValueSelector(t.value, k)
}

func (t ComposeModule) Has(k string) bool {
	switch k {
	case "ID":
		return true
	case "namespaceID":
		return true
	case "name":
		return true
	case "handle":
		return true
	case "labels":
		return true
	case "createdAt":
		return true
	case "updatedAt":
		return true
	case "deletedAt":
		return true
	}
	return false
}

// composeModuleGValSelector is field accessor for *types.Module
func composeModuleGValSelector(res *types.Module, k string) (interface{}, error) {
	switch k {
	case "ID":
		return res.ID, nil
	case "namespaceID":
		return res.NamespaceID, nil
	case "name":
		return res.Name, nil
	case "handle":
		return res.Handle, nil
	case "labels":
		return res.Labels, nil
	case "createdAt":
		return res.CreatedAt, nil
	case "updatedAt":
		return res.UpdatedAt, nil
	case "deletedAt":
		return res.DeletedAt, nil
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// composeModuleTypedValueSelector is field accessor for *types.Module
func composeModuleTypedValueSelector(res *types.Module, k string) (TypedValue, error) {
	switch k {
	case "ID":
		return NewID(res.ID)
	case "namespaceID":
		return NewID(res.NamespaceID)
	case "name":
		return NewString(res.Name)
	case "handle":
		return NewHandle(res.Handle)
	case "labels":
		return NewKV(res.Labels)
	case "createdAt":
		return NewDateTime(res.CreatedAt)
	case "updatedAt":
		return NewDateTime(res.UpdatedAt)
	case "deletedAt":
		return NewDateTime(res.DeletedAt)
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// assignToComposeModule is field value setter for *types.Module
func assignToComposeModule(res *types.Module, k string, val interface{}) error {
	switch k {
	case "ID":
		return fmt.Errorf("field '%s' is read-only", k)
	case "namespaceID":
		return fmt.Errorf("field '%s' is read-only", k)
	case "name":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Name = aux
		return nil
	case "handle":
		aux, err := CastToHandle(val)
		if err != nil {
			return err
		}

		res.Handle = aux
		return nil
	case "labels":
		aux, err := CastToKV(val)
		if err != nil {
			return err
		}

		res.Labels = aux
		return nil
	case "createdAt":
		return fmt.Errorf("field '%s' is read-only", k)
	case "updatedAt":
		return fmt.Errorf("field '%s' is read-only", k)
	case "deletedAt":
		return fmt.Errorf("field '%s' is read-only", k)
	}

	return fmt.Errorf("unknown field '%s'", k)
}

// ComposeNamespace is an expression type, wrapper for *types.Namespace type
type ComposeNamespace struct{ value *types.Namespace }

// NewComposeNamespace creates new instance of ComposeNamespace expression type
func NewComposeNamespace(val interface{}) (*ComposeNamespace, error) {
	if c, err := CastToComposeNamespace(val); err != nil {
		return nil, fmt.Errorf("unable to create ComposeNamespace: %w", err)
	} else {
		return &ComposeNamespace{value: c}, nil
	}
}

// Return underlying value on ComposeNamespace
func (t ComposeNamespace) Get() interface{} { return t.value }

// Return underlying value on ComposeNamespace
func (t ComposeNamespace) GetValue() *types.Namespace { return t.value }

// Return type name
func (ComposeNamespace) Type() string { return "ComposeNamespace" }

// Convert value to *types.Namespace
func (ComposeNamespace) Cast(val interface{}) (TypedValue, error) {
	return NewComposeNamespace(val)
}

// Assign new value to ComposeNamespace
//
// value is first passed through CastToComposeNamespace
func (t *ComposeNamespace) Assign(val interface{}) error {
	if c, err := CastToComposeNamespace(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

func (t *ComposeNamespace) AssignFieldValue(key string, val interface{}) error {
	return assignToComposeNamespace(t.value, key, val)
}

// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access ComposeNamespace's underlying value (*types.Namespace)
// and it's fields
//
func (t ComposeNamespace) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	return composeNamespaceGValSelector(t.value, k)
}

// Select is field accessor for *types.Namespace
//
// Similar to SelectGVal but returns typed values
func (t ComposeNamespace) Select(k string) (TypedValue, error) {
	return composeNamespaceTypedValueSelector(t.value, k)
}

func (t ComposeNamespace) Has(k string) bool {
	switch k {
	case "ID":
		return true
	case "name":
		return true
	case "slug", "handle":
		return true
	case "labels":
		return true
	case "createdAt":
		return true
	case "updatedAt":
		return true
	case "deletedAt":
		return true
	}
	return false
}

// composeNamespaceGValSelector is field accessor for *types.Namespace
func composeNamespaceGValSelector(res *types.Namespace, k string) (interface{}, error) {
	switch k {
	case "ID":
		return res.ID, nil
	case "name":
		return res.Name, nil
	case "slug", "handle":
		return res.Slug, nil
	case "labels":
		return res.Labels, nil
	case "createdAt":
		return res.CreatedAt, nil
	case "updatedAt":
		return res.UpdatedAt, nil
	case "deletedAt":
		return res.DeletedAt, nil
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// composeNamespaceTypedValueSelector is field accessor for *types.Namespace
func composeNamespaceTypedValueSelector(res *types.Namespace, k string) (TypedValue, error) {
	switch k {
	case "ID":
		return NewID(res.ID)
	case "name":
		return NewString(res.Name)
	case "slug", "handle":
		return NewHandle(res.Slug)
	case "labels":
		return NewKV(res.Labels)
	case "createdAt":
		return NewDateTime(res.CreatedAt)
	case "updatedAt":
		return NewDateTime(res.UpdatedAt)
	case "deletedAt":
		return NewDateTime(res.DeletedAt)
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// assignToComposeNamespace is field value setter for *types.Namespace
func assignToComposeNamespace(res *types.Namespace, k string, val interface{}) error {
	switch k {
	case "ID":
		return fmt.Errorf("field '%s' is read-only", k)
	case "name":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Name = aux
		return nil
	case "slug", "handle":
		aux, err := CastToHandle(val)
		if err != nil {
			return err
		}

		res.Slug = aux
		return nil
	case "labels":
		aux, err := CastToKV(val)
		if err != nil {
			return err
		}

		res.Labels = aux
		return nil
	case "createdAt":
		return fmt.Errorf("field '%s' is read-only", k)
	case "updatedAt":
		return fmt.Errorf("field '%s' is read-only", k)
	case "deletedAt":
		return fmt.Errorf("field '%s' is read-only", k)
	}

	return fmt.Errorf("unknown field '%s'", k)
}

// ComposeRecord is an expression type, wrapper for *types.Record type
type ComposeRecord struct{ value *types.Record }

// NewComposeRecord creates new instance of ComposeRecord expression type
func NewComposeRecord(val interface{}) (*ComposeRecord, error) {
	if c, err := CastToComposeRecord(val); err != nil {
		return nil, fmt.Errorf("unable to create ComposeRecord: %w", err)
	} else {
		return &ComposeRecord{value: c}, nil
	}
}

// Return underlying value on ComposeRecord
func (t ComposeRecord) Get() interface{} { return t.value }

// Return underlying value on ComposeRecord
func (t ComposeRecord) GetValue() *types.Record { return t.value }

// Return type name
func (ComposeRecord) Type() string { return "ComposeRecord" }

// Convert value to *types.Record
func (ComposeRecord) Cast(val interface{}) (TypedValue, error) {
	return NewComposeRecord(val)
}

// Assign new value to ComposeRecord
//
// value is first passed through CastToComposeRecord
func (t *ComposeRecord) Assign(val interface{}) error {
	if c, err := CastToComposeRecord(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

func (t *ComposeRecord) AssignFieldValue(key string, val interface{}) error {
	return assignToComposeRecord(t.value, key, val)
}

// Select is field accessor for *types.Record
//
// Similar to SelectGVal but returns typed values
func (t ComposeRecord) Select(k string) (TypedValue, error) {
	return composeRecordTypedValueSelector(t.value, k)
}

func (t ComposeRecord) Has(k string) bool {
	switch k {
	case "ID":
		return true
	case "moduleID":
		return true
	case "namespaceID":
		return true
	case "values":
		return true
	case "labels":
		return true
	case "ownedBy":
		return true
	case "createdAt":
		return true
	case "createdBy":
		return true
	case "updatedAt":
		return true
	case "updatedBy":
		return true
	case "deletedAt":
		return true
	case "deletedBy":
		return true
	}
	return false
}

// composeRecordGValSelector is field accessor for *types.Record
func composeRecordGValSelector(res *types.Record, k string) (interface{}, error) {
	switch k {
	case "ID":
		return res.ID, nil
	case "moduleID":
		return res.ModuleID, nil
	case "namespaceID":
		return res.NamespaceID, nil
	case "values":
		return res.Values, nil
	case "labels":
		return res.Labels, nil
	case "ownedBy":
		return res.OwnedBy, nil
	case "createdAt":
		return res.CreatedAt, nil
	case "createdBy":
		return res.CreatedBy, nil
	case "updatedAt":
		return res.UpdatedAt, nil
	case "updatedBy":
		return res.UpdatedBy, nil
	case "deletedAt":
		return res.DeletedAt, nil
	case "deletedBy":
		return res.DeletedBy, nil
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// composeRecordTypedValueSelector is field accessor for *types.Record
func composeRecordTypedValueSelector(res *types.Record, k string) (TypedValue, error) {
	switch k {
	case "ID":
		return NewID(res.ID)
	case "moduleID":
		return NewID(res.ModuleID)
	case "namespaceID":
		return NewID(res.NamespaceID)
	case "values":
		return NewComposeRecordValues(res.Values)
	case "labels":
		return NewKV(res.Labels)
	case "ownedBy":
		return NewID(res.OwnedBy)
	case "createdAt":
		return NewDateTime(res.CreatedAt)
	case "createdBy":
		return NewID(res.CreatedBy)
	case "updatedAt":
		return NewDateTime(res.UpdatedAt)
	case "updatedBy":
		return NewID(res.UpdatedBy)
	case "deletedAt":
		return NewDateTime(res.DeletedAt)
	case "deletedBy":
		return NewID(res.DeletedBy)
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// assignToComposeRecord is field value setter for *types.Record
func assignToComposeRecord(res *types.Record, k string, val interface{}) error {
	switch k {
	case "ID":
		return fmt.Errorf("field '%s' is read-only", k)
	case "moduleID":
		return fmt.Errorf("field '%s' is read-only", k)
	case "namespaceID":
		return fmt.Errorf("field '%s' is read-only", k)
	case "values":
		aux, err := CastToComposeRecordValues(val)
		if err != nil {
			return err
		}

		res.Values = aux
		return nil
	case "labels":
		aux, err := CastToKV(val)
		if err != nil {
			return err
		}

		res.Labels = aux
		return nil
	case "ownedBy":
		aux, err := CastToID(val)
		if err != nil {
			return err
		}

		res.OwnedBy = aux
		return nil
	case "createdAt":
		return fmt.Errorf("field '%s' is read-only", k)
	case "createdBy":
		return fmt.Errorf("field '%s' is read-only", k)
	case "updatedAt":
		return fmt.Errorf("field '%s' is read-only", k)
	case "updatedBy":
		return fmt.Errorf("field '%s' is read-only", k)
	case "deletedAt":
		return fmt.Errorf("field '%s' is read-only", k)
	case "deletedBy":
		return fmt.Errorf("field '%s' is read-only", k)
	}

	return fmt.Errorf("unknown field '%s'", k)
}

// ComposeRecordValueErrorSet is an expression type, wrapper for *types.RecordValueErrorSet type
type ComposeRecordValueErrorSet struct{ value *types.RecordValueErrorSet }

// NewComposeRecordValueErrorSet creates new instance of ComposeRecordValueErrorSet expression type
func NewComposeRecordValueErrorSet(val interface{}) (*ComposeRecordValueErrorSet, error) {
	if c, err := CastToComposeRecordValueErrorSet(val); err != nil {
		return nil, fmt.Errorf("unable to create ComposeRecordValueErrorSet: %w", err)
	} else {
		return &ComposeRecordValueErrorSet{value: c}, nil
	}
}

// Return underlying value on ComposeRecordValueErrorSet
func (t ComposeRecordValueErrorSet) Get() interface{} { return t.value }

// Return underlying value on ComposeRecordValueErrorSet
func (t ComposeRecordValueErrorSet) GetValue() *types.RecordValueErrorSet { return t.value }

// Return type name
func (ComposeRecordValueErrorSet) Type() string { return "ComposeRecordValueErrorSet" }

// Convert value to *types.RecordValueErrorSet
func (ComposeRecordValueErrorSet) Cast(val interface{}) (TypedValue, error) {
	return NewComposeRecordValueErrorSet(val)
}

// Assign new value to ComposeRecordValueErrorSet
//
// value is first passed through CastToComposeRecordValueErrorSet
func (t *ComposeRecordValueErrorSet) Assign(val interface{}) error {
	if c, err := CastToComposeRecordValueErrorSet(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// ComposeRecordValues is an expression type, wrapper for types.RecordValueSet type
type ComposeRecordValues struct{ value types.RecordValueSet }

// NewComposeRecordValues creates new instance of ComposeRecordValues expression type
func NewComposeRecordValues(val interface{}) (*ComposeRecordValues, error) {
	if c, err := CastToComposeRecordValues(val); err != nil {
		return nil, fmt.Errorf("unable to create ComposeRecordValues: %w", err)
	} else {
		return &ComposeRecordValues{value: c}, nil
	}
}

// Return underlying value on ComposeRecordValues
func (t ComposeRecordValues) Get() interface{} { return t.value }

// Return underlying value on ComposeRecordValues
func (t ComposeRecordValues) GetValue() types.RecordValueSet { return t.value }

// Return type name
func (ComposeRecordValues) Type() string { return "ComposeRecordValues" }

// Convert value to types.RecordValueSet
func (ComposeRecordValues) Cast(val interface{}) (TypedValue, error) {
	return NewComposeRecordValues(val)
}

// Assign new value to ComposeRecordValues
//
// value is first passed through CastToComposeRecordValues
func (t *ComposeRecordValues) Assign(val interface{}) error {
	if c, err := CastToComposeRecordValues(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}
