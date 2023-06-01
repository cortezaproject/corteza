package service

import (
	"context"
	"fmt"
	htpl "html/template"
	"io/ioutil"
	"net/url"

	intAuth "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/logger"
	"github.com/cortezaproject/corteza/server/pkg/mail"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/cortezaproject/corteza/server/system/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gomail "gopkg.in/mail.v2"
)

type (
	authNotification struct {
		logger   *zap.Logger
		settings *types.AppSettings
		ts       TemplateService
		opt      options.AuthOpt
	}

	AuthNotificationService interface {
		EmailOTP(ctx context.Context, emailAddress string, otp string) error
		EmailConfirmation(ctx context.Context, emailAddress string, url string) error
		PasswordReset(ctx context.Context, emailAddress string, url string) error
		PasswordCreate(url string) (string, error)
		InviteEmail(ctx context.Context, emailAddress string, token string) error
	}
)

func AuthNotification(s *types.AppSettings, ts TemplateService, opt options.AuthOpt) *authNotification {
	return &authNotification{
		logger:   DefaultLogger.Named("auth-notification"),
		settings: s,
		ts:       ts,
		opt:      opt,
	}
}

func (svc authNotification) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}

func (svc authNotification) EmailOTP(ctx context.Context, emailAddress string, code string) error {
	return svc.send(ctx, "auth_email_mfa", emailAddress, map[string]interface{}{
		"Code": code,
	})
}

func (svc authNotification) EmailConfirmation(ctx context.Context, emailAddress string, token string) error {
	return svc.send(ctx, "auth_email_confirm", emailAddress, map[string]interface{}{
		"URL": fmt.Sprintf("%s/confirm-email?token=%s", svc.opt.BaseURL, url.QueryEscape(token)),
	})
}

func (svc authNotification) PasswordReset(ctx context.Context, emailAddress string, token string) error {
	return svc.send(ctx, "auth_email_password_reset", emailAddress, map[string]interface{}{
		"URL": fmt.Sprintf("%s/reset-password?token=%s", svc.opt.BaseURL, url.QueryEscape(token)),
	})
}

func (svc authNotification) PasswordCreate(token string) (string, error) {
	return fmt.Sprintf("%s/create-password?token=%s", svc.opt.BaseURL, url.QueryEscape(token)), nil
}

func (svc authNotification) InviteEmail(ctx context.Context, emailAddress string, token string) error {
	return svc.send(ctx, "auth_email_user_invite", emailAddress, map[string]interface{}{
		"URL": fmt.Sprintf("%s/accept-invite?token=%s", svc.opt.BaseURL, url.QueryEscape(token)),
	})
}

func (svc authNotification) newMail() *gomail.Message {
	return mail.New()
}

func (svc authNotification) send(ctx context.Context, name, sendTo string, payload map[string]interface{}) error {
	var (
		err error
		tmp []byte
		ntf = svc.newMail()
		hdl string

		// context with service user
		// we need this for retrieving & rendering email templates
		suCtx = intAuth.SetIdentityToContext(ctx, intAuth.ServiceUser())
	)

	// Fetch parts
	hdl = name + "_subject"
	st, err := svc.ts.FindByHandle(suCtx, hdl)
	if err != nil {
		return fmt.Errorf("cannot generate auth email with template %s: %w", hdl, err)
	}

	hdl = name + "_content"
	ct, err := svc.ts.FindByHandle(suCtx, hdl)
	if err != nil {
		return fmt.Errorf("cannot generate auth email with template %s: %w", hdl, err)
	}

	// Prepare payload
	payload["Logo"] = htpl.URL(svc.settings.General.Mail.Logo)
	payload["BaseURL"] = svc.opt.BaseURL
	payload["EmailAddress"] = sendTo

	// Render document
	subject, err := svc.ts.Render(suCtx, st.ID, "text/plain", payload, nil)
	if err != nil {
		return err
	}

	content, err := svc.ts.Render(suCtx, ct.ID, "text/plain", payload, nil)
	if err != nil {
		return err
	}

	tmp, err = ioutil.ReadAll(subject)
	if err != nil {
		return err
	}
	ntf.SetAddressHeader("To", sendTo, "")
	ntf.SetHeader("Subject", string(tmp))
	tmp, err = ioutil.ReadAll(content)
	if err != nil {
		return err
	}
	ntf.SetBody("text/html", string(tmp))

	err = mail.Send(ntf)

	if err != nil {
		svc.log(ctx).Error(
			"auth notification send failed",
			zap.String("name", name),
			zap.String("email", sendTo),
			zap.Error(err),
		)

		return fmt.Errorf("could not send email, contact your administrator")
	}

	svc.log(ctx).Debug(
		"auth notification sent",
		zap.String("name", name),
		zap.String("email", sendTo),
	)

	return nil
}
