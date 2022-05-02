package drivers

import (
	"github.com/cortezaproject/corteza-server/pkg/data"
	"github.com/doug-martin/goqu/v9/exp"
)

var (
	CheckID          = exp.NewLiteralExpression(`'^[0-9]+$'`)
	CheckNumber      = exp.NewLiteralExpression(`'^[0-9]+(\.[0-9])*$'`)
	CheckFullISO8061 = exp.NewLiteralExpression(`'^([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2}(?:\.[0-9]*)?)((-([0-9]{2}):([0-9]{2})|Z)?)$'`)
	CheckDateISO8061 = exp.NewLiteralExpression(`'^([0-9]{4})-([0-9]{2})-([0-9]{2})(T([0-9]{2}):([0-9]{2}):([0-9]{2}(?:\.[0-9]*)?)((-([0-9]{2}):([0-9]{2})|Z)?))?$'`)
	CheckTimeISO8061 = exp.NewLiteralExpression(`'^(([0-9]{4})-([0-9]{2})-([0-9]{2}))?T([0-9]{2}):([0-9]{2}):([0-9]{2}(?:\.[0-9]*)?)((-([0-9]{2}):([0-9]{2})|Z)?)$'`)
	LiteralNULL      = exp.NewLiteralExpression(`NULL`)
	LiteralFALSE     = exp.NewLiteralExpression(`FALSE`)
	LiteralTRUE      = exp.NewLiteralExpression(`TRUE`)
)

func AttributeCast(attr *data.Attribute, val exp.LiteralExpression) (exp.LiteralExpression, error) {
	var (
		c exp.CastExpression
	)

	switch attr.Type.(type) {
	case *data.TypeID, *data.TypeRef:
		ce := exp.NewCaseExpression().
			When(val.RegexpLike(CheckID), val).
			Else(LiteralNULL)

		c = exp.NewCastExpression(ce, "BIGINT")

	case *data.TypeNumber:
		ce := exp.NewCaseExpression().
			When(val.RegexpLike(CheckNumber), val).
			Else(LiteralNULL)

		c = exp.NewCastExpression(ce, "NUMERIC")

	case *data.TypeTimestamp:
		ce := exp.NewCaseExpression().
			When(val.RegexpLike(CheckFullISO8061), val).
			Else(LiteralNULL)

		c = exp.NewCastExpression(ce, "TIMESTAMPTZ")

	case *data.TypeDate:
		ce := exp.NewCaseExpression().
			When(val.RegexpLike(CheckDateISO8061), val).
			Else(LiteralNULL)

		c = exp.NewCastExpression(ce, "DATE")

	case *data.TypeTime:
		ce := exp.NewCaseExpression().
			When(val.RegexpLike(CheckTimeISO8061), val).
			Else(LiteralNULL)

		c = exp.NewCastExpression(ce, "TIMETZ")

	case *data.TypeBoolean:
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
