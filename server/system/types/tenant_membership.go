package types

import (
	"time"

	"github.com/crusttech/human/server/pkg/filter"
)

type (
	// TenantMembership is the join record between a user and a tenant.
	//
	// One membership record per user per tenant. A user may hold at most one
	// active membership across all tenants (one user = one tenant), enforced at
	// the service layer.
	TenantMembership struct {
		ID        uint64             `json:"tenantMembershipID,string"`
		TenantID  uint64             `json:"tenantID,string"`
		UserID    uint64             `json:"userID,string"`
		Role      TenantMemberRole   `json:"role"`
		Status    TenantMemberStatus `json:"status"`
		InvitedBy uint64             `json:"invitedBy,string,omitempty"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	}

	TenantMembershipFilter struct {
		TenantMembershipID []string           `json:"tenantMembershipID"`
		TenantID           uint64             `json:"tenantID,string,omitempty"`
		UserID             uint64             `json:"userID,string,omitempty"`
		Role               TenantMemberRole   `json:"role,omitempty"`
		Status             TenantMemberStatus `json:"status,omitempty"`

		LabeledIDs []uint64 `json:"-"`

		Check func(*TenantMembership) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}

	TenantMemberRole   string
	TenantMemberStatus string
)

const (
	TenantRoleAdmin  TenantMemberRole = "admin"
	TenantRoleMember TenantMemberRole = "member"

	TenantMemberStatusActive    TenantMemberStatus = "active"
	TenantMemberStatusSuspended TenantMemberStatus = "suspended"
	TenantMemberStatusInvited   TenantMemberStatus = "invited"
)

// Valid reports whether r is one of the known roles.
func (r TenantMemberRole) Valid() bool {
	switch r {
	case TenantRoleAdmin, TenantRoleMember:
		return true
	}
	return false
}

// Valid reports whether s is one of the known statuses.
func (s TenantMemberStatus) Valid() bool {
	switch s {
	case TenantMemberStatusActive, TenantMemberStatusSuspended, TenantMemberStatusInvited:
		return true
	}
	return false
}
