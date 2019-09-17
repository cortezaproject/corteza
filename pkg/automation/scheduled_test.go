package automation

import (
	"reflect"
	"testing"
	"time"
)

func ss2tt(ss ...string) []time.Time {
	tt := make([]time.Time, len(ss))
	for i, s := range ss {
		tt[i], _ = time.Parse(time.RFC3339, s)
	}
	return tt
}

func TestScheduleBuilder(t *testing.T) {
	tests := []struct {
		name string
		ss   ScriptSet
		sch  scheduledSet
	}{
		{name: "basics",
			ss: ScriptSet{
				&Script{ID: 1, Enabled: true, triggers: TriggerSet{
					&Trigger{Enabled: true, Event: EVENT_TYPE_DEFERRED, Condition: "2000-01-01T00:02:00+02:00"},
					&Trigger{Enabled: true, Event: EVENT_TYPE_DEFERRED, Condition: "2000-01-01T00:03:00+02:00"},
				}},
				&Script{ID: 2, Enabled: true, triggers: TriggerSet{
					&Trigger{Enabled: true, Event: EVENT_TYPE_DEFERRED, Condition: "2000-01-01T00:02:00+02:00"},
					&Trigger{Enabled: true, Event: EVENT_TYPE_DEFERRED, Condition: "2000-01-01T00:03:00+02:00"},
				}},
			},

			sch: scheduledSet{
				schedule{scriptID: 1, timestamps: ss2tt("2000-01-01T00:02:00+02:00", "2000-01-01T00:03:00+02:00")},
				schedule{scriptID: 2, timestamps: ss2tt("2000-01-01T00:02:00+02:00", "2000-01-01T00:03:00+02:00")},
			},
		},
	}

	n, _ := time.Parse(time.RFC3339, "2000-01-01T00:00:00+02:00")
	now = func() time.Time { return n }

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := buildScheduleList(tt.ss)
			if !reflect.DeepEqual(out, tt.sch) {
				t.Errorf("Result do not match %v %v", out, tt.sch)
			}
		})
	}
}

func TestSchedulePicker(t *testing.T) {
	tests := []struct {
		name string
		sch  scheduledSet
		ids  []uint64
	}{
		{name: "one",
			ids: []uint64{1},
			sch: scheduledSet{
				schedule{scriptID: 1, timestamps: ss2tt("2000-01-01T00:01:00+02:00", "2000-01-01T00:03:00+02:00")},
				schedule{scriptID: 2, timestamps: ss2tt("2000-01-01T00:04:00+02:00", "2000-01-01T00:05:00+02:00")},
			},
		},
		{name: "two",
			ids: []uint64{1, 2},
			sch: scheduledSet{
				schedule{scriptID: 1, timestamps: ss2tt("2000-01-01T00:01:00+02:00")},
				schedule{scriptID: 2, timestamps: ss2tt("2000-01-01T00:01:00+02:00", "2000-01-01T00:05:00+02:00")},
			},
		},
	}

	n, _ := time.Parse(time.RFC3339, "2000-01-01T00:01:50+02:00")
	now = func() time.Time { return n }

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.sch.pick()

			if !reflect.DeepEqual(out, tt.ids) {
				t.Errorf("Result do not match %v %v", out, tt.ids)
			}
		})
	}
}
