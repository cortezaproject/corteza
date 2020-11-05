package rest

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/pkg/payload/outgoing"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	authInternalValidUserResponse struct {
		JWT  string           `json:"jwt"`
		User *authUserPayload `json:"user"`
	}

	authPasswordResetTokenExchangeResponse struct {
		Token string         `json:"token"`
		User  *outgoing.User `json:"user"`
	}

	AuthInternal struct {
		tokenEncoder auth.TokenEncoder
		authSvc      authInternalAuthService
	}

	authInternalAuthService interface {
		InternalSignUp(ctx context.Context, input *types.User, password string) (*types.User, error)
		InternalLogin(ctx context.Context, email string, password string) (*types.User, error)
		SetPassword(ctx context.Context, userID uint64, AuthActionPassword string) error
		ChangePassword(ctx context.Context, userID uint64, oldPassword, AuthActionPassword string) error
		LoadRoleMemberships(ctx context.Context, user *types.User) error
		ValidateEmailConfirmationToken(ctx context.Context, token string) (user *types.User, err error)
		ExchangePasswordResetToken(ctx context.Context, token string) (user *types.User, exchangedToken string, err error)
		ValidatePasswordResetToken(ctx context.Context, token string) (user *types.User, err error)
		SendPasswordResetToken(ctx context.Context, email string) (err error)
	}
)

func (AuthInternal) New() *AuthInternal {
	return &AuthInternal{
		tokenEncoder: auth.DefaultJwtHandler,
		authSvc:      service.DefaultAuth,
	}
}

func (ctrl *AuthInternal) Login(ctx context.Context, r *request.AuthInternalLogin) (interface{}, error) {
	u, err := ctrl.authSvc.InternalLogin(ctx, r.Email, r.Password)
	if err != nil {
		return nil, err
	}

	return ctrl.authInternalValidUserResponse(ctx, u)
}

func (ctrl *AuthInternal) Signup(ctx context.Context, r *request.AuthInternalSignup) (interface{}, error) {
	newUser := &types.User{
		Email:    r.Email,
		Handle:   r.Handle,
		Username: r.Username,
		Name:     r.Name,
	}

	u, err := ctrl.authSvc.InternalSignUp(ctx, newUser, r.Password)
	if err != nil {
		return nil, err
	}

	if !u.EmailConfirmed {
		// When email is not confirmed, do not send back JWT
		return authInternalValidUserResponse{User: &authUserPayload{User: payload.User(u)}}, nil
	}

	if err = ctrl.authSvc.LoadRoleMemberships(ctx, u); err != nil {
		return nil, err
	}

	return ctrl.authInternalValidUserResponse(ctx, u)
}

func (ctrl *AuthInternal) RequestPasswordReset(ctx context.Context, r *request.AuthInternalRequestPasswordReset) (interface{}, error) {
	return true, ctrl.authSvc.SendPasswordResetToken(ctx, r.Email)
}

func (ctrl *AuthInternal) ExchangePasswordResetToken(ctx context.Context, r *request.AuthInternalExchangePasswordResetToken) (interface{}, error) {
	u, token, err := ctrl.authSvc.ExchangePasswordResetToken(ctx, r.Token)
	if err != nil {
		return nil, err
	}

	return authPasswordResetTokenExchangeResponse{
		Token: token,
		User:  payload.User(u),
	}, nil
}

func (ctrl *AuthInternal) ResetPassword(ctx context.Context, r *request.AuthInternalResetPassword) (interface{}, error) {
	var u, err = ctrl.authSvc.ValidatePasswordResetToken(ctx, r.Token)
	if err != nil {
		return nil, err
	}

	err = ctrl.authSvc.SetPassword(ctx, u.ID, r.Password)
	if err != nil {
		return nil, err
	}

	return ctrl.authInternalValidUserResponse(ctx, u)
}

func (ctrl *AuthInternal) ConfirmEmail(ctx context.Context, r *request.AuthInternalConfirmEmail) (interface{}, error) {
	var u, err = ctrl.authSvc.ValidateEmailConfirmationToken(ctx, r.Token)
	if err != nil {
		return nil, err
	}

	return ctrl.authInternalValidUserResponse(ctx, u)
}

func (ctrl *AuthInternal) ChangePassword(ctx context.Context, r *request.AuthInternalChangePassword) (interface{}, error) {
	var identity = auth.GetIdentityFromContext(ctx)

	if !identity.Valid() {
		return nil, errors.New("invalid user (not authenticated)")
	}

	err := ctrl.authSvc.ChangePassword(ctx, identity.Identity(), r.OldPassword, r.NewPassword)
	if err != nil {
		return nil, err
	} else {
		return true, nil
	}
}

func (ctrl AuthInternal) authInternalValidUserResponse(ctx context.Context, u *types.User) (*authInternalValidUserResponse, error) {
	if err := ctrl.authSvc.LoadRoleMemberships(ctx, u); err != nil {
		return nil, err
	}

	return &authInternalValidUserResponse{
		JWT: ctrl.tokenEncoder.Encode(u),
		User: &authUserPayload{
			User:  payload.User(u),
			Roles: payload.Uint64stoa(u.Roles()),
		},
	}, nil
}
