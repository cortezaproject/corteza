package types

type (
	// ProjectGovernance tracks per-step workflow state of the project build
	// pipeline, keyed by step key (e.g. "summary", "data-sensitivity"). The
	// pipeline shape itself (steps, gates, ordering) is client-side config;
	// the server only stores and transitions step state.
	ProjectGovernance map[string]*ProjectGovernanceStep

	// ProjectGovernanceStep is one step's form values plus its approval state.
	ProjectGovernanceStep struct {
		Values     map[string]any          `json:"values,omitempty"`
		Status     ProjectGovernanceStatus `json:"status"`
		ReviewNote string                  `json:"reviewNote,omitempty"`
	}

	ProjectGovernanceStatus string
	ProjectGovernanceAction string
)

const (
	ProjectGovernanceStatusDraft            ProjectGovernanceStatus = "draft"
	ProjectGovernanceStatusSubmitted        ProjectGovernanceStatus = "submitted"
	ProjectGovernanceStatusApproved         ProjectGovernanceStatus = "approved"
	ProjectGovernanceStatusChangesRequested ProjectGovernanceStatus = "changes-requested"

	// ProjectGovernanceActionSubmit sends a draft/changes-requested step for approval.
	ProjectGovernanceActionSubmit ProjectGovernanceAction = "submit"
	// ProjectGovernanceActionApprove approves a submitted step.
	ProjectGovernanceActionApprove ProjectGovernanceAction = "approve"
	// ProjectGovernanceActionRequestChanges returns a submitted step with a note.
	ProjectGovernanceActionRequestChanges ProjectGovernanceAction = "request-changes"
	// ProjectGovernanceActionReopen unlocks an approved step back to draft.
	ProjectGovernanceActionReopen ProjectGovernanceAction = "reopen"
	// ProjectGovernanceActionRecall withdraws a pending submission back to draft.
	ProjectGovernanceActionRecall ProjectGovernanceAction = "recall"
)

// Step returns the entry for the given step key, initialising a draft entry
// when missing.
func (g *ProjectGovernance) Step(key string) *ProjectGovernanceStep {
	if *g == nil {
		*g = ProjectGovernance{}
	}
	if (*g)[key] == nil {
		(*g)[key] = &ProjectGovernanceStep{Status: ProjectGovernanceStatusDraft}
	}
	return (*g)[key]
}

// RequiresGrant reports whether the action needs the grant-approval capability
// (approve, request-changes, reopen); the others need request-approval.
func (a ProjectGovernanceAction) RequiresGrant() bool {
	switch a {
	case ProjectGovernanceActionApprove,
		ProjectGovernanceActionRequestChanges,
		ProjectGovernanceActionReopen:
		return true
	}
	return false
}

// Transition applies the action to the step state machine:
//
//	submit:          draft|changes-requested → submitted (clears note)
//	approve:         submitted               → approved  (clears note)
//	request-changes: submitted               → changes-requested (sets note)
//	reopen:          approved                → draft     (sets note)
//	recall:          submitted               → draft     (clears note)
//
// Returns false for invalid transitions, leaving the step untouched.
func (s *ProjectGovernanceStep) Transition(action ProjectGovernanceAction, note string) bool {
	switch {
	case action == ProjectGovernanceActionSubmit &&
		(s.Status == ProjectGovernanceStatusDraft || s.Status == ProjectGovernanceStatusChangesRequested):
		s.Status, s.ReviewNote = ProjectGovernanceStatusSubmitted, ""
	case action == ProjectGovernanceActionApprove && s.Status == ProjectGovernanceStatusSubmitted:
		s.Status, s.ReviewNote = ProjectGovernanceStatusApproved, ""
	case action == ProjectGovernanceActionRequestChanges && s.Status == ProjectGovernanceStatusSubmitted:
		s.Status, s.ReviewNote = ProjectGovernanceStatusChangesRequested, note
	case action == ProjectGovernanceActionReopen && s.Status == ProjectGovernanceStatusApproved:
		s.Status, s.ReviewNote = ProjectGovernanceStatusDraft, note
	case action == ProjectGovernanceActionRecall && s.Status == ProjectGovernanceStatusSubmitted:
		s.Status, s.ReviewNote = ProjectGovernanceStatusDraft, ""
	default:
		return false
	}
	return true
}

// Editable reports whether the step's form values may currently be changed.
func (s *ProjectGovernanceStep) Editable() bool {
	return s.Status == ProjectGovernanceStatusDraft || s.Status == ProjectGovernanceStatusChangesRequested
}
