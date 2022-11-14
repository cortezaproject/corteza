package types

import (
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/jmoiron/sqlx/types"
)

type (
	Credential struct {
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

	CredentialFilter struct {
		OwnerID     uint64       `json:"ownerID"`
		Kind        string       `json:"kind"`
		Credentials string       `json:"credentials"`
		Deleted     filter.State `json:"deleted"`
		Limit       uint
	}
)

func (u *Credential) Valid() bool {
	return u.ID > 0 && (u.ExpiresAt == nil || u.ExpiresAt.After(time.Now())) && u.DeletedAt == nil
}
