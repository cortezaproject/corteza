package handlers

import (
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
)

func (h *AuthHandlers) signupForm(req *request.AuthReq) error {
	req.Template = TmplSignup
	req.Data["form"] = req.PopKV()
	return nil
}

func (h *AuthHandlers) signupProc(req *request.AuthReq) error {
	req.RedirectTo = GetLinks().Signup
	req.SetKV(nil)

	payload := &types.User{
		Email:  req.Request.PostFormValue("email"),
		Handle: req.Request.PostFormValue("handle"),
		Name:   req.Request.PostFormValue("name"),
	}

	newUser, err := h.AuthService.InternalSignUp(
		req.Context(),
		payload,
		req.Request.PostFormValue("password"),
	)

	if err == nil {
		if newUser.EmailConfirmed {
			req.NewAlerts = append(req.NewAlerts, request.Alert{
				Type: "primary",
				Text: "Sign-up successful.",
			})

			h.Log.Info(
				"signup successful",
				zap.String("email", newUser.Email),
				zap.Uint64s("roles", newUser.Roles()),
			)
			req.RedirectTo = GetLinks().Profile

			req.AuthUser = request.NewAuthUser(h.Settings, newUser, false, 0)

			// auto-complete EmailOTP
			req.AuthUser.CompleteEmailOTP()

			req.AuthUser.Save(req.Session)

		} else {
			req.RedirectTo = GetLinks().PendingEmailConfirmation
		}

		return nil
	}

	switch {
	case service.AuthErrInternalSignupDisabledByConfig().Is(err):
		h.signupDisabledAlert(req)
		return nil
	case service.AuthErrInvalidEmailFormat().Is(err),
		service.AuthErrInvalidHandle().Is(err),
		service.AuthErrPasswordNotSecure().Is(err),
		service.AuthErrInvalidCredentials().Is(err):
		req.SetKV(map[string]string{
			"error":  err.Error(),
			"email":  payload.Email,
			"handle": payload.Handle,
			"name":   payload.Name,
		})

		h.Log.Warn("handled error", zap.Error(err))
		return nil

	default:
		h.Log.Error("unhandled error", zap.Error(err))
		return err
	}
}

func (h *AuthHandlers) pendingEmailConfirmation(req *request.AuthReq) error {
	req.Template = TmplPendingEmailConfirmation

	if _, has := req.Request.URL.Query()["resend"]; has && req.AuthUser.User != nil {
		err := h.AuthService.SendEmailAddressConfirmationToken(req.Context(), req.AuthUser.User)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *AuthHandlers) confirmEmail(req *request.AuthReq) (err error) {
	if token := req.Request.URL.Query().Get("token"); len(token) > 0 {
		var user *types.User
		user, err = h.AuthService.ValidateEmailConfirmationToken(req.Context(), token)
		if err == nil {
			// redirect back to self (but without token and with user in session
			h.Log.Debug("valid email confirmation token found, redirecting to profile")
			req.NewAlerts = append(req.NewAlerts, request.Alert{
				Type: "primary",
				Text: "Email address confirmed, you're now logged-in.",
			})

			req.RedirectTo = GetLinks().Profile

			req.AuthUser = request.NewAuthUser(h.Settings, user, false, 0)

			// auto-complete EmailOTP
			req.AuthUser.CompleteEmailOTP()

			req.AuthUser.Save(req.Session)

			return nil
		}
	}

	h.Log.Warn("invalid email confirmation token used", zap.Error(err))

	// redirect to the right page
	// not doing this here and relying on handler on subseq. request
	// will cause alerts to be removed
	if req.AuthUser == nil && req.AuthUser.User == nil {
		req.RedirectTo = GetLinks().Login
	} else {
		req.RedirectTo = GetLinks().Profile
	}

	req.NewAlerts = append(req.NewAlerts, request.Alert{
		Type: "warning",
		Text: "Invalid or expired email confirmation token, please resend confirmation request.",
	})

	return nil
}

func (h *AuthHandlers) onlyIfSignupEnabled(fn handlerFn) handlerFn {
	return func(req *request.AuthReq) error {
		if !h.Settings.SignupEnabled || !h.Settings.LocalEnabled {
			h.signupDisabledAlert(req)
			return nil
		}

		return fn(req)
	}
}

func (h *AuthHandlers) signupDisabledAlert(req *request.AuthReq) {
	req.RedirectTo = GetLinks().Login
	req.NewAlerts = append(req.NewAlerts, request.Alert{
		Type: "danger",
		Text: "Signup disabled",
	})
}
