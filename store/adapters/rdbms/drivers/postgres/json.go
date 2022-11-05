package postgres

import (
	"github.com/doug-martin/goqu/v9/exp"
	"strconv"
	"strings"
)

// DeepIdentJSON constructs expression with chain of JSON operators
// that point value inside JSON document
//
// https://www.postgresql.org/docs/9.3/functions-json.html
func DeepIdentJSON(asJSON bool, jsonDoc exp.Expression, pp ...any) exp.LiteralExpression {
	var (
		sql  strings.Builder
		last = len(pp) - 1
	)

	sql.WriteString("?")

	for i, p := range pp {
		sql.WriteString("-")

		sql.WriteString(">")
		if i == last && !asJSON {
			sql.WriteString(">")
		}

		switch path := p.(type) {
		case string:
			sql.WriteString("'")
			sql.WriteString(strings.ReplaceAll(path, "'", `\'`))
			sql.WriteString("'")
		case int:
			sql.WriteString(strconv.Itoa(path))
		default:
			panic("invalid type")
		}
	}

	return exp.NewLiteralExpression(sql.String(), jsonDoc)
}
