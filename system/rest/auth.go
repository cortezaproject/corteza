package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	tokenGenerator interface {
		Generate(ctx context.Context, i auth.Identifiable, clientID uint64, scope ...string) (token []byte, err error)
	}

	Auth struct {
		tokenHandler tokenGenerator
		settings     *types.AppSettings
		authSvc      authUserService
	}

	authUserResponse struct {
		JWT  string           `json:"jwt"`
		User *authUserPayload `json:"user"`
	}

	authUserPayload struct {
		*userPayload
		Roles []string `json:"roles"`
	}

	authUserService interface {
		Impersonate(ctx context.Context, userID uint64) (*types.User, error)
		LoadRoleMemberships(ctx context.Context, user *types.User) error
	}

	userPayload struct {
		// Channel to part (nil) for ALL channels
		ID       uint64 `json:"userID,string"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Handle   string `json:"handle"`
	}
)

func (Auth) New() *Auth {
	return &Auth{
		tokenHandler: auth.JWT(),
		settings:     service.CurrentSettings,
		authSvc:      service.DefaultAuth,
	}
}

// Impersonate implements impersonation functionality
//
// This is experimental and internals will most likely change in the future:
func (ctrl *Auth) Impersonate(ctx context.Context, r *request.AuthImpersonate) (interface{}, error) {
	u, err := ctrl.authSvc.Impersonate(ctx, r.UserID)
	if err != nil {
		return nil, err
	}

	return ctrl.makePayload(ctx, u)
}

func (ctrl *Auth) makePayload(ctx context.Context, user *types.User) (*authUserResponse, error) {
	if err := ctrl.authSvc.LoadRoleMemberships(ctx, user); err != nil {
		return nil, err
	}

	// Generate and save the token
	t, err := ctrl.tokenHandler.Generate(ctx, user, 0)
	if err != nil {
		return nil, nil
	}

	return &authUserResponse{
		JWT: string(t),
		User: &authUserPayload{
			userPayload: &userPayload{
				ID:       user.ID,
				Name:     user.Name,
				Handle:   user.Handle,
				Username: user.Username,
				Email:    user.Email,
			},
			Roles: payload.Uint64stoa(user.Roles()),
		},
	}, nil
}
