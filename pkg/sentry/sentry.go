package sentry

import (
	"github.com/getsentry/sentry-go"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/options"
)

func Init(sentryOpt options.SentryOpt) error {
	if sentryOpt.DSN == "" {
		return nil
	}

	return sentry.Init(sentry.ClientOptions{
		Dsn:              sentryOpt.DSN,
		Debug:            sentryOpt.Debug,
		AttachStacktrace: sentryOpt.AttachStacktrace,
		SampleRate:       float64(sentryOpt.SampleRate),
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

func Recover() {
	// Check if client is configured
	if sentry.CurrentHub().Client() == nil {
		// We do not have Sentry client configured, that means we do not want to
		// recover from panic here as it will only suppress it
		logger.Default()
		return
	}

	if err := recover(); err != nil {
		hub := sentry.CurrentHub()
		hub.Recover(err)
	}
}
