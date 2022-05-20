package drivers

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/doug-martin/goqu/v9/exp"
)

type (
	// TableCodec is an RDBMS representation of data.Model structure and its arguments
	TableCodec interface {
		Columns() []Column
		Ident() exp.IdentifierExpression
		MakeScanBuffer() []any
		Encode(r dal.ValueGetter) (_ []any, err error)
		Decode(buf []any, r dal.ValueSetter) (err error)
		AttributeExpression(string) (exp.LiteralExpression, error)
	}

	// GenericTableCodec is a generic implementation of TableCodec
	GenericTableCodec struct {
		// table identifier (name)
		ident exp.IdentifierExpression

		// all columns we're selecting from when
		// we're selecting from all columns
		columns []Column

		model   *dal.Model
		dialect Dialect
	}
)

var (
	_ TableCodec = &GenericTableCodec{}
)

func NewTableCodec(m *dal.Model, d Dialect) *GenericTableCodec {
	gtc := &GenericTableCodec{
		dialect: d,
		model:   m,
		ident:   exp.NewIdentifierExpression("", m.Ident, ""),
	}

	var (
		colIdent string
		attr     *dal.Attribute
		done     = make(map[string]bool)
		cols     = make([]Column, 0, len(m.Attributes))
	)

	for a := range m.Attributes {
		attr = m.Attributes[a]
		colIdent = attrColumnIdent(attr)

		if done[colIdent] {
			continue
		}

		switch attr.Store.(type) {
		case *dal.CodecRecordValueSetJSON:
			// when dealing with encoded types there is probably
			// a different column that can handle the encoded payload
			cols = append(cols, &SimpleJsonDocColumn{
				name:       colIdent,
				attributes: collectStdRecordValueJSONColumns(colIdent, m.Attributes...),
			})
		default:
			cols = append(cols, NewSingleValueColumn(d, attr))
		}

		done[colIdent] = true
	}

	gtc.columns = cols

	return gtc
}

func (t *GenericTableCodec) Ident() exp.IdentifierExpression {
	return t.ident
}

func (t *GenericTableCodec) Columns() []Column {
	return t.columns
}

func (t *GenericTableCodec) MakeScanBuffer() []any {
	out := make([]any, len(t.columns))

	for c := range t.columns {
		out[c] = t.columns[c].Type().MakeScanBuffer()
	}

	return out
}

func (t *GenericTableCodec) Encode(r dal.ValueGetter) (_ []any, err error) {
	enc := make([]any, len(t.columns))

	for c := range t.columns {
		enc[c], err = t.columns[c].Encode(r)
		if err != nil {
			return
		}
	}

	return enc, nil
}

func (t *GenericTableCodec) Decode(buf []any, r dal.ValueSetter) (err error) {
	if len(buf) != len(t.columns) {
		return fmt.Errorf("columns incompatilbe with scan buffer")
	}

	for i, c := range t.columns {
		if err = c.Decode(buf[i], r); err != nil {
			return err
		}
	}

	return
}

func (t *GenericTableCodec) AttributeExpression(ident string) (exp.LiteralExpression, error) {
	attr := t.model.Attributes.FindByIdent(ident)

	if attr == nil {
		return nil, fmt.Errorf("unknown attribute %q", ident)
	}

	switch s := attr.Store.(type) {
	case *dal.CodecAlias:
		// using column directly
		return exp.NewLiteralExpression("?", exp.NewIdentifierExpression("", t.model.Ident, s.Ident)), nil

	case *dal.CodecRecordValueSetJSON:
		// using JSON to handle embedded values
		lit, err := t.dialect.DeepIdentJSON(exp.NewIdentifierExpression("", t.model.Ident, s.Ident), attr.Ident, 0)
		if err != nil {
			return nil, err
		}

		return t.dialect.AttributeCast(attr, lit)
	}

	return exp.NewLiteralExpression("?", exp.NewIdentifierExpression("", t.model.Ident, ident)), nil
}

func attrColumnIdent(att *dal.Attribute) string {
	switch ss := att.Store.(type) {
	case *dal.CodecRecordValueSetJSON:
		return ss.Ident

	case *dal.CodecAlias:
		return ss.Ident

	default:
		return att.Ident
	}
}

func collectStdRecordValueJSONColumns(ident string, aa ...*dal.Attribute) []*dal.Attribute {
	filtered := make([]*dal.Attribute, 0)
	for _, a := range aa {
		storeType, is := a.Store.(*dal.CodecRecordValueSetJSON)
		if !is {
			continue
		}

		if ident != storeType.Ident {
			continue
		}

		filtered = append(filtered, a)
	}

	return filtered
}
