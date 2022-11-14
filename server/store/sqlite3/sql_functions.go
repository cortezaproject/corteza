package sqlite3

import (
	"fmt"
	"regexp"
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

	supportedSubstitutions = map[string]bool{
		"d": true,
		"H": true,
		"j": true,
		"m": true,
		"M": true,
		"S": true,
		"w": true,
		"W": true,
		"Y": true,
		"%": true,
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
	case "TIMESTAMP", "DATETIME":
		return ql.MakeFormattedNode("DATETIME(%s)", f.Arguments...), nil
	case "DATE_FORMAT":
		if len(f.Arguments) != 2 {
			return nil, fmt.Errorf("expecting exactly two arguments for DATE_FORMAT function")
		}
		format := f.Arguments[1]
		col := f.Arguments[0]

		err := supportedDateFormatParams(format)
		if err != nil {
			return nil, err
		}

		return ql.MakeReplacedFormattedNode("STRFTIME(%s, %s)", translateDateFormatParams, format, col), nil
	case "DATE":
		// need to convert back to datetime so it can be converted to time.Time
		return ql.MakeFormattedNode("STRFTIME('%%Y-%%m-%%dT00:00:00Z', %s)", f.Arguments...), nil
	case "DATE_ADD", "DATE_SUB", "STD":
		return nil, fmt.Errorf("%q function is currently unsupported in SQLite store backend", f.Name)
	}

	return f, nil
}

func supportedDateFormatParams(fmtNode ql.ASTNode) error {
	format := translateDateFormatParams(fmtNode.String())

	r := regexp.MustCompile(`%(?P<sub>.)`)

	for _, m := range r.FindAllStringSubmatch(format, -1) {
		if len(m) == 0 {
			continue
		}

		if _, ok := supportedSubstitutions[m[1]]; !ok {
			return fmt.Errorf("format substitution not supported: %%%s", m[1])
		}
	}

	return nil
}

func translateDateFormatParams(format string) string {
	return strings.NewReplacer(
		`%i`, `%M`,
		`%U`, `%W`,
	).Replace(format)
}
