package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/eventbus.yaml

import (
	"time"
)

type (
	EventbusOpt struct {
		SchedulerEnabled  bool          `env:"EVENTBUS_SCHEDULER_ENABLED"`
		SchedulerInterval time.Duration `env:"EVENTBUS_SCHEDULER_INTERVAL"`
	}
)

// Eventbus initializes and returns a EventbusOpt with default values
func Eventbus() (o *EventbusOpt) {
	o = &EventbusOpt{
		SchedulerEnabled:  true,
		SchedulerInterval: time.Minute,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *Eventbus) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
