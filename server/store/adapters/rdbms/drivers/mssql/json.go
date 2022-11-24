package mssql

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9/exp"
)

func jsonPathExpr(pp ...any) (exp.LiteralExpression, error) {
	if path, err := jsonPath(pp...); err != nil {
		return nil, err
	} else {
		return exp.NewLiteralExpression("'" + path + "'"), nil
	}
}

func jsonPath(pp ...any) (string, error) {
	var (
		sql strings.Builder
	)

	sql.WriteString(`$`)

	for _, p := range pp {
		switch path := p.(type) {
		case string:
			sql.WriteString(".")
			sql.WriteString(strings.ReplaceAll(path, "'", `\'`))
		case int:
			sql.WriteString("[")
			sql.WriteString(strconv.Itoa(path))
			sql.WriteString("]")
		default:
			return "", fmt.Errorf("unexpected path part (%q) type: %T", p, p)
		}
	}

	return sql.String(), nil
}
