package drivers

import (
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/ql"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type (
	Dialect interface {
		// GOQU returns goqu's dialect wrapper struct
		GOQU() goqu.DialectWrapper

		// DeepIdentJSON returns expression that allows us (read) access to a particular
		// value inside JSON document:
		//
		// DeepIdentJSON(exp.ParseExpression("some_column"), "a", "b")
		// should result in something like:
		// "some_column"->'a'->>'b'
		// (depending on what is supported in the underlying database)
		DeepIdentJSON(exp.IdentifierExpression, ...any) (exp.LiteralExpression, error)

		// AttributeCast prepares complex SQL expression that verifies
		// arbitrary string value in the db and casts it to b used in
		// comparison or soring expression
		AttributeCast(*dal.Attribute, exp.LiteralExpression) (exp.LiteralExpression, error)

		// TableCodec returns table codec (encodes & decodes data to/from db table)
		TableCodec(*dal.Model) TableCodec

		// TypeWrap returns driver's type implementation for a particular attribute type
		TypeWrap(dal.Type) Type

		// NativeColumnType converts column type to type that can be used in the underlying rdbms
		NativeColumnType(columnType dal.Type) (*ddl.ColumnType, error)

		// ExprHandler returns driver specific expression handling
		ExprHandler(*ql.ASTNode, ...exp.Expression) (exp.Expression, error)
	}
)
