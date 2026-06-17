package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"time"
)

type ProjectMember struct {
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
