package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	labelTypes "github.com/cortezaproject/corteza/server/pkg/label/types"
	"github.com/cortezaproject/corteza/server/pkg/sql"
)

type (
	NgAutomation struct {
		ID     uint64                           `json:"automationID,string"`
		Handle string                           `json:"handle"`
		Labels map[string]labelTypes.LabelValue `json:"labels,omitempty"`
		Meta   *NgAutomationMeta                `json:"meta,omitempty"`

		Enabled bool `json:"enabled"`

		Scope *expr.Vars `json:"scope"`

		Triggers NgAutomationTriggerSet `json:"triggers"`
		Steps    NgAutomationStepSet    `json:"steps"`
		Paths    NgAutomationPathSet    `json:"paths"`

		RunAs     uint64     `json:"runAs,string"`
		OwnedBy   uint64     `json:"ownedBy,string"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
	}

	NgAutomationFilter struct {
		AutomationID []string `json:"automationID"`

		Handle string `json:"handle"`

		Query string `json:"query"`

		Deleted  filter.State `json:"deleted"`
		Disabled filter.State `json:"disabled"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*NgAutomation) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	NgAutomationTrigger struct {
		ID      uint64                           `json:"triggerID,string"`
		Labels  map[string]labelTypes.LabelValue `json:"labels,omitempty"`
		Meta    *NgTriggerMeta                   `json:"meta,omitempty"`
		Enabled bool                             `json:"enabled"`

		// Resource type that can trigger the automation
		ResourceType string `json:"resourceType"`

		// Event type that can trigger the automation
		EventType string `json:"eventType"`

		// Trigger constraints
		Constraints TriggerConstraintSet `json:"constraints"`

		// Initial input scope,
		// will be merged merged with automation variables
		Input *expr.Vars `json:"input"`

		OwnedBy   uint64     `json:"ownedBy,string"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
	}

	NgTriggerMeta struct {
		Description string             `json:"description"`
		Visual      NgAutomationVisual `json:"visual"`
	}

	NgAutomationMeta struct {
		Name        string             `json:"name"`
		Description string             `json:"description"`
		Visual      NgAutomationVisual `json:"visual"`
	}

	NgAutomationVisual struct {
		// @todo...
	}

	NgAutomationExecParams struct {
		// @todo

		EventType    string
		ResourceType string

		// Wait for workflow to be executed even if it's deferred
		Wait bool

		Input *expr.Vars
	}

	NgAutomationStep struct {
		ID   uint64               `json:"stepID"`
		Meta NgAutomationStepMeta `json:"meta"`

		Kind string `json:"kind"`
		Ref  string

		Arguments []*Expr
		Results   []*Expr
	}

	NgAutomationPath struct {
		ParentID uint64
		ChildID  uint64

		Meta NgAutomationPathMeta
	}

	NgAutomationStepMeta struct {
		Name        string
		Description string
		Visual      NgAutomationVisual
	}

	NgAutomationPathMeta struct {
		Name        string
		Description string
		Visual      NgAutomationVisual
	}
)

func (set *NgAutomationMeta) Scan(src any) error          { return sql.ParseJSON(src, set) }
func (set NgAutomationMeta) Value() (driver.Value, error) { return json.Marshal(set) }

func (set *NgAutomationTriggerSet) Scan(src any) error          { return sql.ParseJSON(src, set) }
func (set NgAutomationTriggerSet) Value() (driver.Value, error) { return json.Marshal(set) }

func (set *NgAutomationStepSet) Scan(src any) error          { return sql.ParseJSON(src, set) }
func (set NgAutomationStepSet) Value() (driver.Value, error) { return json.Marshal(set) }

func (set *NgAutomationPathSet) Scan(src any) error          { return sql.ParseJSON(src, set) }
func (set NgAutomationPathSet) Value() (driver.Value, error) { return json.Marshal(set) }
