package mysql

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/doug-martin/goqu/v9/exp"

	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/jmoiron/sqlx"
)

type (
	// dataDefiner for MySQL
	dataDefiner struct {
		dbName string
		conn   *sqlx.DB
		is     *informationSchema
		d      *mysqlDialect
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

	// Sadly, MySQL does not support conditional indexes
	// We'll solve that on an app level.
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

func (dd *dataDefiner) TableCreate(ctx context.Context, t *ddl.Table) error {
	return ddl.Exec(ctx, dd.conn, &ddl.CreateTable{
		Table:                 t,
		Dialect:               dd.d,
		OmitIfNotExistsClause: true,
	})
}

func (dd *dataDefiner) TableLookup(ctx context.Context, t string) (*ddl.Table, error) {
	return dd.is.TableLookup(ctx, t, dd.dbName)
}

func (dd *dataDefiner) ColumnAdd(ctx context.Context, t string, c *ddl.Column) error {
	return ddl.Exec(ctx, dd.conn, &ddl.AddColumn{
		Table:  exp.NewIdentifierExpression("", t, ""),
		Column: c,
	})
}

func (dd *dataDefiner) ColumnDrop(ctx context.Context, t, col string) error {
	return ddl.Exec(ctx, dd.conn, &ddl.DropColumn{
		Table:  exp.NewIdentifierExpression("", t, ""),
		Column: exp.NewIdentifierExpression("", "", col),
	})
}

func (dd *dataDefiner) ColumnRename(ctx context.Context, t string, o string, n string) error {
	return ddl.Exec(ctx, dd.conn, &ddl.RenameColumn{
		Table: exp.NewIdentifierExpression("", t, ""),
		Old:   exp.NewIdentifierExpression("", "", o),
		New:   exp.NewIdentifierExpression("", "", n),
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
		Index:                 i,
		Dialect:               dd.d,
		OmitIfNotExistsClause: true,
	})
}

func (dd *dataDefiner) IndexDrop(ctx context.Context, t, i string) error {
	return ddl.Exec(ctx, dd.conn, &ddl.DropIndex{
		Ident:   exp.NewIdentifierExpression("", t, i),
		Dialect: dd.d,
	})
}

//
//// TableExists  checks if table exists in the MySQL database
//func (s *schema) TableExists(ctx context.Context, db sqlx.QueryerContext, table string) (bool, error) {
//	return ddl.TableExists(ctx, db, Dialect(), table, s.dbName)
//}
//
//// ColumnExists  checks if column exists in the MySQL table
//func (s *schema) ColumnExists(ctx context.Context, db sqlx.QueryerContext, column, table string) (bool, error) {
//	return ddl.ColumnExists(ctx, db, Dialect(), column, table, s.dbName)
//}
//
//// CreateTable
////
//// MySQL does not hav CREATE-INDEX-IF-NOT-EXISTS; we need to check index existence manually
//func (s *schema) CreateTable(ctx context.Context, db sqlx.ExtContext, t *ddl.Table) (err error) {
//
//	tc := &ddl.CreateTable{
//		Dialect:      Dialect(),
//		Table:        t,
//		SuffixClause: "ENGINE=InnoDB DEFAULT CHARSET=utf8",
//	}
//
//	if err = ddl.Exec(ctx, db, tc); err != nil {
//		return
//	}
//
//	for _, index := range t.Indexes {
//		if index.Condition != "" {
//			// MySQL, sad little DB does not support
//			// conditional indexes
//			//
//			// We'll solve this on an application level
//			continue
//		}
//
//		var doesIt bool
//		if doesIt, err = ddl.IndexExists(ctx, db, Dialect(), index.Name, index.Table, s.dbName); err != nil {
//			return
//		} else if doesIt {
//			continue
//		}
//
//		ic := &ddl.CreateIndex{
//			Dialect:               Dialect(),
//			Index:                 index,
//			OmitIfNotExistsClause: true,
//		}
//
//		if err = ddl.Exec(ctx, db, ic); err != nil {
//			return
//		}
//
//	}
//
//	return
//}
//
//func (s *schema) AddColumn(ctx context.Context, db sqlx.ExtContext, t *ddl.Table, cc ...*ddl.Column) (err error) {
//	var (
//		aux    []any
//		exists bool
//	)
//
//	for _, c := range cc {
//		// check column existence
//		if exists, err = s.ColumnExists(ctx, db, c.Name, t.Name); err != nil {
//			return
//		} else if exists {
//			// column exists
//			continue
//		}
//
//		// Sadly, some column types in MySQL can not have default values
//		if c.Type.Type == ddl.ColumnTypejsonb|| c.Type.Type == ddl.ColumnTypeBinary || c.Type.Type == ddl.ColumnTypeText {
//			c.DefaultValue = ""
//		}
//
//		aux = append(aux, &ddl.AddColumn{
//			Dialect: dialect,
//			Table:   t,
//			Column:  c,
//		})
//	}
//
//	return ddl.Exec(ctx, db, aux...)
//}
//
//func (s *schema) DropColumn(ctx context.Context, db sqlx.ExtContext, t *ddl.Table, cc ...string) (err error) {
//	var (
//		aux    []any
//		exists bool
//	)
//
//	for _, c := range cc {
//		// check column existence
//		if exists, err = s.ColumnExists(ctx, db, c, t.Name); err != nil {
//			return
//		} else if !exists {
//			// column exists
//			continue
//		}
//
//		aux = append(aux, &ddl.DropColumn{
//			Dialect: dialect,
//			Table:   t,
//			Column:  c,
//		})
//	}
//
//	return ddl.Exec(ctx, db, aux...)
//}
//
//func (s *schema) RenameColumn(ctx context.Context, db sqlx.ExtContext, t *ddl.Table, old, new string) (err error) {
//	var (
//		exists bool
//	)
//
//	if exists, err = s.ColumnExists(ctx, db, old, t.Name); err != nil || !exists {
//		// error or old column does not exist
//		return
//	}
//
//	if exists, err = s.ColumnExists(ctx, db, new, t.Name); err != nil || exists {
//		// error or new column already exists
//		return
//	}
//
//	return ddl.Exec(ctx, db, &ddl.RenameColumn{
//		Dialect: dialect,
//		Table:   t,
//		Old:     old,
//		New:     new,
//	})
//}
