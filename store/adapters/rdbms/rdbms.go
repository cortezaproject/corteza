package rdbms

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cortezaproject/corteza-server/store"
	"github.com/doug-martin/goqu/v9/exec"
	"go.uber.org/zap"
)

type (
	sqlizer interface {
		ToSQL() (string, []interface{}, error)
	}

	scalable interface {
		Scan(...interface{}) error
	}

	Store struct {
		db dbLayer

		// rdbms store configuration
		config *Config

		// Logger for connection
		logger *zap.Logger
	}
)

func (s Store) Exec(ctx context.Context, q sqlizer) error {
	var (
		query, args, err = q.ToSQL()
	)

	if err != nil {
		return fmt.Errorf("could not build query: %w", err)
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	return store.HandleError(err, s.config.ErrorHandler)
}

func (s Store) Query(ctx context.Context, q sqlizer) (*sql.Rows, error) {
	var (
		rr *sql.Rows

		query, args, err = q.ToSQL()
	)

	if err != nil {
		return nil, fmt.Errorf("could not build query: %w", err)
	}

	rr, err = s.db.QueryContext(ctx, query, args...)
	if err = store.HandleError(err, s.config.ErrorHandler); err != nil {
		return nil, err
	}

	return rr, nil
}

func (s Store) QueryOne(ctx context.Context, q sqlizer, dst interface{}) (err error) {
	var (
		rows *sql.Rows
	)

	rows, err = s.Query(ctx, q)
	if err != nil {
		return
	}

	defer func() {
		closeError := rows.Close()
		if err == nil {
			// return error from close
			err = closeError
		}
	}()

	if err = rows.Err(); err != nil {
		return
	}

	if !rows.Next() {
		return store.ErrNotFound.Stack(1)
	}

	return exec.NewScanner(rows).ScanStruct(dst)
}

// Config returns store's config as a value
// as protection against changes
func (s Store) Config() Config {
	return *s.config
}
