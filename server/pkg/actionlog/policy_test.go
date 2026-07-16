package actionlog

import (
	"testing"
)

func TestPolicyMatchers(t *testing.T) {
	tests := []struct {
		name string
		actn *Action
		mtch policyMatcher
		want bool
	}{
		{
			"simple any",
			&Action{},
			NewPolicyAny(),
			false,
		},
		{
			"simple all",
			&Action{},
			NewPolicyAll(),
			true,
		},
		{
			"should match severity",
			&Action{Severity: Emergency},
			NewPolicyMatchSeverity(Emergency),
			true,
		},
		{
			"should not match severity",
			&Action{Severity: Debug},
			NewPolicyMatchSeverity(Alert),
			false,
		},
		{
			"should match one of severity",
			&Action{Severity: Warning},
			NewPolicyMatchSeverity(Emergency, Debug, Warning),
			true,
		},
		{
			"should not match one of severity",
			&Action{Severity: Info},
			NewPolicyNegate(NewPolicyMatchSeverity(Emergency, Debug, Warning)),
			true,
		},
		{
			"complex match",
			&Action{Severity: Warning, Resource: "foo", Action: "do"},
			NewPolicyAll(
				NewPolicyMatchSeverity(Warning),
			),
			true,
		},
		{
			"complex miss",
			&Action{Severity: Warning, Resource: "bar", Action: "do"},
			NewPolicyAll(
				NewPolicyMatchSeverity(Emergency),
			),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != tt.mtch.Match(tt.actn) {
				if tt.want {
					t.Errorf("expecting to match")
				} else {
					t.Errorf("expecting not to match ")
				}
			}
		})
	}
}
