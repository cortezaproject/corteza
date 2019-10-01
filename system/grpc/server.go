package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/system/proto"
	"github.com/cortezaproject/corteza-server/system/service"
)

// @todo when we extend gRPC-server capabilities to compose & messaging
//       this needs to be refactored and generalized
func NewServer() *grpc.Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(authCheck(auth.DefaultJwtHandler)),
	)

	proto.RegisterUsersServer(s, NewUserService(
		service.DefaultUser,
		service.DefaultAuth,
		auth.DefaultJwtHandler,
		service.DefaultAccessControl,
	))

	proto.RegisterRolesServer(s, NewRoleService(
		service.DefaultRole,
	))

	return s
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
