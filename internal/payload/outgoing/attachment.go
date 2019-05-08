package outgoing

import (
	"time"
)

type (
	Attachment struct {
		ID         string      `json:"attachmentID"`
		UserID     string      `json:"userID"`
		Url        string      `json:"url"`
		PreviewUrl string      `json:"previewUrl,omitempty"`
		Meta       interface{} `json:"meta"`
		Name       string      `json:"name"`
		CreatedAt  time.Time   `json:"createdAt,omitempty"`
		UpdatedAt  *time.Time  `json:"updatedAt,omitempty"`
	}

	AttachmentSet []*Attachment
)
