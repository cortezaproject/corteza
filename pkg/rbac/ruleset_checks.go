package rbac

import (
	"sort"
)

func check(indexedRules OptRuleSet, rolesByKind partRoles, op, res string) Access {
	a, _, _ := evaluate(indexedRules, rolesByKind, op, res, false)
	return a
}

func evaluate(indexedRules OptRuleSet, rolesByKind partRoles, op, res string, parent bool) (Access, *Rule, explanation) {
	if member(rolesByKind, AnonymousRole) && len(rolesByKind) > 1 {
		// Integrity check; when user is member of anonymous role
		// should not be member of any other type of role

		return Deny, nil, stepIntegrity
	}

	if member(rolesByKind, BypassRole) {
		// if user has at least one bypass role, we allow access
		return Allow, nil, stepBypass
	}

	if len(indexedRules) == 0 {
		// no rules no access
		return Inherit, nil, stepRuleless
	}

	var rules RuleSet

	// Priority is important here. We want to have
	// stable RBAC check behaviour and ability
	// to override allow/deny depending on how niche the role (type) is:
	//  - context (eg owners) are more niche than common
	//  - rules for common roles are more important than
	for _, kind := range []roleKind{ContextRole, CommonRole, AuthenticatedRole, AnonymousRole} {
		// no roles if this kind
		if len(rolesByKind[kind]) == 0 {
			continue
		}

		rules = nil
		for roleID, r := range indexedRules[op] {
			if !rolesByKind[kind][roleID] {
				continue
			}
			rules = append(rules, r...)
		}

		// When evaluating access for parent, omit the exact tule
		if parent {
			nr := make(RuleSet, 0, len(rules))
			for _, r := range rules {
				if r.Resource != res {
					nr = append(nr, r)
				}
			}
			rules = nr
		}

		r, access := checkRulesByResource(rules, op, res)
		if access != Inherit {
			return access, r, stepEvaluated
		}
	}

	return Inherit, nil, stepRuleless
}

// Check given resource match and operation on all given rules
func checkRulesByResource(set RuleSet, op, res string) (*Rule, Access) {
	// Make sure rules are always sorted (by level)
	// to avoid any kind of unstable behaviour
	sort.Sort(set)

	for _, r := range set {
		if !matchResource(r.Resource, res) {
			continue
		}

		if op != r.Operation {
			continue
		}

		if r.Access != Inherit {
			return r, r.Access
		}
	}

	return nil, Inherit
}

// at least one of the roles must be set to true
func member(r partRoles, k roleKind) bool {
	for _, is := range r[k] {
		if is {
			return true
		}
	}

	return false
}
