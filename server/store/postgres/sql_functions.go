package postgres

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
		"quarter": makeGenericExtrFncHandler("QUARTER"),
		"year":    makeGenericExtrFncHandler("YEAR"),
		"month":   makeGenericExtrFncHandler("MONTH"),
		"date":    makeGenericExtrFncHandler("DAY"),
	}
)

func makeGenericExtrFncHandler(extr string) rdbms.HandlerSig {
	return func(aa ...rdbms.FormattedASTArgs) (out string, args []interface{}, selfEnclosed bool, err error) {
		if len(aa) != 1 {
			err = fmt.Errorf("expecting 1 arguments, got %d", len(aa))
			return
		}

		out = fmt.Sprintf("EXTRACT(%s FROM %s)", extr, aa[0].S)
		args = aa[0].Args
		return
	}
}

func sqlASTFormatter(n *qlng.ASTNode) rdbms.HandlerSig {
	return sqlExprRegistry[n.Ref]
}

func sqlFunctionHandler(f ql.Function) (ql.ASTNode, error) {
	switch strings.ToUpper(f.Name) {
	case "QUARTER", "YEAR":
		return ql.MakeFormattedNode(fmt.Sprintf("EXTRACT(%s FROM %%s::date)", f.Name), f.Arguments...), nil
	case "DATE_FORMAT":
		return ql.MakeReplacedFormattedNode("TO_CHAR(%s::TIMESTAMPTZ, %s::TEXT)", translateDateFormatParams, f.Arguments...), nil
	case "TIMESTAMP":
		return ql.MakeFormattedNode("TIMESTAMPTZ(%s::TIMESTAMPTZ)", f.Arguments...), nil
	case "DATE":
		return ql.MakeFormattedNode("%s::DATE", f.Arguments...), nil
	case "TIME":
		return ql.MakeFormattedNode("DATE_TRUNC('second', %s::TIME)::TIME", f.Arguments...), nil
	case "DATE_ADD", "DATE_SUB", "STD":
		return nil, fmt.Errorf("%q function is currently unsupported in PostgreSQL store backend", f.Name)
	}

	return f, nil
}

func translateDateFormatParams(format string) string {
	return strings.NewReplacer(
		// @todo Doing ...%dT%H... (for iso tomestamp) pgsql doesn't format it correctly
		// so I'm covering this edgecase.
		// We should fix this properly when we redo record storage.
		`%dT%H`, `DD"T"HH24`,

		`%a`, `Dy`,
		`%b`, `Mon`,
		`%c`, `FMMM`,
		`%d`, `DD`,
		`%e`, `FMDD`,
		`%f`, `US`,
		`%H`, `HH24`,
		`%h`, `HH12`,
		`%I`, `HH12`,
		`%i`, `MI`,
		`%j`, `DDD`,
		`%k`, `FMHH24`,
		`%l`, `FMHH12`,
		`%M`, `FMMonth`,
		`%m`, `MM`,
		`%p`, `AM`,
		`%r`, `HH12:MI:SS AM`,
		`%S`, `SS`,
		`%s`, `SS`,
		`%T`, `HH24:MI:SS`,
		`%W`, `FMDay`,
		`%Y`, `YYYY`,
		`%y`, `YY`,
		`%%`, `%`,
	).Replace(format)
}
