package rbac

import "sort"

type (
	resolution string

	Trace struct {
		Resource   string     `json:"resource"`
		Operation  string     `json:"operation"`
		Access     Access     `json:"access"`
		Roles      []uint64   `json:"roles"`
		Rules      []*Rule    `json:"rules,omitempty"`
		Resolution resolution `json:"resolution,omitempty"`
	}
)

const (
	failedIntegrityCheck resolution = "failed-integrity-check"
	bypassRoleMembership resolution = "bypass-role-membership"
	noRules              resolution = "no-rules"
	noMatch              resolution = "no-match"
)

// baseTraceInfo updates given check trace struct
//
// If Trace is nil, function terminates early and does not panic
func baseTraceInfo(t *Trace, res, op string, rolesByKind partRoles) {
	if t == nil {
		return
	}

	t.Resource = res
	t.Operation = op

	for _, rr := range rolesByKind {
		for r := range rr {
			t.Roles = append(t.Roles, r)
		}
	}

	sort.Slice(t.Roles, func(i, j int) bool {
		return t.Roles[i] < t.Roles[j]
	})
}

// resolve updates access and resolution and returns access info unmodified
// if nil is passed for Trace, function silently ignores that
func resolve(t *Trace, access Access, res resolution) Access {
	if t != nil {
		t.Access = access
		t.Resolution = res
	}

	return access
}

// ensure stable order of trace rules
func sortTraceRules(trace *Trace) {
	// using custom sorting (and not RulSet implementation) because
	// we might have multiple roles per resource.
	sort.Slice(trace.Rules, func(i, j int) bool {
		var (
			rli = level(trace.Rules[i].Resource)
			rlj = level(trace.Rules[j].Resource)
		)

		if rli != rlj {
			return rli < rlj
		}

		return trace.Rules[i].RoleID < trace.Rules[j].RoleID
	})
}
