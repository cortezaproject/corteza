package options

import (
	"time"
)

type (
	GRPCServerOpt struct {
		Network string `env:"GRPC_SERVER_NETWORK"`
		Addr    string `env:"GRPC_SERVER_ADDR"`

		ClientMaxBackoffDelay time.Duration `env:"GRPC_CLIENT_BACKOFF_DELAY"`
		ClientLog             bool          `env:"GRPC_CLIENT_LOG"`
	}
)

func GRPCServer(pfix string) (o *GRPCServerOpt) {
	o = &GRPCServerOpt{
		Network: "tcp",
		Addr:    ":50051",

		ClientMaxBackoffDelay: time.Minute,
		ClientLog:             false,
	}

	fill(o, pfix)

	return
}
