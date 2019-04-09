package db

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
)

const (
	maxTries = 100
	delay    = 5 * time.Second
	timeout  = 1 * time.Minute
)

func TryToConnect(ctx context.Context, name, dsn, profiler string) (db *factory.DB, err error) {
	factory.Database.Add(name, dsn)

	var (
		connErrCh = make(chan error, 1)
	)

	defer close(connErrCh)

	log.Printf("Connecting to the DB (%q, %q)", name, dsn)

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
				log.Printf(
					"could not connect to %q, try %d, error: %v, retry in %.0fs",
					name,
					try,
					err,
					delay.Seconds(),
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

	// @todo: profiling as an external service?
	switch profiler {
	case "stdout":
		db.Profiler = &factory.Database.ProfilerStdout
	default:
		log.Println("No database query profiler selected")
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}
