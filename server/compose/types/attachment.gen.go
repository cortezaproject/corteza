package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/crusttech/human/server/pkg/sql"
	"time"
)

type Attachment struct {
	ID          uint64         `json:"attachmentID,string"`
	TenantID    uint64         `json:"tenantID,string,omitempty"`
	ProjectID   uint64         `json:"projectID,string,omitempty"`
	NamespaceID uint64         `json:"namespaceID,string"`
	OwnerID     uint64         `json:"ownerID,string"`
	Kind        string         `json:"-"`
	Url         string         `json:"url,omitempty"`
	PreviewUrl  string         `json:"previewUrl,omitempty"`
	Name        string         `json:"name,omitempty"`
	Meta        AttachmentMeta `json:"meta"`
	CreatedAt   time.Time      `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time     `json:"updatedAt,omitempty"`
	DeletedAt   *time.Time     `json:"deletedAt,omitempty"`
}

func (m *AttachmentMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m AttachmentMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseAttachmentMeta(ss []string) (p AttachmentMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
