package permissions

// Check verifies if role has access to perform an operation on a resource
//
// Overall flow:
//  - invalid resource, no access
//  - can this combination of roles perform an operation on this specific resource
//  - can this combination of roles perform an operation on any resource of the type (wildcard)
//  - can anyone/everyone perform an operation on this specific resource
//  - can anyone/everyone perform an operation on any resource of the type (wildcard)
func (set RuleSet) Check(res Resource, op Operation, roles ...uint64) (v Access) {

	if !res.IsValid() {
		return Deny
	}

	if len(roles) > 0 {
		if v = set.checkResource(res, op, roles...); v != Inherit {
			return
		}
	}

	if v = set.checkResource(res, op, EveryoneRoleID); v != Inherit {
		return
	}

	return
}

// Check ability to perform an operation on a specific and wildcard resource
func (set RuleSet) checkResource(res Resource, op Operation, roles ...uint64) (v Access) {
	if v = set.check(res, op, roles...); v != Inherit {
		return
	}

	if res.IsAppendable() {
		// Is this a specific resource and can we turn it ito a wildcarded-resource?
		if v = set.check(res.AppendWildcard(), op, roles...); v != Inherit {
			return
		}
	}

	return
}

// Check verifies if any of given roles has permission to perform an operation over a resource
//
// Will return Inherit when:
//  - no roles are given
//  - more than 1 role is given and one of the given roles is Everyone
//
// Will return Deny when:
//  - there is one rule with Deny value
//
// Will return Allow when:
//  - there is at least one rule with Allow value (and no Deny rules)
func (set RuleSet) check(res Resource, op Operation, roles ...uint64) (v Access) {
	v = Inherit

	for i := range set {
		// Ignore resources & operations that do not match
		if set[i].Resource != res || set[i].Operation != op {
			continue
		}

		// Check for every role
		for _, roleID := range roles {
			// Skip rules that do not match
			if set[i].RoleID != roleID || set[i].Access == Inherit {
				continue
			}

			v = set[i].Access // set to Allow

			// Return on first Deny
			if v == Deny {
				return
			}
		}
	}

	// If none of the rules matched, return Inherit (see 1st line)
	// if at least one of the rules allowed this op over a resource,
	// return Allow.
	return v
}
