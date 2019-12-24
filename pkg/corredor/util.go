package corredor

import (
	"github.com/cortezaproject/corteza-server/pkg/slice"
)

// removes onManual event type from trigger
// returns true if event type was removed or
// false if there was no onManual event
func popOnManualEventType(trigger *Trigger) (found bool) {
	for i := len(trigger.Events) - 1; i >= 0; i-- {
		if trigger.Events[i] == onManualEventType {
			found = true

			// remove from the list
			trigger.Events = append(trigger.Events[:i], trigger.Events[i+1:]...)
		}
	}

	return
}

func makeScriptFilter(f ManualScriptFilter) func(s *Script) (b bool, err error) {
	return func(s *Script) (b bool, err error) {
		b = true
		if len(f.ResourceTypes) > 0 {
			// Filtering by resource type,
			// at least one of the script's triggers should match
			b = false
			for _, t := range s.Triggers {
				if len(slice.IntersectStrings(f.ResourceTypes, t.Resources)) == 0 {
					b = true
				}
			}

			if !b {
				// No match by resource type, break
				return
			}
		}

		if len(f.EventTypes) > 0 {
			// Filtering by event type,
			// at least one of the script's triggers should match
			b = false
			for _, t := range s.Triggers {
				if len(slice.IntersectStrings(f.EventTypes, t.Events)) == 0 {
					b = true
				}
			}

			if !b {
				// No match by event type, break
				return
			}
		}

		// Not explicitly filtered
		return
	}
}
