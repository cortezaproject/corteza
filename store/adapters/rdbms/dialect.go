package rdbms

// dialect.go
//
// Generic SQL functions used by majority of RDBMS drivers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

func init() {
	goqu.SetDefaultPrepared(true)
}

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
	// strings.ReplaceAll(str, "'", `\'`)

	var (
		//sql  = "?" + strings.Repeat("->?", len(pp))
		//args = []any{exp.ParseIdentifier(ident)}
		sql = strconv.Quote(ident)
	)

	for _, p := range pp {
		switch path := p.(type) {
		case string:
			sql += "->'" + strings.ReplaceAll(path, "'", `\'`) + "'"
		case int:
			sql += fmt.Sprintf("->%d", path)
		default:
			panic("invalid type")
		}
	}

	return exp.NewLiteralExpression(sql)
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
