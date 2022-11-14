package options

import "strings"

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

func (e EnvironmentOpt) IsDevelopment() bool {
	return strings.HasPrefix(e.Environment, "dev")
}

func (e EnvironmentOpt) IsTest() bool {
	return strings.HasPrefix(e.Environment, "test")
}

func (e EnvironmentOpt) IsProduction() bool {
	return !e.IsDevelopment() &&
		!e.IsTest()
}
