package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/payload"
	"github.com/cortezaproject/corteza-server/internal/payload/outgoing"
	"github.com/cortezaproject/corteza-server/system/internal/service"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/types"
)

var _ = errors.Wrap

type (
	authInternalValidUserResponse struct {
		JWT  string         `json:"jwt"`
		User *outgoing.User `json:"user"`
	}

	authPasswordResetTokenExchangeResponse struct {
		Token string         `json:"token"`
		User  *outgoing.User `json:"user"`
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
	var svc = ctrl.authSvc.With(ctx)
	u, err := svc.InternalLogin(r.Email, r.Password)
	if err != nil {
		return nil, err
	}

	return ctrl.authInternalValidUserResponse(svc, u)
}

func (ctrl *AuthInternal) Signup(ctx context.Context, r *request.AuthInternalSignup) (interface{}, error) {
	var svc = ctrl.authSvc.With(ctx)

	newUser := &types.User{
		Email:    r.Email,
		Handle:   r.Handle,
		Username: r.Username,
		Name:     r.Name,
	}

	u, err := svc.InternalSignUp(newUser, r.Password)
	if err != nil {
		return nil, err
	}

	if !u.EmailConfirmed {
		// When email is not confirmed, do not send back JWT
		return authInternalValidUserResponse{User: payload.User(u)}, nil
	}

	if err = svc.LoadRoleMemberships(u); err != nil {
		return nil, err
	}

	return ctrl.authInternalValidUserResponse(svc, u)
}

func (ctrl *AuthInternal) RequestPasswordReset(ctx context.Context, r *request.AuthInternalRequestPasswordReset) (interface{}, error) {
	return true, ctrl.authSvc.With(ctx).SendPasswordResetToken(r.Email)
}

func (ctrl *AuthInternal) ExchangePasswordResetToken(ctx context.Context, r *request.AuthInternalExchangePasswordResetToken) (interface{}, error) {
	var svc = ctrl.authSvc.With(ctx)
	u, token, err := svc.ExchangePasswordResetToken(r.Token)
	if err != nil {
		return nil, err
	}

	return authPasswordResetTokenExchangeResponse{
		Token: token,
		User:  payload.User(u),
	}, nil
}

func (ctrl *AuthInternal) ResetPassword(ctx context.Context, r *request.AuthInternalResetPassword) (interface{}, error) {
	var svc = ctrl.authSvc.With(ctx)
	var u, err = svc.ValidatePasswordResetToken(r.Token)
	if err != nil {
		return nil, err
	}

	err = svc.SetPassword(u.ID, r.Password)
	if err != nil {
		return nil, err
	}

	return ctrl.authInternalValidUserResponse(svc, u)
}

func (ctrl *AuthInternal) ConfirmEmail(ctx context.Context, r *request.AuthInternalConfirmEmail) (interface{}, error) {
	var svc = ctrl.authSvc.With(ctx)
	var u, err = svc.ValidateEmailConfirmationToken(r.Token)
	if err != nil {
		return nil, err
	}

	return ctrl.authInternalValidUserResponse(svc, u)
}

func (ctrl *AuthInternal) ChangePassword(ctx context.Context, r *request.AuthInternalChangePassword) (interface{}, error) {
	var svc = ctrl.authSvc.With(ctx)
	var identity = auth.GetIdentityFromContext(ctx)

	if !identity.Valid() {
		return nil, errors.New("invalid user (not authenticated)")
	}

	err := svc.ChangePassword(identity.Identity(), r.OldPassword, r.NewPassword)
	if err != nil {
		return nil, err
	} else {
		return true, nil
	}
}

func (ctrl AuthInternal) authInternalValidUserResponse(svc interface{ LoadRoleMemberships(*types.User) error }, u *types.User) (*authInternalValidUserResponse, error) {
	if err := svc.LoadRoleMemberships(u); err != nil {
		return nil, err
	}

	return &authInternalValidUserResponse{
		JWT:  ctrl.tokenEncoder.Encode(u),
		User: payload.User(u),
	}, nil
}
