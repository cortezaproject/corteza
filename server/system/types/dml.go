package types

import (
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/id"
)

type (
	DmlConnectionInput struct {
		Handle string               `json:"handle"`
		Label  string               `json:"label"`
		Params *DmlConnectionParams `json:"params"`
	}

	DmlConnectionFilter struct {
		ConnectionID id.Uint64s   `json:"connectionID"`
		Handle       string       `json:"handle"`
		Type         string       `json:"type"`
		Deleted      filter.State `json:"deleted"`

		Check func(*DmlConnection) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}

	DmlMappingFilter struct {
		MappingID    id.Uint64s   `json:"mappingID"`
		ConnectionID uint64       `json:"connectionID,string"`
		SourceIdent  string       `json:"sourceIdent"`
		Deleted      filter.State `json:"deleted"`

		Check func(*DmlMapping) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}

	DmlImportRunFilter struct {
		ImportRunID id.Uint64s `json:"importRunID"`
		MappingID   uint64     `json:"mappingID,string"`
		Status      []string   `json:"status"`

		Check func(*DmlImportRun) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}

	DmlModelFilter struct {
		Ident []string `json:"ident"`
	}

	DmlModel struct {
		ConnectionID uint64          `json:"connectionID,string"`
		Ident        string          `json:"ident"`
		Label        string          `json:"label"`
		ResourceType string          `json:"resourceType"`
		Attributes   []*DmlAttribute `json:"attributes"`
	}

	DmlAttribute struct {
		Ident      string `json:"ident"`
		Label      string `json:"label"`
		PrimaryKey bool   `json:"primaryKey"`
		Sortable   bool   `json:"sortable"`
		Filterable bool   `json:"filterable"`
		Type       string `json:"type"`
		Store      string `json:"store"`
	}
)
