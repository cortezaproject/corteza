package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/monitor.yaml

import (
	"time"
)

type (
	MonitorOpt struct {
		Interval time.Duration `env:"MONITOR_INTERVAL"`
	}
)

// Monitor initializes and returns a MonitorOpt with default values
func Monitor() (o *MonitorOpt) {
	o = &MonitorOpt{
		Interval: 300 * time.Second,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *Monitor) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
