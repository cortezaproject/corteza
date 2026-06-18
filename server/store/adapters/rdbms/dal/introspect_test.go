package dal

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/store/adapters/rdbms/ddl"
)

func TestSqlTypeToDalType(t *testing.T) {
	tcs := []struct {
		name string
		sql  string
		null bool
		// assert returns true if the produced dal.Type is the expected one
		assert func(t *testing.T, got dal.Type)
	}{
		{name: "numeric", sql: "numeric(10,2)", assert: func(t *testing.T, got dal.Type) {
			n, ok := got.(*dal.TypeNumber)
			require.True(t, ok, "want *dal.TypeNumber, got %T", got)
			require.Equal(t, 10, n.Precision)
			require.Equal(t, 2, n.Scale)
		}},
		{name: "bigint", sql: "int8", assert: func(t *testing.T, got dal.Type) {
			_, ok := got.(*dal.TypeNumber)
			require.True(t, ok, "want *dal.TypeNumber, got %T", got)
		}},
		{name: "varchar length", sql: "varchar(255)", assert: func(t *testing.T, got dal.Type) {
			tt, ok := got.(*dal.TypeText)
			require.True(t, ok, "want *dal.TypeText, got %T", got)
			require.Equal(t, uint(255), tt.Length)
		}},
		{name: "text", sql: "text", assert: func(t *testing.T, got dal.Type) {
			_, ok := got.(*dal.TypeText)
			require.True(t, ok, "want *dal.TypeText, got %T", got)
		}},
		{name: "timestamptz", sql: "timestamptz", assert: func(t *testing.T, got dal.Type) {
			ts, ok := got.(*dal.TypeTimestamp)
			require.True(t, ok, "want *dal.TypeTimestamp, got %T", got)
			require.True(t, ts.Timezone)
		}},
		{name: "timestamp", sql: "timestamp", assert: func(t *testing.T, got dal.Type) {
			ts, ok := got.(*dal.TypeTimestamp)
			require.True(t, ok, "want *dal.TypeTimestamp, got %T", got)
			require.False(t, ts.Timezone)
		}},
		{name: "mysql datetime", sql: "datetime", assert: func(t *testing.T, got dal.Type) {
			_, ok := got.(*dal.TypeTimestamp)
			require.True(t, ok, "want *dal.TypeTimestamp, got %T", got)
		}},
		{name: "date", sql: "date", assert: func(t *testing.T, got dal.Type) {
			_, ok := got.(*dal.TypeDate)
			require.True(t, ok, "want *dal.TypeDate, got %T", got)
		}},
		{name: "time", sql: "time", assert: func(t *testing.T, got dal.Type) {
			_, ok := got.(*dal.TypeTime)
			require.True(t, ok, "want *dal.TypeTime, got %T", got)
		}},
		{name: "bool", sql: "boolean", assert: func(t *testing.T, got dal.Type) {
			_, ok := got.(*dal.TypeBoolean)
			require.True(t, ok, "want *dal.TypeBoolean, got %T", got)
		}},
		// the bug this guards: only tinyint(1) is boolean
		{name: "tinyint(1) is bool", sql: "tinyint(1)", assert: func(t *testing.T, got dal.Type) {
			_, ok := got.(*dal.TypeBoolean)
			require.True(t, ok, "tinyint(1) must be boolean, got %T", got)
		}},
		{name: "tinyint(4) is number", sql: "tinyint(4)", assert: func(t *testing.T, got dal.Type) {
			_, ok := got.(*dal.TypeNumber)
			require.True(t, ok, "tinyint(4) must be number, got %T", got)
		}},
		{name: "bare tinyint is number", sql: "tinyint", assert: func(t *testing.T, got dal.Type) {
			_, ok := got.(*dal.TypeNumber)
			require.True(t, ok, "bare tinyint must be number, got %T", got)
		}},
		{name: "jsonb", sql: "jsonb", assert: func(t *testing.T, got dal.Type) {
			_, ok := got.(*dal.TypeJSON)
			require.True(t, ok, "want *dal.TypeJSON, got %T", got)
		}},
		{name: "uuid", sql: "uuid", assert: func(t *testing.T, got dal.Type) {
			_, ok := got.(*dal.TypeUUID)
			require.True(t, ok, "want *dal.TypeUUID, got %T", got)
		}},
		{name: "bytea", sql: "bytea", assert: func(t *testing.T, got dal.Type) {
			_, ok := got.(*dal.TypeBlob)
			require.True(t, ok, "want *dal.TypeBlob, got %T", got)
		}},
		{name: "unknown falls back to text", sql: "geometry", assert: func(t *testing.T, got dal.Type) {
			_, ok := got.(*dal.TypeText)
			require.True(t, ok, "unknown type must fall back to text, got %T", got)
		}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := sqlTypeToDalType(&ddl.ColumnType{Name: tc.sql, Null: tc.null})
			require.NotNil(t, got)
			tc.assert(t, got)
		})
	}
}

func TestSqlTypeToDalType_nullable(t *testing.T) {
	got := sqlTypeToDalType(&ddl.ColumnType{Name: "text", Null: true})
	tt, ok := got.(*dal.TypeText)
	require.True(t, ok)
	require.True(t, tt.Nullable)
}

func TestPrimaryKeyColumns(t *testing.T) {
	t.Run("from primary index", func(t *testing.T) {
		tbl := &ddl.Table{
			Indexes: []*ddl.Index{{
				Ident:  "PRIMARY",
				Fields: []*ddl.IndexField{{Column: "id"}},
			}},
			Columns: []*ddl.Column{{Ident: "id"}, {Ident: "name"}},
		}
		pk := primaryKeyColumns(tbl)
		require.True(t, pk["id"])
		require.False(t, pk["name"])
	})

	t.Run("fallback to id column", func(t *testing.T) {
		tbl := &ddl.Table{
			Columns: []*ddl.Column{{Ident: "id"}, {Ident: "name"}},
		}
		pk := primaryKeyColumns(tbl)
		require.True(t, pk["id"])
	})
}
