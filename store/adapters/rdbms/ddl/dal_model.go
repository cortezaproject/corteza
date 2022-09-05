package ddl

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/dal"
)

// CreateModel creates table with columns and indexes that match Model definition
func CreateModel(ctx context.Context, dd DataDefiner, m *dal.Model) (err error) {
	var (
		t *Table
	)

	if t, err = dd.ConvertModel(m); err != nil {
		return err
	}

	if err = dd.TableCreate(ctx, t); err != nil {
		return err
	}

	if err = EnsureIndexes(ctx, dd, m.Indexes...); err != nil {
		return
	}

	return
}

func DeleteModel(ctx context.Context, dd DataDefiner, m *dal.Model) (err error) {
	return fmt.Errorf("not implemented")
}

// UpdateModel alters existing table's columns and indexes to match Model definition
func UpdateModel(ctx context.Context, dd DataDefiner, m *dal.Model) (err error) {
	var (
		t *Table
	)

	if t, err = dd.TableLookup(ctx, m.Ident); err != nil {
		return err
	}

	// @todo check model against table structure
	for _, attr := range m.Attributes {
		// iterate over attributes and check if they exist in the table
		col := t.ColumnByIdent(attr.Ident)
		if col != nil {
			// @todo check if column type matches
			// @todo check if column is nullable
			break
		}

		// @todo add column to table
		return fmt.Errorf("column %q on table %q not found; adding columns is not jet supported", attr.Ident, m.Ident)
	}

	if err = EnsureIndexes(ctx, dd, m.Indexes...); err != nil {
		return
	}

	return
}

func EnsureIndexes(ctx context.Context, dd DataDefiner, ii ...*dal.Index) (err error) {
	return nil
}
