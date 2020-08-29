package sqlite

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/ql"
	"strings"
)

func sqlFunctionHandler(f ql.Function) (ql.ASTNode, error) {
	switch strings.ToUpper(f.Name) {
	case "QUARTER":
		return ql.MakeFormattedNode("(CAST(STRFTIME('%%m', %s) AS INTEGER) + 2) / 3", f.Arguments...), nil
	case "YEAR":
		return ql.MakeFormattedNode("STRFTIME('%%y', %s)", f.Arguments...), nil
	case "NOW":
		return ql.MakeFormattedNode("DATE('now')", f.Arguments...), nil
	case "DATE_FORMAT":
		if len(f.Arguments) != 2 {
			return nil, fmt.Errorf("expecting exactly two arguments for DATE_FORMAT function")
		}

		return ql.MakeFormattedNode("STRFTIME('%s', %s)", f.Arguments[0], f.Arguments[1]), nil
	case "DATE_ADD", "DATE_SUB", "STD":
		return nil, fmt.Errorf("%q function is currently unsupported in SQLite store backend", f.Name)
	}

	return f, nil
}
