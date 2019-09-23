package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/system/proto"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	userService struct {
		ac    userServiceAccessControl
		users service.UserService
		auth  service.AuthService
		jwt   auth.TokenEncoder
	}

	userServiceAccessControl interface {
		CanGrant(ctx context.Context) bool
	}
)

func NewUserService(users service.UserService, auth service.AuthService, jwt auth.TokenEncoder, ac userServiceAccessControl) *userService {
	return &userService{
		ac:    ac,
		users: users,
		auth:  auth,
		jwt:   jwt,
	}
}

func (gs userService) MakeJWT(ctx context.Context, req *proto.MakeJWTRequest) (rsp *proto.MakeJWTResponse, err error) {
	var (
		u *types.User
	)

	if !gs.ac.CanGrant(ctx) {
		return nil, status.Error(codes.PermissionDenied, "no permissions to issue jwt for other users")
	}

	if u, err = gs.users.FindByID(req.UserID); err != nil {
		return
	}

	if err = gs.auth.LoadRoleMemberships(u); err != nil {
		return
	}

	rsp = &proto.MakeJWTResponse{
		JWT: gs.jwt.Encode(u),
	}

	return
}

func (gs userService) FindByID(ctx context.Context, req *proto.FindByIDRequest) (rsp *proto.FindByIDResponse, err error) {
	var (
		u *types.User
	)

	if u, err = gs.users.FindByID(req.UserID); err != nil {
		return
	}

	rsp = &proto.FindByIDResponse{
		User: &proto.User{
			ID:     u.ID,
			Email:  u.Email,
			Handle: u.Handle,
			Name:   u.Name,
			Kind:   string(u.Kind),
		},
	}

	return
}
