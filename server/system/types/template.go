package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/cortezaproject/corteza/server/pkg/sql"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	DocumentType string

	Template struct {
		ID     uint64 `json:"templateID,string"`
		Handle string `json:"handle"`
		// Language specifies the language the template is written for; leave empty for default
		Language string `json:"language"`

		Type DocumentType `json:"type"`
		// Partial templates can be used to construct larger templates; for example headers and footers
		Partial bool `json:"partial"`
		// use int so JS can handle it normally
		//
		// @todo We'll handle this at a later point
		// Revision int           `json:"revision,string"`
		Meta TemplateMeta `json:"meta"`

		Template string `json:"template"`

		Labels map[string]string `json:"labels,omitempty"`

		OwnerID    uint64     `json:"ownerID,string"`
		CreatedAt  time.Time  `json:"createdAt,omitempty"`
		UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
		DeletedAt  *time.Time `json:"deletedAt,omitempty"`
		LastUsedAt *time.Time `json:"lastUsedAt,omitempty"`
	}

	TemplateMeta struct {
		Short       string `json:"short"`
		Description string `json:"description,omitempty"`
	}

	TemplateFilter struct {
		TemplateID []uint64 `json:"templateID"`
		Query      string   `json:"query"`
		Handle     string   `json:"handle"`
		Type       string   `json:"type"`
		OwnerID    uint64   `json:"ownerID,string"`
		Partial    bool     `json:"partial"`

		LabeledIDs []uint64          `json:"-"`
		Labels     map[string]string `json:"labels,omitempty"`

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

func (t *TemplateMeta) Scan(src any) error          { return sql.ParseJSON(src, t) }
func (t TemplateMeta) Value() (driver.Value, error) { return json.Marshal(t) }

func (t Template) Clone() *Template {
	c := &t
	return c
}
