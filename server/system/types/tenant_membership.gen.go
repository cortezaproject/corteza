package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"time"
)

type TenantMembership struct {
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
