package schema

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_generator_CreateTable(t *testing.T) {
	var (
		tests = []struct {
			name string
			in   *Table
			out  string
		}{
			{
				name: "2col",
				in: TableDef("tbl",
					ColumnDef("col1", ColumnTypeInteger, Null),
					ColumnDef("col2", ColumnTypeVarchar, ColumnTypeLength(2)),
				),
				out: `
CREATE TABLE tbl (
  col1 INTEGER
, col2 VARCHAR(2) NOT NULL

)`,
			},
			{
				name: "primary key",
				in: TableDef("tbl",
					ID,
				),
				out: `
CREATE TABLE tbl (
  id BIGINT NOT NULL

, PRIMARY KEY (id)
)`,
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.New(t).Equal(strings.TrimSpace(tt.out), NewCommonDialect(zap.NewNop()).CreateTable(tt.in))
		})
	}
}

func Test_generator_CreateIndex(t *testing.T) {
	var (
		tests = []struct {
			name string
			in   *Table
			out  string
		}{
			{
				name: "columns",
				in: TableDef("tbl",
					AddIndex("mix", IColumn("c1", "c2")),
				),
				out: "CREATE INDEX tbl_mix ON tbl (c1, c2)",
			},
			{
				name: "expression",
				in: TableDef("tbl",
					AddIndex("mix", IExpr("LOWER(exp1)")),
				),
				out: "CREATE INDEX tbl_mix ON tbl ((LOWER(exp1)))",
			},
			{
				name: "conditional",
				in: TableDef("tbl",
					AddIndex("mix", IColumn("c1"), IWhere("cnd IS NULL")),
				),
				out: "CREATE INDEX tbl_mix ON tbl (c1) WHERE (cnd IS NULL)",
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.New(t).Equal(tt.out, NewCommonDialect(zap.NewNop()).CreateIndex(tt.in.Indexes[0]))
		})
	}
}
