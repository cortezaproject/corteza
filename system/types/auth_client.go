package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/sql"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	AuthClient struct {
		ID uint64 `json:"authClientID,string"`

		// Client's handle
		Handle string `json:"handle"`

		// Client's meta data, see comments on AuthClientMeta
		Meta *AuthClientMeta `json:"meta,omitempty"`

		// Client secret
		Secret string `json:"secret,omitempty"`

		Scope string `json:"scope"`

		// valid grant for this client (only one)
		//  - authorization_code
		//  - client_credentials
		ValidGrant string `json:"validGrant"`

		// Valid redirection URIs
		RedirectURI string `json:"redirectURI"`

		// Users will not be prompted to confirm the client after login
		Trusted bool `json:"trusted"`

		// Can client be used for authentication
		Enabled bool `json:"enabled"`

		// Is client valid yet?
		ValidFrom *time.Time `json:"validFrom,omitempty"`

		// Is client still valid
		ExpiresAt *time.Time `json:"expiresAt,omitempty"`

		// Role-specific settings, see comments on AuthClientSecurity
		Security *AuthClientSecurity `json:"security"`

		// Auth client labels
		Labels map[string]string `json:"labels,omitempty"`

		OwnedBy   uint64     `json:"ownedBy"`
		CreatedBy uint64     `json:"createdBy"`
		CreatedAt time.Time  `json:"createdAt"`
		UpdatedBy uint64     `json:"updatedBy,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	AuthClientMeta struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	AuthClientSecurity struct {
		// Impersonates a specific user;
		// ignored when non client-credentials grant is used
		ImpersonateUser uint64 `json:"impersonateUser,string,omitempty"`

		// Subset of roles, permitted to be used with this client
		// IDs are intentionally stored as strings to support JS (int64 only)
		PermittedRoles []string `json:"permittedRoles,omitempty"`

		// Subset of roles, prohibited to be used with this client
		// IDs are intentionally stored as strings to support JS (int64 only)
		ProhibitedRoles []string `json:"prohibitedRoles,omitempty"`

		// Set of additional roles that are forced on this user
		// IDs are intentionally stored as strings to support JS (int64 only)
		ForcedRoles []string `json:"forcedRoles,omitempty"`
	}

	AuthClientFilter struct {
		ClientID []uint64 `json:"authClientID"`

		Handle string `json:"handle"`

		Deleted filter.State `json:"deleted"`

		LabeledIDs []uint64          `json:"-"`
		Labels     map[string]string `json:"labels,omitempty"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*AuthClient) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	AuthConfirmedClient struct {
		UserID      uint64
		ClientID    uint64
		ConfirmedAt time.Time
	}

	AuthConfirmedClientFilter struct {
		UserID uint64
		Limit  uint
	}
)

func (r *AuthClient) String() string {
	switch {
	case r.Meta != nil && r.Meta.Name != "":
		return r.Meta.Name
	case r.Handle != "":
		return r.Handle
	default:
		return fmt.Sprintf("%d", r.ID)
	}
}

func (r AuthClient) Dict() map[string]interface{} {
	dict := map[string]interface{}{
		"ID":          r.ID,
		"labels":      r.Labels,
		"scope":       r.Scope,
		"validGrant":  r.ValidGrant,
		"redirectURI": r.RedirectURI,
		"trusted":     r.Trusted,
		"enabled":     r.Enabled,
		"validFrom":   r.ValidFrom,
		"expiresAt":   r.ExpiresAt,
		"ownedBy":     r.OwnedBy,
		"createdAt":   r.CreatedAt,
		"createdBy":   r.CreatedBy,
		"updatedAt":   r.UpdatedAt,
		"updatedBy":   r.UpdatedBy,
		"deletedAt":   r.DeletedAt,
		"deletedBy":   r.DeletedBy,
	}

	return dict
}

// FindByHandle finds authClient by it's handle
func (set AuthClientSet) FindByHandle(handle string) *AuthClient {
	for i := range set {
		if set[i].Handle == handle {
			return set[i]
		}
	}

	return nil
}

func (r *AuthClient) Verify() error {
	switch {
	case r == nil || !r.Enabled:
		return fmt.Errorf("disabled")
	case r.ExpiresAt != nil && r.ExpiresAt.After(time.Now()):
		return fmt.Errorf("expired")
	case r.ValidFrom != nil && r.ValidFrom.Before(time.Now()):
		return fmt.Errorf("not yet valid")
	}

	return nil
}

func (vv *AuthClientMeta) Scan(src any) error           { return sql.ParseJSON(src, vv) }
func (vv *AuthClientMeta) Value() (driver.Value, error) { return json.Marshal(vv) }

func (vv *AuthClientSecurity) Scan(src any) error           { return sql.ParseJSON(src, vv) }
func (vv *AuthClientSecurity) Value() (driver.Value, error) { return json.Marshal(vv) }
