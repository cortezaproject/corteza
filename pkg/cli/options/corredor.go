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

		// @todo do autodiscovery & prefill these with values we know
		ApiBaseURLSystem    string `env:"CORREDOR_API_BASE_URL_SYSTEM"`
		ApiBaseURLMessaging string `env:"CORREDOR_API_BASE_URL_MESSAGING"`
		ApiBaseURLCompose   string `env:"CORREDOR_API_BASE_URL_COMPOSE"`
	}
)

func Corredor(pfix string) (o *CorredorOpt) {
	o = &CorredorOpt{
		Enabled:         false,
		Addr:            "corredor:80",
		MaxBackoffDelay: time.Minute,
		Log:             false,
	}

	fill(o, pfix)

	return
}
