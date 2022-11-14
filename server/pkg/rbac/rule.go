package rbac

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/resource"
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

	ruleIndexWrap struct {
		index *resource.IndexNode
		rules RuleSet
	}

	// OptRuleSet RBAC rule index (operation / role ID / rules)
	OptRuleSet map[string]map[uint64]*ruleIndexWrap
)

func (r Rule) Clone() *Rule {
	return &r
}

func (r Rule) String() string {
	return fmt.Sprintf("%s %d to %s on %s", r.Access, r.RoleID, r.Operation, r.Resource)
}

func indexRules(rules []*Rule) OptRuleSet {
	i := make(OptRuleSet)
	for _, r := range rules {
		if i[r.Operation] == nil {
			i[r.Operation] = make(map[uint64]*ruleIndexWrap)
		}

		if i[r.Operation][r.RoleID] == nil {
			i[r.Operation][r.RoleID] = &ruleIndexWrap{
				index: resource.NewIndex(),
			}
		}

		i[r.Operation][r.RoleID].index.Add(r, r.IndexPath()...)
		i[r.Operation][r.RoleID].rules = append(i[r.Operation][r.RoleID].rules, r)
	}

	return i
}

func (r Rule) IndexPath() (out [][]string) {
	pts := strings.Split(r.Resource, pathSep)

	for _, p := range pts {
		out = append(out, []string{p})
	}

	return
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
		if !matchResource(r.Resource, res) {
			continue
		}
		out = append(out, r)
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
