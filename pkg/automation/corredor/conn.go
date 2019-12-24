package corredor

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/cortezaproject/corteza-server/pkg/app/options"
)

// Corredor standard connector to Corredor service via gRPC
func NewConnection(ctx context.Context, opt options.CorredorOpt, logger *zap.Logger) (c *grpc.ClientConn, err error) {
	if !opt.Enabled {
		// Do not connect when script runner is not enabled
		return
	}

	if opt.Log {
		// Send logs to zap
		//
		// waiting for https://github.com/uber-go/zap/pull/538
		grpclog.SetLogger(zapgrpc.NewLogger(logger.Named("grpc")))
	}

	var dopts = []grpc.DialOption{
		// @todo insecure?
		grpc.WithInsecure(),
	}

	if opt.MaxBackoffDelay > 0 {
		dopts = append(dopts, grpc.WithBackoffMaxDelay(opt.MaxBackoffDelay))
	}

	return grpc.DialContext(ctx, opt.Addr, dopts...)
}
