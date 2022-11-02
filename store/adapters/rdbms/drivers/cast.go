package drivers

import (
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/doug-martin/goqu/v9/exp"
)

var (
	CheckID          = exp.NewLiteralExpression(`'^[0-9]+$'`)
    CheckNumber      = exp.NewLiteralExpression(`'^[0-9]+(\.[0-9]+)?$'`)
	CheckFullISO8061 = exp.NewLiteralExpression(`'^([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2}(?:\.[0-9]*)?)((-([0-9]{2}):([0-9]{2})|Z)?)$'`)
	CheckDateISO8061 = exp.NewLiteralExpression(`'^([0-9]{4})-([0-9]{2})-([0-9]{2})(T([0-9]{2}):([0-9]{2}):([0-9]{2}(?:\.[0-9]*)?)((-([0-9]{2}):([0-9]{2})|Z)?))?$'`)
	CheckTimeISO8061 = exp.NewLiteralExpression(`'^(([0-9]{4})-([0-9]{2})-([0-9]{2}))?T([0-9]{2}):([0-9]{2}):([0-9]{2}(?:\.[0-9]*)?)((-([0-9]{2}):([0-9]{2})|Z)?)$'`)
	LiteralNULL      = exp.NewLiteralExpression(`NULL`)
	LiteralFALSE     = exp.NewLiteralExpression(`FALSE`)
	LiteralTRUE      = exp.NewLiteralExpression(`TRUE`)
)

func AttributeCast(attr *dal.Attribute, val exp.LiteralExpression) (exp.LiteralExpression, error) {
	var (
		c exp.CastExpression
	)

	switch attr.Type.(type) {
	case *dal.TypeID, *dal.TypeRef:
		ce := exp.NewCaseExpression().
			When(val.RegexpLike(CheckID), val).
			Else(LiteralNULL)

		c = exp.NewCastExpression(ce, "BIGINT")

	case *dal.TypeNumber:
		ce := exp.NewCaseExpression().
			When(val.RegexpLike(CheckNumber), val).
			Else(LiteralNULL)

		c = exp.NewCastExpression(ce, "NUMERIC")

	case *dal.TypeTimestamp:
		ce := exp.NewCaseExpression().
			When(val.RegexpLike(CheckFullISO8061), val).
			Else(LiteralNULL)

		c = exp.NewCastExpression(ce, "TIMESTAMPTZ")

	case *dal.TypeDate:
		ce := exp.NewCaseExpression().
			When(val.RegexpLike(CheckDateISO8061), val).
			Else(LiteralNULL)

		c = exp.NewCastExpression(ce, "DATE")

	case *dal.TypeTime:
		ce := exp.NewCaseExpression().
			When(val.RegexpLike(CheckTimeISO8061), val).
			Else(LiteralNULL)

		c = exp.NewCastExpression(ce, "TIMETZ")

	case *dal.TypeBoolean:
		ce := exp.NewCaseExpression().
			When(val.In(LiteralTRUE, exp.NewLiteralExpression(`'true'`)), LiteralTRUE).
			When(val.In(LiteralFALSE, exp.NewLiteralExpression(`'false'`)), LiteralFALSE).
			Else(LiteralNULL)

		c = exp.NewCastExpression(ce, "BOOLEAN")

	default:
		return val, nil
	}

	return exp.NewLiteralExpression("?", c), nil
}
