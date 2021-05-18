package rbac

func (set RuleSet) Check(rolesByKind partRoles, res, op string) Access {
	for _, kind := range roleKindsByPriority() {
		if len(rolesByKind[kind]) == 0 {
			continue
		}

		if kind == BypassRole {
			return Allow
		}

		access := checkRulesByResource(filterRules(set, rolesByKind[kind], op), res, op)
		if access != Inherit {
			return access
		}
	}

	return Inherit
}

func checkOptimised(indexedRules OptRuleSet, rolesByKind partRoles, res, op string) Access {
	if len(rolesByKind) == 0 || len(indexedRules) == 0 {
		return Inherit
	}

	var rules []*Rule

	// looping through all role kinds
	for _, kind := range roleKindsByPriority() {
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

		access := checkRulesByResource(rules, res, op)
		if access != Inherit {
			return access
		}
	}

	return Inherit
}

// Check given resource match and operation on all given rules
//
// Function expects rules, sorted by level!
func checkRulesByResource(set []*Rule, res, op string) Access {
	for _, r := range set {
		if !matchResource(res, r.Resource) {
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
