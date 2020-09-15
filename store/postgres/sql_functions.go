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
	}

	return f, nil
}
