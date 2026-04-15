package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	labelTypes "github.com/cortezaproject/corteza/server/pkg/label/types"
	"github.com/cortezaproject/corteza/server/pkg/sql"
)

type (
	TriggerDefinition struct {
		ID      uint64 `json:"triggerDefinitionID,string"`
		AgentID uint64 `json:"agentID,string,omitempty"`
		Handle  string `json:"handle"`

		Meta *TriggerDefinitionMeta `json:"meta,omitempty"`

		// Input parameters this trigger accepts from callers
		InputSchema TriggerDefinitionSchema `json:"inputSchema"`

		// Output parameters this trigger returns to callers
		OutputSchema TriggerDefinitionSchema `json:"outputSchema"`

		// When true, triggers using this definition are NOT registered
		// on the eventbus. They are invoked directly via Exec/ExecAndWait.
		SkipEventBus bool `json:"skipEventBus"`

		// Labels, ownership, timestamps...
		Labels    map[string]labelTypes.LabelValue `json:"labels,omitempty"`
		OwnedBy   uint64                           `json:"ownedBy,string"`
		CreatedAt time.Time                        `json:"createdAt,omitempty"`
		CreatedBy uint64                           `json:"createdBy,string"`
		UpdatedAt *time.Time                       `json:"updatedAt,omitempty"`
		UpdatedBy uint64                           `json:"updatedBy,string,omitempty"`
		DeletedAt *time.Time                       `json:"deletedAt,omitempty"`
		DeletedBy uint64                           `json:"deletedBy,string,omitempty"`
	}

	TriggerDefinitionMeta struct {
		Short       string            `json:"short"`
		Description string            `json:"description"`
		Icon        *NgAutomationIcon `json:"icon,omitempty"`
	}

	TriggerDefinitionSchema []TriggerDefinitionParam

	TriggerDefinitionParam struct {
		Name        string `json:"name"`
		Type        string `json:"type"` // "String" default, extensible later
		Required    bool   `json:"required"`
		Description string `json:"description"`
	}

	TriggerDefinitionFilter struct {
		AgentID             []string `json:"agentID"`
		TriggerDefinitionID []string `json:"triggerDefinitionID"`
		Handle              string   `json:"handle"`
		Query               string   `json:"query"`

		Deleted filter.State `json:"deleted"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*TriggerDefinition) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}
)

func (set *TriggerDefinitionMeta) Scan(src any) error          { return sql.ParseJSON(src, set) }
func (set TriggerDefinitionMeta) Value() (driver.Value, error) { return json.Marshal(set) }

func (set *TriggerDefinitionSchema) Scan(src any) error          { return sql.ParseJSON(src, set) }
func (set TriggerDefinitionSchema) Value() (driver.Value, error) { return json.Marshal(set) }

func ParseTriggerDefinitionSchema(ss []string) (p TriggerDefinitionSchema, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func ParseTriggerDefinitionMeta(ss []string) (p *TriggerDefinitionMeta, err error) {
	if len(ss) == 0 {
		return
	}
	err = json.Unmarshal([]byte(ss[0]), &p)
	return
}
