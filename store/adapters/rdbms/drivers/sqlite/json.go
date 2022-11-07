package sqlite

import (
	"fmt"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/valyala/fastjson"
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

func sqliteFuncJsonArrayContains(needle, haystack []byte) (_ bool, err error) {
	var n, h, i *fastjson.Value

	if h, err = fastjson.ParseBytes(haystack); err != nil {
		return
	}

	if h.Type() != fastjson.TypeArray {
		err = fmt.Errorf("haystack is %s, expecting array", h.Type())
		return
	}

	if n, err = fastjson.ParseBytes(needle); err != nil {
		return
	}

	for _, i = range h.GetArray() {
		if i.Type() != n.Type() {
			continue
		}

		switch i.Type() {
		case fastjson.TypeFalse, fastjson.TypeTrue:
			if i.GetBool() == n.GetBool() {
				return true, nil
			}
		case fastjson.TypeNumber, fastjson.TypeString:
			if i.String() == n.String() {
				return true, nil
			}
		}
	}

	return
}
