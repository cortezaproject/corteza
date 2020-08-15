package types

import (
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"time"

	"github.com/jmoiron/sqlx/types"
	"golang.org/x/crypto/bcrypt"
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

	CredentialsFilter struct {
		OwnerID     uint64         `json:"ownerID"`
		Kind        string         `json:"kind"`
		Credentials string         `json:"credentials"`
		Deleted     rh.FilterState `json:"deleted"`
	}
)

func (u *Credentials) Valid() bool {
	return u.ID > 0 && (u.ExpiresAt == nil || u.ExpiresAt.After(time.Now())) && u.DeletedAt == nil
}

// CompareHashAndPassword returns first valid credentials with matching hash
func (cc CredentialsSet) CompareHashAndPassword(password string) *Credentials {
	// We need only valid credentials (skip deleted, expired)
	for _, c := range cc {
		if !c.Valid() {
			continue
		}

		if len(c.Credentials) == 0 {
			continue
		}

		if bcrypt.CompareHashAndPassword([]byte(c.Credentials), []byte(password)) == nil {
			return c
		}
	}

	return nil
}
