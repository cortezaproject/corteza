package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/log.yaml

type (
	LogOpt struct {
		Debug           bool   `env:"LOG_DEBUG"`
		Level           string `env:"LOG_LEVEL"`
		Filter          string `env:"LOG_FILTER"`
		IncludeCaller   bool   `env:"LOG_INCLUDE_CALLER"`
		StacktraceLevel string `env:"LOG_STACKTRACE_LEVEL"`
	}
)

// Log initializes and returns a LogOpt with default values
func Log() (o *LogOpt) {
	o = &LogOpt{
		Level:           "warn",
		IncludeCaller:   false,
		StacktraceLevel: "dpanic",
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *Log) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
