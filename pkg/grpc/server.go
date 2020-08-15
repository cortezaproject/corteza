package grpc

import (
	"context"
	"net"

	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/options"
)

type (
	server struct {
		log *zap.Logger
		opt options.GRPCServerOpt

		s *grpc.Server
	}
)

func New(log *zap.Logger, opt options.GRPCServerOpt) *server {
	return &server{
		log: log.Named("grpc-server"),
		opt: opt,

		s: grpc.NewServer(
			grpc.UnaryInterceptor(authCheck(auth.DefaultJwtHandler)),
		),
	}
}

func (srv *server) Serve(ctx context.Context) {
	ln, err := net.Listen(srv.opt.Network, srv.opt.Addr)
	if err != nil {
		srv.log.Error("could not start gRPC server", zap.Error(err))
	}

	go func() {
		select {
		case <-ctx.Done():
			srv.log.Debug("shutting down")
			srv.s.GracefulStop()
			_ = ln.Close()
		}
	}()

	srv.log.Info("Starting gRPC server", zap.String("address", srv.opt.Addr))
	err = srv.s.Serve(ln)

	if err == nil {
		err = ctx.Err()
		if err == context.Canceled {
			err = nil
		}
	}

	srv.log.Info("Server stopped", zap.Error(err))
}

func (srv *server) RegisterServices(reg func(*grpc.Server)) {
	reg(srv.s)
}

// Creates auth-checking interceptor function
//
// Interceptor expects a valid 'jwt' in meta-data
func authCheck(h auth.TokenDecoder) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if meta, ok := metadata.FromIncomingContext(ctx); !ok {
			return nil, status.Error(codes.Unauthenticated, "could not read metadata")
		} else if len(meta["jwt"]) != 1 {
			return nil, status.Error(codes.Unauthenticated, "metadata without jwt")
		} else if identity, err := h.Decode(meta["jwt"][0]); err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid jwt")
		} else if identity == nil || !identity.Valid() {
			return nil, status.Error(codes.Unauthenticated, "invalid identity")
		} else {
			// Append identity to context and procede
			ctx = auth.SetIdentityToContext(ctx, identity)
		}

		// Serve the request
		return handler(ctx, req)
	}
}
