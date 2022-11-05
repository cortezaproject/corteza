package mysql

import (
	"fmt"
	"github.com/doug-martin/goqu/v9/exp"
	"strconv"
	"strings"
)

func JsonPath(asJSON bool, jsonDoc exp.Expression, pp ...any) (_ exp.LiteralExpression, err error) {
	var (
		sql  strings.Builder
		path string
	)

	if path, err = jsonPath(pp...); err != nil {
		return nil, err
	}

	sql.WriteString(`?->'`)
	if !asJSON {
		// interested in un-encoded value
		sql.WriteString(`>'`)
	}
	sql.WriteString(path)
	sql.WriteString(`'`)

	return exp.NewLiteralExpression(sql.String(), jsonDoc), nil
}

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

// JSONArrayContains prepares MySQL compatible comparison of value and JSON array
//
// # literal value = multi-value field / plain
// # multi-value field = single-value field / plain
// JSON_CONTAINS(v, JSON_EXTRACT(V, '$.f3'), '$.f2')
//
// # single-value field = multi-value field / plain
// # multi-value field = single-value field / plain
// JSON_CONTAINS(v, '"aaa"', '$.f2')
//
// This approach is not optimal, but it is the only way to make it work
func JSONArrayContains(left exp.Expression, ident exp.IdentifierExpression, pp ...any) (jc exp.Expression, err error) {
	var (
		pathAux string
		path    exp.Expression
	)

	// this is the least painful way how to encode
	// an unknown value as JSON
	left = exp.NewSQLFunctionExpression(
		"JSON_EXTRACT",
		exp.NewSQLFunctionExpression("JSON_ARRAY", left),
		exp.NewLiteralExpression(`'$[0]'`),
	)

	if pathAux, err = jsonPath(pp...); err != nil {
		return nil, err
	} else {
		path = exp.NewLiteralExpression(`'` + pathAux + `'`)
	}

	return exp.NewSQLFunctionExpression("JSON_CONTAINS", ident, left, path), nil
}
