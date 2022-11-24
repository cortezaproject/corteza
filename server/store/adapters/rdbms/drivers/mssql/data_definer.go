package mssql

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ddl"
	"github.com/jmoiron/sqlx"
)

type (
	// dataDefiner for MySQL
	dataDefiner struct {
		dbName string
		conn   *sqlx.DB
		is     *informationSchema
		d      *mssqlDialect
	}
)

var (
	_ ddl.DataDefiner = new(dataDefiner)
)

func DataDefiner(dbName string, conn *sqlx.DB) *dataDefiner {
	return &dataDefiner{
		dbName: dbName,
		conn:   conn,
		is:     InformationSchema(conn),
		d:      Dialect(),
	}
}

func (dd *dataDefiner) ConvertModel(m *dal.Model) (tbl *ddl.Table, err error) {
	tbl, err = ddl.ConvertModel(m, dd.d)
	if err != nil {
		return
	}

	// We'll solve conditional indexes on the app level
	// @todo check if we can use what they provide; the code to replace upsert
	// wasn't working for me.
	//
	// We need to prevent these indexes from adding
	//
	// loop through indexes and remove all with predicate

	for i := len(tbl.Indexes) - 1; i >= 0; i-- {
		if tbl.Indexes[i].Predicate != "" {
			tbl.Indexes = append(tbl.Indexes[:i], tbl.Indexes[i+1:]...)
		}
	}

	return
}

func (dd *dataDefiner) ConvertAttribute(attr *dal.Attribute) (*ddl.Column, error) {
	return ddl.ConvertAttribute(attr, dd.d)
}

func (dd *dataDefiner) TableCreate(ctx context.Context, t *ddl.Table) error {
	return ddl.Exec(ctx, dd.conn, &ddl.CreateTable{
		Dialect:               dd.d,
		Table:                 t,
		OmitIfNotExistsClause: true,
	})
}

func (dd *dataDefiner) TableDrop(ctx context.Context, t string) error {
	return ddl.Exec(ctx, dd.conn, &ddl.DropTable{
		Dialect: dd.d,
		Table:   t,
	})
}

func (dd *dataDefiner) TableLookup(ctx context.Context, t string) (*ddl.Table, error) {
	return dd.is.TableLookup(ctx, t, dd.dbName)
}

func (dd *dataDefiner) ColumnAdd(ctx context.Context, t string, c *ddl.Column) error {
	return ddl.Exec(ctx, dd.conn, &ddl.AddColumn{
		Dialect: dd.d,
		Table:   t,
		Column:  c,
	})
}

func (dd *dataDefiner) ColumnDrop(ctx context.Context, t, col string) error {
	return ddl.Exec(ctx, dd.conn, &ddl.DropColumn{
		Dialect: dd.d,
		Table:   t,
		Column:  col,
	})
}

func (dd *dataDefiner) ColumnRename(ctx context.Context, t string, o string, n string) error {
	return ddl.Exec(ctx, dd.conn, &ddl.RenameColumn{
		Dialect: dd.d,
		Table:   t,
		Old:     o,
		New:     n,
	})
}

func (dd *dataDefiner) IndexLookup(ctx context.Context, i, t string) (*ddl.Index, error) {
	if index, err := dd.is.IndexLookup(ctx, i, t, dd.dbName); err != nil {
		return nil, err
	} else {
		return index, nil
	}
}

func (dd *dataDefiner) IndexCreate(ctx context.Context, t string, i *ddl.Index) error {
	return ddl.Exec(ctx, dd.conn, &ddl.CreateIndex{
		Dialect:               dd.d,
		Index:                 i,
		OmitIfNotExistsClause: true,
	})
}

func (dd *dataDefiner) IndexDrop(ctx context.Context, t, i string) error {
	return ddl.Exec(ctx, dd.conn, &ddl.DropIndex{
		Dialect: dd.d,
		Ident:   i,
	})
}
