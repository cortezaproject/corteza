package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"

	"github.com/pkg/errors"
)

type (
	Attachment struct {
		ID         uint64         `json:"attachmentID,string"`
		OwnerID    uint64         `json:"ownerID,string"`
		Kind       string         `json:"-"`
		Url        string         `json:"url,omitempty"`
		PreviewUrl string         `json:"previewUrl,omitempty"`
		Name       string         `json:"name,omitempty"`
		Meta       AttachmentMeta `json:"meta"`

		NamespaceID uint64 `json:"namespaceID,string"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
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

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Attachment) (bool, error)

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	AttachmentImageMeta struct {
		Width    int  `json:"width,omitempty"`
		Height   int  `json:"height,omitempty"`
		Animated bool `json:"animated"`
	}

	AttachmentFileMeta struct {
		Size      int64                `json:"size"`
		Extension string               `json:"ext"`
		Mimetype  string               `json:"mimetype"`
		Image     *AttachmentImageMeta `json:"image,omitempty"`
	}

	AttachmentMeta struct {
		Original AttachmentFileMeta  `json:"original"`
		Preview  *AttachmentFileMeta `json:"preview,omitempty"`
	}
)

const (
	PageAttachment      string = "page"
	RecordAttachment    string = "record"
	NamespaceAttachment string = "namespace"
)

func (a *Attachment) SetOriginalImageMeta(width, height int, animated bool) *AttachmentFileMeta {
	a.imageMeta(&a.Meta.Original, width, height, animated)
	return &a.Meta.Original
}

func (a *Attachment) SetPreviewImageMeta(width, height int, animated bool) *AttachmentFileMeta {
	if a.Meta.Preview == nil {
		a.Meta.Preview = &AttachmentFileMeta{}
	}

	a.imageMeta(a.Meta.Preview, width, height, animated)
	return a.Meta.Preview
}

func (a *Attachment) imageMeta(in *AttachmentFileMeta, width, height int, animated bool) {
	if in.Image == nil {
		in.Image = &AttachmentImageMeta{}
	}

	if width > 0 && height > 0 {
		in.Image.Animated = animated
		in.Image.Width = width
		in.Image.Height = height
	}
}

func (meta *AttachmentMeta) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*meta = AttachmentMeta{}
	case []uint8:
		if err := json.Unmarshal(value.([]byte), meta); err != nil {
			return errors.Wrapf(err, "cannot scan '%v' into attachmentMeta", value)
		}
	}

	return nil
}

func (meta AttachmentMeta) Value() (driver.Value, error) {
	return json.Marshal(meta)
}
