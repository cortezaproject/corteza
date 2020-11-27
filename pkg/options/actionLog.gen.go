package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/actionLog.yaml

type (
	ActionLogOpt struct {
		Enabled bool `env:"ACTIONLOG_ENABLED"`
		Debug   bool `env:"ACTIONLOG_DEBUG"`
	}
)

// ActionLog initializes and returns a ActionLogOpt with default values
func ActionLog() (o *ActionLogOpt) {
	o = &ActionLogOpt{
		Enabled: true,
		Debug:   false,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *ActionLog) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
