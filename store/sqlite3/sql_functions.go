package sqlite3

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/ql"
	"github.com/cortezaproject/corteza-server/pkg/qlng"
	"github.com/cortezaproject/corteza-server/store/rdbms"
)

var (
	sqlExprRegistry = map[string]rdbms.HandlerSig{
		// functions
		// - filtering
		"now": func(aa ...rdbms.FormattedASTArgs) (out string, args []interface{}, selfEnclosed bool, err error) {
			if len(aa) != 0 {
				err = fmt.Errorf("expecting 0 arguments, got %d", len(aa))
				return
			}

			out = "DATE('now')"
			return
		},
		"quarter": func(aa ...rdbms.FormattedASTArgs) (out string, args []interface{}, selfEnclosed bool, err error) {
			if len(aa) != 1 {
				err = fmt.Errorf("expecting 1 arguments, got %d", len(aa))
				return
			}

			out = fmt.Sprintf("(CAST(STRFTIME('%%m', %s) AS INTEGER) + 2) / 3", aa[0].S)
			args = aa[0].Args
			return
		},
		"year": func(aa ...rdbms.FormattedASTArgs) (out string, args []interface{}, selfEnclosed bool, err error) {
			if len(aa) != 1 {
				err = fmt.Errorf("expecting 1 arguments, got %d", len(aa))
				return
			}

			out = fmt.Sprintf("STRFTIME('%%Y', %s)", aa[0].S)
			args = aa[0].Args
			return
		},
		"month": func(aa ...rdbms.FormattedASTArgs) (out string, args []interface{}, selfEnclosed bool, err error) {
			if len(aa) != 1 {
				err = fmt.Errorf("expecting 1 arguments, got %d", len(aa))
				return
			}

			out = fmt.Sprintf("STRFTIME('%%m', %s)", aa[0].S)
			args = aa[0].Args
			return
		},
		"date": func(aa ...rdbms.FormattedASTArgs) (out string, args []interface{}, selfEnclosed bool, err error) {
			if len(aa) != 1 {
				err = fmt.Errorf("expecting 1 arguments, got %d", len(aa))
				return
			}

			out = fmt.Sprintf("STRFTIME('%%d', %s)", aa[0].S)
			args = aa[0].Args
			return
		},

		// - strings
		"concat": func(aa ...rdbms.FormattedASTArgs) (out string, args []interface{}, selfEnclosed bool, err error) {
			selfEnclosed = true

			params := make([]string, len(aa))
			for i, a := range aa {
				params[i] = a.S
				args = append(args, a.Args...)
			}

			out = fmt.Sprintf("(%s)", strings.Join(params, "||"))
			return
		},

		// - typecast
		"float": func(aa ...rdbms.FormattedASTArgs) (out string, args []interface{}, selfEnclosed bool, err error) {
			selfEnclosed = true

			if len(aa) != 1 {
				err = fmt.Errorf("expecting 1 argument, got %d", len(aa))
				return
			}

			out = fmt.Sprintf("CAST(%s AS FLOAT)", aa[0].S)
			args = aa[0].Args
			return
		},
		"string": func(aa ...rdbms.FormattedASTArgs) (out string, args []interface{}, selfEnclosed bool, err error) {
			selfEnclosed = true

			if len(aa) != 1 {
				err = fmt.Errorf("expecting 1 argument, got %d", len(aa))
				return
			}

			out = fmt.Sprintf("CAST(%s AS TEXT)", aa[0].S)
			args = aa[0].Args
			return
		},
	}
)

func sqlASTFormatter(n *qlng.ASTNode) rdbms.HandlerSig {
	return sqlExprRegistry[n.Ref]
}

func sqlFunctionHandler(f ql.Function) (ql.ASTNode, error) {
	switch strings.ToUpper(f.Name) {
	case "QUARTER":
		return ql.MakeFormattedNode("(CAST(STRFTIME('%%m', %s) AS INTEGER) + 2) / 3", f.Arguments...), nil
	case "YEAR":
		return ql.MakeFormattedNode("STRFTIME('%%Y', %s)", f.Arguments...), nil
	case "NOW":
		return ql.MakeFormattedNode("DATE('now')", f.Arguments...), nil
	case "DATE_FORMAT":
		if len(f.Arguments) != 2 {
			return nil, fmt.Errorf("expecting exactly two arguments for DATE_FORMAT function")
		}
		return ql.MakeFormattedNode("STRFTIME('%s', %s)", f.Arguments[0], f.Arguments[1]), nil
	case "DATE":
		// need to convert back to datetime so it can be converted to time.Time
		return ql.MakeFormattedNode("STRFTIME('%%Y-%%m-%%dT00:00:00Z', %s)", f.Arguments...), nil
	case "DATE_ADD", "DATE_SUB", "STD":
		return nil, fmt.Errorf("%q function is currently unsupported in SQLite store backend", f.Name)
	}

	return f, nil
}
