package service

import (
	"fmt"
	"sort"

	"github.com/cortezaproject/corteza/server/automation/types"
)

// analyzeWorkflow performs non-fatal static analysis on a workflow
// definition and returns a WorkflowIssueSet suitable for assignment to
// wf.Warnings. Analysis is intentionally conservative: it only reports
// on patterns that can be decided purely from the workflow graph and
// the declared step metadata, without evaluating expressions.
//
// Current checks:
//
//   - Parallel-join variable conflicts: for every join gateway, walk
//     every upstream branch back to the matching fork (or the workflow
//     entry) and collect the set of variable names each branch can
//     write. Any variable name written in more than one branch is
//     reported as a potential conflict at the join step. Matches the
//     runtime detectConflicts check in wfexec.joinGateway.Exec but
//     fires at save time so authors see the warning before they run.
//
// Assumptions (documented so future readers know the limits):
//
//   - Only statically declared writes are considered. An Expressions
//     step's writes are its argument target names; a Function or
//     ExecWorkflow step's writes are its result target names. Dynamic
//     writes inside expression bodies (e.g. vars[computed] = ...) are
//     not detected.
//   - Subworkflow contents are NOT walked. A subworkflow invocation's
//     writes into the calling scope are exactly the ExecWorkflow step's
//     own declared results, so the step's static declaration already
//     captures the full contract. This is how the runtime wires
//     subworkflow results back too (see convExecWorkflowStep).
//   - The analyzer never fails the workflow. If an unsupported topology
//     is encountered (e.g. a broken path with no reachable fork) it
//     silently skips that join and moves on.
func analyzeWorkflow(wf *types.Workflow) types.WorkflowIssueSet {
	if wf == nil || len(wf.Steps) == 0 || len(wf.Paths) == 0 {
		return nil
	}

	// step index lookup: stepID -> position in wf.Steps (used as
	// culprit index for per-step badge placement in the editor)
	stepIdx := make(map[uint64]int, len(wf.Steps))
	stepByID := make(map[uint64]*types.WorkflowStep, len(wf.Steps))
	for i, s := range wf.Steps {
		stepIdx[s.ID] = i
		stepByID[s.ID] = s
	}

	// parents[stepID] = list of parent step IDs feeding into stepID
	parents := make(map[uint64][]uint64, len(wf.Steps))
	for _, p := range wf.Paths {
		parents[p.ChildID] = append(parents[p.ChildID], p.ParentID)
	}

	var warnings types.WorkflowIssueSet

	for _, s := range wf.Steps {
		if s.Kind != types.WorkflowStepKindGateway || s.Ref != "join" {
			continue
		}

		// For each incoming branch, collect the set of variables that
		// branch can write. A "branch" is defined as the set of steps
		// strictly between a fork (or entry) and the join, following
		// parent edges backwards from each direct parent of the join.
		// The home join ID is passed in so the walker can stop at any
		// *other* join it encounters — otherwise, in nested parallels,
		// the outer branch's writes get unioned into every inner
		// branch and produce spurious inner-join conflicts.
		branchWrites := make([]map[string]bool, 0, len(parents[s.ID]))
		for _, parentID := range parents[s.ID] {
			writes := collectBranchWrites(parentID, s.ID, stepByID, parents)
			if len(writes) > 0 {
				branchWrites = append(branchWrites, writes)
			}
		}

		if len(branchWrites) < 2 {
			continue
		}

		// For each variable name, count how many branches wrote it.
		counts := make(map[string]int)
		for _, bw := range branchWrites {
			for name := range bw {
				counts[name]++
			}
		}

		var conflicts []string
		for name, c := range counts {
			if c >= 2 {
				conflicts = append(conflicts, name)
			}
		}
		if len(conflicts) == 0 {
			continue
		}

		sort.Strings(conflicts)
		for _, name := range conflicts {
			warnings = warnings.Append(
				fmt.Errorf(
					"variable %q is written by more than one parallel branch joining here; the value from the first-configured parent wins at runtime and the others are discarded. Rename the conflicting outputs (e.g. %s_a, %s_b) if you need them all.",
					name, name, name,
				),
				map[string]int{"step": stepIdx[s.ID]},
			)
		}
	}

	return warnings
}

