package outgoing

import (
	"time"
)

type (
	Attachment struct {
		ID         string     `json:"ID"`
		UserID     string     `json:"userID"`
		Url        string     `json:"url"`
		PreviewUrl string     `json:"previewUrl"`
		Size       int64      `json:"size"`
		Mimetype   string     `json:"mimetype"`
		Name       string     `json:"name"`
		CreatedAt  time.Time  `json:"createdAt,omitempty"`
		UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
	}

	AttachmentSet []*Attachment
)
