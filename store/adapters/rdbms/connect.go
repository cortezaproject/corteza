package rdbms

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/healthcheck"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type (
	dbLayer interface {
		sqlx.ExecerContext
		SelectContext(context.Context, interface{}, string, ...interface{}) error
		GetContext(context.Context, interface{}, string, ...interface{}) error
		QueryRowContext(context.Context, string, ...interface{}) *sql.Row
		QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	}

	dbTransactionMaker interface {
		BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	}
)

const (
	// TxRetryHardLimit is the absolute maximum retries we'll allow
	TxRetryHardLimit = 100

	DefaultSliceCapacity = 1000

	MinEnsureFetchLimit = 10
	MaxRefetches        = 100

	MaxLimit = 1000
)

// Connect is called from the adapter's Connect  function
//
// It is intentionally not compatible with store.ConnectorFn
// and can not be used to register
func Connect(ctx context.Context, cfg *Config) (s *Store, err error) {
	return s, func() error {
		if err = cfg.ParseExtra(); err != nil {
			return err
		}

		cfg.SetDefaults()
		s = &Store{
			config: cfg,
		}

		return s.Connect(ctx)
	}()
}

// WithTx spins up new store instance with transaction
func (s *Store) withTx(tx dbLayer) *Store {
	return &Store{
		config: s.config,
		//sug:    s.sug,
		db: tx,
	}
}

func (s Store) DB() dbLayer {
	return s.db
}

func (s *Store) Connect(ctx context.Context) error {
	s.log(ctx).Debug("opening connection", zap.String("driver", s.config.DriverName), zap.String("dsn", s.config.MaskedDSN()))

	db, err := sql.Open(s.config.DriverName, s.config.DataSourceName)
	if err != nil {
		return err
	}

	healthcheck.Defaults().Add(dbHealthcheck(db), "Store/RDBMS/"+s.config.DriverName)

	dbx := sqlx.NewDb(db, s.config.DriverName)
	s.log(ctx).Debug(
		"setting connection parameters",
		zap.Int("MaxOpenConns", s.config.MaxOpenConns),
		zap.Duration("MaxLifetime", s.config.ConnMaxLifetime),
		zap.Int("MaxIdleConns", s.config.MaxIdleConns),
	)

	dbx.SetMaxOpenConns(s.config.MaxOpenConns)
	dbx.SetConnMaxLifetime(s.config.ConnMaxLifetime)
	dbx.SetMaxIdleConns(s.config.MaxIdleConns)

	if err = s.tryToConnect(ctx, dbx); err != nil {
		return err
	}

	s.db = dbx
	return err
}

func (s Store) tryToConnect(ctx context.Context, db *sqlx.DB) error {
	var (
		connErrCh = make(chan error, 1)
		patience  = time.Now().Add(s.config.ConnTryPatience)
	)

	go func() {
		defer sentry.Recover()

		var (
			err error
			try = 0

			log = s.log(ctx).
				// Make a small adjustment when
				// collecting callers from the callstack for this
				WithOptions(zap.AddCallerSkip(-2))
		)

		for {
			try++

			if s.config.ConnTryMax <= try {
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
						zap.Float64("delay", s.config.ConnTryBackoffDelay.Seconds()),
					)
				}

				select {
				case <-ctx.Done():
					// Forced break
					break
				case <-time.After(s.config.ConnTryBackoffDelay):
					//	Wait before next try
					continue
				}
			}

			log.Debug("connected to the database")
			break
		}

		connErrCh <- err
	}()

	to := s.config.ConnTryTimeout * time.Duration(s.config.ConnTryMax*2)
	select {
	case err := <-connErrCh:
		return err
	case <-time.After(to):
		// Wait before next try
		return fmt.Errorf("timedout after %.2fs", to.Seconds())
	case <-ctx.Done():
		return fmt.Errorf("connection cancelled")
	}
}

func dbHealthcheck(db *sql.DB) func(ctx context.Context) error {
	return db.PingContext
}
