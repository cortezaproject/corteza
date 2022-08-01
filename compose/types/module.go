package types

import (
	"database/sql/driver"
	"encoding/json"
	discovery "github.com/cortezaproject/corteza-server/discovery/types"
	"github.com/cortezaproject/corteza-server/pkg/sql"
	"github.com/jmoiron/sqlx/types"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/locale"
)

type (
	Module struct {
		ID     uint64 `json:"moduleID,string"`
		Handle string `json:"handle"`

		// collection of configurations for various subsystems that
		// use this module and how it affects their behaviour
		Config ModuleConfig `json:"config"`

		// @todo should be removed and placed into a separate subsystem
		//       mostly because we want to allow client apps to store
		//       application configs away from the module config
		//       using separate access-control
		Meta types.JSONText `json:"meta"`

		Fields ModuleFieldSet `json:"fields"`

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

	SystemFieldEncoding struct {
		ID *EncodingStrategy `json:"id"`

		ModuleID    *EncodingStrategy `json:"moduleID"`
		NamespaceID *EncodingStrategy `json:"namespaceID"`

		Revision *EncodingStrategy `json:"revision"`

		OwnedBy *EncodingStrategy `json:"ownedBy"`

		CreatedAt *EncodingStrategy `json:"createdAt"`
		CreatedBy *EncodingStrategy `json:"createdBy"`

		UpdatedAt *EncodingStrategy `json:"updatedAt"`
		UpdatedBy *EncodingStrategy `json:"updatedBy"`

		DeletedAt *EncodingStrategy `json:"deletedAt"`
		DeletedBy *EncodingStrategy `json:"deletedBy"`
	}

	ModuleConfig struct {
		// How and where the records of this module are stored in the database
		DAL ModuleConfigDAL `json:"dal"`

		// Record data privacy settings
		Privacy ModuleConfigDataPrivacy `json:"privacy"`

		// @todo we need to transfer this from meta!!
		Discovery discovery.ModuleMeta `json:"discovery"`

		RecordRevisions ModuleConfigRecordRevisions `json:"recordRevisions"`
	}

	ModuleConfigDAL struct {
		ConnectionID uint64           `json:"connectionID,string"`
		Capabilities capabilities.Set `json:"capabilities"`

		Issues []string `json:"issues,omitempty"`

		Constraints map[string][]any `json:"constraints"`

		Partitioned     bool   `json:"partitioned"`
		PartitionFormat string `json:"partitionFormat"`

		SystemFieldEncoding SystemFieldEncoding `json:"systemFieldEncoding"`
	}

	ModuleConfigRecordRevisions struct {
		// enable or disable revisions
		Enabled bool `json:"enabled"`

		// where are record revisions stored
		Ident string `json:"ident"`
	}

	ModuleConfigDataPrivacy struct {
		// Define the highest sensitivity level which
		// can be configured on the module fields
		SensitivityLevel uint64 `json:"sensitivityLevel,string,omitempty"`

		UsageDisclosure string `json:"usageDisclosure"`
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

func (m Module) HasIssues() bool {
	return len(m.Config.DAL.Issues) > 0
}

// We won't worry about fields at this point
func (m *Module) decodeTranslations(tt locale.ResourceTranslationIndex) {
	return
}

// We won't worry about fields at this point
func (m *Module) encodeTranslations() (out locale.ResourceTranslationSet) {
	return
}

func (m *Module) ModelRef() dal.ModelRef {
	return dal.ModelRef{
		ConnectionID: m.Config.DAL.ConnectionID,

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

func (c *ModuleConfig) Scan(src any) error          { return sql.ParseJSON(src, c) }
func (c ModuleConfig) Value() (driver.Value, error) { return json.Marshal(c) }

func ParseModuleConfig(ss []string) (m ModuleConfig, err error) {
	if len(ss) == 0 {
		return
	}

	err = json.Unmarshal([]byte(ss[0]), &m)
	return
}
