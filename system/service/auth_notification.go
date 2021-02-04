package service

import (
	"bytes"
	"context"
	"fmt"
	htpl "html/template"

	"github.com/cortezaproject/corteza-server/system/types"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gomail "gopkg.in/mail.v2"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/mail"
)

type (
	authNotification struct {
		logger   *zap.Logger
		settings *types.AppSettings
	}

	AuthNotificationService interface {
		EmailConfirmation(ctx context.Context, lang string, emailAddress string, url string) error
		PasswordReset(ctx context.Context, lang string, emailAddress string, url string) error
	}

	authNotificationPayload struct {
		EmailAddress   string
		URL            string
		BaseURL        string
		Logo           htpl.URL
		SignatureName  string
		SignatureEmail string
		EmailHeaderEn  htpl.HTML
		EmailFooterEn  htpl.HTML
	}
)

func AuthNotification(s *types.AppSettings) AuthNotificationService {
	return &authNotification{
		logger:   DefaultLogger.Named("auth-notification"),
		settings: s,
	}
}

func (svc authNotification) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}

func (svc authNotification) EmailConfirmation(ctx context.Context, lang string, emailAddress string, token string) error {
	return svc.send(ctx, "email-confirmation", lang, authNotificationPayload{
		EmailAddress: emailAddress,
		URL:          svc.settings.Auth.Frontend.Url.EmailConfirmation + token,
	})
}

func (svc authNotification) PasswordReset(ctx context.Context, lang string, emailAddress string, token string) error {
	return svc.send(ctx, "password-reset", lang, authNotificationPayload{
		EmailAddress: emailAddress,
		URL:          svc.settings.Auth.Frontend.Url.PasswordReset + token,
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

func (svc authNotification) send(ctx context.Context, name, lang string, payload authNotificationPayload) error {
	var (
		err error
		tmp string
		ntf = svc.newMail()
	)

	payload.Logo = htpl.URL(svc.settings.General.Mail.Logo)
	payload.BaseURL = svc.settings.Auth.Frontend.Url.Base
	payload.SignatureName = svc.settings.Auth.Mail.FromName
	payload.SignatureEmail = svc.settings.Auth.Mail.FromAddress

	// @todo translations
	if tmp, err = svc.render(svc.settings.General.Mail.Header, payload); err != nil {
		return fmt.Errorf("failed to render svc.settings.General.Mail.Header: %w", err)
	}
	payload.EmailHeaderEn = htpl.HTML(tmp)
	if tmp, err = svc.render(svc.settings.General.Mail.Footer, payload); err != nil {
		return fmt.Errorf("failed to render svc.settings.General.Mail.Footer: %w", err)
	}
	payload.EmailFooterEn = htpl.HTML(tmp)

	ntf.SetAddressHeader("To", payload.EmailAddress, "")
	// @todo translations
	switch name {
	case "email-confirmation":
		if tmp, err = svc.render(svc.settings.Auth.Mail.EmailConfirmation.Subject, payload); err != nil {
			return fmt.Errorf("failed to render svc.settings.Auth.Mail.EmailConfirmation.Subject: %w", err)
		}
		ntf.SetHeader("Subject", tmp)
		if tmp, err = svc.render(svc.settings.Auth.Mail.EmailConfirmation.Body, payload); err != nil {
			return fmt.Errorf("failed to render svc.settings.Auth.Mail.EmailConfirmation.Body: %w", err)
		}
		ntf.SetBody("text/html", tmp)

	case "password-reset":
		if tmp, err = svc.render(svc.settings.Auth.Mail.PasswordReset.Subject, payload); err != nil {
			return fmt.Errorf("failed to render svc.settings.Auth.Mail.PasswordReset.Subject: %w", err)
		}
		ntf.SetHeader("Subject", tmp)
		if tmp, err = svc.render(svc.settings.Auth.Mail.PasswordReset.Body, payload); err != nil {
			return fmt.Errorf("failed to render svc.settings.Auth.Mail.PasswordReset.Body: %w", err)
		}
		ntf.SetBody("text/html", tmp)

	default:
		return fmt.Errorf("unknown notification email template %q", name)
	}

	svc.log(ctx).Debug(
		"sending auth notification",
		zap.String("name", name),
		zap.String("language", lang),
		zap.String("email", payload.EmailAddress),
	)

	return mail.Send(ntf)
}

func (svc authNotification) render(source string, payload interface{}) (string, error) {
	var (
		err error
		tpl *htpl.Template
		buf = bytes.Buffer{}
	)

	tpl, err = htpl.New("").Parse(source)
	if err != nil {
		return "", fmt.Errorf("could not parse template: %w", err)
	}

	err = tpl.Execute(&buf, payload)
	if err != nil {
		return "", fmt.Errorf("could not render template: %w", err)
	}

	return buf.String(), nil
}