// collectBranchWrites walks parent edges backwards from the given
// starting step, collecting the set of variable names that any step
// on the branch can write. Walking stops when it hits:
//
//   - A fork gateway (upstream boundary of the branch).
//   - Any join gateway other than `homeJoinID`. In nested parallels,
//     an inner branch's upstream chain passes through the inner join,
//     then the inner fork, then the outer branch — walking through
//     the inner join would pull the outer branch's writes into every
//     inner branch and produce spurious inner-join conflicts. Stopping
//     at any *other* join cuts that bleed.
//   - No more parent edges (workflow entry).
//
// The starting step itself is included.
//
// Cycles are possible if the workflow has iterators (loop-back edges).
// A visited set prevents infinite walks; iterator bodies still get
// their writes collected on the first visit.
func collectBranchWrites(
	startID uint64,
	homeJoinID uint64,
	stepByID map[uint64]*types.WorkflowStep,
	parents map[uint64][]uint64,
) map[string]bool {
	writes := make(map[string]bool)
	visited := make(map[uint64]bool)

	var walk func(id uint64)
	walk = func(id uint64) {
		if visited[id] {
			return
		}
		visited[id] = true

		s := stepByID[id]
		if s == nil {
			return
		}

		// Stop at a fork — that's the upstream boundary of the branch.
		// Do not collect the fork's writes; it doesn't produce any.
		if s.Kind == types.WorkflowStepKindGateway && s.Ref == "fork" {
			return
		}

		// Stop at any join other than the one we're analyzing. The
		// inner join's own analysis will handle its incoming branches
		// independently; we must not cross it.
		if s.Kind == types.WorkflowStepKindGateway && s.Ref == "join" && id != homeJoinID {
			return
		}

		for name := range stepWrites(s) {
			writes[name] = true
		}

		for _, pid := range parents[id] {
			walk(pid)
		}
	}

	walk(startID)
	return writes
}

// stepWrites returns the set of variable names a step can write into
// the shared scope based on its declared metadata. Returns an empty
// map for step kinds that do not produce scope writes.
func stepWrites(s *types.WorkflowStep) map[string]bool {
	out := make(map[string]bool)
	if s == nil {
		return out
	}

	switch s.Kind {
	case types.WorkflowStepKindExpressions:
		// Expressions step declares its assignments via Arguments,
		// each with an explicit Target name that becomes the scope
		// key. This is how convExpressionStep / ExpressionsStep wire
		// outputs at runtime.
		for _, a := range s.Arguments {
			if a != nil && a.Target != "" {
				out[a.Target] = true
			}
		}

	case types.WorkflowStepKindFunction,
		types.WorkflowStepKindExecWorkflow:
		// Function / subworkflow call: step results are the declared
		// outputs. For ExecWorkflow specifically, the runtime evaluates
		// s.Results against the subworkflow's return value, so the
		// declared targets ARE the contract.
		for _, r := range s.Results {
			if r != nil && r.Target != "" {
				out[r.Target] = true
			}
		}

	case types.WorkflowStepKindIterator:
		// Iterator step results (item, counter, isFirst, ...) are
		// loop-local: they're injected into the body scope on each
		// iteration and do NOT persist past the iterator's exit path.
		// Treating them as branch writes produces false positives
		// (two parallel iterators both called "item" are not actually
		// in conflict at the downstream join). Steps inside the
		// iterator body still get their writes collected through the
		// normal parent-edge walk.

	case types.WorkflowStepKindErrHandler:
		// Error handler declares author-chosen variable names for
		// error / errorMessage / errorStepID via Results targets.
		// These become scope variables when a downstream step errors.
		for _, r := range s.Results {
			if r != nil && r.Target != "" {
				out[r.Target] = true
			}
		}

	default:
		// Error, Termination, Break, Continue, Debug, Delay, Prompt,
		// Visual and plain Gateway steps do not write to scope.
	}

	return out
}
