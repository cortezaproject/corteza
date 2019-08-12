package automation

import (
	"testing"
)

func TestTriggerSet_HasMatch(t *testing.T) {
	type args struct {
		m Trigger
	}
	tests := []struct {
		name string
		set  TriggerSet
		args args
		want bool
	}{
		{
			name: "simple match",
			set:  TriggerSet{nil, &Trigger{}, &Trigger{Event: "e", Enabled: true}, nil, &Trigger{}},
			args: args{m: Trigger{Event: "e"}},
			want: true,
		}, {
			name: "simple miss",
			set:  TriggerSet{nil, &Trigger{}, &Trigger{Event: "e", Enabled: true}, nil, &Trigger{}},
			args: args{m: Trigger{Event: "E"}},
			want: false,
		}, {
			name: "specific",
			set:  TriggerSet{nil, &Trigger{}, &Trigger{ID: 2, Enabled: true}, nil, &Trigger{}},
			args: args{m: Trigger{ID: 2}},
			want: true,
		}, {
			name: "invalid",
			set:  TriggerSet{nil, &Trigger{}, &Trigger{Event: "e", Enabled: false}, nil, &Trigger{}},
			args: args{m: Trigger{Event: "e"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.HasMatch(tt.args.m); got != tt.want {
				t.Errorf("HasMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
