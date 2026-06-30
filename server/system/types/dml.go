package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type DmlImportMethod string

const (
	DmlImportMethodBackground DmlImportMethod = "background"
)

type (
	DmlConnectionParams struct {
		Type            string         `json:"type"`
		Params          map[string]any `json:"params"`
		ModelIdent      string         `json:"modelIdent"`
		ModelIdentCheck []string       `json:"modelIdentCheck"`
	}

	DmlConnectionInput struct {
		Handle string               `json:"handle"`
		Label  string               `json:"label"`
		Params *DmlConnectionParams `json:"params"`
	}

	DmlConnectionFilter struct {
		ConnectionID []uint64     `json:"connectionID"`
		Handle       string       `json:"handle"`
		Type         string       `json:"type"`
		Deleted      filter.State `json:"deleted"`

		Check func(*DmlConnection) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}

	DmlMappingFilter struct {
		MappingID    []uint64     `json:"mappingID"`
		ConnectionID uint64       `json:"connectionID"`
		SourceIdent  string       `json:"sourceIdent"`
		Deleted      filter.State `json:"deleted"`

		Check func(*DmlMapping) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}

	DmlImportRunFilter struct {
		ImportRunID []uint64 `json:"importRunID"`
		MappingID   uint64   `json:"mappingID"`
		Status      []string `json:"status"`

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

	DmlColumnMap struct {
		SourceIdent string `json:"sourceIdent"`
		FieldName   string `json:"fieldName"`
		Label       string `json:"label,omitempty"`
		FieldKind   string `json:"fieldKind"`
		Skip        bool   `json:"skip"`
	}
)
