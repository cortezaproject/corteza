package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/limit.yaml

type (
	LimitOpt struct {
		SystemUsers int `env:"LIMIT_SYSTEM_USERS"`
	}
)

// Limit initializes and returns a LimitOpt with default values
func Limit() (o *LimitOpt) {
	o = &LimitOpt{}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *Limit) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
