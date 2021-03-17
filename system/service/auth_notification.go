package service

import (
	"context"
	"fmt"
	htpl "html/template"
	"io/ioutil"
	"net/url"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/mail"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/system/types"
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

func (svc authNotification) newMail() *gomail.Message {
	var (
		m    = mail.New()
		addr = svc.settings.Auth.Mail.FromAddress
		name = svc.settings.Auth.Mail.FromName
	)

	if addr != "" {
		m.SetAddressHeader("From", addr, name)
	}

	return m
}

func (svc authNotification) send(ctx context.Context, name, sendTo string, payload map[string]interface{}) error {
	var (
		err error
		tmp []byte
		ntf = svc.newMail()
		hdl string
	)

	// Fetch parts
	hdl = name + "_subject"
	st, err := svc.ts.FindByHandle(ctx, hdl)
	if err != nil {
		return fmt.Errorf("cannot generate auth email with template %s: %w", hdl, err)
	}

	hdl = name + "_content"
	ct, err := svc.ts.FindByHandle(ctx, hdl)
	if err != nil {
		return fmt.Errorf("cannot generate auth email with template %s: %w", hdl, err)
	}

	// Prepare payload
	payload["Logo"] = htpl.URL(svc.settings.General.Mail.Logo)
	payload["BaseURL"] = svc.opt.BaseURL
	payload["SignatureName"] = svc.settings.Auth.Mail.FromName
	payload["SignatureEmail"] = svc.settings.Auth.Mail.FromAddress
	payload["EmailAddress"] = sendTo

	// Render document
	subject, err := svc.ts.Render(ctx, st.ID, "text/plain", payload, nil)
	if err != nil {
		return err
	}

	content, err := svc.ts.Render(ctx, ct.ID, "text/plain", payload, nil)
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

	svc.log(ctx).Debug(
		"sending auth notification",
		zap.String("name", name),
		zap.String("email", sendTo),
	)

	return mail.Send(ntf)
}
