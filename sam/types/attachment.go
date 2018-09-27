package types

import (
	"fmt"
	"net/url"
	"time"
)

const (
	attachmentURL        = "/attachment/%d/%s"
	attachmentPreviewURL = "/attachment/%d/%s/preview"
)

type (
	Attachment struct {
		ID         uint64     `db:"id"         json:"ID,omitempty"`
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

	MessageAttachment struct {
		Attachment
		MessageID uint64 `db:"rel_message" json:"-"`
	}

	MessageAttachmentSet []*MessageAttachment
)

func (aa MessageAttachmentSet) Walk(w func(*MessageAttachment) error) (err error) {
	for i := range aa {
		if err = w(aa[i]); err != nil {
			return
		}
	}

	return
}

func (a *Attachment) GenerateURLs() {
	a.Url = fmt.Sprintf(attachmentURL, a.ID, url.PathEscape(a.Name))
	a.PreviewUrl = fmt.Sprintf(attachmentPreviewURL, a.ID, url.PathEscape(a.Name))
}
