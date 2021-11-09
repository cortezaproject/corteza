package rbac

import (
	"sort"
)

func check(indexedRules OptRuleSet, rolesByKind partRoles, op, res string) Access {
	if member(rolesByKind, AnonymousRole) && len(rolesByKind) > 1 {
		// Integrity check; when user is member of anonymous role
		// should not be member of any other type of role
		return Deny
	}

	if member(rolesByKind, BypassRole) {
		// if user has at least one bypass role, we allow access
		return Allow
	}

	if len(indexedRules) == 0 {
		// no rules no access
		return Inherit
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

		// user has at least one bypass role
		if kind == BypassRole {
			return Allow
		}

		rules = nil
		for roleID, r := range indexedRules[op] {
			if !rolesByKind[kind][roleID] {
				continue
			}
			rules = append(rules, r...)
		}

		access := checkRulesByResource(rules, op, res)
		if access != Inherit {
			return access
		}
	}

	return Inherit
}

// Check given resource match and operation on all given rules
func checkRulesByResource(set RuleSet, op, res string) Access {
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
			return r.Access
		}
	}

	return Inherit
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
