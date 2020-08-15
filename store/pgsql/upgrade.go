package pgsql

// PostgreSQL specific prefixes, sql
// templates, functions and other helpers

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cortezaproject/corteza-server/store/rdbms"
	"github.com/cortezaproject/corteza-server/store/rdbms/ddl"
	"go.uber.org/zap"
)

type (
	upgrader struct {
		s   *Store
		log *zap.Logger
		ddl *ddl.Generator
	}
)

// NewUpgrader returns PostgreSQL schema upgrader
func NewUpgrader(log *zap.Logger, store *Store) *upgrader {
	var g = &upgrader{store, log, ddl.NewGenerator(log)}

	// All modifications we need for the DDL generator
	// to properly support PostgreSQL dialect

	g.ddl.AddTemplate("create-table-suffix", "WITHOUT OIDS")

	return g
}

// Before runs before all tables are upgraded
func (u upgrader) Before(ctx context.Context) error {
	return rdbms.GenericUpgrades(u.log, u).Before(ctx)
}

// After runs after all tables are upgraded
func (u upgrader) After(ctx context.Context) error {
	return rdbms.GenericUpgrades(u.log, u).After(ctx)
}

// CreateTable is triggered for every table defined in the rdbms package
//
// It checks if table is missing and creates it, otherwise
// it runs
func (u upgrader) CreateTable(ctx context.Context, t *ddl.Table) (err error) {
	var exists bool
	if exists, err = u.TableExists(ctx, t.Name); err != nil {
		return
	}

	if !exists {
		if err = u.Exec(ctx, u.ddl.CreateTable(t)); err != nil {
			return err
		}

		for _, i := range t.Indexes {
			if err = u.Exec(ctx, u.ddl.CreateIndex(i)); err != nil {
				return fmt.Errorf("could not create index %s on table %s: %w", i.Name, i.Table, err)
			}
		}
	} else {
		return u.upgradeTable(ctx, t)
	}

	return nil
}

func (u upgrader) Exec(ctx context.Context, sql string, aa ...interface{}) error {
	_, err := u.s.DB().ExecContext(ctx, sql, aa...)
	return err
}

// upgradeTable applies any necessary changes connected to that specific table
func (u upgrader) upgradeTable(ctx context.Context, t *ddl.Table) error {
	g := rdbms.GenericUpgrades(u.log, u)

	switch t.Name {
	default:
		return g.Upgrade(ctx, t)
	}
}

func (u upgrader) TableExists(ctx context.Context, table string) (bool, error) {
	var exists bool

	if err := u.s.DB().GetContext(ctx, &exists, "SELECT TO_REGCLASS($1) IS NOT NULL", "public."+table); err != nil {
		return false, fmt.Errorf("could not check if table exists: %w", err)
	}

	return exists, nil
}

func (u upgrader) AddColumn(ctx context.Context, table string, col *ddl.Column) (added bool, err error) {
	var (
		lookup = `SELECT is_nullable = 'YES' AS is_nullable,
                         data_type
                    FROM information_schema.columns 
                   WHERE table_catalog = $1 
                     AND table_name = $2
                     AND column_name = $3`

		tmp struct {
			IsNullable bool   `db:"is_nullable"`
			DataType   string `db:"data_type"`
		}
	)

	if err = u.s.DB().GetContext(ctx, &tmp, lookup, u.s.Config().DBName, table, col.Name); err == sql.ErrNoRows {
		if err = u.Exec(ctx, u.ddl.AddColumn(table, col)); err != nil {
			return false, fmt.Errorf("could not add column %s to table %s: %w", table, col.Name, err)
		}

		return true, nil
	} else if err != nil {
		return false, fmt.Errorf("could not check if column exists: %w", err)
	}

	return false, nil
}

func (u upgrader) AddPrimaryKey(ctx context.Context, table string, ind *ddl.Index) (added bool, err error) {
	if err = u.Exec(ctx, u.ddl.AddPrimaryKey(table, ind)); err != nil {
		return false, fmt.Errorf("could not add primary key to table %s: %w", table, err)
	}

	return true, nil
}
