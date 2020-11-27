package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/upgrade.yaml

type (
	UpgradeOpt struct {
		Debug  bool `env:"UPGRADE_DEBUG"`
		Always bool `env:"UPGRADE_ALWAYS"`
	}
)

// Upgrade initializes and returns a UpgradeOpt with default values
func Upgrade() (o *UpgradeOpt) {
	o = &UpgradeOpt{
		Debug:  false,
		Always: true,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *Upgrade) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
