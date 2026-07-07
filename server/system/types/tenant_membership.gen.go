package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/crusttech/human/server/pkg/revisions"
	"reflect"
	"time"
)

type (
	TenantMembership struct {
		ID        uint64             `json:"tenantMembershipID,string"`
		TenantID  uint64             `json:"tenantID,string"`
		ProjectID uint64             `json:"projectID,string,omitempty"`
		UserID    uint64             `json:"userID,string"`
		Role      TenantMemberRole   `json:"role"`
		Status    TenantMemberStatus `json:"status"`
		InvitedBy uint64             `json:"invitedBy,string,omitempty"`
		CreatedAt time.Time          `json:"createdAt,omitempty"`
		UpdatedAt *time.Time         `json:"updatedAt,omitempty"`
	}

	TenantMemberRole string

	TenantMemberStatus string
)

func (r TenantMembership) Clone() *TenantMembership {
	dup := r
	if r.UpdatedAt != nil {
		v := *r.UpdatedAt
		dup.UpdatedAt = &v
	}

	return &dup
}

func (r TenantMembership) Diff(cmp *TenantMembership) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &TenantMembership{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "tenantMembershipID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.UserID != cmp.UserID {
		out = append(out, &revisions.Change{Key: "userID", Old: []any{cmp.UserID}, New: []any{r.UserID}})
	}

	if !reflect.DeepEqual(r.Role, cmp.Role) {
		out = append(out, &revisions.Change{Key: "role", Old: []any{cmp.Role}, New: []any{r.Role}})
	}

	if !reflect.DeepEqual(r.Status, cmp.Status) {
		out = append(out, &revisions.Change{Key: "status", Old: []any{cmp.Status}, New: []any{r.Status}})
	}

	if r.InvitedBy != cmp.InvitedBy {
		out = append(out, &revisions.Change{Key: "invitedBy", Old: []any{cmp.InvitedBy}, New: []any{r.InvitedBy}})
	}

	if !reflect.DeepEqual(r.CreatedAt, cmp.CreatedAt) {
		out = append(out, &revisions.Change{Key: "createdAt", Old: []any{cmp.CreatedAt}, New: []any{r.CreatedAt}})
	}

	if !reflect.DeepEqual(r.UpdatedAt, cmp.UpdatedAt) {
		out = append(out, &revisions.Change{Key: "updatedAt", Old: []any{cmp.UpdatedAt}, New: []any{r.UpdatedAt}})
	}

	return out
}

const (
	TenantMemberRoleTenantRoleAdmin  TenantMemberRole = "admin"
	TenantMemberRoleTenantRoleMember TenantMemberRole = "member"
)

const (
	TenantMemberStatusActive    TenantMemberStatus = "active"
	TenantMemberStatusSuspended TenantMemberStatus = "suspended"
	TenantMemberStatusInvited   TenantMemberStatus = "invited"
)
