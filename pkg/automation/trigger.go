package automation

import (
	"errors"
	"strconv"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	Event string

	Trigger struct {
		ID uint64 `json:"triggerID,string" db:"id"`

		// Resource that triggered the event
		//  - "compose:" (unspec, general)
		//  - "compose:record"
		//  - "compose:namespace"
		Resource string `json:"resource" db:"resource"`

		// Event name, arbitrary string
		//  - "manual"
		//  - "interval"
		//  - "deferred"
		//  - "before"
		//  - "after"
		Event string `json:"event" db:"event"`

		// Arbitrary data for trigger condition
		//
		// It is caller's responsibility to encode, decode and verify conditions
		Condition string `json:"condition" db:"event_condition"`

		ScriptID uint64 `json:"scriptID,string" db:"rel_script"`

		// Is trigger enabled or disabled?
		Enabled bool `json:"enabled" db:"enabled"`

		Weight int `json:"weight" db:"weight"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt"`
		CreatedBy uint64     `db:"created_by" json:"createdBy,string" `
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		UpdatedBy uint64     `db:"updated_by" json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
		DeletedBy uint64     `db:"deleted_by" json:"deletedBy,string,omitempty" `
	}

	TriggerFilter struct {
		Resource  string `json:"resource"`
		Event     string `json:"event"`
		Condition string `json:"condition"`
		ScriptID  uint64 `json:"scriptID"`

		IncDeleted bool `json:"incDeleted"`

		// Standard paging fields & helpers
		rh.PageFilter
	}

	TriggerConditionChecker func(string) bool
)

const (
	EVENT_TYPE_INTERVAL = "interval"
	EVENT_TYPE_DEFERRED = "deferred"
)

var (
	ErrAutomationTriggerInvalidResource  = errors.New("AutomationTriggerInvalidResource")
	ErrAutomationTriggerInvalidCondition = errors.New("AutomationTriggerInvalidCondition")
	ErrAutomationTriggerInvalidEvent     = errors.New("AutomationTriggerInvalidEvent")
)

// IsValid checks if trigger is enabled and not deleted
func (t *Trigger) IsValid() bool {
	return t != nil && t.Enabled && t.DeletedAt == nil
}

// IsDeferred - not called as consequence of a user's action (create, delete, update)
func (t Trigger) IsInterval() bool {
	return t.Event == EVENT_TYPE_INTERVAL
}

// IsDeferred - not called as consequence of a user's action (create, delete, update)
func (t Trigger) IsDeferred() bool {
	return t.Event == EVENT_TYPE_DEFERRED
}

// Uint64Condition converts condition to uint64
//
// Errors are ignored
func (t Trigger) Uint64Condition() (o uint64) {
	o, _ = strconv.ParseUint(t.Condition, 10, 64)
	return
}

// HasMatch checks if any og the triggers in a set matches the given parameters
func (set TriggerSet) HasMatch(m Trigger, ff ...TriggerConditionChecker) bool {
withTriggers:
	for _, t := range set {
		if !t.IsValid() {
			// only valid can match
			continue withTriggers
		}

		if m.ID > 0 && m.ID != t.ID {
			// Are we looking for a particular trigger?
			continue withTriggers
		}

		if m.Resource != "" && m.Resource != t.Resource {
			// event should match
			continue withTriggers
		}

		if m.Event != "" && m.Event != t.Event {
			// event should match
			continue withTriggers
		}

		// Go through all condition checking functions
		// All of them should return true for trigger to match
		for _, fn := range ff {
			if !fn(t.Condition) {
				continue withTriggers
			}
		}

		return true
	}

	return false
}
