package service

import (
	"bytes"
	"context"
	"html/template"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gomail "gopkg.in/mail.v2"

	"github.com/cortezaproject/corteza-server/internal/mail"
	"github.com/cortezaproject/corteza-server/pkg/logger"
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
		EmailAddress   string
		URL            string
		BaseURL        string
		Logo           template.URL
		SignatureName  string
		SignatureEmail string
	}
)

var (
	// Crust / Unify your busines (.png)
	defaultLogo = "data:image/png;base64,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"

	emailTemplateHeader = `<div style="width:100%;min-height:100%;margin:0;padding:0;color:#3a393c;font-size:12px;line-height:18px;font-family:Verdana,Arial,sans-serif">
  <table width="100%" align="center" style="width:100%;height:100%;border-collapse:collapse;border:0;padding:60px" border="0" cellspacing="0" cellpadding="0" summary="">
    <tbody>
      <tr>
        <td valign="top" align="center" style="padding: 20px 0;">
          <table width="800" cellspacing="0" cellpadding="0" border="0">
            <tbody>
              <tr>
                <td width="800" bgcolor="#ffffff" style="color:#3a393c;font-size:14px;line-height:20px;font-family:Helvetica Neue,Helvetica,Arial,sans-serif;text-align:left">
                  <table width="800" cellspacing="0" cellpadding="0" border="0">
                    <tbody>
                      <tr style="background-color:#ffffff;height:50px;">
                        <td style="border-bottom:2px solid #1397CB;">
                          <a href="{{ .BaseURL }}" style="text-decoration:none" target="_blank">
						    <img src="{{ .Logo }}" style="display: block;margin: 0 auto;padding: 10px;">
						  </a>
                        </td>
                      </tr>
                      <tr>
                        <td width="800" style="padding:40px 30px">`

	emailTemplateFooter = `</td>
                      </tr>
                      <tr>
                        <td style="padding:30px;border-top: 1px solid #F3F3F5">
                          <p>If you have any questions, please contact <a href="mailto:{{ .SignatureEmail }}" style="color:#1397CB;">{{ .SignatureEmail }}</a>.</p>
                          <p>We hope you enjoy using Corteza!</p>
                          <p>Best regards, <br>
                          {{ .SignatureName }}</p>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </tbody>
          </table>
        </td>
      </tr>
    </tbody>
  </table>
</div>`

	// @todo Temporary email template storage
	emailTemplates = map[string]string{
		"email-confirmation.en.subject": `[Corteza] Confirm your email address`,
		"email-confirmation.en.html": emailTemplateHeader +
			`<h2 style="color: #1397CB;text-align: center;">Confirm your email address</h2>
			  <p>Hello,</p>
			  <p>Follow <a href="{{ .URL }}" style="color:#1397CB;">this link</a> to confirm your email address.</p>
			  <p>You will be logged-in after successful confirmation.</p>` +
			emailTemplateFooter,

		"password-reset.en.subject": `[Corteza] Reset your password`,
		"password-reset.en.html": emailTemplateHeader +
			`<h2 style="color: #1397CB;text-align: center;">Reset your password</h2>
			  <p>Hello,</p>
			  <p>Follow <a href="{{ .URL }}" style="color:#1397CB;">this link</a> and reset your password.</p>
			  <p>You will be logged-in after successful reset.</p>` +
			emailTemplateFooter,
	}
)

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

func (svc authNotification) EmailConfirmation(lang string, emailAddress string, token string) error {
	return svc.send("email-confirmation", lang, authNotificationPayload{
		EmailAddress: emailAddress,
		URL:          svc.settings.frontendUrlEmailConfirmation + token,
	})
}

func (svc authNotification) PasswordReset(lang string, emailAddress string, token string) error {
	return svc.send("password-reset", lang, authNotificationPayload{
		EmailAddress: emailAddress,
		URL:          svc.settings.frontendUrlPasswordReset + token,
	})
}

func (svc authNotification) newMail() *gomail.Message {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", svc.settings.mailFromAddress, svc.settings.mailFromName)
	return m
}

func (svc authNotification) send(name, lang string, payload authNotificationPayload) error {
	ntf := svc.newMail()

	payload.Logo = template.URL(defaultLogo)
	payload.BaseURL = svc.settings.frontendUrlBase
	payload.SignatureName = svc.settings.mailFromName
	payload.SignatureEmail = svc.settings.mailFromAddress

	ntf.SetAddressHeader("To", payload.EmailAddress, "")
	ntf.SetHeader("Subject", svc.render(emailTemplates[name+"."+lang+".subject"], payload))
	ntf.SetBody("text/html", svc.render(emailTemplates[name+"."+lang+".html"], payload))

	svc.log().Debug(
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
