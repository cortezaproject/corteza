package ddl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"strconv"
)

type (
	driverDialect interface {
		QuoteIdent(string) string
		// AttributeToColumn converts attribute to column
		AttributeToColumn(*dal.Attribute) (*Column, error)
		IndexFieldModifiers(*dal.Attribute, ...dal.IndexFieldModifier) (string, error)
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
		Ident string
		Type  *ColumnType

		Default string

		// implementation variations
		Meta map[string]interface{}

		Comment string
	}

	ColumnType struct {
		Name string

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

	// IndexField describes a single field (column or expression) of the SQL index
	IndexField struct {
		// Expression or a single column
		Column string

		Length int

		// Wrap part in parentheses
		Expression string

		// Ascending or descending
		Sort  dal.IndexFieldSort
		Nulls dal.IndexFieldNulls

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
	var (
		col *Column
		idx *Index

		// keeps track of embedded attributes
		// key is storeIdent
		embeds = make(map[string]bool)
	)

	t = &Table{Ident: m.Ident}
	for _, a := range m.Attributes {
		if a.Type == nil {
			continue
		}

		if embeds[a.StoreIdent()] {
			continue
		}

		switch a.Store.(type) {
		case *dal.CodecRecordValueSetJSON:
			// add to embeds and make sure we do not add the column again!
			embeds[a.StoreIdent()] = true

			// throw away the attribute and create a new one
			// to for JSON storage
			a = &dal.Attribute{
				Ident: a.StoreIdent(),
				Type:  &dal.TypeJSON{Nullable: false},
				Store: &dal.CodecPlain{},
			}
		}

		col, err = d.AttributeToColumn(a)
		if err != nil {
			return nil, fmt.Errorf("could not convert attribute %q to column: %w", a.Ident, err)
		}

		t.Columns = append(t.Columns, col)
	}

	for _, i := range m.Indexes {
		if idx, err = ConvertIndex(i, m.Attributes, m.Ident, d); err != nil {
			return nil, fmt.Errorf("could not convert index %q: %w", i.Ident, err)
		}

		t.Indexes = append(t.Indexes, idx)
	}

	return
}

// ConvertIndex converts dal.Index to ddl.Index
func ConvertIndex(i *dal.Index, aa dal.AttributeSet, table string, d driverDialect) (idx *Index, err error) {
	var (
		a        *dal.Attribute
		idxField *IndexField
	)

	idx = &Index{
		TableIdent: table,
		Ident:      i.Ident,
		Type:       i.Type,
		Unique:     i.Unique,
		Predicate:  i.Predicate,
	}

	for _, f := range i.Fields {
		// ensure attribute exists
		if a = aa.FindByIdent(f.AttributeIdent); a == nil {
			return nil, fmt.Errorf("referenced attribute %q does not exist", f.AttributeIdent)
		}

		idxField = &IndexField{
			Sort:  f.Sort,
			Nulls: f.Nulls,
		}

		if len(f.Modifiers) > 0 {
			if idxField.Expression, err = d.IndexFieldModifiers(a, f.Modifiers...); err != nil {
				return nil, fmt.Errorf("could not convert index field modifiers: %w", err)
			}
		} else {
			idxField.Column = a.StoreIdent()
		}

		idx.Fields = append(idx.Fields, idxField)
	}

	return
}

func DefaultValueCurrentTimestamp(set bool) string {
	if !set {
		return ""
	}

	return "CURRENT_TIMESTAMP"
}

func DefaultBoolean(set, value bool) string {
	switch {
	case !set:
		return ""
	case value:
		return "true"
	default:
		return "false"
	}
}

func DefaultID(set bool, value uint64) string {
	switch {
	case !set:
		return ""
	default:
		return strconv.FormatUint(value, 10)
	}
}

func DefaultNumber(set bool, precision int, value float64) string {
	switch {
	case !set:
		return ""
	case precision >= 0:
		return strconv.FormatFloat(value, 'f', precision, 64)
	default:
		return strconv.FormatFloat(value, 'f', -1, 64)
	}
}

func DefaultJSON(set bool, value any) (_ string, err error) {
	if str, is := value.(string); is {
		return "'" + str + "'", nil
	}

	switch {
	case !set:
		return "", nil
	default:
		var aux []byte
		if aux, err = json.Marshal(value); err != nil {
			return "", fmt.Errorf("could not serialize default value for JSON field: %w", err)
		}

		return string(aux), nil
	}
}
