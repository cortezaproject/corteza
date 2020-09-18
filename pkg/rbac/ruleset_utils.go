package rbac

import "github.com/cortezaproject/corteza-server/pkg/slice"

// merge applies new rules (changes) to existing set and mark all changes as dirty
func (set RuleSet) merge(rules ...*Rule) (out RuleSet) {
	var (
		o    int
		olen = len(set)
	)

	if olen == 0 {
		// Nothing exists yet, mark all as dirty
		for r := range rules {
			rules[r].dirty = true
		}

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

// Missing compares cmp with existing set
// and returns rules that exists in set but not in cmp
func (set RuleSet) Diff(cmp RuleSet) RuleSet {
	diff := RuleSet{}
base:
	for _, s := range set {
		for _, c := range cmp {
			if c.Equals(s) {
				continue base
			}
		}

		diff = append(diff, s)
	}

	return diff
}

// Roles returns list of unique id of all roles in the rule set
func (set RuleSet) Roles() []uint64 {
	roles := make([]uint64, 0)
	for _, r := range set {
		if slice.HasUint64(roles, r.RoleID) {
			continue
		}

		roles = append(roles, r.RoleID)
	}

	return roles
}

func (set RuleSet) ByResource(res Resource) RuleSet {
	out, _ := set.Filter(func(r *Rule) (bool, error) {
		return res == r.Resource, nil
	})
	return out
}

func (set RuleSet) AllAllows() RuleSet {
	return set.ByAccess(Allow)
}

func (set RuleSet) AllDenies() RuleSet {
	return set.ByAccess(Deny)
}

func (set RuleSet) ByAccess(a Access) RuleSet {
	out, _ := set.Filter(func(r *Rule) (bool, error) {
		return a == r.Access, nil
	})
	return out
}

func (set RuleSet) ByRole(roleID uint64) RuleSet {
	out, _ := set.Filter(func(r *Rule) (bool, error) {
		return roleID == r.RoleID, nil
	})
	return out
}
