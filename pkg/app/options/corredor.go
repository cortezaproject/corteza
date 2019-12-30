package options

import (
	"time"
)

type (
	CorredorOpt struct {
		Enabled bool `env:"CORREDOR_ENABLED"`

		// Also used by corredor service to configure gRPC server
		Addr string `env:"CORREDOR_ADDR"`

		// Also used by corredor service to enable logging
		Log bool `env:"CORREDOR_LOG_ENABLED"`

		MaxBackoffDelay time.Duration `env:"CORREDOR_MAX_BACKOFF_DELAY"`

		DefaultExecTimeout time.Duration `env:"CORREDOR_DEFAULT_EXEC_TIMEOUT"`

		// Allow scripts to have runner explicitly defined
		RunAsEnabled bool `env:"CORREDOR_RUN_AS_ENABLED"`
	}
)

func Corredor(pfix string) (o *CorredorOpt) {
	o = &CorredorOpt{
		Enabled:            true,
		RunAsEnabled:       true,
		Addr:               "corredor:80",
		MaxBackoffDelay:    time.Minute,
		DefaultExecTimeout: time.Minute,
		Log:                false,
	}

	fill(o, pfix)

	return
}
