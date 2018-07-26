package types

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `user.go`, `user.util.go` or `user_test.go` to
	implement your API calls, helper functions and tests. The file `user.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"time"
)

type (
	// Users -
	User struct {
		ID             uint64      `db:"id"`
		Username       string      `db:"username"`
		Meta           interface{} `json:"-" db:"meta"`
		OrganisationID uint64      `db:"rel_organisation"`
		Password       []byte      `json:"-" db:"password"`
		CreatedAt      time.Time   `json:"created_at,omitempty" db:"created_at"`
		UpdatedAt      *time.Time  `json:"updated_at,omitempty" db:"updated_at"`
		SuspendedAt    *time.Time  `json:"suspended_at,omitempty" db:"suspended_at"`
		DeletedAt      *time.Time  `json:"deleted_at,omitempty" db:"deleted_at"`

		changed []string
	}
)

// New constructs a new instance of User
func (User) New() *User {
	return &User{}
}

// Get the value of ID
func (u *User) GetID() uint64 {
	return u.ID
}

// Set the value of ID
func (u *User) SetID(value uint64) *User {
	if u.ID != value {
		u.changed = append(u.changed, "ID")
		u.ID = value
	}
	return u
}

// Get the value of Username
func (u *User) GetUsername() string {
	return u.Username
}

// Set the value of Username
func (u *User) SetUsername(value string) *User {
	if u.Username != value {
		u.changed = append(u.changed, "Username")
		u.Username = value
	}
	return u
}

// Get the value of Meta
func (u *User) GetMeta() interface{} {
	return u.Meta
}

// Set the value of Meta
func (u *User) SetMeta(value interface{}) *User {
	if u.Meta != value {
		u.changed = append(u.changed, "Meta")
		u.Meta = value
	}
	return u
}

// Get the value of OrganisationID
func (u *User) GetOrganisationID() uint64 {
	return u.OrganisationID
}

// Set the value of OrganisationID
func (u *User) SetOrganisationID(value uint64) *User {
	if u.OrganisationID != value {
		u.changed = append(u.changed, "OrganisationID")
		u.OrganisationID = value
	}
	return u
}

// Get the value of Password
func (u *User) GetPassword() []byte {
	return u.Password
}

// Set the value of Password
func (u *User) SetPassword(value []byte) *User {
	u.changed = append(u.changed, "Password")
	u.Password = value
	return u
}

// Get the value of CreatedAt
func (u *User) GetCreatedAt() time.Time {
	return u.CreatedAt
}

// Set the value of CreatedAt
func (u *User) SetCreatedAt(value time.Time) *User {
	u.changed = append(u.changed, "CreatedAt")
	u.CreatedAt = value
	return u
}

// Get the value of UpdatedAt
func (u *User) GetUpdatedAt() *time.Time {
	return u.UpdatedAt
}

// Set the value of UpdatedAt
func (u *User) SetUpdatedAt(value *time.Time) *User {
	u.changed = append(u.changed, "UpdatedAt")
	u.UpdatedAt = value
	return u
}

// Get the value of SuspendedAt
func (u *User) GetSuspendedAt() *time.Time {
	return u.SuspendedAt
}

// Set the value of SuspendedAt
func (u *User) SetSuspendedAt(value *time.Time) *User {
	u.changed = append(u.changed, "SuspendedAt")
	u.SuspendedAt = value
	return u
}

// Get the value of DeletedAt
func (u *User) GetDeletedAt() *time.Time {
	return u.DeletedAt
}

// Set the value of DeletedAt
func (u *User) SetDeletedAt(value *time.Time) *User {
	u.changed = append(u.changed, "DeletedAt")
	u.DeletedAt = value
	return u
}

// Changes returns the names of changed fields
func (u *User) Changes() []string {
	return u.changed
}
