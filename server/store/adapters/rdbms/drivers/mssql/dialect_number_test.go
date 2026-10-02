package mssql

import (
	"testing"

	"github.com/crusttech/human/server/pkg/dal"
	"github.com/stretchr/testify/assert"
)

func TestAttributeToColumnNumber(t *testing.T) {
	tcc := []struct {
		name     string
		typ      *dal.TypeNumber
		expected string
	}{
		{
			name:     "no precision and scale",
			typ:      &dal.TypeNumber{},
			expected: "DECIMAL",
		},
		{
			name:     "negative precision and scale",
			typ:      &dal.TypeNumber{Precision: -1, Scale: -1},
			expected: "DECIMAL",
		},
		{
			name:     "precision only",
			typ:      &dal.TypeNumber{Precision: 15},
			expected: "DECIMAL(15)",
		},
		{
			name:     "precision and scale",
			typ:      &dal.TypeNumber{Precision: 15, Scale: 3},
			expected: "DECIMAL(15,3)",
		},
	}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {
			col, err := Dialect().AttributeToColumn(&dal.Attribute{Ident: "n", Type: c.typ, Store: &dal.CodecPlain{}})
			assert.NoError(t, err)
			assert.Equal(t, c.expected, col.Type.Name)
		})
	}
}
