package automation

import (
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	Script struct {
		ID uint64 `json:"scriptID,string" db:"id"`

		NamespaceID uint64 `json:"namespaceID,string" db:"rel_namespace"`

		Name string `json:"name" db:"name"`

		// (URL) Where did we get the source from?
		SourceRef string `json:"sourceRef" db:"source_ref"`

		// Code
		Source string `json:"source" db:"source"`

		// No need to wait for script to return the value
		Async bool `json:"async" db:"async"`

		// Who is running this script?
		// Leave it at 0 for the current user (security invoker) or
		// set ID of specific user (security definer)
		RunAs uint64 `json:"runAs,string" db:"rel_runner"`

		// Where can we run this script? user-agent? corredor service?
		RunInUA bool `json:"runInUA" db:"run_in_ua"`

		// Are you doing something that can take more time?
		// specify timeout (in milliseconds)
		Timeout uint `json:"timeout" db:"timeout"`

		// Is it critical to run this script successfully?
		Critical bool `json:"critical" db:"critical"`

		Enabled bool `json:"enabled" db:"enabled"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt"`
		CreatedBy uint64     `db:"created_by" json:"createdBy,string" `
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		UpdatedBy uint64     `db:"updated_by" json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
		DeletedBy uint64     `db:"deleted_by" json:"deletedBy,string,omitempty" `

		// Serves as container for valid triggers for runnable scripts internal cache
		// and for as a transport on create/update operations
		//
		// Node: on c/u op., we currently just merge current state with the given list,
		// w/o updating the rest of the script's triggers
		triggers TriggerSet

		// How are we merging?
		tms triggersMergeStrategy

		// Script running credentials
		//
		// We'll store credentials for security-defined scripts
		// here when (runnable) scripts are loaded
		credentials string
	}

	ScriptFilter struct {
		NamespaceID uint64 `json:"namespaceID,string"`

		Name       string
		Query      string
		Resource   string
		IncDeleted bool `json:"incDeleted"`

		permissions.AccessCheck `json:"-"`

		// Standard paging fields & helpers
		rh.PageFilter
	}

	triggersMergeStrategy int
)

const (
	// Ignore the given triggers
	STMS_IGNORE triggersMergeStrategy = iota

	// Create triggers, no pre-checks
	STMS_FRESH

	// Update existing with new
	STMS_UPDATE

	// Replace existing with new
	STMS_REPLACE
)

var (
	ErrAutomationScriptInvalid     = errors.New("AutomationScriptInvalid")
	ErrAutomationScriptMissingUser = errors.New("AutomationScriptMissingUser")
)

// IsValid - enabled, deleted?
func (s *Script) IsValid() bool {
	return s != nil && s.Enabled && s.DeletedAt == nil
}

// Verify - sanity check of script's properties
func (s Script) Verify() error {
	if s.RunAsDefined() && s.RunInUA {
		return errors.New("user-agent engine does not support run-as-defined scripts")
	}

	if s.Critical && s.RunInUA {
		return errors.New("user-agent engine scripts can not be critical")
	}

	return nil
}

// IsCompatible verifies if trigger can be added to a script
func (s *Script) CheckCompatibility(t *Trigger) error {
	if s == nil {
		return errors.New("not compatible with nil script")
	}
	if s == nil || t == nil {
		return errors.New("not compatible with nil trigger")
	}

	if t.IsDeferred() || t.IsInterval() {
		if s.RunInUA {
			return errors.New("deferred triggers are not compatible with user-agent scripts")
		}

		if s.RunAsInvoker() {
			return errors.New("deferred triggers are not compatible with run-as-invoker scripts")
		}
	}

	return nil
}

// FilterByTrigger
//
// Filters non-UA scripts that match event and resource + all extra conditions
func (set ScriptSet) FilterByTrigger(event, resource string, cc ...TriggerConditionChecker) (out ScriptSet) {
	out, _ = set.Filter(func(s *Script) (bool, error) {
		return s.IsValid() && s.triggers.HasMatch(Trigger{Event: event, Resource: resource}, cc...), nil
	})

	return
}

// FindByName finds undeleted script in a cetain namespace
func (set ScriptSet) FindByName(name string, ids ...uint64) *Script {
	var namespaceID uint64
	if len(ids) > 0 {
		namespaceID = ids[0]
	}

	for i := range set {
		if set[i].NamespaceID == namespaceID && set[i].Name == name && set[i].DeletedAt == nil {
			return set[i]
		}
	}

	return nil
}

// RunAsDefined - script should be run with pre-defined privileges (user)
func (s Script) RunAsDefined() bool {
	return s.RunAs > 0
}

// RunAsInvoker - this script should run with invoker's privileges (user)
func (s Script) RunAsInvoker() bool {
	return s.RunAs == 0
}

// AddTrigger appends one or more triggers to internal list of triggers on script struct
//
// We do not do any compatibility check (See Script.CheckCompatibility());
// this is only an utility func that helps us pass data along
func (s *Script) AddTrigger(strategy triggersMergeStrategy, tt ...*Trigger) {
	s.tms = strategy

	if s.tms == STMS_IGNORE {
		return
	}

	if s.tms == STMS_REPLACE {
		s.triggers = TriggerSet{}
	}

	// Make sure all your trigger belong to us (ref same script or no ref):
	for _, t := range tt {
		if t.ScriptID == 0 || t.ScriptID == s.ID {
			s.triggers = append(s.triggers, t)
		}
	}
}

func (s *Script) Triggers() TriggerSet {
	return s.triggers
}

func (s Script) HasEvent(event string) bool {
	return s.triggers.HasMatch(Trigger{Event: event})
}

func (s Script) Credentials() string {
	return s.credentials
}

func MakeMatcherIDCondition(id uint64) TriggerConditionChecker {
	// We'll be comparing strings, not uint64!
	var s = strconv.FormatUint(id, 10)

	return func(c string) bool {
		return s == c
	}
}
