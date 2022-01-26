package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/sentry.yaml

import (
	"github.com/cortezaproject/corteza-server/pkg/version"
)

type (
	SentryOpt struct {
		DSN              string  `env:"SENTRY_DSN"`
		Debug            bool    `env:"SENTRY_DEBUG"`
		AttachStacktrace bool    `env:"SENTRY_ATTACH_STACKTRACE"`
		SampleRate       float32 `env:"SENTRY_SAMPLE_RATE"`
		MaxBreadcrumbs   int     `env:"SENTRY_MAX_BREADCRUMBS"`
		ServerName       string  `env:"SENTRY_SERVERNAME"`
		Release          string  `env:"SENTRY_RELEASE"`
		Dist             string  `env:"SENTRY_DIST"`
		Environment      string  `env:"SENTRY_ENVIRONMENT"`
	}
)

// Sentry initializes and returns a SentryOpt with default values
func Sentry() (o *SentryOpt) {
	o = &SentryOpt{
		AttachStacktrace: true,
		MaxBreadcrumbs:   0,
		Release:          version.Version,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *Sentry) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
