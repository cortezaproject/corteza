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
	ModelConfig struct {
		ConnectionID uint64           `json:"connectionID,string"`
		Capabilities capabilities.Set `json:"capabilities"`

		Issues []string `json:"issues,omitempty"`

		Constraints map[string][]any `json:"constraints"`

		Partitioned     bool   `json:"partitioned"`
		PartitionFormat string `json:"partitionFormat"`

		SystemFieldEncoding SystemFieldEncoding `json:"systemFieldEncoding"`
	}

	DataPrivacyConfig struct {
		SensitivityLevel uint64 `json:"sensitivityLevel,string,omitempty"`
		UsageDisclosure  string `json:"usageDisclosure"`
	}

	Module struct {
		ID     uint64         `json:"moduleID,string"`
		Handle string         `json:"handle"`
		Meta   types.JSONText `json:"meta"`

		ModelConfig ModelConfig       `json:"modelConfig"`
		Privacy     DataPrivacyConfig `json:"privacy"`

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

		OwnedBy *EncodingStrategy `json:"ownedBy"`

		CreatedAt *EncodingStrategy `json:"createdAt"`
		CreatedBy *EncodingStrategy `json:"createdBy"`

		UpdatedAt *EncodingStrategy `json:"updatedAt"`
		UpdatedBy *EncodingStrategy `json:"updatedBy"`

		DeletedAt *EncodingStrategy `json:"deletedAt"`
		DeletedBy *EncodingStrategy `json:"deletedBy"`
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

func (m Module) HasIssues() bool {
	return len(m.ModelConfig.Issues) > 0
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
		ConnectionID: m.ModelConfig.ConnectionID,

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

func (nm *ModelConfig) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*nm = ModelConfig{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, nm); err != nil {
			return errors.Wrapf(err, "cannot scan '%v' into ModelConfig", string(b))
		}
	}

	return nil
}

func (nm ModelConfig) Value() (driver.Value, error) {
	return json.Marshal(nm)
}

func ParseModelConfig(ss []string) (m ModelConfig, err error) {
	if len(ss) == 0 {
		return
	}

	err = json.Unmarshal([]byte(ss[0]), &m)
	return
}

func (nm *DataPrivacyConfig) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*nm = DataPrivacyConfig{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, nm); err != nil {
			return errors.Wrapf(err, "cannot scan '%v' into DataPrivacyConfig", string(b))
		}
	}

	return nil
}

func (nm DataPrivacyConfig) Value() (driver.Value, error) {
	return json.Marshal(nm)
}

func ParseDataPrivacyConfig(ss []string) (dpc DataPrivacyConfig, err error) {
	if len(ss) == 0 {
		return
	}

	err = json.Unmarshal([]byte(ss[0]), &dpc)
	return
}
