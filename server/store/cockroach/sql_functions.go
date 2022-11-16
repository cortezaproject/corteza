package cockroach

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/ql"
)

func sqlFunctionHandler(f ql.Function) (ql.ASTNode, error) {
	switch strings.ToUpper(f.Name) {
	case "QUARTER", "YEAR":
		return ql.MakeFormattedNode(fmt.Sprintf("EXTRACT(%s FROM %%s::date)", f.Name), f.Arguments...), nil
	case "DATE_FORMAT":
		return ql.MakeFormattedNode("TO_CHAR(%s::TIMESTAMPTZ, %s::TEXT)", f.Arguments...), nil
	case "DATE":
		return ql.MakeFormattedNode("%s::DATE", f.Arguments...), nil
	case "DATE_ADD", "DATE_SUB", "STD":
		return nil, fmt.Errorf("%q function is currently unsupported in PostgreSQL store backend", f.Name)
	}

	return f, nil
}
