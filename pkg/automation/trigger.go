package automation

import (
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
		//  - "before"
		//  - "after"
		//  - "on"
		//  - "at"
		Event string `json:"event" db:"event"`

		// Arbitrary data for trigger condition
		//
		// It is caller's responsibility to encode, decode and verify conditions
		Condition string `json:"condition" db:"condition"`

		ScriptID uint64 `json:"scriptID,string" db:"rel_script"`

		// Is trigger enabled or disabled?
		Enabled bool `json:"enabled" db:"enabled"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt"`
		CreatedBy uint64     `db:"created_by" json:"createdBy,string" `
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		UpdatedBy uint64     `db:"updated_by" json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
		DeletedBy uint64     `db:"deleted_by" json:"deletedBy,string,omitempty" `
	}

	TriggerFilter struct {
		Resource string
		Event    string
		ScriptID uint64

		IncDeleted bool

		// Standard paging fields & helpers
		rh.PageFilter
	}

	TriggerConditionChecker func(string) bool
)

const (
	EVENT_TYPE_INTERVAL  = "interval"
	EVENT_TYPE_TIMESTAMP = "at"
)

// IsValid checks if trigger is enabled and not deleted
func (t *Trigger) IsValid() bool {
	return t != nil && t.Enabled && t.DeletedAt == nil
}

// IsDeferred - not called as consequence of a user's action (create, delete, update)
func (t Trigger) IsDeferred() bool {
	return t.Event == EVENT_TYPE_INTERVAL || t.Event == EVENT_TYPE_TIMESTAMP
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

		if m.Resource != t.Resource {
			// event should match
			continue withTriggers
		}

		if m.Event != t.Event {
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
