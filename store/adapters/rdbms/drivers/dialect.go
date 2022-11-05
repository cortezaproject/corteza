package drivers

import (
	"fmt"
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

		// JsonExtract returns expression that returns a value from  inside JSON document
		//
		// Use this when you want use JSON encoded value
		JsonExtract(exp.Expression, ...any) (exp.Expression, error)

		// JsonExtractUnquote returns expression that returns a value from  inside JSON document:
		//
		// Use this when you want to use unencoded value!
		JsonExtractUnquote(exp.Expression, ...any) (exp.Expression, error)

		// JsonArrayContains generates expression JSON array containment check expression
		//
		// Literal values need to be JSON docs!
		JsonArrayContains(needle, haystack exp.Expression) (exp.Expression, error)

		// AttributeCast prepares complex SQL expression that verifies
		// arbitrary string value in the db and casts it to b used in
		// comparison or soring expression
		AttributeCast(*dal.Attribute, exp.Expression) (exp.Expression, error)

		// TableCodec returns table codec (encodes & decodes data to/from db table)
		TableCodec(*dal.Model) TableCodec

		// TypeWrap returns driver's type implementation for a particular attribute type
		TypeWrap(dal.Type) Type

		QuoteIdent(string) string

		// AttributeToColumn converts attribute to column defunition
		AttributeToColumn(*dal.Attribute) (*ddl.Column, error)

		// ExprHandler returns driver specific expression handling
		ExprHandler(*ql.ASTNode, ...exp.Expression) (exp.Expression, error)
	}
)

func init() {
	goqu.SetDefaultPrepared(true)
}

func IndexFieldModifiers(attr *dal.Attribute, quoteIdent func(i string) string, mm ...dal.IndexFieldModifier) (string, error) {
	var (
		modifier string
		out      = quoteIdent(attr.StoreIdent())
	)

	for _, m := range mm {
		switch m {
		case dal.IndexFieldModifierLower:
			modifier = "LOWER"
		default:
			return "", fmt.Errorf("unknown index field modifier: %s", m)
		}

		out = fmt.Sprintf("%s(%s)", modifier, out)
	}

	return out, nil
}
