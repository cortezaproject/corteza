package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/locale.yaml

type (
	LocaleOpt struct {
		Path string `env:"LOCALE_PATH"`
		Log  bool   `env:"LOCALE_LOG"`
	}
)

// Locale initializes and returns a LocaleOpt with default values
func Locale() (o *LocaleOpt) {
	o = &LocaleOpt{}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *Locale) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
