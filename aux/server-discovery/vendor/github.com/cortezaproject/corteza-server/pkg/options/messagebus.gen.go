package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/messagebus.yaml

type (
	MessagebusOpt struct {
		Enabled    bool `env:"MESSAGEBUS_ENABLED"`
		LogEnabled bool `env:"MESSAGEBUS_LOG_ENABLED"`
	}
)

// Messagebus initializes and returns a MessagebusOpt with default values
func Messagebus() (o *MessagebusOpt) {
	o = &MessagebusOpt{
		Enabled:    true,
		LogEnabled: false,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *Messagebus) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
