package sqlite

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/ql"
	"github.com/cortezaproject/corteza-server/store/rdbms"
)

func fieldToColumnTypeCaster(field rdbms.ModuleFieldTypeDetector, i ql.Ident) (ql.Ident, error) {
	switch true {
	case field.IsBoolean():
		i.Value = fmt.Sprintf("(rv_%s.value NOT IN ('', '0', 'false', 'f',  'FALSE', 'F', false))", i.Value)
	case field.IsNumeric():
		i.Value = fmt.Sprintf("CAST(rv_%s.value AS SIGNED)", i.Value)
	case field.IsRef():
		i.Value = fmt.Sprintf("rv_%s.ref ", i.Value)
	default:
		i.Value = fmt.Sprintf("rv_%s.value ", i.Value)
	}

	return i, nil
}
