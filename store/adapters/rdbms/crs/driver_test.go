package crs

import (
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/data"
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

		req  = require.New(t)
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
