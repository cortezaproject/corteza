package types

import (
	"encoding/json"
	"time"
)

type (
	Attachment struct {
		ID         uint64          `db:"id"`
		UserID     uint64          `db:"rel_user"`
		MessageID  uint64          `db:"rel_message"`
		ChannelID  uint64          `db:"rel_channel"`
		Attachment json.RawMessage `db:"attachment"`
		Url        string          `db:"url"`
		PreviewUrl string          `db:"preview_url"`
		Size       uint            `db:"size"`
		Mimetype   string          `db:"mimetype"`
		Name       string          `db:"name"`
		CreatedAt  time.Time       `json:"created_at,omitempty" db:"created_at"`
		UpdatedAt  *time.Time      `json:"updated_at,omitempty" db:"updated_at"`
		DeletedAt  *time.Time      `json:"deleted_at,omitempty" db:"deleted_at"`
	}
)
