package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/crusttech/human/server/pkg/sql"

	"github.com/crusttech/human/server/pkg/filter"
)

type (
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

)

const (
	AttachmentKindSettings       string = "settings"
	AttachmentKindAvatar         string = "avatar"
	AttachmentKindAvatarInitials string = "avatar-initials"
	AttachmentKindChatbot        string = "chatbot"
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
