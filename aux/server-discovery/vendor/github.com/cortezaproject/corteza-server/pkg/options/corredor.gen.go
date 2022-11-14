package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/corredor.yaml

import (
	"time"
)

type (
	CorredorOpt struct {
		Enabled               bool          `env:"CORREDOR_ENABLED"`
		Addr                  string        `env:"CORREDOR_ADDR"`
		MaxBackoffDelay       time.Duration `env:"CORREDOR_MAX_BACKOFF_DELAY"`
		MaxReceiveMessageSize int           `env:"CORREDOR_MAX_RECEIVE_MESSAGE_SIZE"`
		DefaultExecTimeout    time.Duration `env:"CORREDOR_DEFAULT_EXEC_TIMEOUT"`
		ListTimeout           time.Duration `env:"CORREDOR_LIST_TIMEOUT"`
		ListRefresh           time.Duration `env:"CORREDOR_LIST_REFRESH"`
		RunAsEnabled          bool          `env:"CORREDOR_RUN_AS_ENABLED"`
		TlsCertEnabled        bool          `env:"CORREDOR_CLIENT_CERTIFICATES_ENABLED"`
		TlsCertPath           string        `env:"CORREDOR_CLIENT_CERTIFICATES_PATH"`
		TlsCertCA             string        `env:"CORREDOR_CLIENT_CERTIFICATES_CA"`
		TlsCertPrivate        string        `env:"CORREDOR_CLIENT_CERTIFICATES_PRIVATE"`
		TlsCertPublic         string        `env:"CORREDOR_CLIENT_CERTIFICATES_PUBLIC"`
		TlsServerName         string        `env:"CORREDOR_CLIENT_CERTIFICATES_SERVER_NAME"`
	}
)

// Corredor initializes and returns a CorredorOpt with default values
func Corredor() (o *CorredorOpt) {
	o = &CorredorOpt{
		Enabled:               false,
		Addr:                  "localhost:50051",
		MaxBackoffDelay:       time.Minute,
		MaxReceiveMessageSize: 2 << 23,
		DefaultExecTimeout:    time.Minute,
		ListTimeout:           time.Second * 2,
		ListRefresh:           time.Second * 5,
		RunAsEnabled:          true,
		TlsCertEnabled:        false,
		TlsCertPath:           "/certs/corredor/client",
		TlsCertCA:             "ca.crt",
		TlsCertPrivate:        "private.key",
		TlsCertPublic:         "public.crt",
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *Corredor) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
