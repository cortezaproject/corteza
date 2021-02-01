package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/automation/expr_types.yaml

import (
	"context"
	"fmt"
	. "github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/system/types"
)

var _ = context.Background
var _ = fmt.Errorf

// Role is an expression type, wrapper for *types.Role type
type Role struct{ value *types.Role }

// NewRole creates new instance of Role expression type
func NewRole(val interface{}) (*Role, error) {
	if c, err := CastToRole(val); err != nil {
		return nil, fmt.Errorf("unable to create Role: %w", err)
	} else {
		return &Role{value: c}, nil
	}
}

// Return underlying value on Role
func (t Role) Get() interface{} { return t.value }

// Return type name
func (Role) Type() string { return "Role" }

// Convert value to *types.Role
func (Role) Cast(val interface{}) (TypedValue, error) {
	return NewRole(val)
}

// Assign new value to Role
//
// value is first passed through CastToRole
func (t *Role) Assign(val interface{}) error {
	if c, err := CastToRole(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

func (t *Role) AssignFieldValue(key string, val interface{}) error {
	return assignToRole(t.value, key, val)
}

// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access Role's underlying value (*types.Role)
// and it's fields
//
func (t Role) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	return roleGValSelector(t.value, k)
}

// Select is field accessor for *types.Role
//
// Similar to SelectGVal but returns typed values
func (t Role) Select(k string) (TypedValue, error) {
	return roleTypedValueSelector(t.value, k)
}

func (t Role) Has(k string) bool {
	switch k {
	case "ID":
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
	case "archivedAt":
		return true
	case "deletedAt":
		return true
	}
	return false
}

// roleGValSelector is field accessor for *types.Role
func roleGValSelector(res *types.Role, k string) (interface{}, error) {
	switch k {
	case "ID":
		return res.ID, nil
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
	case "archivedAt":
		return res.ArchivedAt, nil
	case "deletedAt":
		return res.DeletedAt, nil
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// roleTypedValueSelector is field accessor for *types.Role
func roleTypedValueSelector(res *types.Role, k string) (TypedValue, error) {
	switch k {
	case "ID":
		return NewID(res.ID)
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
	case "archivedAt":
		return NewDateTime(res.ArchivedAt)
	case "deletedAt":
		return NewDateTime(res.DeletedAt)
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// assignToRole is field value setter for *types.Role
func assignToRole(res *types.Role, k string, val interface{}) error {
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
	case "archivedAt":
		return fmt.Errorf("field '%s' is read-only", k)
	case "deletedAt":
		return fmt.Errorf("field '%s' is read-only", k)
	}

	return fmt.Errorf("unknown field '%s'", k)
}

// User is an expression type, wrapper for *types.User type
type User struct{ value *types.User }

// NewUser creates new instance of User expression type
func NewUser(val interface{}) (*User, error) {
	if c, err := CastToUser(val); err != nil {
		return nil, fmt.Errorf("unable to create User: %w", err)
	} else {
		return &User{value: c}, nil
	}
}

// Return underlying value on User
func (t User) Get() interface{} { return t.value }

// Return type name
func (User) Type() string { return "User" }

// Convert value to *types.User
func (User) Cast(val interface{}) (TypedValue, error) {
	return NewUser(val)
}

// Assign new value to User
//
// value is first passed through CastToUser
func (t *User) Assign(val interface{}) error {
	if c, err := CastToUser(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

func (t *User) AssignFieldValue(key string, val interface{}) error {
	return assignToUser(t.value, key, val)
}

// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access User's underlying value (*types.User)
// and it's fields
//
func (t User) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	return userGValSelector(t.value, k)
}

// Select is field accessor for *types.User
//
// Similar to SelectGVal but returns typed values
func (t User) Select(k string) (TypedValue, error) {
	return userTypedValueSelector(t.value, k)
}

func (t User) Has(k string) bool {
	switch k {
	case "ID":
		return true
	case "username":
		return true
	case "email":
		return true
	case "name":
		return true
	case "handle":
		return true
	case "emailConfirmed":
		return true
	case "labels":
		return true
	case "createdAt":
		return true
	case "updatedAt":
		return true
	case "suspendedAt":
		return true
	case "deletedAt":
		return true
	}
	return false
}

// userGValSelector is field accessor for *types.User
func userGValSelector(res *types.User, k string) (interface{}, error) {
	switch k {
	case "ID":
		return res.ID, nil
	case "username":
		return res.Username, nil
	case "email":
		return res.Email, nil
	case "name":
		return res.Name, nil
	case "handle":
		return res.Handle, nil
	case "emailConfirmed":
		return res.EmailConfirmed, nil
	case "labels":
		return res.Labels, nil
	case "createdAt":
		return res.CreatedAt, nil
	case "updatedAt":
		return res.UpdatedAt, nil
	case "suspendedAt":
		return res.SuspendedAt, nil
	case "deletedAt":
		return res.DeletedAt, nil
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// userTypedValueSelector is field accessor for *types.User
func userTypedValueSelector(res *types.User, k string) (TypedValue, error) {
	switch k {
	case "ID":
		return NewID(res.ID)
	case "username":
		return NewString(res.Username)
	case "email":
		return NewString(res.Email)
	case "name":
		return NewString(res.Name)
	case "handle":
		return NewHandle(res.Handle)
	case "emailConfirmed":
		return NewBoolean(res.EmailConfirmed)
	case "labels":
		return NewKV(res.Labels)
	case "createdAt":
		return NewDateTime(res.CreatedAt)
	case "updatedAt":
		return NewDateTime(res.UpdatedAt)
	case "suspendedAt":
		return NewDateTime(res.SuspendedAt)
	case "deletedAt":
		return NewDateTime(res.DeletedAt)
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// assignToUser is field value setter for *types.User
func assignToUser(res *types.User, k string, val interface{}) error {
	switch k {
	case "ID":
		return fmt.Errorf("field '%s' is read-only", k)
	case "username":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Username = aux
		return nil
	case "email":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Email = aux
		return nil
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
	case "emailConfirmed":
		aux, err := CastToBoolean(val)
		if err != nil {
			return err
		}

		res.EmailConfirmed = aux
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
	case "suspendedAt":
		return fmt.Errorf("field '%s' is read-only", k)
	case "deletedAt":
		return fmt.Errorf("field '%s' is read-only", k)
	}

	return fmt.Errorf("unknown field '%s'", k)
}
