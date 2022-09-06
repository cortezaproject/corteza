package drivers

// dialect.go
//
// Generic SQL functions used by majority of RDBMS drivers

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/dal"
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
func DeepIdentJSON(ident exp.IdentifierExpression, pp ...any) exp.LiteralExpression {
	var (
		sql  strings.Builder
		last = len(pp) - 1
	)

	sql.WriteString("?")

	for i, p := range pp {
		sql.WriteString("-")

		sql.WriteString(">")
		if i == last {
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

	return exp.NewLiteralExpression(sql.String(), ident)
}

// JsonPath constructs json-path string from the slice of path parts.
//
func JsonPath(pp ...any) string {
	var (
		path strings.Builder
	)

	path.WriteString("$")
	for i := range pp {
		switch part := pp[i].(type) {
		case string:
			path.WriteString(".")
			path.WriteString(part)
		case int:
			path.WriteString("[")
			path.WriteString(strconv.Itoa(i))
			path.WriteString("]")
		default:
			panic(fmt.Errorf("JsonPath expect string or int, got %T", i))
		}
	}

	return path.String()
}

func IndexFieldModifiers(attr *dal.Attribute, quoteIdent func(i string) string, mm ...dal.IndexFieldModifier) (string, error) {
	var (
		out = quoteIdent(attr.StoreIdent())
	)

	for _, m := range mm {
		out = fmt.Sprintf("%s(%s)", m, out)
	}

	return out, nil
}
