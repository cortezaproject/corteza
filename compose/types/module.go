package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	discovery "github.com/cortezaproject/corteza-server/discovery/types"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"
)

type (
	DALConfig struct {
		ConnectionID uint64           `json:"connectionID,string"`
		Capabilities capabilities.Set `json:"capabilities"`

		Constraints map[string][]any `json:"constraints"`

		Partitioned     bool   `json:"partitioned"`
		PartitionFormat string `json:"partitionFormat"`
	}

	Module struct {
		ID     uint64         `json:"moduleID,string"`
		Handle string         `json:"handle"`
		Meta   types.JSONText `json:"meta"`

		DALConfig DALConfig `json:"DALConfig"`

		Fields       ModuleFieldSet `json:"fields"`
		SystemFields SystemFieldSet `json:"systemFields"`

		Labels map[string]string `json:"labels,omitempty"`

		NamespaceID uint64 `json:"namespaceID,string"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`

		// Warning: value of this field is now handled via resource-translation facility
		//          struct field is kept for the convenience for now since it allows us
		//          easy encoding/decoding of the outgoing/incoming values
		Name string `json:"name"`
	}

	SystemFieldSet struct {
		ID *EncodingStrategy

		ModuleID    *EncodingStrategy
		NamespaceID *EncodingStrategy

		OwnedBy *EncodingStrategy

		CreatedAt *EncodingStrategy
		CreatedBy *EncodingStrategy

		UpdatedAt *EncodingStrategy
		UpdatedBy *EncodingStrategy

		DeletedAt *EncodingStrategy
		DeletedBy *EncodingStrategy
	}

	ModuleMeta struct {
		Discovery discovery.ModuleMeta `json:"discovery"`
	}

	ModuleFilter struct {
		ModuleID    []uint64 `json:"moduleID"`
		NamespaceID uint64   `json:"namespaceID,string"`
		Query       string   `json:"query"`
		Handle      string   `json:"handle"`
		Name        string   `json:"name"`

		LabeledIDs []uint64          `json:"-"`
		Labels     map[string]string `json:"labels,omitempty"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Module) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}
)

func (m Module) Clone() *Module {
	c := &m
	c.Fields = m.Fields.Clone()
	return c
}

// We won't worry about fields at this point
func (m *Module) decodeTranslations(tt locale.ResourceTranslationIndex) {
	return
}

// We won't worry about fields at this point
func (m *Module) encodeTranslations() (out locale.ResourceTranslationSet) {
	return
}

func (m *Module) ModelFilter() dal.ModelFilter {
	return dal.ModelFilter{
		ConnectionID: m.DALConfig.ConnectionID,

		ResourceID: m.ID,

		ResourceType: ModuleResourceType,
		// @todo will use this for now but should probably change
		Resource: m.RbacResource(),
	}
}

// FindByHandle finds module by it's handle
func (set ModuleSet) FindByHandle(handle string) *Module {
	for i := range set {
		if set[i].Handle == handle {
			return set[i]
		}
	}

	return nil
}

func (nm *ModuleMeta) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*nm = ModuleMeta{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, nm); err != nil {
			return errors.Wrapf(err, "cannot scan '%v' into ModuleMeta", string(b))
		}
	}

	return nil
}

func (nm ModuleMeta) Value() (driver.Value, error) {
	return json.Marshal(nm)
}
