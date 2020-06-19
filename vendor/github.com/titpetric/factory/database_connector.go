package factory

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/pkg/errors"
	"github.com/titpetric/factory/logger"
)

// DatabaseConnectionOptions is a configuration struct for connection retry
type DatabaseConnectionOptions struct {
	DSN        string
	DriverName string
	Logger     string

	Retries        int
	RetryTimeout   time.Duration
	ConnectTimeout time.Duration
}

var (
	dsnMasker = regexp.MustCompile("(.)(?:.*)(.):(.)(?:.*)(.)@")
)

func (df *DatabaseFactory) TryToConnect(ctx context.Context, name string, options *DatabaseConnectionOptions) (db *DB, err error) {
	df.Add(name, DatabaseCredential{
		DSN:        options.DSN,
		DriverName: options.DriverName,
	})

	var (
		connErrCh = make(chan error, 1)
	)

	// We'll not add this to the general log because we do not want to carry it with us for every query.
	dsnField := fmt.Sprintf("dsn=%s", dsnMasker.ReplaceAllString(options.DSN, "$1****$2:$3****$4@"))

	defer close(connErrCh)

	log.Println("connecting to database", dsnField)

	go func() {
		var (
			try = 0
		)

		for {
			try++

			if options.Retries <= try {
				err = errors.Errorf("could not connect to %q, in %d tries", name, try)
				return
			}

			db, err = df.Get(name)
			if err != nil {
				log.Println(
					"could not connect to the database",
					err,
					dsnField,
					fmt.Sprintf("try=%d", try),
				)

				select {
				case <-ctx.Done():
					// Forced break
					break
				case <-time.After(options.RetryTimeout):
					// Wait before next try
					continue
				}
			}

			break
		}

		connErrCh <- err
	}()

	select {
	case err = <-connErrCh:
		break
	case <-time.After(options.ConnectTimeout):
		// Wait before next try
		return nil, errors.Errorf("db init for %q timed out", name)
	case <-ctx.Done():
		return nil, errors.Errorf("db connection for %q cancelled", name)
	}

	switch options.Logger {
	case "stdout":
		db.SetLogger(logger.Default{})
	default:
		db.SetLogger(logger.Silent{})
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}
