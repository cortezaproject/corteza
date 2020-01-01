package options

import (
	"path"
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

		ListTimeout time.Duration `env:"CORREDOR_LIST_TIMEOUT"`

		// Allow scripts to have runner explicitly defined
		RunAsEnabled bool `env:"CORREDOR_RUN_AS_ENABLED"`

		TlsCertEnabled bool   `env:"CORREDOR_CLIENT_CERTIFICATES_ENABLED"`
		TlsCertPath    string `env:"CORREDOR_CLIENT_CERTIFICATES_PATH"`
		TlsCertCA      string `env:"CORREDOR_CLIENT_CERTIFICATES_CA"`
		TlsCertPrivate string `env:"CORREDOR_CLIENT_CERTIFICATES_PUBLIC"`
		TlsCertPublic  string `env:"CORREDOR_CLIENT_CERTIFICATES_PRIVATE"`
		TlsServerName  string `env:"CORREDOR_CLIENT_CERTIFICATES_SERVER_NAME"`
	}
)

func Corredor() (o *CorredorOpt) {
	o = &CorredorOpt{
		Enabled:            true,
		RunAsEnabled:       true,
		Addr:               "corredor:80",
		MaxBackoffDelay:    time.Minute,
		DefaultExecTimeout: time.Minute,
		ListTimeout:        time.Second * 2,
		Log:                false,

		TlsCertEnabled: true,
		TlsCertPath:    "/corredor-certificates",
		TlsCertCA:      "ca.crt",
		TlsCertPublic:  "public.crt",
		TlsCertPrivate: "private.key",
	}

	fill(o, "")

	o.TlsCertCA = path.Join(o.TlsCertPath, o.TlsCertCA)
	o.TlsCertPrivate = path.Join(o.TlsCertPath, o.TlsCertPrivate)
	o.TlsCertPublic = path.Join(o.TlsCertPath, o.TlsCertPublic)

	return
}
