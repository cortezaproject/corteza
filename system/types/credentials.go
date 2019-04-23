package types

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

type (
	Credentials struct {
		ID          uint64         `json:"credentialsID,string" db:"id"`
		OwnerID     uint64         `json:"ownerID,string" db:"rel_owner"`
		Label       string         `json:"label" db:"label"`
		Kind        string         `json:"kind" db:"kind"`
		Credentials string         `json:"-" db:"credentials"`
		Meta        types.JSONText `json:"-" db:"meta"`
		LastUsedAt  *time.Time     `json:"lastUsedAt,omitempty" db:"last_used_at"`
		ExpiresAt   *time.Time     `json:"expiresAt,omitempty"  db:"expires_at"`
		CreatedAt   time.Time      `json:"createdAt,omitempty"  db:"created_at"`
		UpdatedAt   *time.Time     `json:"updatedAt,omitempty"  db:"updated_at"`
		DeletedAt   *time.Time     `json:"deletedAt,omitempty"  db:"deleted_at"`
	}
)

func (u *Credentials) Valid() bool {
	return u.ID > 0 && (u.ExpiresAt == nil || u.ExpiresAt.After(time.Now())) && u.DeletedAt == nil
}
