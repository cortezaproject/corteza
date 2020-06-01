package actionlog

func MakeDebugPolicy() policyMatcher {
	return NewPolicyAll()
}

func MakeProductionPolicy() policyMatcher {
	return NewPolicyAll(
		// Ignore debug actions
		NewPolicyNegate(NewPolicyMatchSeverity(Debug)),
	)
}
