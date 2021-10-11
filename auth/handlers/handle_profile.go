package handlers

import (
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
)

func (h *AuthHandlers) profileForm(req *request.AuthReq) error {
	req.Template = TmplProfile

	var (
		preferredLanguage string

		u = req.AuthUser.User
	)

	if langList := h.Locale.LocalizedList(req.Context()); len(langList) > 0 {
		req.Data["languages"] = langList

		preferredLanguage = langList[0].Tag.String()
		if u.Meta != nil && u.Meta.PreferredLanguage != "" {
			preferredLanguage = u.Meta.PreferredLanguage
		}
	}

	if form := req.PopKV(); len(form) > 0 {
		req.Data["form"] = form
	} else {
		req.Data["form"] = map[string]string{
			"email":             u.Email,
			"handle":            u.Handle,
			"name":              u.Name,
			"preferredLanguage": preferredLanguage,
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

	if pl := req.Request.PostFormValue("preferredLanguage"); pl != "" {
		if u.Meta == nil {
			u.Meta = &types.UserMeta{}
		}

		u.Meta.PreferredLanguage = pl
	}

	// a little workaround to inject current user as authenticated identity into context
	// this way user service will pass us through.
	user, err := h.UserService.Update(req.Context(), u)

	if err == nil {
		err = h.AuthService.LoadRoleMemberships(req.Context(), user)
	}

	if err == nil {
		req.AuthUser.User = user
		req.AuthUser.Save(req.Session)

		t := translator(req, "auth")
		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "primary",
			Text: t("profile.alerts.profile-updated"),
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

		t := translator(req, "auth")
		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "danger",
			Text: t("profile.alerts.profile-update-fail"),
		})

		h.Log.Warn("handled error", zap.Error(err))
		return nil

	default:
		h.Log.Error("unhandled error", zap.Error(err))
		return err
	}
}
