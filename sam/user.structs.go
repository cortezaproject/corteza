package sam

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Users
type User struct {
	ID       uint64
	Username string
	Password []byte `json:"-"`

	SuspendedAt *time.Time `json:",omitempty"`
	DeletedAt   *time.Time `json:",omitempty"`

	changed []string
}

func (User) new() *User {
	return &User{}
}

// Basic user validation
func (u *User) CanLogin() bool {
	return u.ID > 0
}

func (u *User) GetID() uint64 {
	return u.ID
}

func (u *User) SetID(value uint64) *User {
	if u.ID != value {
		u.changed = append(u.changed, "id")
		u.ID = value
	}
	return u
}
func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) SetUsername(value string) *User {
	if u.Username != value {
		u.changed = append(u.changed, "username")
		u.Username = value
	}
	return u
}

func (u *User) SetPassword(value string) error {
	if !u.ValidatePassword(value) {
		if encrypted, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost); err != nil {
			return err
		} else {
			u.Password = encrypted
			u.changed = append(u.changed, "password")
		}
	}

	return nil
}

func (u User) ValidatePassword(value string) bool {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(value)) == nil
}
