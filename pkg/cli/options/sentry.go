package options

import (
	"github.com/cortezaproject/corteza-server/internal/version"
)

type (
	SentryOpt struct {
		DSN string `env:"SENTRY_DSN"`

		Debug            bool    `env:"SENTRY_DEBUG"`
		AttachStacktrace bool    `env:"SENTRY_ATTACH_STACKTRACE"`
		SampleRate       float32 `env:"SENTRY_SAMPLE_RATE"`
		MaxBreadcrumbs   int     `env:"SENTRY_MAX_BREADCRUMBS"`

		ServerName  string `env:"SENTRY_SERVERNAME"`
		Release     string `env:"SENTRY_RELEASE"`
		Dist        string `env:"SENTRY_DIST"`
		Environment string `env:"SENTRY_ENVIRONMENT"`
	}
)

func Sentry(pfix string) (o *SentryOpt) {
	o = &SentryOpt{
		AttachStacktrace: true,
		MaxBreadcrumbs:   0,

		Release: version.Version,
	}

	fill(o, pfix)

	return
}
