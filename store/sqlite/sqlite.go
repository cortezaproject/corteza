package sqlite

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/rdbms"
	"github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"strings"
)

type (
	Store struct {
		*rdbms.Store
	}
)

func New(ctx context.Context, dsn string) (s *Store, err error) {
	var cfg *rdbms.Config

	if cfg, err = ProcDataSourceName(dsn); err != nil {
		return nil, err
	}

	cfg.PlaceholderFormat = squirrel.Dollar
	cfg.TxRetryErrHandler = txRetryErrHandler
	cfg.ErrorHandler = errorHandler
	//cfg.TxDisabled = true

	s = new(Store)
	if s.Store, err = rdbms.New(ctx, cfg); err != nil {
		return nil, err
	}

	return s, nil
}

func NewInMemory(ctx context.Context) (s *Store, err error) {
	return New(ctx, "sqlite3://file::memory:?cache=shared&mode=rwc")
}

func (s *Store) Upgrade(ctx context.Context, log *zap.Logger) (err error) {
	if err = (&rdbms.Schema{}).Upgrade(ctx, NewUpgrader(log, s)); err != nil {
		return fmt.Errorf("can not upgrade sqlite schema: %w", err)
	}

	return nil
}

// ProcDataSourceName validates given DSN and ensures
// params are present and correct
func ProcDataSourceName(in string) (*rdbms.Config, error) {
	const (
		schemeDel   = "://"
		validScheme = "sqlite3"
	)

	var (
		endOfSchema = strings.Index(in, schemeDel)
		c           = &rdbms.Config{}
	)

	if endOfSchema > 0 && (in[:endOfSchema] == validScheme || strings.HasPrefix(in[:endOfSchema], validScheme+"+")) {
		c.DriverName = in[:endOfSchema]
		c.DataSourceName = in[endOfSchema+len(schemeDel):]
	} else {
		return nil, fmt.Errorf("expecting valid schema (sqlite3://) at the beginning of the DSN (%s)", in)
	}

	return c, nil
}

func txRetryErrHandler(try int, err error) bool {
	for errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}

	var sqliteErr, ok = err.(sqlite3.Error)
	if !ok {
		return false
	}

	switch sqliteErr.Code {
	case sqlite3.ErrLocked:
		return true

	}

	return false
}

func errorHandler(err error) error {
	if err != nil {
		if implErr, ok := err.(sqlite3.Error); ok {
			switch implErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				return store.ErrNotUnique.Wrap(implErr)
			}
		}
	}

	return err
}
