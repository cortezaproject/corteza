package types

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

type (
	Credentials struct {
		ID          uint64          `json:"credentialsID,string" db:"id"`
		OwnerID     uint64          `json:"ownerID,string" db:"rel_owner"`
		Label       string          `json:"label" db:"label"`
		Kind        CredentialsKind `json:"kind" db:"kind"`
		Credentials string          `json:"-" db:"credentials"`
		Meta        types.JSONText  `json:"-" db:"meta"`
		ExpiresAt   *time.Time      `json:"expiresAt,omitempty" db:"expires_at"`
		CreatedAt   time.Time       `json:"createdAt,omitempty" db:"created_at"`
		DeletedAt   *time.Time      `json:"deletedAt,omitempty" db:"deleted_at"`
	}

	CredentialsKind string
)

const (
	// Use as a password for users or as API secret for bots (and credentials-id as a key) as a value for "credentials"
	CredentialsKindHash CredentialsKind = "hash"

	// Identity (profile-id) stored under "credentials"
	CredentialsKindFacebook CredentialsKind = "facebook"
	CredentialsKindGPlus    CredentialsKind = "gplus"
	CredentialsKindGitHub   CredentialsKind = "github"
	CredentialsKindLinkedin CredentialsKind = "linkedin"
	// CredentialsKindSatosa   CredentialsKind = "satosa"
)

func (u *Credentials) Valid() bool {
	return u.ID > 0 && (u.ExpiresAt == nil || u.ExpiresAt.Before(time.Now())) && u.DeletedAt == nil
}
