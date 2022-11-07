package tests

import (
	"encoding/json"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJSONOp(t *testing.T) {
	var (
		tbl = ddl.Table{
			Ident: "test_json_path_test",
			Columns: []*ddl.Column{
				{Ident: "c", Type: &ddl.ColumnType{Name: "JSON"}},
			},
			Temporary: true,
		}

		req = require.New(t)

		// test utility function
		// counts how many rows have a value in the json path
		count = func(req *require.Assertions, expr ...exp.Expression) int {
			var (
				out = struct {
					Count int `db:"count"`
				}{}
			)

			query := conn.dialect.GOQU().
				Select(goqu.COUNT(goqu.Star()).As("count")).
				From(tbl.Ident).
				Where(expr...)

			err := conn.store.QueryOne(ctx, query, &out)
			req.NoError(err)

			return out.Count
		}

		// test utility function
		// counts how many rows have a value in the json path
		countDeepJSON = func(req *require.Assertions, val string) int {
			col := exp.NewIdentifierExpression("", tbl.Ident, "c")
			diJSON, err := conn.dialect.JsonExtractUnquote(col, "a", "b", "c")
			req.NoError(err)

			return count(req, exp.NewLiteralExpression("?", diJSON).Eq(val))
		}

		a2e = func(attr *dal.Attribute) exp.Expression {
			expr, err := conn.dialect.JsonExtract(
				exp.NewIdentifierExpression("", tbl.Ident, attr.StoreIdent()),
				attr.Ident,
			)

			req.NoError(err)
			return expr
		}

		// test utility function
		// counts how many rows have a value in the json path
		countContains = func(req *require.Assertions, attr *dal.Attribute, val any) int {
			contains, err := conn.dialect.JsonArrayContains(
				exp.NewLiteralExpression("?", val),
				a2e(attr),
			)
			req.NoError(err)

			return count(req, contains)
		}

		// test utility function
		insert = func(vv ...string) {
			for _, val := range vv {
				insert := conn.dialect.GOQU().
					Insert(tbl.Ident).
					Cols("c").
					Vals(goqu.Vals{val})

				req.NoError(conn.store.Exec(ctx, insert))
			}
		}

		asJSON = func(val any) string {
			enc, err := json.Marshal(val)
			req.NoError(err)
			return string(enc)
		}

		jsonStore = &dal.CodecRecordValueSetJSON{Ident: "c"}
		numAttr   = &dal.Attribute{Ident: "n", Type: &dal.TypeNumber{}, Store: jsonStore}
		textAttr  = &dal.Attribute{Ident: "s", Type: &dal.TypeText{}, Store: jsonStore}
		boolAttr  = &dal.Attribute{Ident: "b", Type: &dal.TypeBoolean{}, Store: jsonStore}

		numCheck0Attr  = &dal.Attribute{Ident: "nCheck0", Type: &dal.TypeNumber{}, Store: jsonStore}
		textCheck0Attr = &dal.Attribute{Ident: "sCheck0", Type: &dal.TypeText{}, Store: jsonStore}
		boolCheck0Attr = &dal.Attribute{Ident: "bCheck0", Type: &dal.TypeBoolean{}, Store: jsonStore}

		numCheck1Attr  = &dal.Attribute{Ident: "nCheck1", Type: &dal.TypeNumber{}, Store: jsonStore}
		textCheck1Attr = &dal.Attribute{Ident: "sCheck1", Type: &dal.TypeText{}, Store: jsonStore}
		boolCheck1Attr = &dal.Attribute{Ident: "bCheck1", Type: &dal.TypeBoolean{}, Store: jsonStore}
	)

	{
		// check if table exists, drop it and create a new, temporary one

		dd := conn.store.DataDefiner

		// utilize DDL and create a table with a json column
		_, err := dd.TableLookup(ctx, tbl.Ident)
		if errors.IsNotFound(err) {
			err = nil
		} else if err == nil {
			err = dd.TableDrop(ctx, tbl.Ident)
		}

		req.NoError(err)
		req.NoError(dd.TableCreate(ctx, &tbl))
	}

	t.Run("deep ident json", func(t *testing.T) {
		req = require.New(t)
		req.NoError(conn.store.Exec(ctx, conn.dialect.GOQU().Truncate(tbl.Ident)))

		insert(`{"a": {"b": {"c": "match"}}}`)

		req.Equal(1, countDeepJSON(req, "match"))
		req.Equal(0, countDeepJSON(req, "nope"))
	})

	t.Run("json array contains", func(t *testing.T) {
		req = require.New(t)
		req.NoError(conn.store.Exec(ctx, conn.dialect.GOQU().Truncate(tbl.Ident)))

		insert(`{"n": [1,2,3], "b": [true], "s": ["foo", "bar"], "nCheck0": 0, "bCheck0": false, "sCheck0": "baz", "nCheck1": 1, "bCheck1": true, "sCheck1": "foo"}`)

		t.Log("Validating contains check with numeric value")
		req.Equal(1, countContains(req, numAttr, asJSON(1)))
		req.Equal(0, countContains(req, numAttr, asJSON(0)))

		t.Log("Validating contains check with string value")
		req.Equal(1, countContains(req, textAttr, asJSON(`foo`)))
		req.Equal(1, countContains(req, textAttr, asJSON(`bar`)))
		req.Equal(0, countContains(req, textAttr, asJSON(`baz`)))

		t.Log("Validating contains check with boolean value")
		req.Equal(1, countContains(req, boolAttr, asJSON(true)))
		req.Equal(0, countContains(req, boolAttr, asJSON(false)))

		t.Log("Validating contains check with numeric field")
		req.Equal(1, countContains(req, numAttr, a2e(numCheck1Attr)))
		req.Equal(0, countContains(req, numAttr, a2e(numCheck0Attr)))

		t.Log("Validating contains check with string field")
		req.Equal(1, countContains(req, textAttr, a2e(textCheck1Attr)))
		req.Equal(1, countContains(req, textAttr, a2e(textCheck1Attr)))
		req.Equal(0, countContains(req, textAttr, a2e(textCheck0Attr)))

		t.Log("Validating contains check with boolean field")
		req.Equal(1, countContains(req, boolAttr, a2e(boolCheck1Attr)))
		req.Equal(0, countContains(req, boolAttr, a2e(boolCheck0Attr)))
	})
}
