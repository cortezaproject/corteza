package service

import (
	"bytes"
	"context"
	"html/template"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gomail "gopkg.in/mail.v2"

	"github.com/crusttech/crust/internal/logger"
	"github.com/crusttech/crust/internal/mail"
)

type (
	authNotification struct {
		ctx    context.Context
		logger *zap.Logger

		settings authSettings
	}

	AuthNotificationService interface {
		With(ctx context.Context) AuthNotificationService

		EmailConfirmation(lang string, emailAddress string, url string) error
		PasswordReset(lang string, emailAddress string, url string) error
	}

	authNotificationPayload struct {
		EmailAddress string
		URL          string
	}
)

// @todo Temporary email template storage
var emailTemplates = map[string]string{
	"email-confirmation.en.subject": `[Crust] Confirm your email address`,
	"email-confirmation.en.html":    `<p><a href="{{ .URL }}">Confirm your email address ({{ .EmailAddress }})</a></p>`,

	"password-reset.en.subject": `[Crust] Change your password`,
	"password-reset.en.html":    `<p><a href="{{ .URL }}">Change your password</a></p>`,
}

func AuthNotification(ctx context.Context) AuthNotificationService {
	return (&authNotification{
		logger: DefaultLogger.Named("auth-notification"),
	}).With(ctx)
}

func (svc authNotification) With(ctx context.Context) AuthNotificationService {
	return &authNotification{
		ctx:    ctx,
		logger: svc.logger,

		settings: DefaultAuthSettings,
	}
}

func (svc authNotification) log(fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(svc.ctx, svc.logger).With(fields...)
}

func (svc authNotification) EmailConfirmation(lang string, emailAddress string, url string) error {
	return svc.send("email-confirmation", lang, authNotificationPayload{
		EmailAddress: emailAddress,
		URL:          url,
	})
}

func (svc authNotification) PasswordReset(lang string, emailAddress string, url string) error {
	return svc.send("password-reset", lang, authNotificationPayload{
		EmailAddress: emailAddress,
		URL:          url,
	})
}

func (svc authNotification) newMail() *gomail.Message {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", svc.settings.mailFromAddress, svc.settings.mailFromName)
	return m
}

func (svc authNotification) send(name, lang string, payload authNotificationPayload) error {
	ntf := svc.newMail()

	ntf.SetAddressHeader("To", payload.EmailAddress, "")
	ntf.SetHeader("Subject", svc.render(emailTemplates[name+"."+lang+".subject"], payload))
	ntf.SetBody("text/html", svc.render(emailTemplates[name+"."+lang+".html"], payload))

	svc.log().Debug(
		"sending auth notification (%s.%s) to %q",
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
		svc.log(zap.Error(err)).Error("could not parse template")
		return
	}

	err = tpl.Execute(&buf, payload)
	if err != nil {
		svc.log(zap.Error(err)).Error("could not render template")
		return
	}

	out = buf.String()
	return
}
