package types

// RevisionRef implementations for the project work item types.
//
// Work items (unlike compose namespaces, workflows, agents, ...) file
// against the chain ROOT project via ProjectRef(), and separately carry
// their own nullable RevisionID naming the revision they're assigned to.
// That second reference can't be expressed by the generic, per-type
// single-field ProjectRef()/project_ref.gen.go mechanism, so it gets its own
// small interface here.
//
// Hand-maintained (not *.gen.go): unlike project_ref.gen.go, there is no cue
// codegen source driving this file yet, so a future `make codegen` run will
// not touch or regenerate it.

// RevisionRef returns the ID of the revision ProjectIncident is assigned to
// (0 = unassigned).
//
// It implements actionlog.RevisionResourcer, letting the action log
// attribute an event to the revision the affected work item is assigned to,
// independently of its chain-root ProjectRef().
func (r ProjectIncident) RevisionRef() uint64 {
	return r.RevisionID
}

// RevisionRef returns the ID of the revision ProjectFeature is assigned to
// (0 = unassigned).
//
// It implements actionlog.RevisionResourcer, letting the action log
// attribute an event to the revision the affected work item is assigned to,
// independently of its chain-root ProjectRef().
func (r ProjectFeature) RevisionRef() uint64 {
	return r.RevisionID
}

// RevisionRef returns the ID of the revision ProjectPrivacy is assigned to
// (0 = unassigned).
//
// It implements actionlog.RevisionResourcer, letting the action log
// attribute an event to the revision the affected work item is assigned to,
// independently of its chain-root ProjectRef().
func (r ProjectPrivacy) RevisionRef() uint64 {
	return r.RevisionID
}

// RevisionRef returns the ID of the revision ProjectTask is assigned to
// (0 = unassigned).
//
// It implements actionlog.RevisionResourcer, letting the action log
// attribute an event to the revision the affected work item is assigned to,
// independently of its chain-root ProjectRef().
func (r ProjectTask) RevisionRef() uint64 {
	return r.RevisionID
}

// RevisionRef returns the ID of the revision ProjectReview is assigned to
// (0 = unassigned).
//
// It implements actionlog.RevisionResourcer, letting the action log
// attribute an event to the revision the affected work item is assigned to,
// independently of its chain-root ProjectRef().
func (r ProjectReview) RevisionRef() uint64 {
	return r.RevisionID
}

// RevisionRef returns the ID of the revision ProjectBacklogItem is assigned
// to (0 = unassigned).
//
// It implements actionlog.RevisionResourcer, letting the action log
// attribute an event to the revision the affected work item is assigned to,
// independently of its chain-root ProjectRef().
func (r ProjectBacklogItem) RevisionRef() uint64 {
	return r.RevisionID
}
