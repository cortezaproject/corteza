package automation

import (
	"time"
)

type (
	// List of scheduled scripts with schedule time
	scheduledSet []schedule

	// ref to script & list of timestamps when we should run this
	schedule struct {
		scriptID   uint64
		timestamps []time.Time
		// intervals []interface{}
	}
)

var (
	now = func() time.Time { return time.Now() }
)

// scans runnables and builds a set of scheduled scripts
func buildScheduleList(runables ScriptSet) scheduledSet {
	set := scheduledSet{}

	_ = runables.Walk(func(s *Script) error {
		sch := schedule{scriptID: s.ID}
		for _, t := range s.triggers {
			if !t.IsDeferred() {
				// only interested in deferred scripts
				continue
			}

			if t.Condition == "" {
				continue
			}

			if ts, err := time.Parse(time.RFC3339, t.Condition); err == nil {
				ts = ts.Truncate(time.Minute)
				if ts.Before(now()) {
					// in the past...
					continue
				}

				sch.timestamps = append(sch.timestamps, ts)
				continue
			}

			// @todo parse cron format and fill intervals
		}

		// If there is anything useful in the schedule,
		// add it to the list
		if len(sch.timestamps) > 0 {
			set = append(set, sch)
		}

		return nil
	})

	return set
}

// scans scheduled set and picks out all candidates that
// are scheduled this minute
//
// Script is executed only once, even if
// is scheduled multiple times,
func (set scheduledSet) pick() []uint64 {
	uu := []uint64{}
	thisMinute := now().Truncate(time.Minute)
	for _, s := range set {
		for _, t := range s.timestamps {
			if thisMinute.Equal(t) {
				uu = append(uu, s.scriptID)
				break
			}
		}

		// @todo pick from intervals
	}

	return uu
}
