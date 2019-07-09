package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/http"
	"github.com/cortezaproject/corteza-server/internal/mail"
	"github.com/cortezaproject/corteza-server/pkg/cli/options"
)

func InitGeneralServices(smtpOpt *options.SMTPOpt, jwtOpt *options.JWTOpt, httpClientOpt *options.HttpClientOpt) {
	auth.SetupDefault(jwtOpt.Secret, int(jwtOpt.Expiry/time.Minute))
	mail.SetupDialer(smtpOpt.Host, smtpOpt.Port, smtpOpt.User, smtpOpt.Pass, smtpOpt.From)
	http.SetupDefaults(
		httpClientOpt.HttpClientTimeout,
		httpClientOpt.ClientTSLInsecure,
	)
}

func InitSentry(sentryOpt *options.SentryOpt) error {
	if sentryOpt.DSN == "" {
		return nil
	}

	return sentry.Init(sentry.ClientOptions{
		Dsn:              sentryOpt.DSN,
		Debug:            sentryOpt.Debug,
		AttachStacktrace: sentryOpt.AttachStacktrace,
		SampleRate:       sentryOpt.SampleRate,
		MaxBreadcrumbs:   sentryOpt.MaxBreadcrumbs,
		IgnoreErrors:     nil,
		BeforeSend:       nil,
		BeforeBreadcrumb: nil,
		Integrations:     nil,
		Transport:        nil,
		ServerName:       sentryOpt.ServerName,
		Release:          sentryOpt.Release,
		Dist:             sentryOpt.Dist,
		Environment:      sentryOpt.Environment,
	})
}

func HandleError(err error) {
	if err == nil {
		return
	}

	_, _ = fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
