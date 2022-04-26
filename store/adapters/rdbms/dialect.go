package rdbms

// dialect.go
//
// Generic SQL functions used by majority of RDBMS drivers

import (
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9/exp"
)

// DeepIdentJSON constructs expression with chain of JSON operators
// that point value inside JSON document
//
// Supported in databases:
//
// PostgreSQL
// https://www.postgresql.org/docs/9.3/functions-json.html
//
// MySQL
// https://dev.mysql.com/doc/expman/5.7/en/json-function-experence.html
//
// SQLite
// https://www.sqlite.org/json1.html#jptr
func DeepIdentJSON(ident string, pp ...any) exp.LiteralExpression {
	var (
		sql  = "?" + strings.Repeat("->?", len(pp))
		args = []any{exp.ParseIdentifier(ident)}
	)

	for _, p := range pp {
		switch p.(type) {
		case string, int:
			args = append(args, exp.NewLiteralExpression("?", p))
		default:
			panic("invalid type")
		}
	}

	return exp.NewLiteralExpression(sql, args...)
}

// JsonPath constructs json-path string from the slice of path parts.
//
func JsonPath(pp ...any) string {
	var (
		path = "$"
	)

	for i := range pp {
		switch part := pp[i].(type) {
		case string:
			path = path + "." + part
		case int:
			path = path + fmt.Sprintf("[%d]", pp[i])
		default:
			panic(fmt.Errorf("JsonPath expect string or int, got %T", i))
		}
	}

	return path
}
