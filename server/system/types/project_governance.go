package types

type (
	// ProjectGovernance tracks per-step workflow state on a project, keyed by
	// step key. Most steps (e.g. "summary", "data-sensitivity") only ever use
	// it to persist form values via SaveGovernanceStep and stay in the draft
	// status forever — per-section approval gates were removed, so those
	// steps never lock via their own submit/approve cycle. The one step that
	// drives the full submit → approve/request-changes state machine is the
	// well-known "publish" key (see ProjectGovernanceStepPublish): in
	// gated-mode projects, Publish() requires it to be approved and resets it
	// back to draft on success, so every publish needs its own fresh approval
	// cycle. The pipeline shape itself (steps, ordering) is client-side
	// config; the server only stores and transitions step state.
	//
	// The one other way a non-publish step's status moves is the "flag"
	// escalation (see Flag): a granter can request changes on any step, at
	// any time, as an out-of-band review comment — this bypasses the normal
	// submitted-only request-changes rule. The cross-step orchestration that
	// comes with it (sending "publish" back for review, clearing flags on
	// resubmit) lives in service.TransitionGovernanceStep, not here — this
	// type only models a single step's own state.
	ProjectGovernance map[string]*ProjectGovernanceStep

	ProjectGovernanceAction string
)

// ProjectGovernanceStepPublish is the well-known governance step key that
// gates publishing of a gated-mode project: Publish() requires this step to
// be ProjectGovernanceStatusApproved and resets it back to draft once the
// publish succeeds.
const ProjectGovernanceStepPublish = "publish"

const (
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

// Reset unconditionally returns the step to its initial draft state, clearing
// any review note. Unlike the reopen action, this is not a user-facing
// transition gated by CanGrantApproval — it's the system-driven counterpart
// used by:
//
//   - the publish flow, to require a fresh approval cycle for the next
//     publish once one succeeds, regardless of the step's current status;
//   - the auto-clear-on-resubmit side effect, which resets every flagged
//     (changes-requested) step other than "publish" back to draft when
//     "publish" is resubmitted, since resubmission asserts the feedback was
//     addressed.
func (s *ProjectGovernanceStep) Reset() {
	s.Status = ProjectGovernanceStatusDraft
	s.ReviewNote = ""
}

// Flag unconditionally moves the step to changes-requested with the given
// note, regardless of its current status — including a step with no prior
// entry, which the *ProjectGovernance.Step accessor auto-creates as a draft
// before Flag overwrites it. This is the counterpart to Transition's own
// request-changes case (which only fires from "submitted"): it models an
// out-of-band review flag that a granter can raise on any step at any time,
// not the normal submitted → changes-requested review response. Like Reset,
// it is not itself gated by a capability check — callers (see
// service.TransitionGovernanceStep) are responsible for requiring
// CanGrantApproval and restricting it to gated-mode projects before calling
// it.
func (s *ProjectGovernanceStep) Flag(note string) {
	s.Status = ProjectGovernanceStatusChangesRequested
	s.ReviewNote = note
}
