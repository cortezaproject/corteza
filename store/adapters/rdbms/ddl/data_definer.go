package ddl

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/dal"
)

type (
	driverDialect interface {
		// NativeColumnType converts column type to type that can be used in the underlying rdbms
		NativeColumnType(columnType dal.Type) (*ColumnType, error)
	}

	// DataDefiner describes an interface for all DDL commands
	DataDefiner interface {
		ConvertModel(*dal.Model) (*Table, error)

		// Tables(ctx context.Context) ([]*Table, error)
		TableLookup(context.Context, string) (*Table, error)
		TableCreate(context.Context, *Table) error

		ColumnAdd(context.Context, string, *Column) error
		ColumnDrop(context.Context, string, string) error
		ColumnRename(context.Context, string, string, string) error

		IndexLookup(context.Context, string, string) (*Index, error)
		IndexCreate(context.Context, string, *Index) error
		IndexDrop(context.Context, string, string) error
	}

	// Table describes structure of the SQL table
	Table struct {
		Ident   string
		Columns []*Column
		Indexes []*Index
		Comment string

		// implementation variations
		Meta map[string]interface{}
	}

	Column struct {
		Ident        string
		Type         *ColumnType
		DefaultValue string

		// implementation variations
		Meta map[string]interface{}

		Comment string
	}

	ColumnType struct {
		Name string

		// additional params for type
		Parameters []any

		Null bool

		// implementation variations
		Meta map[string]interface{}
	}

	// Index describes structure of the SQL index
	Index struct {
		TableIdent string
		Ident      string
		Type       string
		Fields     []*IndexField
		Unique     bool
		Predicate  string
		Comment    string

		// implementation variations
		Meta map[string]interface{}
	}

	IndexFieldSorted int

	// IndexField describes a single field (column or expression) of the SQL index
	IndexField struct {
		// Expression or a single column
		Column string

		Length int

		// Wrap part in parentheses
		Expression string

		// Ascending or descending
		Sorted IndexFieldSorted

		Statistics *IndexFieldStatistics

		// implementation variations
		Meta map[string]interface{}
	}

	IndexFieldStatistics struct {
		// Cardinality is an indicator that refers to the uniqueness
		// of all values in a column. Low cardinality means a lot
		// of duplicate values in that column. For example, a column
		// that stores the gender values has low cardinality.
		// In contrast, high cardinality means that there are many distinct values.
		Cardinality int64
	}
)

const (
	PRIMARY_KEY = "PRIMARY"

	IndexFieldSortDesc = -1
	IndexFieldUnsorted = 0
	IndexFieldSortAsc  = 1
)

func (t *Table) ColumnByIdent(i string) *Column {
	for _, c := range t.Columns {
		if c.Ident == i {
			return c
		}
	}

	return nil
}

// ConvertModel is generic model converter
func ConvertModel(m *dal.Model, d driverDialect) (t *Table, err error) {
	t = &Table{Ident: m.Ident}
	for _, a := range m.Attributes {
		if a.Type == nil {
			continue
		}

		col := &Column{
			Ident:   a.StoreIdent(),
			Comment: a.Label,
		}

		col.Type, err = d.NativeColumnType(a.Type)
		if err != nil {
			return nil, fmt.Errorf("could not convert attribute %q to column: %w", a.Ident, err)
		}

		t.Columns = append(t.Columns, col)
	}

	// @todo indexes

	return
}
