package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/slice"
	"strconv"
	"time"
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
		ValidGrant string `json:"grant"`

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
	}
)

// Resource returns a resource ID for this type
func (r *AuthClient) RBACResource() rbac.Resource {
	return AuthClientRBACResource.AppendID(r.ID)
}

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

func (vv *AuthClientMeta) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*vv = AuthClientMeta{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, vv); err != nil {
			return fmt.Errorf("cannot scan '%v' into AuthClientMeta: %w", string(b), err)
		}
	}

	return nil
}

// Scan on AuthClientMeta gracefully handles conversion from NULL
func (vv *AuthClientMeta) Value() (driver.Value, error) {
	if vv == nil {
		return []byte("null"), nil
	}

	return json.Marshal(vv)
}

func (vv *AuthClientSecurity) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*vv = AuthClientSecurity{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, vv); err != nil {
			return fmt.Errorf("cannot scan '%v' into AuthClientSecurity: %w", string(b), err)
		}
	}

	return nil
}

// Scan on AuthClientSecurity gracefully handles conversion from NULL
func (vv *AuthClientSecurity) Value() (driver.Value, error) {
	if vv == nil {
		return []byte("null"), nil
	}

	return json.Marshal(vv)
}

// Takes user's roles, filter out only allowed roles (when set), remove denied and add all forced
func (vv *AuthClientSecurity) ProcessRoles(rr ...uint64) (out []uint64) {
	var (
		permitted  = slice.ToStringBoolMap(vv.PermittedRoles)
		prohibited = slice.ToStringBoolMap(vv.ProhibitedRoles)
		forced     = slice.ToStringBoolMap(vv.ForcedRoles)
		aux        string
		roleID     uint64
	)

	// iterate over user's roles and just append them (obeying allow&deny rules)
	// to list of forced roles
	for _, r := range rr {
		aux = strconv.FormatUint(r, 10)
		if (len(vv.PermittedRoles) == 0 || permitted[aux]) && !prohibited[aux] {
			forced[aux] = true
		}
	}

	out = make([]uint64, 0, len(forced))
	for i := range forced {
		if roleID, _ = strconv.ParseUint(i, 10, 64); roleID > 0 {
			out = append(out, roleID)
		}
	}

	return
}
