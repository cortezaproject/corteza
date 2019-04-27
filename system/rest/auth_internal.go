package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/system/internal/service"
	"github.com/crusttech/crust/system/rest/request"
	"github.com/crusttech/crust/system/types"
)

var _ = errors.Wrap

type (
	authInternalValidUserResponse struct {
		JWT  string      `json:"jwt"`
		User *types.User `json:"user"`
	}

	authPasswordResetTokenExchangeResponse struct {
		Token string      `json:"token"`
		User  *types.User `json:"user"`
	}

	AuthInternal struct {
		tokenEncoder auth.TokenEncoder
		authSvc      service.AuthService
	}
)

func (AuthInternal) New() *AuthInternal {
	return &AuthInternal{
		tokenEncoder: auth.DefaultJwtHandler,
		authSvc:      service.DefaultAuth,
	}
}

func (ctrl *AuthInternal) Login(ctx context.Context, r *request.AuthInternalLogin) (interface{}, error) {
	u, err := ctrl.authSvc.InternalLogin(r.Email, r.Password)
	if err != nil {
		return nil, err
	}

	return authInternalValidUserResponse{
		JWT:  ctrl.tokenEncoder.Encode(u),
		User: u,
	}, nil
}

func (ctrl *AuthInternal) Signup(ctx context.Context, r *request.AuthInternalSignup) (interface{}, error) {
	newUser := &types.User{
		Email:    r.Email,
		Handle:   r.Handle,
		Username: r.Username,
		Name:     r.Name,
	}

	u, err := ctrl.authSvc.InternalSignUp(newUser, r.Password)
	if err != nil {
		return nil, err
	}

	if !u.EmailConfirmed {
		return nil, errors.New("user email pending confirmation")
	}

	return authInternalValidUserResponse{
		JWT:  ctrl.tokenEncoder.Encode(u),
		User: u,
	}, nil
}

func (ctrl *AuthInternal) RequestPasswordReset(ctx context.Context, r *request.AuthInternalRequestPasswordReset) (interface{}, error) {
	return true, ctrl.authSvc.SendPasswordResetToken(r.Email)
}

func (ctrl *AuthInternal) ExchangePasswordResetToken(ctx context.Context, r *request.AuthInternalExchangePasswordResetToken) (interface{}, error) {
	u, token, err := ctrl.authSvc.ExchangePasswordResetToken(r.Token)
	if err != nil {
		return nil, err
	}

	return authPasswordResetTokenExchangeResponse{
		Token: token,
		User:  u,
	}, nil
}

func (ctrl *AuthInternal) ResetPassword(ctx context.Context, r *request.AuthInternalResetPassword) (interface{}, error) {
	var u, err = ctrl.authSvc.ValidatePasswordResetToken(r.Token)
	if err != nil {
		return nil, err
	}

	err = ctrl.authSvc.SetPassword(u.ID, r.Password)
	if err != nil {
		return nil, err
	}

	return authInternalValidUserResponse{
		JWT:  ctrl.tokenEncoder.Encode(u),
		User: u,
	}, nil
}

func (ctrl *AuthInternal) ConfirmEmail(ctx context.Context, r *request.AuthInternalConfirmEmail) (interface{}, error) {
	var u, err = ctrl.authSvc.ValidateEmailConfirmationToken(r.Token)
	if err != nil {
		return nil, err
	}

	return authInternalValidUserResponse{
		JWT:  ctrl.tokenEncoder.Encode(u),
		User: u,
	}, nil
}

func (ctrl *AuthInternal) ChangePassword(ctx context.Context, r *request.AuthInternalChangePassword) (interface{}, error) {
	var identity = auth.GetIdentityFromContext(ctx)

	if !identity.Valid() {
		return nil, errors.New("invalid user (not authenticated)")
	}

	err := ctrl.authSvc.ChangePassword(identity.Identity(), r.OldPassword, r.NewPassword)
	if err != nil {
		return nil, err
	} else {
		return true, nil
	}
}
