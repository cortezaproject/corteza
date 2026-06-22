package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"encoding/json"
	"time"
)

type Attachment struct {
	ID         uint64         `json:"attachmentID,string"`
	TenantID   uint64         `json:"tenantID,string,omitempty"`
	ProjectID  uint64         `json:"projectID,string,omitempty"`
	OwnerID    uint64         `json:"ownerID,string"`
	Kind       string         `json:"-"`
	Url        string         `json:"url,omitempty"`
	PreviewUrl string         `json:"previewUrl,omitempty"`
	Name       string         `json:"name,omitempty"`
	Meta       AttachmentMeta `json:"meta"`
	CreatedAt  time.Time      `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time     `json:"updatedAt,omitempty"`
	DeletedAt  *time.Time     `json:"deletedAt,omitempty"`
}

func (r Attachment) Clone() *Attachment {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}
