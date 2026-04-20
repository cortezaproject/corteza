package mssql

import (
	"testing"

	"github.com/crusttech/human/server/store/adapters/rdbms/ddl"
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
			name: "exact match (varchar)",
			target: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "varchar",
				},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "varchar",
				},
			},
			expected: true,
		},
		{
			name: "fits somewhere",
			target: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "varchar(max)",
				},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "decimal(1,2)",
				},
			},
			expected: true,
		},
		{
			name: "doesn't fit",
			target: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "decimal(1,2)",
				},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "varchar(max)",
				},
			},
			expected: false,
		},

		{
			name: "decimal fits",
			target: &ddl.Column{
				Type: &ddl.ColumnType{Name: "decimal(1,2)"},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{Name: "decimal(1,2)"},
			},
			expected: true,
		},

		{
			name: "decimal doesn't fit",
			target: &ddl.Column{
				Type: &ddl.ColumnType{Name: "decimal(1,2)"},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{Name: "decimal(2,3)"},
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
			name: "varchar(max) fits",
			target: &ddl.Column{
				Type: &ddl.ColumnType{Name: "varchar(max)"},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{Name: "varchar(42)"},
			},
			expected: true,
		},

		{
			name: "varchar(max) doesn't fit",
			target: &ddl.Column{
				Type: &ddl.ColumnType{Name: "varchar(42)"},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{Name: "varchar(max)"},
			},
			expected: false,
		},

		{
			name: "datetime exact match",
			target: &ddl.Column{
				Type: &ddl.ColumnType{Name: "datetime"},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{Name: "datetime"},
			},
			expected: true,
		},

		{
			name: "time exact match",
			target: &ddl.Column{
				Type: &ddl.ColumnType{Name: "time"},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{Name: "time"},
			},
			expected: true,
		},

		{
			name: "datetime doesn't fit time",
			target: &ddl.Column{
				Type: &ddl.ColumnType{Name: "time"},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{Name: "datetime"},
			},
			expected: false,
		},

		{
			name: "time doesn't fit datetime",
			target: &ddl.Column{
				Type: &ddl.ColumnType{Name: "datetime"},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{Name: "time"},
			},
			expected: false,
		},

		{
			name: "datetimeoffset exact match",
			target: &ddl.Column{
				Type: &ddl.ColumnType{Name: "datetimeoffset"},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{Name: "datetimeoffset"},
			},
			expected: true,
		},
	}

	d := mssqlDialect{}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {
			out := d.ColumnFits(c.target, c.assert)
			assert.Equal(t, c.expected, out)
		})
	}

}
