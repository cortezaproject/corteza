package postgres

import (
	"fmt"

	"github.com/cortezaproject/corteza/server/store/rdbms"
)

// fieldToColumnTypeCaster handles special ComposeModule field query representations
// @todo Not as elegant as it should be but it'll do the trick until the #2 store iteration
//
// Return parameters:
//   * full cast: query column + datatype cast
//   * field cast tpl: fmt template to get query column
//   * type cast tpl: fmt template to cast the compared to value
func fieldToColumnTypeCaster(field rdbms.ModuleFieldTypeDetector, ident string) (string, string, string, error) {
	fcp := "rv_%s.value"
	fcpRef := "rv_%s.ref"

	switch true {
	case field.IsBoolean():
		tcp := "%s NOT IN ('', '0', 'false', 'f',  'FALSE', 'F')"
		fc := fmt.Sprintf(fcp, ident)
		return fmt.Sprintf(tcp, fc), fcp, tcp, nil
	case field.IsNumeric():
		tcp := "CAST('0' || %s AS NUMERIC)"
		fc := fmt.Sprintf(fcp, ident)
		return fmt.Sprintf(tcp, fc), fcp, tcp, nil
	case field.IsDateTime():
		tcp := "to_timestamp(NULLIF(%s, ''),'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"') "
		fc := fmt.Sprintf(fcp, ident)
		return fmt.Sprintf(tcp, fc), fcp, tcp, nil
	case field.IsDateOnly():
		tcp := "CAST(NULLIF(%s, '') AS DATE)"
		fc := fmt.Sprintf(fcp, ident)
		return fmt.Sprintf(tcp, fc), fcp, tcp, nil
	case field.IsTimeOnly():
		tcp := "CAST(NULLIF(%s, '') AS TIME)"
		fc := fmt.Sprintf(fcp, ident)
		return fmt.Sprintf(tcp, fc), fcp, tcp, nil
	case field.IsRef():
		tcp := "%s"
		fc := fmt.Sprintf(fcpRef, ident)
		return fmt.Sprintf(tcp, fc), fcpRef, tcp, nil
	}

	tcp := "%s"
	fc := fmt.Sprintf(fcp, ident)
	return fmt.Sprintf(tcp, fc), fcp, tcp, nil
}
