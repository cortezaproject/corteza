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
	// Users
	User struct {
		ID          uint64     `db:"id"`
		Username    string     `db:"username"`
		Password    []byte     `json:"-" db:"password"`
		SuspendedAt *time.Time `json:",omitempty" db:"suspended_at"`
		DeletedAt   *time.Time `json:",omitempty" db:"deleted_at"`

		changed []string
	}
)

/* Constructors */
func (User) New() *User {
	return &User{}
}

/* Getters/setters */
func (u *User) GetID() uint64 {
	return u.ID
}

func (u *User) SetID(value uint64) *User {
	if u.ID != value {
		u.changed = append(u.changed, "ID")
		u.ID = value
	}
	return u
}
func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) SetUsername(value string) *User {
	if u.Username != value {
		u.changed = append(u.changed, "Username")
		u.Username = value
	}
	return u
}
func (u *User) GetPassword() []byte {
	return u.Password
}

func (u *User) SetPassword(value []byte) *User {
	u.Password = value
	return u
}
func (u *User) GetSuspendedAt() *time.Time {
	return u.SuspendedAt
}

func (u *User) SetSuspendedAt(value *time.Time) *User {
	u.SuspendedAt = value
	return u
}
func (u *User) GetDeletedAt() *time.Time {
	return u.DeletedAt
}

func (u *User) SetDeletedAt(value *time.Time) *User {
	u.DeletedAt = value
	return u
}
