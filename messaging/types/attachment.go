package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type (
	Attachment struct {
		ID         uint64         `json:"ID,omitempty"`
		OwnerID    uint64         `json:"userID,omitempty"`
		Url        string         `json:"url,omitempty"`
		PreviewUrl string         `json:"previewUrl,omitempty"`
		Name       string         `json:"name,omitempty"`
		Meta       attachmentMeta `json:"meta"`
		CreatedAt  time.Time      `json:"createdAt,omitempty"`
		UpdatedAt  *time.Time     `json:"updatedAt,omitempty"`
		DeletedAt  *time.Time     `json:"deletedAt,omitempty"`
	}

	MessageAttachment struct {
		AttachmentID uint64
		MessageID    uint64
	}

	// AttachmentFilter is used for filtering and as a return value from Find
	AttachmentFilter struct {
		AttachmentID []uint64
	}

	// MessageAttachmentFilter is used for filtering and as a return value from Find
	MessageAttachmentFilter struct {
		MessageID []uint64
	}

	attachmentImageMeta struct {
		Width    int  `json:"width,omitempty"`
		Height   int  `json:"height,omitempty"`
		Animated bool `json:"animated"`
	}

	attachmentFileMeta struct {
		Size      int64                `json:"size"`
		Extension string               `json:"ext"`
		Mimetype  string               `json:"mimetype"`
		Image     *attachmentImageMeta `json:"image,omitempty"`
	}

	attachmentMeta struct {
		Original attachmentFileMeta  `json:"original"`
		Preview  *attachmentFileMeta `json:"preview,omitempty"`
	}
)

func (a *Attachment) SetOriginalImageMeta(width, height int, animated bool) *attachmentFileMeta {
	a.imageMeta(&a.Meta.Original, width, height, animated)
	return &a.Meta.Original
}

func (a *Attachment) SetPreviewImageMeta(width, height int, animated bool) *attachmentFileMeta {
	if a.Meta.Preview == nil {
		a.Meta.Preview = &attachmentFileMeta{}
	}

	a.imageMeta(a.Meta.Preview, width, height, animated)
	return a.Meta.Preview
}

func (a *Attachment) imageMeta(in *attachmentFileMeta, width, height int, animated bool) {
	if in.Image == nil {
		in.Image = &attachmentImageMeta{}
	}

	if width > 0 && height > 0 {
		in.Image.Animated = animated
		in.Image.Width = width
		in.Image.Height = height
	}
}

func (meta *attachmentMeta) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*meta = attachmentMeta{}
	case []uint8:
		if err := json.Unmarshal(value.([]byte), meta); err != nil {
			return errors.Wrapf(err, "Can not scan '%v' into attachmentMeta", value)
		}
	}

	return nil
}

func (meta attachmentMeta) Value() (driver.Value, error) {
	return json.Marshal(meta)
}
