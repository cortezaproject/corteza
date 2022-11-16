package actionlog

import "github.com/cortezaproject/corteza/server/pkg/slice"

type (
	policyMatcher interface {
		Match(*Action) bool
	}

	logPolicyAny struct {
		mm []policyMatcher
	}

	logPolicyAll struct {
		mm []policyMatcher
	}

	logPolicyNegate struct {
		m policyMatcher
	}

	logPolicyMatchRequestOrigin struct {
		origin map[string]bool
	}

	logPolicyMatchResource struct {
		resources map[string]bool
	}

	logPolicyMatchAction struct {
		actions map[string]bool
	}

	logPolicyMatchSeverity struct {
		severities map[Severity]bool
	}

	logPolicyNoop struct {
		v bool
	}
)

// NewPolicyNone ignores all action logs
func NewPolicyNone() policyMatcher {
	return &logPolicyNoop{v: false}
}

// NewPolicyAny returns policy where at least one of the sub-policies should match
func NewPolicyAny(mm ...policyMatcher) policyMatcher {
	return &logPolicyAny{mm: mm}
}

func (p logPolicyAny) Match(a *Action) bool {
	for _, m := range p.mm {
		if m.Match(a) {
			return true
		}
	}

	return false
}

// NewPolicyAll returns policy where all sub-policies should match
func NewPolicyAll(mm ...policyMatcher) policyMatcher {
	return &logPolicyAll{mm: mm}
}

func (p logPolicyAll) Match(a *Action) bool {
	for _, m := range p.mm {
		if !m.Match(a) {
			return false
		}
	}

	return true
}

// NewPolicyNegate negates passed policy
func NewPolicyNegate(m policyMatcher) policyMatcher {
	return &logPolicyNegate{m: m}
}

func (p logPolicyNegate) Match(a *Action) bool {
	return !p.m.Match(a)
}

// NewPolicyMatchResource matches resources
func NewPolicyMatchResource(rr ...string) policyMatcher {
	return &logPolicyMatchResource{resources: slice.ToStringBoolMap(rr)}
}

func (p logPolicyMatchResource) Match(a *Action) bool {
	return p.resources[a.Resource]
}

// NewPolicyMatchAction matches action
func NewPolicyMatchAction(aa ...string) policyMatcher {
	return &logPolicyMatchAction{actions: slice.ToStringBoolMap(aa)}
}

func (p logPolicyMatchAction) Match(a *Action) bool {
	return p.actions[a.Action]
}

// NewPolicyMatchSeverity matches severity
func NewPolicyMatchSeverity(ss ...Severity) policyMatcher {
	var p = &logPolicyMatchSeverity{severities: make(map[Severity]bool)}

	for _, s := range ss {
		p.severities[s] = true
	}

	return p
}

func (p logPolicyMatchSeverity) Match(a *Action) bool {
	return p.severities[a.Severity]
}

// NewPolicyMatchRequestOrigin matches resources
func NewPolicyMatchRequestOrigin(rr ...string) policyMatcher {
	return &logPolicyMatchRequestOrigin{origin: slice.ToStringBoolMap(rr)}
}

func (p logPolicyMatchRequestOrigin) Match(a *Action) bool {
	return p.origin[a.RequestOrigin]
}

// Match Internal policy
func (p logPolicyNoop) Match(*Action) bool {
	return p.v
}
