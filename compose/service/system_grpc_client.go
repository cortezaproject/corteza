package service

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/cortezaproject/corteza-server/pkg/app/options"
)

// Connects to system gRPC server
func NewSystemGRPCClient(ctx context.Context, opt options.GRPCServerOpt, logger *zap.Logger) (c *grpc.ClientConn, err error) {
	if opt.ClientLog {
		// Send logs to zap
		//
		// waiting for https://github.com/uber-go/zap/pull/538
		grpclog.SetLogger(zapgrpc.NewLogger(logger.Named("grpc-client-system")))
	}

	var dopts = []grpc.DialOption{
		// @todo insecure?
		grpc.WithInsecure(),
	}

	if opt.ClientMaxBackoffDelay > 0 {
		dopts = append(dopts, grpc.WithBackoffMaxDelay(opt.ClientMaxBackoffDelay))
	}

	return grpc.DialContext(ctx, opt.Addr, dopts...)
}
