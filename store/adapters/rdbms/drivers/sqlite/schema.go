package sqlite

//
//import (
//	"context"
//	"fmt"
//	"github.com/cortezaproject/corteza-server/pkg/dal"
//
//	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
//	"github.com/jmoiron/sqlx"
//)
//
//type (
//	schema struct{}
//)
//
//func (s *schema) TableExists(ctx context.Context, db sqlx.QueryerContext, table string) (bool, error) {
//	var (
//		exists bool
//		sql    = `SELECT COUNT(*) > 0 FROM sqlite_master WHERE type = 'table' AND name = ?`
//	)
//
//	if err := sqlx.GetContext(ctx, db, &exists, sql, table); err != nil {
//		return false, fmt.Errorf("could not check if table exists: %w", err)
//	}
//
//	return exists, nil
//}
//
//func (s *schema) ColumnExists(ctx context.Context, db sqlx.QueryerContext, column, table string) (bool, error) {
//	var (
//		exists bool
//		sql    = `SELECT COUNT(*) > 0 FROM pragma_table_info(?) WHERE name = ?;`
//	)
//
//	if err := sqlx.GetContext(ctx, db, &exists, sql, table, column); err != nil {
//		return false, fmt.Errorf("could not check if column exists: %w", err)
//	}
//
//	return exists, nil
//}
//
//func (s *schema) CreateTable(ctx context.Context, db sqlx.ExtContext, m *dal.Model) (err error) {
//	tt := append(
//		[]any{&ddl.CreateTable{
//			Dialect: Dialect(),
//			Model:   m,
//		}},
//
//		// @todo
//		// ddl.CreateIndexTemplates(&ddl.CreateIndex{OmitFieldLength: true}, m.Indexes...)...,
//	)
//
//	return ddl.Exec(ctx, db, tt...)
//}
//
//func (s *schema) AddColumn(ctx context.Context, db sqlx.ExtContext, m *dal.Model, aa ...*dal.Attribute) (err error) {
//	var (
//		aux    []any
//		exists bool
//	)
//
//	for _, a := range aa {
//		// check column existence
//		if exists, err = s.ColumnExists(ctx, db, a.Ident, m.Ident); err != nil {
//			return
//		} else if exists {
//			// column exists
//			continue
//		}
//
//		aux = append(aux, &ddl.AddColumn{
//			Dialect:   dialect,
//			Model:     m,
//			Attribute: a,
//		})
//	}
//
//	return ddl.Exec(ctx, db, aux...)
//}
//
//func (s *schema) DropColumn(ctx context.Context, db sqlx.ExtContext, m *dal.Model, cc ...string) (err error) {
//	var (
//		aux    []any
//		exists bool
//	)
//
//	for _, c := range cc {
//		// check column existence
//		if exists, err = s.ColumnExists(ctx, db, c, m.Ident); err != nil {
//			return
//		} else if !exists {
//			// column exists
//			continue
//		}
//
//		aux = append(aux, &ddl.DropColumn{
//			Dialect: dialect,
//			Model:   m,
//			Column:  c,
//		})
//	}
//
//	return ddl.Exec(ctx, db, aux...)
//}
//
//func (s *schema) RenameColumn(ctx context.Context, db sqlx.ExtContext, m *dal.Model, old, new string) (err error) {
//	var (
//		exists bool
//	)
//
//	if exists, err = s.ColumnExists(ctx, db, old, m.Ident); err != nil || !exists {
//		// error or old column does not exist
//		return
//	}
//
//	if exists, err = s.ColumnExists(ctx, db, new, m.Ident); err != nil || exists {
//		// error or new column already exists
//		return
//	}
//
//	return ddl.Exec(ctx, db, &ddl.RenameColumn{
//		Dialect: dialect,
//		Model:   m,
//		Old:     old,
//		New:     new,
//	})
//}
