package mysql

// MySQL specific prefixes, sql
// templates, functions and other helpers

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cortezaproject/corteza-server/store/rdbms"
	"github.com/cortezaproject/corteza-server/store/rdbms/ddl"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

type (
	upgrader struct {
		s   *Store
		log *zap.Logger
		ddl *ddl.Generator
	}
)

// NewUpgrader returns MySQL schema upgrader
func NewUpgrader(log *zap.Logger, store *Store) *upgrader {
	var u = &upgrader{store, log, ddl.NewGenerator(log)}

	// All modifications we need for the DDL generator
	// to properly support MySQL dialect:

	u.ddl.AddTemplate("create-table-suffix", "ENGINE=InnoDB DEFAULT CHARSET=utf8")

	// Sadly, MySQL does not support partial indexes
	//
	// To work around this, we'll ignore partial indexes
	// and solve this on application level
	u.ddl.AddTemplate(
		"create-index",
		`{{ if not .Condition }}CREATE {{ if .Unique }}UNIQUE {{ end }}INDEX {{ template "index-name" . }} ON {{ .Table }} {{ template "index-fields" .Fields }}{{ else }}SELECT 1 -- dummy sql, just to prevent "empty query" errors...{{ end }}`,
	)

	// Cover mysql exceptions
	u.ddl.AddTemplateFunc("columnType", func(ct *ddl.ColumnType) string {
		switch ct.Type {
		case ddl.ColumnTypeIdentifier:
			return "BIGINT UNSIGNED"
		case ddl.ColumnTypeText:
			if y, has := ct.Flags["mysqlLongText"].(bool); has && y {
				return "LONGTEXT"
			}

			return "TEXT"
		case ddl.ColumnTypeTimestamp:
			return "DATETIME"
		case ddl.ColumnTypeBoolean:
			return "TINYINT(1)"
		default:
			return ddl.GenColumnType(ct)
		}
	})

	return u
}

// Before runs before all tables are upgraded
func (u upgrader) Before(ctx context.Context) error {
	err := func() error {
		const migrations = "migrations"
		if exists, err := u.TableExists(ctx, migrations); err != nil || !exists {
			return err
		}

		if _, err := u.s.DB().ExecContext(ctx, fmt.Sprintf("DROP TABLE %s", migrations)); err != nil {
			return err
		}

		u.log.Debug(fmt.Sprintf("%s table Deleted", migrations))

		return nil
	}()

	if err != nil {
		return err
	}

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
	}

	if err = u.upgradeTable(ctx, t); err != nil {
		return
	}

	return nil
}

func (u upgrader) DropTable(ctx context.Context, table string) (dropped bool, err error) {
	var exists bool
	exists, err = u.TableExists(ctx, table)
	if err != nil || !exists {
		return false, err
	}

	err = u.Exec(ctx, fmt.Sprintf(`DROP TABLE %s`, table))
	if err != nil {
		return false, err
	}

	return true, nil
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
	var tmp interface{}
	if err := u.s.DB().GetContext(ctx, &tmp, fmt.Sprintf("SHOW TABLES LIKE %q", table)); err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("could not check if table exists: %w", err)
	}

	return true, nil
}

func (u upgrader) TableSchema(ctx context.Context, table string) (ddl.Columns, error) {
	return nil, fmt.Errorf("pending implementation")
}

// AddColumn adds column to table
// @todo extract column lookup
func (u upgrader) AddColumn(ctx context.Context, table string, col *ddl.Column) (added bool, err error) {
	var (
		lookup = `SELECT IS_NULLABLE = 'YES' AS IS_NULLABLE,
                         DATA_TYPE
                    FROM INFORMATION_SCHEMA.COLUMNS 
                   WHERE TABLE_SCHEMA = ? 
                     AND TABLE_NAME = ?
                     AND COLUMN_NAME = ?`

		tmp struct {
			IsNullable bool   `db:"IS_NULLABLE"`
			DataType   string `db:"DATA_TYPE"`
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

// DropColumn drops column from table
// @todo extract column lookup
func (u upgrader) DropColumn(ctx context.Context, table, column string) (dropped bool, err error) {
	var (
		lookup = `SELECT IS_NULLABLE = 'YES' AS IS_NULLABLE,
                         DATA_TYPE
                    FROM INFORMATION_SCHEMA.COLUMNS 
                   WHERE TABLE_SCHEMA = ? 
                     AND TABLE_NAME = ?
                     AND COLUMN_NAME = ?`

		tmp struct {
			IsNullable bool   `db:"IS_NULLABLE"`
			DataType   string `db:"DATA_TYPE"`
		}
	)

	if err = u.s.DB().GetContext(ctx, &tmp, lookup, u.s.Config().DBName, table, column); err == nil {
		if err = u.Exec(ctx, u.ddl.DropColumn(table, column)); err != nil {
			return false, fmt.Errorf("could not add column %s to table %s: %w", table, column, err)
		}

		return true, nil
	} else if err != sql.ErrNoRows {
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
