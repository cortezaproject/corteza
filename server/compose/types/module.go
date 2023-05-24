package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/sql"
	"github.com/jmoiron/sqlx/types"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/locale"
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

		Issues []string `json:"issues,omitempty"`

		NamespaceID uint64 `json:"namespaceID,string"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`

		// Warning: value of this field is now handled via resource-translation facility
		//          struct field is kept for the convenience for now since it allows us
		//          easy encoding/decoding of the outgoing/incoming values
		Name string `json:"name"`
	}

	ModuleConfig struct {
		// How and where the records of this module are stored in the database
		DAL ModuleConfigDAL `json:"dal"`

		// Record data privacy settings
		Privacy ModuleConfigDataPrivacy `json:"privacy"`

		Discovery ModuleConfigDiscovery `json:"discovery"`

		RecordRevisions ModuleConfigRecordRevisions `json:"recordRevisions"`

		// RecordDeDup value duplicate detection settings
		RecordDeDup ModuleConfigRecordDeDup `json:"recordDeDup"`
	}

	ModuleConfigDAL struct {
		ConnectionID uint64 `json:"connectionID,string"`

		Constraints map[string][]any `json:"constraints"`

		// model identifier (table, collection on the database)
		// can contain {{placeholders}}
		Ident string `json:"ident"`

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
		SensitivityLevelID uint64 `json:"sensitivityLevelID,string,omitempty"`

		UsageDisclosure string `json:"usageDisclosure"`
	}

	ModuleConfigRecordDeDup struct {
		// strictly restrict record saving
		// 		otherwise show a warning with list of duplicated records
		Strict bool `json:"-"`

		// list of duplicate detection rules applied to module's fields
		Rules DeDupRuleSet `json:"rules,omitempty"`
	}

	ModuleConfigDiscovery struct {
		Public    DiscoveryResult `json:"public"`
		Private   DiscoveryResult `json:"private"`
		Protected DiscoveryResult `json:"protected"`
	}

	DiscoveryResult struct {
		Result []struct {
			Lang   string   `json:"lang"`
			Fields []string `json:"fields"`
		} `json:"result"`
	}

	ModuleFilter struct {
		ModuleID    []string `json:"moduleID"`
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
	return len(m.Issues) > 0
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
