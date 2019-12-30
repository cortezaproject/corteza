package corredor

import (
	"encoding/json"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/eventbus"
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

// pluckManualTriggers removes all manual triggers from the list of script's triggers
//
// and returns a hash map with resources from these manual triggers
func pluckManualTriggers(script *ServerScript) map[string]string {
	var (
		hash = make(map[string]string)
	)

	for i := range script.Triggers {
		// We're modifying trigger in the loop,
		// so let's make a copy we can play with
		trigger := script.Triggers[i]

		if popOnManualEventType(trigger) {
			for _, res := range trigger.Resources {
				hash[res] = trigger.RunAs
			}
		}
	}

	return hash
}

// converts trigger's constraint to eventbus' constraint options
func makeTriggerOpts(t *Trigger) (oo []eventbus.TriggerRegOp, err error) {
	if len(t.Events) == 0 {
		return nil, fmt.Errorf("can not generate trigger without at least one events")
	}
	if len(t.Resources) == 0 {
		return nil, fmt.Errorf("can not generate trigger without at least one resource")
	}

	oo = append(oo, eventbus.On(t.Events...))
	oo = append(oo, eventbus.For(t.Resources...))

	for i := range t.Constraints {
		oo = append(oo, eventbus.Constraint(
			t.Constraints[i].Name,
			t.Constraints[i].Op,
			t.Constraints[i].Value...,
		))
	}

	return
}

// encode adds entry (with json encoded value) to hash map
// used to prepare data for transmission
func encodeArguments(args map[string]string, key string, val interface{}) (err error) {
	var tmp []byte

	if tmp, err = json.Marshal(val); err != nil {
		return
	}

	args[key] = string(tmp)
	return
}

// Creates a filter fn for script filtering
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
