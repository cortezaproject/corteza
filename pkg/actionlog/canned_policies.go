package actionlog

func MakeDebugPolicy() policyMatcher {
	return NewPolicyAll()
}

func MakeProductionPolicy() policyMatcher {
	return NewPolicyAll(
		NewPolicyAny(
			// Match all actions from automation
			NewPolicyMatchRequestOrigin(RequestOrigin_Automation),

			// Ignore debug actions
			NewPolicyNegate(NewPolicyMatchSeverity(Debug, Info)),
		),
	)
}

func MakeDisabledPolicy() policyMatcher {
	return NewPolicyNone()
}
