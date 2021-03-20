package handlers

import (
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"go.uber.org/zap"
)

func (h *AuthHandlers) profileForm(req *request.AuthReq) error {
	req.Template = TmplProfile
	u := req.AuthUser.User

	if form := req.PopKV(); len(form) > 0 {
		req.Data["form"] = form
	} else {
		req.Data["form"] = map[string]string{
			"email":  u.Email,
			"handle": u.Handle,
			"name":   u.Name,
		}
	}

	req.Data["emailConfirmationRequired"] = !u.EmailConfirmed && h.Settings.EmailConfirmationRequired
	return nil
}

func (h *AuthHandlers) profileProc(req *request.AuthReq) error {
	req.RedirectTo = GetLinks().Profile
	req.SetKV(nil)
	u := req.AuthUser.User

	u.Handle = req.Request.PostFormValue("handle")
	u.Name = req.Request.PostFormValue("name")

	// a little workaround to inject current user as authenticated identity into context
	// this way user service will pass us through.
	user, err := h.UserService.Update(req.Context(), u)

	if err == nil {
		req.AuthUser.User = user
		req.AuthUser.Save(req.Session)

		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "primary",
			Text: "Profile successfully updated.",
		})

		req.RedirectTo = GetLinks().Profile
		return nil
	}

	switch {
	case
		service.UserErrInvalidID().Is(err),
		service.UserErrInvalidHandle().Is(err),
		service.UserErrInvalidEmail().Is(err),
		service.UserErrHandleNotUnique().Is(err),
		service.UserErrNotAllowedToUpdate().Is(err):
		req.SetKV(map[string]string{
			"error":  err.Error(),
			"email":  u.Email,
			"handle": u.Handle,
			"name":   u.Name,
		})

		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "danger",
			Text: "Could not update profile due to input errors",
		})

		h.Log.Warn("handled error", zap.Error(err))
		return nil

	default:
		h.Log.Error("unhandled error", zap.Error(err))
		return err
	}
}
