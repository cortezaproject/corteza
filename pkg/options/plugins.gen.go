package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/plugins.yaml

type (
	PluginsOpt struct {
		Enabled bool   `env:"PLUGINS_ENABLED"`
		Paths   string `env:"PLUGINS_PATHS"`
	}
)

// Plugins initializes and returns a PluginsOpt with default values
func Plugins() (o *PluginsOpt) {
	o = &PluginsOpt{
		Enabled: true,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *Plugins) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
