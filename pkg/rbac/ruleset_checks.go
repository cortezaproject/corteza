package rbac

import (
	"sort"
)

// function checks all given rules
//
//  - indexRules are rules optimized for quick lookup rules ar grouped in 2 levels:
//    by operation, and role, containing slice of rules (for specific operation and role)
//    at the deepest level
//
//  - rolesByKind are roles optimized for quick lookup
//    roles are grouped by kind and each kind contains fast-lookup (map[role-id]bool)
//
//  - op and res represent operation and resource that are checked
//
//  - trace is optional; when not nil, function will update trace struct
//    with information as it traverses and checks the rules
//
func check(indexedRules OptRuleSet, rolesByKind partRoles, op, res string, trace *Trace) Access {
	baseTraceInfo(trace, res, op, rolesByKind)

	if member(rolesByKind, AnonymousRole) && len(rolesByKind) > 1 {
		// Integrity check; when user is member of anonymous role
		// should not be member of any other type of role
		return resolve(trace, Deny, failedIntegrityCheck)
	}

	if member(rolesByKind, BypassRole) {
		// if user has at least one bypass role, we allow access
		return resolve(trace, Allow, bypassRoleMembership)
	}

	if len(indexedRules) == 0 {
		// no rules to check
		return resolve(trace, Inherit, noRules)
	}

	var (
		match   *Rule
		allowed bool
	)

	//
	if trace != nil {
		// from this point on, there is a chance trace (if set)
		// will contain some rules.
		//
		// Stable order needs to be ensured: there is no production
		// code that relies on that but tests might fail and API
		// response would be flaky.
		defer sortTraceRules(trace)
	}

	// Priority is important here. We want to have
	// stable RBAC check behaviour and ability
	// to override allow/deny depending on how niche the role (type) is:
	//  - context (eg owners) are more niche than common
	//  - rules for common roles are more important than authenticated and anonymous role types
	//
	// Note that bypass roles are intentionally ignored here; if user is member of
	// bypass role there is no need to check any other rule
	for _, kind := range []roleKind{ContextRole, CommonRole, AuthenticatedRole, AnonymousRole} {
		// not a member of any role of this kind
		if len(rolesByKind[kind]) == 0 {
			continue
		}

		// reset allowed to false
		// for each role kind
		allowed = false

		for roleID, rr := range indexedRules[op] {
			if !rolesByKind[kind][roleID] {
				continue
			}

			if len(rr) == 0 {
				// no rules found
				continue
			}

			// check all rules for each role the security-context
			if match = findRuleByResOp(rr, op, res); match == nil {
				// no rules match
				continue
			}

			if trace != nil {
				// if trace is enabled, append
				// each matching rule
				trace.Rules = append(trace.Rules, match)
			}

			if match.Access == Deny {
				// if we stumble upon Deny we short-circuit the check
				return resolve(trace, Deny, "")
			}

			if match.Access == Allow {
				// allow rule found, we need to check rules on other roles
				// before we allow it
				allowed = true
			}
		}

		if allowed {
			// at least one of the roles (per role type) in the security context
			// allows operation on a resource
			return resolve(trace, Allow, "")
		}
	}

	// No rule matched
	return resolve(trace, Inherit, noMatch)
}

// Check given resource match and operation on all given rules
func findRuleByResOp(set RuleSet, op, res string) *Rule {
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
			return r
		}
	}

	return nil
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
