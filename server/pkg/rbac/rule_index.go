package rbac

import (
	"sort"
)

type (
	// ruleIndex indexes all given RBAC rules to optimize lookup times
	ruleIndex struct {
		// children map[uint64]*ruleIndexNode
		bits map[string]map[string]*Rule
	}
)

// buildRuleIndex indexes the given rules for optimal lookups
//
// The build isn't that cleanned up but the lookup is good, I promise <3
func buildRuleIndex(rules []*Rule) (index *ruleIndex) {
	index = &ruleIndex{}
	index.add(rules...)
	return index
}

// add adds a new Rule to the index
func (index *ruleIndex) add(rules ...*Rule) {
	if index.bits == nil {
		index.bits = make(map[string]map[string]*Rule, len(rules)/2)
	}

	for _, r := range rules {
		if _, ok := index.bits[r.Operation]; !ok {
			index.bits[r.Operation] = make(map[string]*Rule, 4)
		}

		index.bits[r.Operation][r.Resource] = r
	}
}

// has checks if the rule is already in there
func (t *ruleIndex) has(r *Rule) bool {
	return len(t.collect(true, r.Operation, r.Resource)) > 0
}

// get returns the matching rules
func (t *ruleIndex) get(op, res string) (out []*Rule) {
	return t.collect(false, op, res)
}

// get returns all RBAC rules matching these constraints
//
// The get operation's lookup complexity is the longest RBAC key + 1 for
// the operation + 1 for the role.
//
// Our longest bit will be 6 so this is essentially constant time.
func (t *ruleIndex) collect(exact bool, op, res string) (out []*Rule) {
	if t.bits[op] == nil {
		return
	}

	for _, res := range permuteResource(res) {
		aux := t.bits[op][res]
		if aux == nil {
			continue
		}

		out = append(out, t.bits[op][res])
	}

	return out
}

// empty returns true if the index is empty
func (t *ruleIndex) empty() bool {
	return t == nil || t.bits == nil || len(t.bits) == 0
}

// matchingRule returns the first matching rule for the role, op, res
func (t *ruleIndex) matchingRule(op, res string) (out *Rule) {
	set := RuleSet(t.get(op, res))
	sort.Sort(set)

	for _, s := range set {
		if s.Access == Inherit {
			continue
		}

		return s
	}

	return nil
}
