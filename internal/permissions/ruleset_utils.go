package permissions

// merge applies new rules (changes) to existing set and mark  all changes as dirty
func (set RuleSet) merge(rules ...*Rule) (out RuleSet) {
	var (
		o    int
		olen = len(set)
	)

	if olen == 0 {
		// Nothing exists yet
		return rules
	} else {
		out = set

	newRules:
		for _, rule := range rules {
			// Never go beyond the last old rule (olen)
			for o = 0; o < olen; o++ {
				if out[o].Equals(rule) {
					out[o].dirty = out[o].Access != rule.Access
					out[o].Access = rule.Access

					// only one rule can match so proceed with next new rule
					continue newRules
				}
			}

			// none of the old rules matched, append
			var c = *rule
			c.dirty = true

			out = append(out, &c)
		}

	}

	return
}

// dirty returns list of changed (dirty==true) and deleted (Access==Inherit) rules
func (set RuleSet) dirty() (inherited, rest RuleSet) {
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

// reset dirty flag
func (set RuleSet) clear() {
	_ = set.Walk(func(rule *Rule) error {
		rule.dirty = false
		return nil
	})
}
