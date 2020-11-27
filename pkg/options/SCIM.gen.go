package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/SCIM.yaml

type (
	SCIMOpt struct {
		Enabled bool   `env:"SCIM_ENABLED"`
		BaseURL string `env:"SCIM_BASE_URL"`
		Secret  string `env:"SCIM_SECRET"`
	}
)

// SCIM initializes and returns a SCIMOpt with default values
func SCIM() (o *SCIMOpt) {
	o = &SCIMOpt{
		BaseURL: "/scim",
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *SCIM) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
