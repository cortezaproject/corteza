package types

type (
	// ProjectGovernance tracks per-step workflow state on a project, keyed by
	// step key. Two shapes of step live here:
	//
	//   - Ordinary Build/Govern steps (e.g. "summary", "data-sensitivity"):
	//     a granter reviews them directly — Approve marks the step approved,
	//     request-changes flags it (see Flag) — at any time, with no submit
	//     stage of their own. SaveGovernanceStep persists their form values
	//     while they are still editable.
	//   - The well-known "publish" key (see ProjectGovernanceStepPublish),
	//     which drives the project-level approval state machine
	//     (submit → approve / request-changes). Publish() requires it to be
	//     approved and resets it back to draft on success, so every publish
	//     needs its own fresh submit → approve cycle.
	//
	// The pipeline shape itself (steps, ordering) is client-side config; the
	// server only stores and transitions step state. The cross-step
	// orchestration (a flagged step sending "publish" back for review,
	// clearing flags on resubmit, blocking project approval while any step is
	// flagged) lives in service.TransitionGovernanceStep, not here — this type
	// only models a single step's own state.
	ProjectGovernance map[string]*ProjectGovernanceStep

	ProjectGovernanceAction string
)

// ProjectGovernanceStepPublish is the well-known governance step key that
// drives the project-level approval state machine and gates publishing:
// Publish() requires this step to be ProjectGovernanceStatusApproved and
// resets it back to draft once the publish succeeds.
const ProjectGovernanceStepPublish = "publish"

const (
	// ProjectGovernanceActionSubmit sends the publish step (a draft or
	// changes-requested project) for approval.
	ProjectGovernanceActionSubmit ProjectGovernanceAction = "submit"
	// ProjectGovernanceActionApprove approves a step: the submitted publish
	// step, or any Build/Govern step directly.
	ProjectGovernanceActionApprove ProjectGovernanceAction = "approve"
	// ProjectGovernanceActionRequestChanges flags a step with a note.
	ProjectGovernanceActionRequestChanges ProjectGovernanceAction = "request-changes"
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
// (approve, request-changes); submit needs only request-approval.
func (a ProjectGovernanceAction) RequiresGrant() bool {
	switch a {
	case ProjectGovernanceActionApprove,
		ProjectGovernanceActionRequestChanges:
		return true
	}
	return false
}

// Transition applies the action to the "publish" step's state machine:
//
//	submit:          draft|changes-requested → submitted (clears note)
//	approve:         submitted               → approved  (clears note)
//	request-changes: submitted               → changes-requested (sets note)
//
// Returns false for invalid transitions, leaving the step untouched. Only the
// publish step uses this machine; ordinary Build/Govern steps are approved or
// flagged directly (see Approve/Flag) with no submit stage.
func (s *ProjectGovernanceStep) Transition(action ProjectGovernanceAction, note string) bool {
	switch {
	case action == ProjectGovernanceActionSubmit &&
		(s.Status == ProjectGovernanceStatusDraft || s.Status == ProjectGovernanceStatusChangesRequested):
		s.Status, s.ReviewNote = ProjectGovernanceStatusSubmitted, ""
	case action == ProjectGovernanceActionApprove && s.Status == ProjectGovernanceStatusSubmitted:
		s.Status, s.ReviewNote = ProjectGovernanceStatusApproved, ""
	case action == ProjectGovernanceActionRequestChanges && s.Status == ProjectGovernanceStatusSubmitted:
		s.Status, s.ReviewNote = ProjectGovernanceStatusChangesRequested, note
	default:
		return false
	}
	return true
}

// Approve unconditionally marks the step approved and clears any review note,
// regardless of its current status. It is the direct-review counterpart to
// Flag used by ordinary Build/Govern steps, which have no submit stage: a
// granter approves the step (clearing a prior changes-requested flag) at any
// time. Like Flag and Reset, it is not itself gated by a capability check —
// callers (see service.TransitionGovernanceStep) require CanGrantApproval
// before calling it.
func (s *ProjectGovernanceStep) Approve() {
	s.Status = ProjectGovernanceStatusApproved
	s.ReviewNote = ""
}

// Editable reports whether the step's form values may currently be changed.
func (s *ProjectGovernanceStep) Editable() bool {
	return s.Status == ProjectGovernanceStatusDraft || s.Status == ProjectGovernanceStatusChangesRequested
}

// Reset unconditionally returns the step to its initial draft state, clearing
// any review note. This is not a user-facing transition gated by a capability
// — it's the system-driven counterpart used by:
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
// before Flag overwrites it. It is the direct request-changes review a
// granter can raise on any step at any time. Like Reset and Approve, it is
// not itself gated by a capability check — callers (see
// service.TransitionGovernanceStep) require CanGrantApproval before calling
// it.
func (s *ProjectGovernanceStep) Flag(note string) {
	s.Status = ProjectGovernanceStatusChangesRequested
	s.ReviewNote = note
}
