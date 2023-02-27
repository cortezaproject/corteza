package rbac

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/cast2"
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

func (set RuleSet) FilterOperation(op string) (out RuleSet) {
	out = make(RuleSet, 0, len(set))

	for _, s := range set {
		if s.Operation == op {
			out = append(out, s)
		}
	}

	return out
}

// FilterResource returns rules that match given list of resources
// Wildcards are not used!
//
// Note that empty resource list will return ALL rules!
func (set RuleSet) FilterResource(rr ...Resource) (out RuleSet) {
	if len(rr) == 0 {
		return set
	}

	out = RuleSet{}
	for _, rule := range set {
		for _, res := range rr {
			if res.RbacResource() != rule.Resource {
				continue
			}

			out = append(out, rule)
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

func (u *Rule) SetValue(name string, pos uint, v any) (err error) {
	switch name {
	case "roleID", "RoleID":
		return cast2.Uint64(v, &u.RoleID)
	case "resource", "Resource":
		return cast2.String(v, &u.Resource)
	default:
		return u.setValue(name, pos, v)
	}
}

func (u *Rule) setValue(name string, pos uint, v any) (err error) {
	pp := strings.Split(name, ".")

	switch pp[0] {
	case "resource", "Resource", "Path", "path":
		ix, err := strconv.ParseUint(pp[1], 10, 64)
		if err != nil {
			return err
		}

		res := strings.Split(u.Resource, "/")

		aux := ""
		err = cast2.String(v, &aux)

		// +1 bacause the first bit is the resource
		res[ix+1] = aux
		u.Resource = strings.Join(res, "/")
		return err
	}

	return
}

func (u *Rule) GetValue(name string, pos uint) (v any, err error) {
	switch name {
	case "roleID", "RoleID":
		return u.RoleID, nil
	case "resource", "Resource":
		return u.Resource, nil
	}

	return
}

func (r Rule) GetID() uint64 {
	// The resource does not define an ID field
	return 0
}
