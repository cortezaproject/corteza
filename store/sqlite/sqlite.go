package sqlite

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/store/rdbms"
	_ "github.com/mattn/go-sqlite3"
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

	s = new(Store)
	if s.Store, err = rdbms.New(ctx, cfg); err != nil {
		return nil, err
	}

	return s, nil
}

func NewInMemory(ctx context.Context) (s *Store, err error) {
	var (
		dsn = "sqlite3://file::memory:?cache=shared"
		cfg *rdbms.Config
	)

	if cfg, err = ProcDataSourceName(dsn); err != nil {
		return nil, err
	}

	cfg.PlaceholderFormat = squirrel.Dollar

	s = new(Store)
	if s.Store, err = rdbms.New(ctx, cfg); err != nil {
		return nil, err
	}

	return s, nil
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
