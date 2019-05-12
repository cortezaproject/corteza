package permissions

func (set RuleSet) merge(rules ...*Rule) (out RuleSet, err error) {
	var (
		o    int
		olen = len(set)

		skipInherited = func(r *Rule) (b bool, e error) {
			return r != nil, nil
		}

		merged = set
	)

	if olen == 0 {
		// Nothing exists yet, just assign
		merged = rules
	} else {
	newRules:
		for _, rule := range rules {
			for ; o < olen; o++ {
				// Never go beyond the last old rule
				if merged[o].Equals(rule) {
					merged[o].Access = rule.Access

					// only one rule can match so proceed with next new rule
					continue newRules
				}
			}

			// none of the old rules matched, append
			merged = append(merged, rule)
		}

	}

	// Filter out all rules with access = inherit
	return merged.Filter(skipInherited)
}

func (set RuleSet) split() (inherited, rest RuleSet) {
	inherited, rest = RuleSet{}, RuleSet{}

	for _, r := range set {
		if r.Access == Inherit {
			inherited = append(inherited, r)
		} else {
			rest = append(rest, r)

		}
	}

	return
}
