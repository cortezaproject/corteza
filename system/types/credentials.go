package types

import (
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/jmoiron/sqlx/types"
	"golang.org/x/crypto/bcrypt"
)

type (
	Credentials struct {
		ID          uint64         `json:"credentialsID,string"`
		OwnerID     uint64         `json:"ownerID,string"`
		Label       string         `json:"label"`
		Kind        string         `json:"kind"`
		Credentials string         `json:"-"`
		Meta        types.JSONText `json:"-"`
		LastUsedAt  *time.Time     `json:"lastUsedAt,omitempty"`
		ExpiresAt   *time.Time     `json:"expiresAt,omitempty"`
		CreatedAt   time.Time      `json:"createdAt,omitempty"`
		UpdatedAt   *time.Time     `json:"updatedAt,omitempty"`
		DeletedAt   *time.Time     `json:"deletedAt,omitempty"`
	}

	CredentialsFilter struct {
		OwnerID     uint64       `json:"ownerID"`
		Kind        string       `json:"kind"`
		Credentials string       `json:"credentials"`
		Deleted     filter.State `json:"deleted"`
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
