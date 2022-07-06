package rbac

import (
	"fmt"
)

type (
	Rule struct {
		RoleID    uint64 `json:"roleID,string"`
		Resource  string `json:"resource"`
		Operation string `json:"operation"`
		Access    Access `json:"access,string"`

		// Do we need to flush it to storage?
		dirty bool
	}

	RuleSet []*Rule

	// OptRuleSet RBAC rule index (operation / role ID / rules)
	OptRuleSet map[string]map[uint64]RuleSet
)

func (r Rule) String() string {
	return fmt.Sprintf("%s %d to %s on %s", r.Access, r.RoleID, r.Operation, r.Resource)
}

func indexRules(rules []*Rule) OptRuleSet {
	i := make(OptRuleSet)
	for _, r := range rules {
		if i[r.Operation] == nil {
			i[r.Operation] = make(map[uint64]RuleSet)
		}

		if i[r.Operation][r.RoleID] == nil {
			i[r.Operation][r.RoleID] = RuleSet{}
		}

		i[r.Operation][r.RoleID] = append(i[r.Operation][r.RoleID], r)
	}

	return i
}

func (set RuleSet) Len() int      { return len(set) }
func (set RuleSet) Swap(i, j int) { set[i], set[j] = set[j], set[i] }
func (set RuleSet) Less(i, j int) bool {
	return level(set[i].Resource) > level(set[j].Resource)
}

func (set RuleSet) FilterAccess(a Access) (out RuleSet) {
	out = make(RuleSet, 0, len(set))

	for _, s := range set {
		if s.Access == a {
			out = append(out, s)
		}
	}

	return out
}

func (set RuleSet) FilterResource(res string) (out RuleSet) {
	for _, r := range set {
		if !matchResource(res, r.Resource) {
			continue
		}
		out = append(out, r)
	}

	return
}

// FilterRules will filter the rules based on given parameter(specific),
//		If params is true then it will return only the specific rules otherwise it will return non-specific rules
func (set RuleSet) FilterRules(specific bool) (out RuleSet) {
	for _, r := range set {
		if specific == isSpecific(r.Resource) {
			out = append(out, r)
		}
	}

	return
}

// AllowRule helper func to create allow rule
func AllowRule(id uint64, r, o string) *Rule {
	return &Rule{id, r, o, Allow, false}
}

// DenyRule helper func to create deny rule
func DenyRule(id uint64, r, o string) *Rule {
	return &Rule{id, r, o, Deny, false}
}

// InheritRule helper func to create inherit rule
func InheritRule(id uint64, r, o string) *Rule {
	return &Rule{id, r, o, Inherit, false}
}
