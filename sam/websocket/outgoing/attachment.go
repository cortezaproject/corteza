package outgoing

import (
	"time"
)

type (
	Attachment struct {
		ID         string     `json:"id"`
		UserID     string     `json:"uid"`
		Url        string     `json:"url"`
		PreviewUrl string     `json:"prw"`
		Size       int64      `json:"sze"`
		Mimetype   string     `json:"typ"`
		Name       string     `json:"nme"`
		CreatedAt  time.Time  `json:"cat,omitempty"`
		UpdatedAt  *time.Time `json:"uat,omitempty"`
	}
)
