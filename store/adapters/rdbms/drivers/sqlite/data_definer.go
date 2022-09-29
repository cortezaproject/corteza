package sqlite

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/jmoiron/sqlx"
)

type (
	// @todo this is unmodified copy of postgres's data definer struct!
	dataDefiner struct {
		dbName string
		conn   *sqlx.DB
		is     *informationSchema
		d      *sqliteDialect
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

func (dd *dataDefiner) ConvertModel(m *dal.Model) (*ddl.Table, error) {
	return ddl.ConvertModel(m, dd.d)
}

func (dd *dataDefiner) ConvertAttribute(attr *dal.Attribute) (*ddl.Column, error) {
	return ddl.ConvertAttribute(attr, dd.d)
}

func (dd *dataDefiner) TableCreate(ctx context.Context, t *ddl.Table) error {
	return ddl.Exec(ctx, dd.conn, &ddl.CreateTable{
		Dialect: dd.d,
		Table:   t,
	})
}

func (dd *dataDefiner) TableLookup(ctx context.Context, t string) (*ddl.Table, error) {
	return dd.is.TableLookup(ctx, t, "public", dd.dbName)
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
		Dialect: dd.d,
		Index:   i,
	})
}

func (dd *dataDefiner) IndexDrop(ctx context.Context, t, i string) error {
	return ddl.Exec(ctx, dd.conn, &ddl.DropIndex{
		Dialect: dd.d,
		Ident:   i,
	})
}
