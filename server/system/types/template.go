package types

import (
	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

type (
	DocumentType string

	TemplateMeta struct {
		Short       string `json:"short"`
		Description string `json:"description,omitempty"`
	}

	TemplateFilter struct {
		TemplateID []string `json:"templateID"`
		Query      string   `json:"query"`
		Handle     string   `json:"handle"`
		Type       string   `json:"type"`
		OwnerID    uint64   `json:"ownerID,string"`
		Partial    bool     `json:"partial"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Template) (bool, error) `json:"-"`

		Deleted filter.State `json:"deleted"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}
)

const (
	DocumentTypePlain DocumentType = "text/plain"
	DocumentTypeHTML  DocumentType = "text/html"
	DocumentTypePDF   DocumentType = "application/pdf"
)

func (t Template) Clone() *Template {
	c := &t
	return c
}
