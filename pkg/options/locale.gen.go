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
		Languages                   string `env:"LOCALE_LANGUAGES"`
		Path                        string `env:"LOCALE_PATH"`
		QueryStringParam            string `env:"LOCALE_QUERY_STRING_PARAM"`
		ResourceTranslationsEnabled bool   `env:"LOCALE_RESOURCE_TRANSLATIONS_ENABLED"`
		Log                         bool   `env:"LOCALE_LOG"`
		DevelopmentMode             bool   `env:"LOCALE_DEVELOPMENT_MODE"`
	}
)

// Locale initializes and returns a LocaleOpt with default values
func Locale() (o *LocaleOpt) {
	o = &LocaleOpt{
		Languages:        "en",
		QueryStringParam: "lng",
	}

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
