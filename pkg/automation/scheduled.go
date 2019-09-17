package automation

import (
	"time"

	"github.com/DestinyWang/cronexpr"
)

type (
	// List of scheduled scripts with schedule time
	scheduledSet []schedule

	// ref to script & list of timestamps when we should run this
	schedule struct {
		scriptID   uint64
		timestamps []time.Time
		intervals  []cronexpr.Expression
	}
)

var (
	now = func() time.Time { return time.Now() }
)

const (
	pickInterval = time.Minute
)

// scans runnables and builds a set of scheduled scripts
func buildScheduleList(runables ScriptSet) scheduledSet {
	set := scheduledSet{}
	n := now()

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
				if ts.Before(n) {
					// in the past...
					continue
				}

				sch.timestamps = append(sch.timestamps, ts)
				continue
			}
		}

		for _, t := range s.triggers {
			if !t.IsInterval() {
				// only interested in interval scripts
				continue
			}

			if t.Condition == "" {
				continue
			}

			if itv, err := cronexpr.Parse(t.Condition); err == nil {
				sch.intervals = append(sch.intervals, *itv)
			}
		}

		// If there is anything useful in the schedule,
		// add it to the list
		if len(sch.timestamps)+len(sch.intervals) > 0 {
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
	n := now()
	thisMinute := n.Truncate(pickInterval)
	for _, s := range set {
		for _, t := range s.timestamps {
			if thisMinute.Equal(t.Truncate(pickInterval)) {
				uu = append(uu, s.scriptID)
				break
			}
		}

		for _, i := range s.intervals {
			// go back 1 interval step to check if script is ran in the next step;
			// since cronexpr can't assert if given time is in the interval.
			nn := thisMinute.Add(-pickInterval)
			if thisMinute.Equal(i.Next(nn).Truncate(pickInterval)) {
				uu = append(uu, s.scriptID)
				break
			}
		}
	}

	return uu
}
