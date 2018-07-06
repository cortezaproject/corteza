package sam

import (
	"time"
)

type (
	// Users
	User struct {
		ID          uint64
		Username    string
		Password    []byte     `json:"-"`
		SuspendedAt *time.Time `json:",omitempty"`
		DeletedAt   *time.Time `json:",omitempty"`

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
	if u.Password != value {
		u.changed = append(u.changed, "Password")
		u.Password = value
	}
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
