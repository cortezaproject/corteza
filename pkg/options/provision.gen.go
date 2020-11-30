package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/provision.yaml

type (
	ProvisionOpt struct {
		Always bool   `env:"PROVISION_ALWAYS"`
		Path   string `env:"PROVISION_PATH"`
	}
)

// Provision initializes and returns a ProvisionOpt with default values
func Provision() (o *ProvisionOpt) {
	o = &ProvisionOpt{
		Always: true,
		Path:   "provision/*",
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *Provision) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
