package crs

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/data"
	"github.com/cortezaproject/corteza-server/pkg/slice"
	"github.com/spf13/cast"
)

// Encode takes a record and encodes it according to the given model
// encoding record values into form acceptable by underlying storage
//
// with omit slice you specify list of attribute idents to be omitted!
func encode(columns []*column, r *types.Record, omit ...string) (enc map[string]any, err error) {
	enc = make(map[string]any)

	for _, col := range columns {
		if slice.ContainsAny(omit, col.attributes[0].Ident) {
			continue
		}

		if enc[col.ident], err = col.encode(col.attributes, r); err != nil {
			return nil, fmt.Errorf("could not encode value for attribute %q: %w", col.ident, err)
		}
	}

	return
}

// @todo this could probably be generalized, nothing special for RDBMS
func encodeStdRecordValueJSON(attributes []*data.Attribute, r *types.Record) (any, error) {
	var (
		aux   = make(map[string][]any)
		rvs   types.RecordValueSet
		value any
		place int
		size  int
	)

	for _, attr := range attributes {
		place = 0

		if isSystemField(attr.Ident) {
			//	for system fields just get the one
			val, err := getSystemFieldValue(r, attr.Ident)
			if err != nil {
				return nil, err
			}

			aux[attr.Ident] = []any{val}
		} else {
			// for the record values, collect all we
			// care about, sort and copy them to the auxiliary var.
			rvs = r.Values.FilterByName(attr.Ident)
			sort.Sort(rvs)

			if size = len(rvs); size == 0 {
				continue
			}

			aux[attr.Ident] = make([]any, size)
			for _, v := range rvs {
				aux[attr.Ident][place] = v.Value
			}
		}

		// now, encode the value according to JSON format constraints

		for place, value = range aux[attr.Ident] {
			switch attr.Type.(type) {
			case data.TypeBoolean:
				aux[attr.Ident][place] = cast.ToBool(value)

			default:
				// by default: as string
				aux[attr.Ident][place] = value
			}

			if !attr.MultiValue {
				// model attribute supports storing of single values only.
				break
			}
		}
	}

	return json.Marshal(aux)
}

// extracts value from a record with system fields as priority
func getSystemFieldValue(r *types.Record, name string) (any, error) {
	switch name {
	// handle system fields
	case sysID:
		return r.ID, nil
	case sysNamespaceID:
		return r.NamespaceID, nil
	case sysModuleID:
		return r.ModuleID, nil
	case sysCreatedAt:
		return r.CreatedAt, nil
	case sysCreatedBy:
		return r.CreatedBy, nil
	case sysUpdatedAt:
		return r.UpdatedAt, nil
	case sysUpdatedBy:
		return r.UpdatedBy, nil
	case sysDeletedAt:
		return r.DeletedAt, nil
	case sysDeletedBy:
		return r.DeletedBy, nil
	case sysOwnedBy:
		return r.OwnedBy, nil
	}

	return nil, fmt.Errorf("unknown compose record system field %q", name)
}
