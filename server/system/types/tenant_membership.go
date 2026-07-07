package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	TenantMembershipFilter struct {
		TenantMembershipID []string           `json:"tenantMembershipID"`
		TenantID           uint64             `json:"tenantID,string,omitempty"`
		ProjectID          uint64             `json:"projectID,string,omitempty"`
		UserID             uint64             `json:"userID,string,omitempty"`
		Role               TenantMemberRole   `json:"role,omitempty"`
		Status             TenantMemberStatus `json:"status,omitempty"`

		LabeledIDs []uint64 `json:"-"`

		Check func(*TenantMembership) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)

const (
	TenantRoleAdmin  TenantMemberRole = "admin"
	TenantRoleMember TenantMemberRole = "member"
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
