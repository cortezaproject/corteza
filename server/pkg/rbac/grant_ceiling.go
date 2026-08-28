package rbac

import (
	"context"
	"fmt"
)

// CanPassOn answers the one question every permission write turns on: does the
// calling user already hold this operation on this resource?
//
// It runs the same evaluation the server runs when the caller performs the
// operation themselves, which is what makes the boundary structural rather than
// a rule someone has to remember. Two entry points because one of them refuses
// half the question:
//
//   - Can is the full evaluation, org tree included, and is what a concrete
//     resource deserves. It answers false for ANY resource carrying a wildcard —
//     checkValidity rejects those outright — so it cannot be the only path.
//   - Trace evaluates a wildcard resource honestly: rules are matched with
//     path.Match against the requested pattern, so a rule at the same or a
//     broader scope matches and a narrower one does not. That is exactly
//     "do you hold this at this breadth". It skips the user-group branch, so it
//     can only under-report, which is the safe direction for a ceiling.
//
// Both err towards refusal. A caller wrongly refused loses a grant they could
// have made by hand; a caller wrongly allowed escalates.
func CanPassOn(ctx context.Context, resource, operation string) bool {
	svc := Global()
	if svc == nil {
		return false
	}

	var (
		ses = ContextToSession(ctx)
		res = NewResource(resource)
	)

	if svc.Can(ses, operation, res) {
		return true
	}

	return svc.Trace(ses, operation, res).Access == Allow
}

// EnforceGrantCeiling refuses a rule set that would hand out more access than
// the caller holds.
//
// Granting is not one permission but two questions: may you write rules at all,
// and may you write THIS one. Only the first was ever asked here, so anyone
// holding a component's grant could give any role every operation that
// component declares, including the ones they could not perform themselves.
//
// It covers inherit as well as allow and deny: an inherit deletes a rule, and
// the rule it deletes may be the deny that was holding an operation shut.
//
// Every rule is judged before any is written. A permission change that
// half-lands leaves a role in a state nobody asked for, and every refusal here
// is one the caller can fix and retry.
func EnforceGrantCeiling(ctx context.Context, rules ...*Rule) error {
	if Global() == nil {
		// Nothing is evaluating access at all, so there is no ceiling to read.
		// Refusing here would break the paths that write rules before the
		// service exists rather than close a hole.
		return nil
	}

	for _, r := range rules {
		if r == nil || CanPassOn(ctx, r.Resource, r.Operation) {
			continue
		}

		return fmt.Errorf(
			"refused: you do not hold %q on %s yourself, and a permission can only be passed on by someone who has it. Nothing was changed",
			r.Operation, r.Resource,
		)
	}

	return nil
}
