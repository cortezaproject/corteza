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
	ProjectMember struct {
		ID         uint64            `json:"projectMemberID,string"`
		TenantID   uint64            `json:"tenantID,string,omitempty"`
		ProjectID  uint64            `json:"projectID,string"`
		UserID     uint64            `json:"userID,string"`
		RolePreset ProjectMemberRole `json:"rolePreset"`
		InvitedBy  uint64            `json:"invitedBy,string,omitempty"`
		CreatedAt  time.Time         `json:"createdAt,omitempty"`
		UpdatedAt  *time.Time        `json:"updatedAt,omitempty"`
		DeletedAt  *time.Time        `json:"deletedAt,omitempty"`
	}
)

func (r ProjectMember) Clone() *ProjectMember {
	dup := r
	if r.UpdatedAt != nil {
		v := *r.UpdatedAt
		dup.UpdatedAt = &v
	}

	if r.DeletedAt != nil {
		v := *r.DeletedAt
		dup.DeletedAt = &v
	}

	return &dup
}

func (r ProjectMember) Diff(cmp *ProjectMember) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ProjectMember{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "projectMemberID", Old: []any{cmp.ID}, New: []any{r.ID}})
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

	if !reflect.DeepEqual(r.RolePreset, cmp.RolePreset) {
		out = append(out, &revisions.Change{Key: "rolePreset", Old: []any{cmp.RolePreset}, New: []any{r.RolePreset}})
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

	if !reflect.DeepEqual(r.DeletedAt, cmp.DeletedAt) {
		out = append(out, &revisions.Change{Key: "deletedAt", Old: []any{cmp.DeletedAt}, New: []any{r.DeletedAt}})
	}

	return out
}
