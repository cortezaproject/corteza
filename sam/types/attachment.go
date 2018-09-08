package types

import (
	"time"
)

type (
	Attachment struct {
		ID         uint64     `db:"id"         json:"id,omitempty"`
		UserID     uint64     `db:"rel_user"   json:"userID,omitempty"`
		Url        string     `db:"url"        json:"url,omitempty"`
		PreviewUrl string     `db:"preview_url"json:"previewUrl,omitempty"`
		Size       int64      `db:"size"       json:"size,omitempty"`
		Mimetype   string     `db:"mimetype"   json:"mimetype,omitempty"`
		Name       string     `db:"name"       json:"name,omitempty"`
		CreatedAt  time.Time  `db:"created_at" json:"createdAt,omitempty"`
		UpdatedAt  *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		DeletedAt  *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
	}
)
