package actionlog

import (
	"testing"
)

func TestCannedPolies(t *testing.T) {
	tests := []struct {
		name string
		actn *Action
		mtch policyMatcher
		want bool
	}{
		{
			"debug policy should pass on anything",
			&Action{},
			MakeDebugPolicy(),
			true,
		},
		{
			"production policy should record notice",
			&Action{Severity: Notice},
			MakeProductionPolicy(),
			true,
		},
		{
			"production policy should not record debug",
			&Action{Severity: Debug},
			MakeProductionPolicy(),
			false,
		},
		{
			"disabled policy should not record anything",
			&Action{Severity: Alert},
			MakeDisabledPolicy(),
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
