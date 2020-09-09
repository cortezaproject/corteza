package sqlite

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/store/rdbms"
)

func fieldToColumnTypeCaster(field rdbms.ModuleFieldTypeDetector, ident string) (string, error) {
	switch true {
	case field.IsBoolean():
		return fmt.Sprintf("(rv_%s.value NOT IN ('', '0', 'false', 'f',  'FALSE', 'F', false))", ident), nil
	case field.IsNumeric():
		return fmt.Sprintf("CAST(rv_%s.value AS SIGNED)", ident), nil
	case field.IsRef():
		return fmt.Sprintf("rv_%s.ref ", ident), nil
	default:
		return fmt.Sprintf("rv_%s.value ", ident), nil
	}

	return ident, nil
}
