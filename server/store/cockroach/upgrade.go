package cockroach

// PostgreSQL specific prefixes, sql
// templates, functions and other helpers

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/store/rdbms"
	"github.com/cortezaproject/corteza/server/store/rdbms/ddl"
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
	// to properly support CockroachDB dialect

	g.ddl.AddTemplate("index-fields", `
	({{ range $n, $f := . -}}
	{{ if $n }}, {{ end }}
	{{- .Field | trimExpression}}
	{{- if .Desc }} DESC{{ end }}
	{{- end }})
	`)
	g.ddl.AddTemplate("if-not-exists-clause", "IF NOT EXISTS")

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
	//const tableSchemaName = "public"
	var exists bool
	// @todo use pq.NewError(regex?!) -> ADDED IF NOT EXISTS
	//if err := u.s.DB().GetContext(ctx, &exists, fmt.Sprintf("SELECT * FROM %[1]s.%[2]s", tableSchemaName, table)); err != nil {
	//	if err.Error() == fmt.Sprintf("pq: relation \"%[1]s.%[2]s\" does not exist", tableSchemaName, table) {
	//		exists = false
	//	} else {
	//		return false, fmt.Errorf("could not check if table exists: %w", err)
	//	}
	//}

	return exists, nil
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

func (u upgrader) TableSchema(ctx context.Context, table string) (ddl.Columns, error) {
	return nil, fmt.Errorf("pending implementation")
}

// AddColumn adds column to table
func (u upgrader) AddColumn(ctx context.Context, table string, col *ddl.Column) (added bool, err error) {
	err = func() error {
		var columns ddl.Columns
		if columns, err = u.getColumns(ctx, table); err != nil {
			return err
		}

		if columns.Get(col.Name) != nil {
			return nil
		}

		if err = u.Exec(ctx, u.ddl.AddColumn(table, col)); err != nil {
			return err
		}

		added = true
		return nil
	}()

	if err != nil {
		return false, fmt.Errorf("could not add column %q to %q: %w", col.Name, table, err)
	}

	return
}

// DropColumn drops column from table
func (u upgrader) DropColumn(ctx context.Context, table, column string) (dropped bool, err error) {
	err = func() error {
		var columns ddl.Columns
		if columns, err = u.getColumns(ctx, table); err != nil {
			return err
		}

		if columns.Get(column) == nil {
			return nil
		}

		if err = u.Exec(ctx, u.ddl.DropColumn(table, column)); err != nil {
			return err
		}

		dropped = true
		return nil
	}()

	if err != nil {
		return false, fmt.Errorf("could not drop column %q from %q: %w", column, table, err)
	}

	return
}

// RenameColumn renames column on a table
func (u upgrader) RenameColumn(ctx context.Context, table, oldName, newName string) (changed bool, err error) {
	err = func() error {
		if oldName == newName {
			return nil
		}

		var columns ddl.Columns
		if columns, err = u.getColumns(ctx, table); err != nil {
			return err
		}

		if columns.Get(oldName) == nil {
			// Old column does not exist anymore

			if columns.Get(newName) == nil {
				return fmt.Errorf("old and new columns are missing")
			}

			return nil
		}

		if columns.Get(newName) != nil {
			return fmt.Errorf("new column already exists")

		}

		if err = u.Exec(ctx, u.ddl.RenameColumn(table, oldName, newName)); err != nil {
			return err
		}

		changed = true
		return nil
	}()

	if err != nil {
		return false, fmt.Errorf("could not rename column %q on table %q to %q: %w", oldName, table, newName, err)
	}

	return
}

func (u upgrader) AddPrimaryKey(ctx context.Context, table string, ind *ddl.Index) (added bool, err error) {
	if err = u.Exec(ctx, u.ddl.AddPrimaryKey(table, ind)); err != nil {
		return false, fmt.Errorf("could not add primary key to table %s: %w", table, err)
	}

	return true, nil
}

// loads and returns all tables columns
func (u upgrader) getColumns(ctx context.Context, table string) (out ddl.Columns, err error) {
	type (
		col struct {
			Name       string `db:"column_name"`
			IsNullable bool   `db:"is_nullable"`
			DataType   string `db:"data_type"`
		}
	)

	var (
		lookup = `SELECT column_name,
                         is_nullable = 'YES' AS is_nullable,
                         data_type
                    FROM information_schema.columns 
                   WHERE table_catalog = $1 
                     AND table_name = $2`

		cols []*col
	)

	if err = u.s.DB().SelectContext(ctx, &cols, lookup, u.s.Config().DBName, table); err != nil {
		return nil, err
	}

	out = make([]*ddl.Column, len(cols))
	for i := range cols {
		out[i] = &ddl.Column{
			Name: cols[i].Name,
			//Type:         ddl.ColumnType{},
			IsNull: cols[i].IsNullable,
			//DefaultValue: "",
		}
	}

	return out, nil
}
