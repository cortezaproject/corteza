package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type (
	Attachment struct {
		ID         uint64         `db:"id"         json:"ID,omitempty"`
		UserID     uint64         `db:"rel_user"   json:"userID,omitempty"`
		Url        string         `db:"url"        json:"url,omitempty"`
		PreviewUrl string         `db:"preview_url"json:"previewUrl,omitempty"`
		Name       string         `db:"name"       json:"name,omitempty"`
		Meta       attachmentMeta `db:"meta"       json:"meta"`
		CreatedAt  time.Time      `db:"created_at" json:"createdAt,omitempty"`
		UpdatedAt  *time.Time     `db:"updated_at" json:"updatedAt,omitempty"`
		DeletedAt  *time.Time     `db:"deleted_at" json:"deletedAt,omitempty"`
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

	MessageAttachment struct {
		Attachment
		MessageID uint64 `db:"rel_message" json:"-"`
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
