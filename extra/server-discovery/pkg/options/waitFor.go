package options

import (
	"time"
)

type (
	WaitForOpt struct {
		Delay      time.Duration `env:"WAIT_FOR"`
		StatusPage bool          `env:"WAIT_FOR_STATUS_PAGE"`
	}
)

// WaitFor initializes and returns a WaitForOpt with default values
func WaitFor() (o *WaitForOpt) {
	o = &WaitForOpt{
		Delay:      0,
		StatusPage: true,
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
