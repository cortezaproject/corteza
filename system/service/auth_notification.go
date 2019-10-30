package service

import (
	"bytes"
	"context"
	"html/template"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gomail "gopkg.in/mail.v2"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/mail"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	authNotification struct {
		ctx    context.Context
		logger *zap.Logger

		// @todo merge auth & system settings
		settings *types.Settings
	}

	AuthNotificationService interface {
		With(ctx context.Context) AuthNotificationService

		EmailConfirmation(lang string, emailAddress string, url string) error
		PasswordReset(lang string, emailAddress string, url string) error
	}

	authNotificationPayload struct {
		EmailAddress   string
		URL            string
		BaseURL        string
		Logo           template.URL
		SignatureName  string
		SignatureEmail string
		EmailHeaderEn  template.HTML
		EmailFooterEn  template.HTML
	}
)

func AuthNotification(ctx context.Context) AuthNotificationService {
	return (&authNotification{
		logger:   DefaultLogger.Named("auth-notification"),
		settings: CurrentSettings,
	}).With(ctx)
}

func (svc authNotification) With(ctx context.Context) AuthNotificationService {
	return &authNotification{
		ctx:      ctx,
		logger:   logger.AddRequestID(ctx, svc.logger),
		settings: svc.settings,
	}
}

func (svc authNotification) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}

func (svc authNotification) EmailConfirmation(lang string, emailAddress string, token string) error {
	return svc.send("email-confirmation", lang, authNotificationPayload{
		EmailAddress: emailAddress,
		URL:          svc.settings.Auth.Frontend.Url.EmailConfirmation + token,
	})
}

func (svc authNotification) PasswordReset(lang string, emailAddress string, token string) error {
	return svc.send("password-reset", lang, authNotificationPayload{
		EmailAddress: emailAddress,
		URL:          svc.settings.Auth.Frontend.Url.PasswordReset + token,
	})
}

func (svc authNotification) newMail() *gomail.Message {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", svc.settings.Auth.Mail.FromAddress, svc.settings.Auth.Mail.FromName)
	return m
}

func (svc authNotification) send(name, lang string, payload authNotificationPayload) error {
	ntf := svc.newMail()

	payload.Logo = template.URL(svc.settings.General.Mail.Logo)
	payload.BaseURL = svc.settings.Auth.Frontend.Url.Base
	payload.SignatureName = svc.settings.Auth.Mail.FromName
	payload.SignatureEmail = svc.settings.Auth.Mail.FromAddress

	// @todo translations
	payload.EmailHeaderEn = template.HTML(svc.render(svc.settings.General.Mail.Header, payload))
	payload.EmailFooterEn = template.HTML(svc.render(svc.settings.General.Mail.Footer, payload))

	ntf.SetAddressHeader("To", payload.EmailAddress, "")
	// @todo translations
	switch name {
	case "email-confirmation":
		ntf.SetHeader("Subject", svc.render(svc.settings.Auth.Mail.EmailConfirmation.Subject, payload))
		ntf.SetBody("text/html", svc.render(svc.settings.Auth.Mail.EmailConfirmation.Body, payload))

	case "password-reset":
		ntf.SetHeader("Subject", svc.render(svc.settings.Auth.Mail.PasswordReset.Subject, payload))
		ntf.SetBody("text/html", svc.render(svc.settings.Auth.Mail.PasswordReset.Body, payload))

	default:
		return ErrNoEmailTemplateForGivenOperation
	}

	svc.log(svc.ctx).Debug(
		"sending auth notification",
		zap.String("name", name),
		zap.String("language", lang),
		zap.String("email", payload.EmailAddress),
	)

	return mail.Send(ntf)
}

func (svc authNotification) render(source string, payload interface{}) (out string) {
	var (
		err error
		tpl *template.Template
		buf = bytes.Buffer{}
	)

	tpl, err = template.New("").Parse(source)
	if err != nil {
		svc.log(svc.ctx, zap.Error(err)).Error("could not parse template")
		return
	}

	err = tpl.Execute(&buf, payload)
	if err != nil {
		svc.log(svc.ctx, zap.Error(err)).Error("could not render template")
		return
	}

	out = buf.String()
	return
}
