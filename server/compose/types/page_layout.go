package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/sql"

	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	PageLayout struct {
		ID       uint64 `json:"pageLayoutID,string"`
		PageID   uint64 `json:"pageID,string"`
		ParentID uint64 `json:"parentID,string"`
		Handle   string `json:"handle"`
		Primary  bool   `json:"primary"`

		NamespaceID uint64 `json:"namespaceID,string"`
		ModuleID    uint64 `json:"moduleID,string"`

		Meta *PageLayoutMeta `json:"meta"`

		Config PageLayoutConfig `json:"config"`
		Blocks PageBlocks       `json:"blocks"`

		Labels map[string]string `json:"labels,omitempty"`

		OwnedBy uint64 `json:"ownedBy,string"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	PageLayoutMeta struct {
		// Warning: value of this field is now handled via resource-translation facility
		//          struct field is kept for the convenience for now since it allows us
		//          easy encoding/decoding of the outgoing/incoming values
		Name string `json:"name"`

		// Warning: value of this field is now handled via resource-translation facility
		//          struct field is kept for the convenience for now since it allows us
		//          easy encoding/decoding of the outgoing/incoming values
		Description string `json:"description"`
	}

	PageLayoutButton struct {
		Label   string `json:"label"`
		Enabled bool   `json:"enabled"`
	}

	// @todo this will probably expand with a more modular action-like system
	PageLayoutButtonConfig struct {
		New    PageLayoutButton `json:"new"`
		Edit   PageLayoutButton `json:"edit"`
		Submit PageLayoutButton `json:"submit"`
		Delete PageLayoutButton `json:"delete"`
		Clone  PageLayoutButton `json:"clone"`
		Back   PageLayoutButton `json:"back"`
	}

	PageLayoutConfig struct {
		Buttons *PageLayoutButtonConfig `json:"buttons,omitempty"`
	}

	PageLayoutFilter struct {
		NamespaceID uint64 `json:"namespaceID,string"`
		PageID      uint64 `json:"pageID,string,omitempty"`
		ModuleID    uint64 `json:"moduleID,string,omitempty"`
		Default     bool   `json:"default,omitempty"`
		Handle      string `json:"handle"`
		Name        string `json:"name"`
		Query       string `json:"query"`

		LabeledIDs []uint64          `json:"-"`
		Labels     map[string]string `json:"labels,omitempty"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*PageLayout) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}
)

func (m PageLayout) Clone() *PageLayout {
	c := &m
	return c
}

// FindByHandle finds pageLayout by it's handle
func (set PageLayoutSet) FindByHandle(handle string) *PageLayout {
	for i := range set {
		if set[i].Handle == handle {
			return set[i]
		}
	}

	return nil
}

func (bb *PageLayoutConfig) Scan(src any) error { return sql.ParseJSON(src, bb) }
func (bb PageLayoutConfig) Value() (driver.Value, error) {
	// We're not saving button config to the DB; no need for it
	bb.Buttons = nil
	return json.Marshal(bb)
}

func (vv *PageLayoutMeta) Scan(src any) error           { return sql.ParseJSON(src, vv) }
func (vv *PageLayoutMeta) Value() (driver.Value, error) { return json.Marshal(vv) }
