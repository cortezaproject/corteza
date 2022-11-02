package drivers

import (
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/doug-martin/goqu/v9/exp"
)

var (
	CheckID          = exp.NewLiteralExpression(`'^[0-9]+$'`)
	CheckNumber      = exp.NewLiteralExpression(`'^[0-9]+(\.[0-9]+)?$'`)
	CheckFullISO8061 = exp.NewLiteralExpression(`'^([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2}(\.[0-9]*)?)((-([0-9]{2}):([0-9]{2})|Z)?)$'`)
	CheckDateISO8061 = exp.NewLiteralExpression(`'^([0-9]{4})-([0-9]{2})-([0-9]{2})(T([0-9]{2}):([0-9]{2}):([0-9]{2}(\.[0-9]*)?)((-([0-9]{2}):([0-9]{2})|Z)?))?$'`)
	CheckTimeISO8061 = exp.NewLiteralExpression(`'^(([0-9]{4})-([0-9]{2})-([0-9]{2}))?T?([0-9]{2}):([0-9]{2}):([0-9]{2}(\.[0-9]*)?)((-([0-9]{2}):([0-9]{2})|Z)?)$'`)
	LiteralNULL      = exp.NewLiteralExpression(`NULL`)
	LiteralFALSE     = exp.NewLiteralExpression(`FALSE`)
	LiteralTRUE      = exp.NewLiteralExpression(`TRUE`)
)

func RegexpLike(format, val exp.Expression) exp.BooleanExpression {
	return exp.NewBooleanExpression(exp.RegexpLikeOp, val, format)
}

func BooleanCheck(val exp.Expression) exp.Expression {
	return exp.NewCaseExpression().
		When(exp.NewBooleanExpression(exp.InOp, val, []any{LiteralTRUE, exp.NewLiteralExpression(`'true'`)}), LiteralTRUE).
		When(exp.NewBooleanExpression(exp.InOp, val, []any{LiteralFALSE, exp.NewLiteralExpression(`'false'`)}), LiteralFALSE).
		Else(LiteralNULL)
}

func AttributeCast(attr *dal.Attribute, val exp.Expression) (exp.Expression, error) {
	var (
		c exp.CastExpression
	)

	switch attr.Type.(type) {
	case *dal.TypeID, *dal.TypeRef:
		ce := exp.NewCaseExpression().
			When(RegexpLike(CheckID, val), val).
			Else(LiteralNULL)

		c = exp.NewCastExpression(ce, "BIGINT")

	case *dal.TypeNumber:
		ce := exp.NewCaseExpression().
			When(RegexpLike(CheckNumber, val), val).
			Else(LiteralNULL)

		c = exp.NewCastExpression(ce, "NUMERIC")

	case *dal.TypeTimestamp:
		ce := exp.NewCaseExpression().
			When(RegexpLike(CheckFullISO8061, val), val).
			Else(LiteralNULL)

		c = exp.NewCastExpression(ce, "TIMESTAMPTZ")

	case *dal.TypeDate:
		ce := exp.NewCaseExpression().
			When(RegexpLike(CheckDateISO8061, val), val).
			Else(LiteralNULL)

		c = exp.NewCastExpression(ce, "DATE")

	case *dal.TypeTime:
		ce := exp.NewCaseExpression().
			When(RegexpLike(CheckTimeISO8061, val), val).
			Else(LiteralNULL)

		c = exp.NewCastExpression(ce, "TIMETZ")

	case *dal.TypeBoolean:
		c = exp.NewCastExpression(BooleanCheck(val), "BOOLEAN")

	default:
		return val, nil
	}

	return exp.NewLiteralExpression("?", c), nil
}
