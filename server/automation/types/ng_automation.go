package types

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/crusttech/human/server/pkg/ast"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/sql"
)

type (
	NgAutomation struct {
		ID        uint64                           `json:"automationID,string"`
		TenantID  uint64                           `json:"tenantID,string,omitempty"`
		ProjectID uint64                           `json:"projectID,string,omitempty"`
		Handle    string                           `json:"handle"`
		Labels map[string]labelTypes.LabelValue `json:"labels,omitempty"`
		Meta   *NgAutomationMeta                `json:"meta,omitempty"`

		Enabled bool `json:"enabled"`

		Scope *expr.Vars `json:"scope"`

		Triggers NgAutomationTriggerSet `json:"triggers"`
		Steps    NgAutomationStepSet    `json:"steps"`
		Paths    NgAutomationPathSet    `json:"paths"`

		// Collection of issues from the last parse
		Issues NgAutomationIssueSet `json:"issues,omitempty"`

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
		TenantID     uint64   `json:"tenantID,string,omitempty"`
		ProjectID    uint64   `json:"projectID,string,omitempty"`

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

	NgAutomationIssue struct {
		// url encoded location of the error:
		Culprit     map[string]int `json:"culprit"`
		Description string         `json:"description"`
	}

	NgAutomationTrigger struct {
		ID      uint64                           `json:"triggerID,string"`
		Labels  map[string]labelTypes.LabelValue `json:"labels,omitempty"`
		Handle  string                           `json:"handle"`
		Meta    *NgTriggerMeta                   `json:"meta,omitempty"`
		Enabled bool                             `json:"enabled"`

		// Resource type that can trigger the automation
		ResourceType string `json:"resourceType"`

		// Event type that can trigger the automation
		EventType string `json:"eventType"`

		// Trigger constraints
		Constraints []NgTriggerConstraint `json:"constraints"`

		// Initial input scope,
		// will be merged merged with automation variables
		Input *expr.Vars `json:"input"`

		// Input parameter schema this trigger accepts from callers
		InputSchema NgAutomationTriggerSchema `json:"inputSchema,omitempty"`
	}

	NgAutomationTriggerSchema []NgAutomationTriggerParam

	NgAutomationTriggerParam struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Required    bool   `json:"required"`
		Description string `json:"description"`
	}

	NgTriggerConstraint struct {
		Name   string                     `json:"name"`
		Op     string                     `json:"op,omitempty"`
		Values []NgTriggerConstraintValue `json:"values,omitempty"`
	}

	NgTriggerConstraintValue struct {
		Type  string `json:"@type"`
		Value string `json:"@value"`
	}

	NgTriggerMeta struct {
		Short       string             `json:"short"`
		Description string             `json:"description"`
		Visual      NgAutomationVisual `json:"visual"`
		Icon        *NgAutomationIcon  `json:"icon,omitempty"`
	}

	NgAutomationMeta struct {
		Short       string             `json:"short"`
		Description string             `json:"description"`
		Visual      NgAutomationVisual `json:"visual"`
		Icon        *NgAutomationIcon  `json:"icon,omitempty"`
	}

	NgAutomationIcon struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}

	NgAutomationVisual struct {
		// @todo...
	}

	NgAutomationExecParams struct {
		EntryPoint string `json:"entryPoint"`
		// @todo

		EventType    string `json:"eventType"`
		ResourceType string `json:"resourceType"`

		// Wait for workflow to be executed even if it's deferred
		Wait bool `json:"wait"`

		Input *expr.Vars `json:"input"`
	}

	NgAutomationStep struct {
		ID     uint64               `json:"stepID,string"`
		Handle string               `json:"handle"`
		Meta   NgAutomationStepMeta `json:"meta"`

		Kind string `json:"kind"`
		Ref  string `json:"ref"`

		Arguments []*Expr `json:"arguments"`
		Results   []*Expr `json:"results"`

		// Recoverable enables pause-and-retry on RecoverableError.
		Recoverable bool `json:"recoverable,omitempty"`

		// MaxRetries is the per-phase retry limit (0 = unlimited pause).
		MaxRetries int `json:"maxRetries,omitempty"`
	}

	NgAutomationPath struct {
		ParentID  uint64        `json:"parentID,string"`
		ChildID   uint64        `json:"childID,string"`
		Condition *ast.ASTNode  `json:"condition,omitempty"`

		// Kind distinguishes special paths.
		// "error" marks the catch handler path; empty = normal flow.
		Kind string `json:"kind,omitempty"`

		Meta NgAutomationPathMeta `json:"meta"`
	}

	NgAutomationStepMeta struct {
		Short       string             `json:"short"`
		Description string             `json:"description"`
		Visual      NgAutomationVisual `json:"visual"`
		Icon        *NgAutomationIcon  `json:"icon,omitempty"`

		// Extra holds step-kind-specific metadata (e.g. "message" and "recoverable" for error steps).
		Extra map[string]any `json:"extra,omitempty"`
	}

	NgAutomationPathMeta struct {
		Short       string             `json:"short"`
		Description string             `json:"description"`
		Visual      NgAutomationVisual `json:"visual"`
		Icon        *NgAutomationIcon  `json:"icon,omitempty"`
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

func (set *NgAutomationIssueSet) Scan(src any) error          { return sql.ParseJSON(src, set) }
func (set NgAutomationIssueSet) Value() (driver.Value, error) { return json.Marshal(set) }

func (set NgAutomationIssueSet) Error() string {
	out := make([]string, 0, 4)
	for _, s := range set {
		out = append(out, s.Description)
	}

	return strings.Join(out, ", ")
}
