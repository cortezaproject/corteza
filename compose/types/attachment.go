package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type (
	Attachment struct {
		ID         uint64         `db:"id"          json:"attachmentID,string"`
		OwnerID    uint64         `db:"rel_owner"   json:"ownerID,string"`
		Kind       string         `db:"kind"        json:"-"`
		Url        string         `db:"url"         json:"url,omitempty"`
		PreviewUrl string         `db:"preview_url" json:"previewUrl,omitempty"`
		Name       string         `db:"name"        json:"name,omitempty"`
		Meta       attachmentMeta `db:"meta"        json:"meta"`

		NamespaceID uint64 `db:"rel_namespace" json:"namespaceID,string"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
	}

	// AttachmentFilter is used for filtering and as a return value from Find
	AttachmentFilter struct {
		NamespaceID uint64 `json:"namespaceID,string"`
		Kind        string `json:"kind,omitempty"`
		PageID      uint64 `json:"pageID,string,omitempty"`
		RecordID    uint64 `json:"recordID,string,omitempty"`
		ModuleID    uint64 `json:"moduleID,string,omitempty"`
		FieldName   string `json:"fieldName,omitempty"`
		Filter      string `json:"filter"`
		Page        uint   `json:"page"`
		PerPage     uint   `json:"perPage"`
		Sort        string `json:"sort"`
		Count       uint   `json:"count"`
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

const (
	PageAttachment   string = "page"
	RecordAttachment string = "record"
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
