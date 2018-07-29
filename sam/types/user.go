package types

import (
	"encoding/json"
	"time"
)

type (
	User struct {
		ID             uint64          `db:"id"`
		Username       string          `db:"username"`
		Meta           json.RawMessage `json:"-" db:"meta"`
		OrganisationID uint64          `db:"rel_organisation"`
		Password       []byte          `json:"-" db:"password"`
		CreatedAt      time.Time       `json:"created_at,omitempty" db:"created_at"`
		UpdatedAt      *time.Time      `json:"updated_at,omitempty" db:"updated_at"`
		SuspendedAt    *time.Time      `json:"suspended_at,omitempty" db:"suspended_at"`
		DeletedAt      *time.Time      `json:"deleted_at,omitempty" db:"deleted_at"`
	}

	UserFilter struct {
		Query            string
		MembersOfChannel uint64
	}
)

func (u *User) Valid() bool {
	return u.ID > 0 && u.SuspendedAt == nil && u.DeletedAt == nil
}

func (u *User) Identity() uint64 {
	return u.ID
}
