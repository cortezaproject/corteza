package types

import (
	"encoding/json"
	"time"
)

type (
	User struct {
		ID             uint64          `json:"id" db:"id"`
		Username       string          `json:"username" db:"username"`
		Meta           json.RawMessage `json:"-" db:"meta"`
		OrganisationID uint64          `json:"organisationId" db:"rel_organisation"`
		Password       []byte          `json:"-" db:"password"`
		CreatedAt      time.Time       `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt      *time.Time      `json:"updatedAt,omitempty" db:"updated_at"`
		SuspendedAt    *time.Time      `json:"suspendedAt,omitempty" db:"suspended_at"`
		DeletedAt      *time.Time      `json:"deletedAt,omitempty" db:"deleted_at"`
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
