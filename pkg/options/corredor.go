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

		MaxBackoffDelay time.Duration `env:"CORREDOR_MAX_BACKOFF_DELAY"`

		MaxReceiveMessageSize int `env:"CORREDOR_MAX_RECEIVE_MESSAGE_SIZE"`

		DefaultExecTimeout time.Duration `env:"CORREDOR_DEFAULT_EXEC_TIMEOUT"`

		ListTimeout time.Duration `env:"CORREDOR_LIST_TIMEOUT"`
		ListRefresh time.Duration `env:"CORREDOR_LIST_REFRESH"`

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
		Enabled:               true,
		RunAsEnabled:          true,
		Addr:                  "localhost:50051",
		MaxBackoffDelay:       time.Minute,
		MaxReceiveMessageSize: 2 << 23, // 16MB
		DefaultExecTimeout:    time.Minute,
		ListTimeout:           time.Second * 2,
		ListRefresh:           time.Second * 5,

		TlsCertEnabled: false,
		TlsCertPath:    "/certs/corredor/client",
		TlsCertCA:      "ca.crt",
		TlsCertPublic:  "public.crt",
		TlsCertPrivate: "private.key",
	}

	fill(o)

	o.TlsCertCA = path.Join(o.TlsCertPath, o.TlsCertCA)
	o.TlsCertPrivate = path.Join(o.TlsCertPath, o.TlsCertPrivate)
	o.TlsCertPublic = path.Join(o.TlsCertPath, o.TlsCertPublic)

	return
}
