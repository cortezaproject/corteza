package rdbms

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/sentry"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	// TxRetryHardLimit is the absolute maximum retries we'll allow
	TxRetryHardLimit = 100

	DefaultSliceCapacity = 1000

	MinEnsureFetchLimit = 10
	MaxRefetches        = 100

	MaxLimit = 1000
)

// Connect is rdbms' package connector
//
// Function is called from (store) driver's required connection function
// to open connection and bind it with
func Connect(ctx context.Context, log *zap.Logger, cfg *ConnConfig) (db *sqlx.DB, err error) {
	var (
		connErrCh = make(chan error, 1)
		patience  = time.Now().Add(cfg.ConnTryPatience)
		base      *sql.DB
	)

	log = log.Named("store")

	if base, err = sql.Open(cfg.DriverName, cfg.DataSourceName); err != nil {
		return
	}

	db = sqlx.NewDb(base, cfg.DriverName)
	log.Debug(
		"setting database connection parameters",
		zap.Int("MaxOpenConns", cfg.MaxOpenConns),
		zap.Duration("MaxLifetime", cfg.ConnMaxLifetime),
		zap.Int("MaxIdleConns", cfg.MaxIdleConns),

		// log DSN with masked username and password
		zap.String("DSN", cfg.MaskedDSN),
	)

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	go func() {
		defer sentry.Recover()

		var (
			err error
			try = 0

			// Make a small adjustment when
			// collecting callers from the callstack for this
			log = log.WithOptions(zap.AddCallerSkip(-2))
		)

		for {
			try++

			if cfg.ConnTryMax <= try {
				connErrCh <- fmt.Errorf("could not connect in %d tries", try)
				return
			}

			if err = db.PingContext(ctx); err != nil {

				if time.Now().After(patience) {
					// don't make too much fuss
					// if we're in patience mode
					log.Warn(
						"could not connect to the database",
						zap.Error(err),
						zap.Int("try", try),
						zap.Float64("delay", cfg.ConnTryBackoffDelay.Seconds()),
					)
				}

				select {
				case <-ctx.Done():
					// Forced break
					break
				case <-time.After(cfg.ConnTryBackoffDelay):
					//	Wait before next try
					continue
				}
			}

			log.Debug("connected to the database")
			break
		}

		connErrCh <- err
	}()

	to := cfg.ConnTryTimeout * time.Duration(cfg.ConnTryMax*2)
	select {
	case err = <-connErrCh:
		return
	case <-time.After(to):
		// Wait before next try
		return nil, fmt.Errorf("timedout after %.2fs", to.Seconds())
	case <-ctx.Done():
		return nil, fmt.Errorf("connection cancelled")
	}
}

// Connect is called from the adapter's Connect function
//
// It is intentionally not compatible with store.ConnectorFn
// and can not be used to register
//func Connect(ctx context.Context, cfg *ConnConfig) (s *Store, err error) {
//	return s, func() error {
//		if err = cfg.ParseExtra(); err != nil {
//			return err
//		}
//
//		cfg.SetDefaults()
//		s = &Store{
//			config: cfg,
//			//schema: ddl.SchemaAPI(s.DB(), ddl.NewCommonDialect()),
//		}
//
//		return s.Connect(ctx)
//	}()
//}

func dbHealthcheck(db *sqlx.DB) func(ctx context.Context) error {
	return db.PingContext
}
