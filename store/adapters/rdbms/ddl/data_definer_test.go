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

func (mockDriver) NativeColumnType(ct dal.Type) (*ColumnType, error) {
	col := &ColumnType{
		Name: string(ct.Type()),
		Null: ct.IsNullable(),
	}

	switch ct.(type) {
	case *dal.TypeID:
		col.Name = "INT"
	}

	return col, nil
}

func TestModelToTable(t *testing.T) {

	tests := []struct {
		name string
		m    *dal.Model
		d    driverDialect
		sql  string
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
				},
			},

			sql: `
CREATE TABLE IF NOT EXISTS "simple" (
  "null_ID" INT     NULL

)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req = require.New(t)

			tbl, err := ConvertModel(tt.m, tt.d)
			req.NoError(err)

			req.Equal(strings.TrimSpace(tt.sql), strings.TrimSpace((&CreateTable{Table: tbl}).String()))
		})
	}
}
