package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	// AttachmentFilter is used for filtering and as a return value from Find
	AttachmentFilter struct {
		TenantID    uint64 `json:"tenantID,string,omitempty"`
		ProjectID   uint64 `json:"projectID,string,omitempty"`
		NamespaceID uint64 `json:"namespaceID,string"`
		Kind        string `json:"kind,omitempty"`
		PageID      uint64 `json:"pageID,string,omitempty"`
		RecordID    uint64 `json:"recordID,string,omitempty"`
		ModuleID    uint64 `json:"moduleID,string,omitempty"`
		FieldName   string `json:"fieldName,omitempty"`
		Filter      string `json:"filter"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Attachment) (bool, error)

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}
)

const (
	PageAttachment      string = "page"
	IconAttachment      string = "icon"
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
