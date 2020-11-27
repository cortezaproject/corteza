package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/waitFor.yaml

import (
	"time"
)

type (
	WaitForOpt struct {
		Delay                 time.Duration `env:"WAIT_FOR"`
		StatusPage            bool          `env:"WAIT_FOR_STATUS_PAGE"`
		Services              string        `env:"WAIT_FOR_SERVICES"`
		ServicesTimeout       time.Duration `env:"WAIT_FOR_SERVICES_TIMEOUT"`
		ServicesProbeTimeout  time.Duration `env:"WAIT_FOR_SERVICES_PROBE_TIMEOUT"`
		ServicesProbeInterval time.Duration `env:"WAIT_FOR_SERVICES_PROBE_INTERVAL"`
	}
)

// WaitFor initializes and returns a WaitForOpt with default values
func WaitFor() (o *WaitForOpt) {
	o = &WaitForOpt{
		Delay:                 0,
		StatusPage:            true,
		ServicesTimeout:       time.Minute,
		ServicesProbeTimeout:  time.Second * 30,
		ServicesProbeInterval: time.Second * 5,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *WaitFor) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
