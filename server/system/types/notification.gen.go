package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"time"
)

type Notification struct {
	ID        uint64             `json:"notificationID,string"`
	TenantID  uint64             `json:"tenantID,string,omitempty"`
	ProjectID uint64             `json:"projectID,string,omitempty"`
	Kind      NotificationKind   `json:"kind"`
	Config    NotificationConfig `json:"config"`
	Recipient uint64             `json:"recipient,string"`
	CreatedBy uint64             `json:"createdBy,string"`
	ReadAt    *time.Time         `json:"readAt"`
	CreatedAt time.Time          `json:"createdAt,omitempty"`
	UpdatedAt *time.Time         `json:"updatedAt,omitempty"`
	DeletedAt *time.Time         `json:"deletedAt,omitempty"`
}
