package rbac

import (
	"github.com/cortezaproject/corteza/server/pkg/slice"
)

func merge(base RuleSet, new ...*Rule) (out RuleSet) {
	var (
		o    int
		blen = len(base)
	)

	if blen == 0 {
		// Nothing exists yet, mark all as dirty
		for r := range new {
			new[r].dirty = true
		}

		return new
	} else {
		out = base

	newRules:
		for _, rule := range new {
			// Never go beyond the last base rule (blen)
			for o = 0; o < blen; o++ {
				if eq(out[o], rule) {
					out[o].dirty = out[o].Access != rule.Access
					out[o].Access = rule.Access

					// only one rule can match so proceed with next new rule
					continue newRules
				}
			}

			// none of the base new matched, append
			var c = *rule
			c.dirty = true

			out = append(out, &c)
		}

	}

	return
}

func eq(a, b *Rule) bool {
	if a == nil || b == nil {
		return false
	}

	return a.RoleID == b.RoleID &&
		a.Resource == b.Resource &&
		a.Operation == b.Operation
}

// Dirty returns list of deleted (Access==Inherit) and changed (dirty) rules
func flushable(set RuleSet) (deletable, updatable, final RuleSet) {
	deletable, updatable, final = RuleSet{}, RuleSet{}, RuleSet{}

	for _, r := range set {
		var c = *r
		if r.Access == Inherit {
			deletable = append(deletable, &c)
			continue
		}

		if r.dirty {
			updatable = append(updatable, &c)
		}

		final = append(final, &c)
	}

	return
}

// Significant roles (sigRoles) returns two list of significant roles.
//
// First slice are roles that are allowed to perform an operation on a specific resource (directly or indirectly),
// 2nd slice are roles that are denied the op.
func (set RuleSet) sigRoles(res string, op string) (aRR, dRR []uint64) {
	if hasWildcards(res) {
		// nothing to do here, we need a direct resource id (level=0)
		return
	}

	const (
		dirRes = 0
		indRes = 1
	)

	var (
		rr = map[Access]map[int][]uint64{
			Allow: {
				dirRes: make([]uint64, 0),
				indRes: make([]uint64, 0),
			},
			Deny: {
				dirRes: make([]uint64, 0),
				indRes: make([]uint64, 0),
			},
		}
	)

	// Extract all relevant rules (by op and resource) and group them by
	// access and distance (rules for direct resources and rules for indirect resources)
	for _, r := range set {
		if r.Operation != op {
			continue
		}

		if r.Resource == res {
			// direct rules
			rr[r.Access][dirRes] = append(rr[r.Access][dirRes], r.RoleID)
		}

		if matchResource(r.Resource, res) {
			// rules on all resources of this type
			rr[r.Access][indRes] = append(rr[r.Access][indRes], r.RoleID)
		}
	}

	// Process all extracted roles and make sure that are filtered and ordered by relevance:
	//  1. list of roles with denied operation directly on the resource
	dRoles := slice.ToUint64BoolMap(rr[Deny][dirRes])

	//  2. list of roles with allowed operation directly on the resource
	aRoles := make(map[uint64]bool)
	for _, r := range rr[Allow][dirRes] {
		aRoles[r] = !dRoles[r]
	}

	//  3. list of roles with denied operation indirectly on the resource
	for _, r := range rr[Deny][indRes] {
		dRoles[r] = true
	}

	//  4. list of roles with allowed operation indirectly on the resource
	for _, r := range rr[Allow][indRes] {
		aRoles[r] = !dRoles[r]
	}

	for r, chk := range aRoles {
		if chk {
			aRR = append(aRR, r)
		}
	}

	for r, chk := range dRoles {
		if chk {
			dRR = append(dRR, r)
		}
	}

	return
}

func clear(set []*Rule) {
	for _, r := range set {
		r.dirty = false
	}
}
