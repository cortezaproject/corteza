package postgres

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/store/rdbms"
)

func fieldToColumnTypeCaster(field rdbms.ModuleFieldTypeDetector, ident string) (string, error) {
	switch true {
	case field.IsBoolean():
		return fmt.Sprintf("CASE WHEN rv_%s.value NOT IN ('', '0', 'false', 'f',  'FALSE', 'F') THEN 1 ELSE 0 END ", ident), nil
	case field.IsNumeric():
		return fmt.Sprintf("CAST('0' || rv_%s.value AS NUMERIC) ", ident), nil
	case field.IsDateTime():
		return fmt.Sprintf(`to_timestamp(rv_%s.value,'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"') `, ident), nil
	case field.IsRef():
		return fmt.Sprintf("rv_%s.ref ", ident), nil
	default:
		return fmt.Sprintf("rv_%s.value ", ident), nil
	}
}
