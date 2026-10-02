package mssql

import (
	"testing"

	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/store/adapters/rdbms/ddl"
	"github.com/stretchr/testify/assert"
)

// Text must be stored as Unicode (NVARCHAR); VARCHAR turns characters
// outside of the database collation into "?"
func TestAttributeToColumnText(t *testing.T) {
	tcc := []struct {
		name     string
		typ      dal.Type
		expected string
	}{
		{name: "text", typ: &dal.TypeText{}, expected: "NVARCHAR(MAX)"},
		{name: "text with length", typ: &dal.TypeText{Length: 64}, expected: "NVARCHAR(64)"},
		{name: "enum", typ: &dal.TypeEnum{}, expected: "NVARCHAR(MAX)"},
		{name: "json", typ: &dal.TypeJSON{}, expected: "NVARCHAR(MAX)"},
		{name: "geometry", typ: &dal.TypeGeometry{}, expected: "NVARCHAR(MAX)"},
	}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {
			col, err := Dialect().AttributeToColumn(&dal.Attribute{Ident: "t", Type: c.typ, Store: &dal.CodecPlain{}})
			assert.NoError(t, err)
			assert.Equal(t, c.expected, col.Type.Name)
		})
	}
}

// Columns created as VARCHAR before the switch to NVARCHAR must still be
// accepted so existing tables are not reported as needing alterations
func TestColumnFits_existingVarchar(t *testing.T) {
	target := &ddl.Column{Type: &ddl.ColumnType{Name: "varchar(max)"}}
	assertCol := &ddl.Column{Type: &ddl.ColumnType{Name: "NVARCHAR(MAX)"}}

	assert.True(t, Dialect().ColumnFits(target, assertCol))
}
