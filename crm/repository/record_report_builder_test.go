package repository

import (
	"testing"
)

func TestRecordReportBuilder_parseExpression(t *testing.T) {
	b := recordReportBuilder{jsonField: "JSONFIELD"}

	tc := []struct {
		exp string
		sql string
		arg []interface{}
		err error
	}{
		{exp: "count(foo)", sql: "COUNT(JSONFIELD)", arg: []interface{}{"foo"}},
		{exp: "sum(count(foo))", sql: "SUM(COUNT(JSONFIELD))", arg: []interface{}{"foo"}},
		{exp: "sum( count( foo))  ", sql: "SUM(COUNT(JSONFIELD))", arg: []interface{}{"foo"}},
	}

	for _, c := range tc {
		sql, arg, err := b.parseExpression(c.exp).ToSql()
		assert(t, sql == c.sql, "Expecting expression SQL to match (%v == %v)", sql, c.sql)
		assert(t, len(arg) == len(c.arg), "Expecting arguments count to match (%v == %v)", arg, c.arg)
		assert(t, err == c.err, "Expecting errors to match (%v == %v)", err, c.err)
	}
}
