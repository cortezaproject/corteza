package handlers

import (
	"fmt"
	"github.com/cortezaproject/corteza/server/auth/request"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
	"go.uber.org/zap"
	"strings"
)

func (h *AuthHandlers) profileForm(req *request.AuthReq) (err error) {
	req.Template = TmplProfile

	var (
		preferredLanguage string
		avatarUrl         string
	)

	u, err := h.UserService.FindByAny(req.Context(), req.AuthUser.User.ID)

	if err != nil {
		h.Log.Error("find user error", zap.Error(err))
		return err
	}

	if langList := h.Locale.LocalizedList(req.Context()); len(langList) > 0 {
		req.Data["languages"] = langList

		preferredLanguage = langList[0].Tag.String()
		if u.Meta != nil && u.Meta.PreferredLanguage != "" {
			preferredLanguage = u.Meta.PreferredLanguage
		}
	}

	avatarUrl = fmt.Sprintf("%s/system/attachment/avatar/%d/original/%s", getApiFullPath(), u.Meta.AvatarID, types.AttachmentKindAvatar)

	if form := req.PopKV(); len(form) > 0 {
		req.Data["form"] = form
	} else {
		req.Data["form"] = map[string]string{
			"email":             u.Email,
			"handle":            u.Handle,
			"name":              u.Name,
			"preferredLanguage": preferredLanguage,
			"avatarUrl":         avatarUrl,
			"initialTextColor":  u.Meta.AvatarColor,
			"initialBgColor":    u.Meta.AvatarBgColor,
		}
	}

	req.Data["emailConfirmationRequired"] = !u.EmailConfirmed && h.Settings.EmailConfirmationRequired
	req.Data["avatarEnabled"] = h.Settings.ProfileAvatarEnabled

	if h.Settings.ProfileAvatarEnabled {
		req.Data["isAvatar"] = false

		if u.Meta.AvatarKind == types.AttachmentKindAvatar {
			req.Data["isAvatar"] = true
		}
	}

	return nil
}

func (h *AuthHandlers) profileProc(req *request.AuthReq) error {
	req.RedirectTo = GetLinks().Profile
	req.SetKV(nil)

	var (
		t = translator(req, "auth")
	)

	u, err := h.UserService.FindByAny(req.Context(), req.AuthUser.User.ID)
	if err != nil {
		h.Log.Error("find user error", zap.Error(err))
		return err
	}

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

	// process avatar initial generation
	bgColor := req.Request.PostFormValue("initial-bg")
	initialColor := req.Request.PostFormValue("initial-color")
	if bgColor != u.Meta.AvatarBgColor || initialColor != u.Meta.AvatarColor {
		err = h.UserService.GenerateAvatar(
			req.Context(),
			u.ID,
			bgColor,
			initialColor,
		)
	}

	// process avatar upload
	_, header, _ := req.Request.FormFile("avatar")
	if header != nil {
		err = h.UserService.UploadAvatar(
			req.Context(),
			u.ID,
			header,
		)
	}

	// Process avatar delete
	if req.Request.PostFormValue("avatar-delete") == "avatar-delete" {
		if err = h.UserService.DeleteAvatar(req.Context(), u.ID); err != nil {
			return err
		}
	}

	if err == nil {
		err = h.AuthService.LoadRoleMemberships(req.Context(), user)
	}

	if err == nil {
		req.AuthUser.User = user
		req.AuthUser.Save(req.Session)

		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "primary",
			Text: t("profile.alerts.profile-updated"),
		})

		req.RedirectTo = GetLinks().Profile
		return nil
	}

	avatarUrl := fmt.Sprintf("%s/system/attachment/avatar/%d/original/%s", getApiFullPath(), u.Meta.AvatarID, types.AttachmentKindAvatar)
	switch {
	case
		service.UserErrInvalidID().Is(err),
		service.UserErrInvalidHandle().Is(err),
		service.UserErrInvalidEmail().Is(err),
		service.UserErrHandleNotUnique().Is(err),
		service.UserErrNotAllowedToUpdate().Is(err),
		service.AttachmentErrInvalidAvatarFileType().Is(err),
		service.AttachmentErrInvalidAvatarFileSize().Is(err),
		service.AttachmentErrInvalidAvatarGenerateFontFile().Is(err):
		req.SetKV(map[string]string{
			"error":            err.Error(),
			"email":            u.Email,
			"handle":           u.Handle,
			"name":             u.Name,
			"avatarUrl":        avatarUrl,
			"initialTextColor": u.Meta.AvatarColor,
			"initialBgColor":   u.Meta.AvatarBgColor,
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

func getApiFullPath() (path string) {
	return fmt.Sprintf("/%s", strings.TrimPrefix(options.CleanBase(options.HttpServer().BaseUrl, options.HttpServer().ApiBaseUrl), "/"))
}
