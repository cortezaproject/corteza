package db

import (
	"context"
	"regexp"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/cli/options"
)

var (
	dsnMasker = regexp.MustCompile("(.)(?:.*)(.):(.)(?:.*)(.)@")
)

func TryToConnect(ctx context.Context, log *zap.Logger, name string, opt options.DBOpt) (db *factory.DB, err error) {
	factory.Database.Add(name, opt.DSN)

	var (
		connErrCh = make(chan error, 1)
	)

	// We'll not add this to the general log because we do not want to carry it with us for every query.
	dsnField := zap.String("dsn", dsnMasker.ReplaceAllString(opt.DSN, "$1****$2:$3****$4@"))

	// This logger is also used inside profiler
	log = log.Named("database").With(zap.String("name", name))

	defer close(connErrCh)

	log.Debug("connecting to the database",
		dsnField,
		zap.Int("tries", opt.MaxTries),
		zap.Duration("delay", opt.Delay),
		zap.Duration("timeout", opt.Timeout))

	go func() {
		defer sentry.Recover()

		var (
			try = 0
		)

		for {
			try++

			if opt.MaxTries <= try {
				err = errors.Errorf("could not connect to %q, in %d tries", name, try)
				return
			}

			db, err = factory.Database.Get(name)
			if err != nil {
				log.Warn(
					"could not connect to the database",
					zap.Error(err),
					zap.Int("try", try),
					dsnField,
					zap.Float64("delay", opt.Delay.Seconds()),
				)

				select {
				case <-ctx.Done():
					// Forced break
					break
				case <-time.After(opt.Delay):
					// Wait before next try
					continue
				}
			}

			log.Info("connected to the database", dsnField)

			// Connected
			break

		}

		connErrCh <- err
	}()

	select {
	case err = <-connErrCh:
		break
	case <-time.After(opt.Timeout):
		// Wait before next try
		return nil, errors.Errorf("db init for %q timedout", name)
	case <-ctx.Done():
		return nil, errors.Errorf("db connection for %q cancelled", name)
	}

	switch opt.Profiler {
	case "stdout":
		db.Profiler = &factory.Database.ProfilerStdout
	case "logger":
		// Skip 3 levels in call stack to get to the actual function used
		db.Profiler = ZapProfiler(log.
			WithOptions(zap.AddCallerSkip(3)),
		)
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}
