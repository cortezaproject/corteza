package db

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"go.uber.org/zap"
)

const (
	maxTries = 100
	delay    = 5 * time.Second
	timeout  = 1 * time.Minute
)

func TryToConnect(ctx context.Context, log *zap.Logger, name, dsn, profiler string) (db *factory.DB, err error) {
	factory.Database.Add(name, dsn)

	var (
		connErrCh = make(chan error, 1)
	)

	// This logger is also used inside profiler
	log = log.Named("database").With(zap.String("name", name))

	defer close(connErrCh)

	log.Debug("connecting to the database",
		zap.String("dsn", dsn),
		zap.Int("tries", maxTries),
		zap.Duration("delay", delay),
		zap.Duration("timeout", timeout))

	go func() {
		var (
			try = 0
		)

		for {
			try++

			if maxTries <= try {
				err = errors.Errorf("could not connect to %q, in %d tries", name, try)
				return
			}

			db, err = factory.Database.Get(name)
			if err != nil {
				log.Warn(
					"could not connect to the database",
					zap.Error(err),
					zap.Int("try", try),
					zap.String("dsn", dsn),
					zap.Float64("delay", delay.Seconds()),
				)

				select {
				case <-ctx.Done():
					// Forced break
					break
				case <-time.After(delay):
					// Wait before next try
					continue
				}
			}

			log.Info("connected to the database", zap.String("dsn", dsn))

			// Connected
			break

		}

		connErrCh <- err
	}()

	select {
	case err = <-connErrCh:
		break
	case <-time.After(timeout):
		// Wait before next try
		return nil, errors.Errorf("db init for %q timedout", name)
	case <-ctx.Done():
		return nil, errors.Errorf("db connection for %q cancelled", name)
	}

	switch profiler {
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
