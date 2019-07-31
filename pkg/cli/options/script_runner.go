package options

import (
	"time"
)

type (
	ScriptRunnerOpt struct {
		Enabled         bool          `env:"SCRIPT_RUNNER_ENABLED"`
		Addr            string        `env:"SCRIPT_RUNNER_ADDR"`
		MaxBackoffDelay time.Duration `env:"SCRIPT_RUNNER_MAX_BACKOFF_DELAY"`
		Log             bool          `env:"SCRIPT_RUNNER_LOG"`
	}
)

func ScriptRunner(pfix string) (o *ScriptRunnerOpt) {
	o = &ScriptRunnerOpt{
		Enabled:         false,
		Addr:            "corredor:80",
		MaxBackoffDelay: time.Minute,
		Log:             false,
	}

	fill(o, pfix)

	return
}
