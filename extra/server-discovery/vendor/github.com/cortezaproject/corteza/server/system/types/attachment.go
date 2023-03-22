package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/cortezaproject/corteza/server/pkg/sql"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/filter"
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

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	// AttachmentFilter is used for filtering and as a return value from Find
	AttachmentFilter struct {
		Kind   string `json:"kind,omitempty"`
		Filter string `json:"filter"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Attachment) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	AttachmentImageMeta struct {
		Width           int    `json:"width,omitempty"`
		Height          int    `json:"height,omitempty"`
		Animated        bool   `json:"animated"`
		Initial         string `json:"initial,omitempty"`
		InitialColor    string `json:"initial-color,omitempty"`
		BackgroundColor string `json:"background-color,omitempty"`
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
		Labels   map[string]string   `json:"labels,omitempty"`
	}
)

const (
	AttachmentKindSettings       string = "settings"
	AttachmentKindAvatar         string = "avatar"
	AttachmentKindAvatarInitials string = "avatar-initials"
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

func (meta *AttachmentMeta) Scan(src any) error { return sql.ParseJSON(src, meta) }
func (meta AttachmentMeta) Value() (driver.Value, error) {

	return json.Marshal(meta)
}
