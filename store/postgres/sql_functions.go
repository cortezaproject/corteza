package postgres

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/ql"
	"strings"
)

func sqlFunctionHandler(f ql.Function) (ql.ASTNode, error) {
	switch strings.ToUpper(f.Name) {
	case "QUARTER", "YEAR":
		return ql.MakeFormattedNode(fmt.Sprintf("EXTRACT(%s FROM %%s::date)", f.Name), f.Arguments...), nil
	case "DATE_FORMAT":
		return ql.MakeFormattedNode("TO_CHAR(%s, %s)", f.Arguments...), nil
	case "DATE":
		return ql.MakeFormattedNode("%s::DATE", f.Arguments...), nil
	case "DATE_ADD", "DATE_SUB", "STD":
		return nil, fmt.Errorf("%q function is currently unsupported in PostgreSQL store backend", f.Name)
	}

	return f, nil
}

func translateDateFormatParams(format string) string {
	return strings.NewReplacer(
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
