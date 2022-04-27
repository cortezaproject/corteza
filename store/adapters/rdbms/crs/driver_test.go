package crs

import (
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/data"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/stretchr/testify/require"
)

func Test_sqlizers(t *testing.T) {
	var (
		m = &data.Model{
			Ident: "test_tbl",
			Attributes: data.AttributeSet{
				&data.Attribute{Ident: sysID, Type: &data.TypeID{}, Store: &data.StoreCodecAlias{Ident: "id"}},
				&data.Attribute{Ident: sysCreatedAt, Type: &data.TypeTimestamp{}, Store: &data.StoreCodecAlias{Ident: "created_at"}},
				&data.Attribute{Ident: sysUpdatedAt, Type: &data.TypeTimestamp{}, Store: &data.StoreCodecAlias{Ident: "updated_at"}},
			},
		}

		ms = CRS(m, nil)

		q    sqlizer
		sql  string
		args []any
		err  error

		tenTenTen, _ = time.Parse("2006-01-02", "2010-10-10")
	)

	cases := []struct {
		fn   func() sqlizer
		sql  string
		args []any
		err  error
	}{
		{
			fn:   func() sqlizer { return ms.lookupByIdSql(10) },
			sql:  `SELECT "test_tbl"."id", "test_tbl"."created_at", "test_tbl"."updated_at" FROM "test_tbl" WHERE ("test_tbl"."id" = 10) LIMIT 1`,
			args: []any{},
			err:  nil,
		},
		{
			fn: func() sqlizer {
				return ms.updateSql(&types.Record{ID: 10, CreatedAt: tenTenTen, UpdatedAt: &tenTenTen})
			},
			sql:  `UPDATE "test_tbl" SET "created_at"='2010-10-10T00:00:00Z',"updated_at"='2010-10-10T00:00:00Z' WHERE ("test_tbl"."id" = 10)`,
			args: []any{},
			err:  nil,
		},
		{
			fn:   func() sqlizer { return ms.insertSql(&types.Record{ID: 10, CreatedAt: tenTenTen}) },
			sql:  `INSERT INTO "test_tbl" ("created_at", "id", "updated_at") VALUES ('2010-10-10T00:00:00Z', 10, NULL)`,
			args: []any{},
			err:  nil,
		},
		{
			fn:   func() sqlizer { return ms.deleteByIdSql(12345) },
			sql:  `DELETE FROM "test_tbl" WHERE ("test_tbl"."id" = 12345)`,
			args: []any{},
			err:  nil,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			req := require.New(t)
			q = c.fn()

			sql, args, err = q.ToSQL()
			if c.args == nil {
				req.NoError(err)
			} else {
				req.ErrorIs(err, c.err)
			}

			req.Equal(c.sql, sql)
			req.Equal(c.args, args)

		})
	}
}

func Test_search(t *testing.T) {
	var (
		m = &data.Model{
			Ident: "test_tbl",
			Attributes: data.AttributeSet{
				&data.Attribute{Ident: sysID, Type: &data.TypeID{}, Store: &data.StoreCodecAlias{Ident: "id"}},
				&data.Attribute{Ident: sysModuleID, Type: &data.TypeID{}, Store: &data.StoreCodecAlias{Ident: "rel_module"}},
				&data.Attribute{Ident: sysCreatedAt, Type: &data.TypeTimestamp{}, Store: &data.StoreCodecAlias{Ident: "created_at"}},
				&data.Attribute{Ident: sysUpdatedAt, Type: &data.TypeTimestamp{}, Store: &data.StoreCodecAlias{Ident: "updated_at"}},
				&data.Attribute{Ident: "foo", Type: &data.TypeText{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "bar", Type: &data.TypeNumber{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "baz", Type: &data.TypeBoolean{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "bbb", Type: &data.TypeUUID{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "phy", Type: &data.TypeText{}, Store: &data.StoreCodecPlain{}},
			},
		}

		sql  string
		args []any
		err  error

		moduleID = uint64(1<<64 - 1)

		prefix = `SELECT "test_tbl"."id", "test_tbl"."rel_module", "test_tbl"."created_at", "test_tbl"."updated_at", "test_tbl"."meta", "test_tbl"."phy" FROM "test_tbl"`
	)

	cases := []struct {
		f    types.RecordFilter
		sql  string
		args []any
		err  error
	}{
		{
			f: types.RecordFilter{
				ModuleID: moduleID,
			},
			sql:  prefix + ` WHERE ("test_tbl"."rel_module" = $1)`,
			args: []any{moduleID},
			err:  nil,
		},
		{
			f: types.RecordFilter{
				ModuleID: moduleID,
				Query:    "moduleID == 1234",
			},
			sql:  prefix + ` WHERE (("test_tbl"."rel_module" = $1) AND ("test_tbl"."rel_module" = $2))`,
			args: []any{moduleID, int64(1234)},
			err:  nil,
		},
		{
			f: types.RecordFilter{
				ModuleID: moduleID,
				Query:    `bar = 1 AND foo = phy`,
			},
			sql:  prefix + ` WHERE (("test_tbl"."rel_module" = $1) AND (("meta"->'bar'->0 = $2) AND ("meta"->'foo'->0 = "test_tbl"."phy")))`,
			args: []any{moduleID, int64(1)},
			err:  nil,
		},
	}

	d := CRS(m, nil)

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			req := require.New(t)
			q := d.searchSql(c.f).WithDialect("postgres")
			sql, args, err = q.ToSQL()
			if c.args == nil {
				req.NoError(err)
			} else {
				req.ErrorIs(err, c.err)
			}

			if len(c.sql) > 0 {
				req.Equal(c.sql, sql)
			}

			req.Equal(c.args, args)
		})
	}
}
