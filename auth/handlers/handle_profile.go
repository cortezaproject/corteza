package handlers

import (
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"go.uber.org/zap"
)

func (h *AuthHandlers) profileForm(req *request.AuthReq) error {
	req.Template = TmplProfile
	if form := req.GetKV(); len(form) > 0 {
		req.Data["form"] = form
		req.SetKV(nil)
	} else {
		req.Data["form"] = map[string]string{
			"email":  req.User.Email,
			"handle": req.User.Handle,
			"name":   req.User.Name,
		}
	}

	req.Data["emailConfirmationRequired"] = !req.User.EmailConfirmed && h.Settings.EmailConfirmationRequired
	return nil
}

func (h *AuthHandlers) profileProc(req *request.AuthReq) error {
	req.RedirectTo = GetLinks().Profile
	req.SetKV(nil)

	req.User.Handle = req.Request.PostFormValue("handle")
	req.User.Name = req.Request.PostFormValue("name")

	// a little workaround to inject current user as authenticated identity into context
	// this way user service will pass us through.
	user, err := h.UserService.Update(req.Context(), req.User)

	if err == nil {
		h.storeUserToSession(req, user)
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
			"email":  req.User.Email,
			"handle": req.User.Handle,
			"name":   req.User.Name,
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
