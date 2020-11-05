package postgres

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/store/rdbms"
)

func fieldToColumnTypeCaster(field rdbms.ModuleFieldTypeDetector, ident string) (string, error) {
	switch true {
	case field.IsBoolean():
		return fmt.Sprintf("rv_%s.value NOT IN ('', '0', 'false', 'f',  'FALSE', 'F', false)", ident), nil
	case field.IsNumeric():
		return fmt.Sprintf("rv_%s.value::NUMERIC", ident), nil
	case field.IsDateTime():
		return fmt.Sprintf("rv_%s.value::TIMESTAMP", ident), nil
	case field.IsRef():
		return fmt.Sprintf("rv_%s.ref ", ident), nil
	default:
		return fmt.Sprintf("rv_%s.value ", ident), nil
	}
}
