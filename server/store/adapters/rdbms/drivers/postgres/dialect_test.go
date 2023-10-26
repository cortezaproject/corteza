package postgres

import (
	"testing"

	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ddl"
	"github.com/stretchr/testify/assert"
)

func TestColumnFits(t *testing.T) {
	tcc := []struct {
		name     string
		target   *ddl.Column
		assert   *ddl.Column
		expected bool
	}{
		{
			name: "exact match (text)",
			target: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "text",
				},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "text",
				},
			},
			expected: true,
		},
		{
			name: "fits somewhere",
			target: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "text",
				},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "numeric(1,2)",
				},
			},
			expected: true,
		},
		{
			name: "doesn't fit",
			target: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "numeric(1,2)",
				},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "text",
				},
			},
			expected: false,
		},

		{
			name: "numeric fits",
			target: &ddl.Column{
				Type: &ddl.ColumnType{Name: "numeric(1,2)"},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{Name: "numeric(1,2)"},
			},
			expected: true,
		},

		{
			name: "numeric doesn't fit",
			target: &ddl.Column{
				Type: &ddl.ColumnType{Name: "numeric(1,2)"},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{Name: "numeric(2,3)"},
			},
			expected: false,
		},

		{
			name: "varchar fits",
			target: &ddl.Column{
				Type: &ddl.ColumnType{Name: "varchar(42)"},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{Name: "varchar(42)"},
			},
			expected: true,
		},

		{
			name: "varchar doesn't fit",
			target: &ddl.Column{
				Type: &ddl.ColumnType{Name: "varchar(42)"},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{Name: "varchar(84)"},
			},
			expected: false,
		},

		{
			name: "timestamp fits into timestamptz",
			target: &ddl.Column{
				Type: &ddl.ColumnType{Name: "timestamptz"},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{Name: "timestamp"},
			},
			expected: true,
		},
		{
			name: "timestamptz doesn't fit into timestamp",
			target: &ddl.Column{
				Type: &ddl.ColumnType{Name: "timestamp"},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{Name: "timestamptz"},
			},
			expected: false,
		},

		{
			name: "time fits into timetz",
			target: &ddl.Column{
				Type: &ddl.ColumnType{Name: "timetz"},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{Name: "time"},
			},
			expected: true,
		},
		{
			name: "timetz doesn't fit into time",
			target: &ddl.Column{
				Type: &ddl.ColumnType{Name: "time"},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{Name: "timetz"},
			},
			expected: false,
		},
	}

	d := postgresDialect{}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {
			out := d.ColumnFits(c.target, c.assert)
			assert.Equal(t, c.expected, out)
		})
	}

}
