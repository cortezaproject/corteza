package rbac

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

func ruleByRole(base RuleSet, roleID uint64) (out RuleSet) {
	for _, r := range base {
		if r.RoleID == roleID {
			out = append(out, r)
		}
	}

	return
}

// Dirty returns list of deleted (Access==Inherit) and changed (dirty) rules
func flushable(set RuleSet) (inherited, rest RuleSet) {
	inherited, rest = RuleSet{}, RuleSet{}

	for _, r := range set {
		var c = *r
		if r.Access == Inherit {
			inherited = append(inherited, &c)
		} else if r.dirty {
			rest = append(rest, &c)
		}
	}

	return
}

func clear(set []*Rule) {
	for _, r := range set {
		r.dirty = false
	}
}
