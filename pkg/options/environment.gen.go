package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/environment.yaml

type (
	EnvironmentOpt struct {
		Environment string `env:"ENVIRONMENT"`
	}
)

// Environment initializes and returns a EnvironmentOpt with default values
func Environment() (o *EnvironmentOpt) {
	o = &EnvironmentOpt{
		Environment: "production",
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *Environment) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
