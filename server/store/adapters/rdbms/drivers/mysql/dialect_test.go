package mysql

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
					Name: "text",
				},
			},
			expected: false,
		},
		{
			name: "decimal fits",
			target: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "decimal(1,2)",
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
			name: "decimal doesn't fit",
			target: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "decimal(1,2)",
				},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "decimal(2,3)",
				},
			},
			expected: false,
		},

		{
			name: "varchar fits",
			target: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "varchar(42)",
				},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "varchar(42)",
				},
			},
			expected: true,
		},
		{
			name: "varchar doesn't fit",
			target: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "varchar(42)",
				},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "varchar(84)",
				},
			},
			expected: false,
		},

		{
			name: "sneaking unsigned into signed",
			target: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "bigint signed",
				},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "bigint unsigned",
				},
			},
			expected: false,
		},
		{
			name: "sneaking signed into unsigned",
			target: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "bigint unsigned",
				},
			},
			assert: &ddl.Column{
				Type: &ddl.ColumnType{
					Name: "bigint signed",
				},
			},
			expected: false,
		},
	}

	d := mysqlDialect{}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {
			out := d.ColumnFits(c.target, c.assert)
			assert.Equal(t, c.expected, out)
		})
	}

}
