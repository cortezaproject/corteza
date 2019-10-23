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
)

type (
	authNotification struct {
		ctx    context.Context
		logger *zap.Logger

		settings       *AuthSettings
		systemSettings *SystemSettings
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
		logger:         DefaultLogger.Named("auth-notification"),
		settings:       DefaultAuthSettings,
		systemSettings: DefaultSystemSettings,
	}).With(ctx)
}

func (svc authNotification) With(ctx context.Context) AuthNotificationService {
	return &authNotification{
		ctx:            ctx,
		logger:         logger.AddRequestID(ctx, svc.logger),
		settings:       svc.settings,
		systemSettings: svc.systemSettings,
	}
}

func (svc authNotification) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}

func (svc authNotification) EmailConfirmation(lang string, emailAddress string, token string) error {
	return svc.send("email-confirmation", lang, authNotificationPayload{
		EmailAddress: emailAddress,
		URL:          svc.settings.FrontendUrlEmailConfirmation + token,
	})
}

func (svc authNotification) PasswordReset(lang string, emailAddress string, token string) error {
	return svc.send("password-reset", lang, authNotificationPayload{
		EmailAddress: emailAddress,
		URL:          svc.settings.FrontendUrlPasswordReset + token,
	})
}

func (svc authNotification) newMail() *gomail.Message {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", svc.settings.MailFromAddress, svc.settings.MailFromName)
	return m
}

func (svc authNotification) send(name, lang string, payload authNotificationPayload) error {
	ntf := svc.newMail()

	payload.Logo = template.URL(svc.systemSettings.DefaultLogo)
	payload.BaseURL = svc.settings.FrontendUrlBase
	payload.SignatureName = svc.settings.MailFromName
	payload.SignatureEmail = svc.settings.MailFromAddress

	// @todo translations
	payload.EmailHeaderEn = template.HTML(svc.render(svc.systemSettings.MailHeader, payload))
	payload.EmailFooterEn = template.HTML(svc.render(svc.systemSettings.MailFooter, payload))

	ntf.SetAddressHeader("To", payload.EmailAddress, "")
	// @todo translations
	switch name {
	case "email-confirmation":
		ntf.SetHeader("Subject", svc.render(svc.settings.MailEmailConfirmationSubject, payload))
		ntf.SetBody("text/html", svc.render(svc.settings.MailEmailConfirmationBody, payload))

	case "password-reset":
		ntf.SetHeader("Subject", svc.render(svc.settings.MailPasswordResetSubject, payload))
		ntf.SetBody("text/html", svc.render(svc.settings.MailPasswordResetBody, payload))

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
