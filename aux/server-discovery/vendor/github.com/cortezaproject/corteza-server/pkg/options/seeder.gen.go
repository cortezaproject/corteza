package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/seeder.yaml

type (
	SeederOpt struct {
		LogEnabled bool `env:"SEEDER_LOG_ENABLED"`
	}
)

// Seeder initializes and returns a SeederOpt with default values
func Seeder() (o *SeederOpt) {
	o = &SeederOpt{}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *Seeder) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
