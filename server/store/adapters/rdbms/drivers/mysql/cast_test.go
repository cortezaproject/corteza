package mysql

import (
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/stretchr/testify/require"
)

// A cast checks its value with REGEXP alone; REGEXP BINARY fails on MySQL 8
// with utf8mb4 columns (error 3995)
func TestAttributeCastRegexpWithoutBinary(t *testing.T) {
	req := require.New(t)
	val := exp.NewIdentifierExpression("", "", "v")

	for _, tt := range []dal.Type{&dal.TypeNumber{}, &dal.TypeTimestamp{}, &dal.TypeDate{}, &dal.TypeTime{}, &dal.TypeID{}} {
		e, err := mysqlDialect{}.AttributeCast(&dal.Attribute{Ident: "v", Type: tt}, val)
		req.NoError(err)

		sql, _, err := goquDialectWrapper.Select(e).ToSQL()
		req.NoError(err)
		req.Contains(sql, "REGEXP", tt.Type())
		req.NotContains(sql, "BINARY", tt.Type())
	}
}
