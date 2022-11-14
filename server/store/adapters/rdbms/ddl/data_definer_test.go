package ddl

import (
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

type (
	mockDriver struct{}
)

func (mockDriver) QuoteIdent(i string) string { return `"` + i + `"` }

func (mockDriver) AttributeToColumn(attr *dal.Attribute) (col *Column, err error) {
	col = &Column{
		Ident:   attr.StoreIdent(),
		Comment: attr.Label,
		Type: &ColumnType{
			Null: attr.Type.IsNullable(),
		},
	}

	switch t := attr.Type.(type) {
	case *dal.TypeID:
		_ = t
		col.Type.Name = "INT"
	case *dal.TypeText:
		_ = t
		col.Type.Name = "TEXT"
	}

	return col, nil
}

func (d mockDriver) IndexFieldModifiers(attr *dal.Attribute, mm ...dal.IndexFieldModifier) (string, error) {
	return IndexFieldModifiers(attr, d.QuoteIdent, mm...)
}

func TestModelToTable(t *testing.T) {

	tests := []struct {
		name string
		m    *dal.Model
		d    driverDialect

		createTableSQL string
		createIndexSQL string
	}{
		{
			name: "simple",
			d:    mockDriver{},
			m: &dal.Model{
				Ident: "simple",
				Attributes: dal.AttributeSet{
					&dal.Attribute{
						Ident: "null_ID",
						Type: &dal.TypeID{
							Nullable: true,
						},
					},
					&dal.Attribute{
						Ident: "some_txt",
						Type:  &dal.TypeText{},
					},
				},

				Indexes: []*dal.Index{
					{
						Ident: PRIMARY_KEY,
						Fields: []*dal.IndexField{
							{
								AttributeIdent: "null_ID",
							},
						},
					},
					{
						Ident: "first_idx",
						Fields: []*dal.IndexField{
							{
								AttributeIdent: "null_ID",
								Modifiers:      []dal.IndexFieldModifier{dal.IndexFieldModifierLower},
								Sort:           dal.IndexFieldSortDesc,
							},
						},
					},
					{
						Ident: "second_idx",
						Fields: []*dal.IndexField{
							{
								AttributeIdent: "null_ID",
							},
							{
								AttributeIdent: "some_txt",
								Modifiers:      []dal.IndexFieldModifier{dal.IndexFieldModifierLower},
								Sort:           dal.IndexFieldSortDesc,
								Nulls:          dal.IndexFieldNullsLast,
							},
						},
					},
				},
			},

			createTableSQL: `
CREATE TABLE IF NOT EXISTS "simple" (
  "null_ID" INT     NULL
, "some_txt" TEXT NOT NULL

, PRIMARY KEY ("null_ID")
)`,

			createIndexSQL: `
CREATE INDEX IF NOT EXISTS "first_idx" ON "simple" ((LOWER("null_ID")) DESC)
CREATE INDEX IF NOT EXISTS "second_idx" ON "simple" ("null_ID", (LOWER("some_txt")) DESC NULLS LAST)
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req = require.New(t)

			tbl, err := ConvertModel(tt.m, tt.d)
			req.NoError(err)

			req.Equal(
				strings.TrimSpace(tt.createTableSQL),
				strings.TrimSpace((&CreateTable{Table: tbl, Dialect: tt.d}).String()),
			)

			idxSQL := "\n"
			for i := 1; i < len(tbl.Indexes); i++ {
				idxSQL += (&CreateIndex{Index: tbl.Indexes[i], Dialect: tt.d}).String() + "\n"
			}

			req.Equal(tt.createIndexSQL, idxSQL)
		})
	}
}
