package automation

import (
	"testing"
)

func TestScript_CheckCompatibility(t *testing.T) {

	tests := []struct {
		name    string
		s       *Script
		t       *Trigger
		wantErr bool
	}{
		{name: "both nil",
			s:       nil,
			t:       nil,
			wantErr: true,
		},
		{name: "both vanilla",
			s:       &Script{},
			t:       &Trigger{},
			wantErr: false,
		},
		{name: "deferred trigger with UA script",
			s:       &Script{RunInUA: true},
			t:       &Trigger{Event: EVENT_TYPE_INTERVAL},
			wantErr: true,
		},
		{name: "deferred trigger with invoker security",
			s:       &Script{RunAs: 0},
			t:       &Trigger{Event: EVENT_TYPE_INTERVAL},
			wantErr: true,
		},
		{name: "deferred trigger with invoker security",
			s:       &Script{RunAs: 1, RunInUA: false},
			t:       &Trigger{Event: EVENT_TYPE_INTERVAL},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.CheckCompatibility(tt.t); (err != nil) != tt.wantErr {
				t.Errorf("CheckCompatibility() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
