package drivers

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/crs"
	"github.com/cortezaproject/corteza-server/pkg/data"
	"github.com/doug-martin/goqu/v9/exp"
)

type (
	// TableCodec is an RDBMS representation of data.Model structure and its arguments
	TableCodec interface {
		Columns() []Column
		Ident() exp.IdentifierExpression
		MakeScanBuffer() []any
		Encode(r crs.ValueGetter) (_ []any, err error)
		Decode(buf []any, r crs.ValueSetter) (err error)
		AttributeExpression(string) (exp.LiteralExpression, error)
	}

	// GenericTableCodec is a generic implementation of TableCodec
	GenericTableCodec struct {
		// table identifier (name)
		ident exp.IdentifierExpression

		// all columns we're selecting from when
		// we're selecting from all columns
		columns []Column

		model   *data.Model
		dialect Dialect
	}
)

var (
	_ TableCodec = &GenericTableCodec{}
)

func NewTableCodec(m *data.Model, d Dialect) *GenericTableCodec {
	gtc := &GenericTableCodec{
		dialect: d,
		model:   m,
		ident:   exp.NewIdentifierExpression("", m.Ident, ""),
	}

	var (
		colIdent string
		att      *data.Attribute
		done     = make(map[string]bool)
		cols     = make([]Column, 0, len(m.Attributes))
	)

	for a := range m.Attributes {
		att = m.Attributes[a]
		colIdent = attrColumnIdent(att)

		if done[colIdent] {
			continue
		}

		if data.StoreCodecStdRecordValueJSONType.Is(att.Store) {
			// when dealing with encoded types there is probably
			// a different column that can handle the encoded payload
			cols = append(cols, &SimpleJsonDocColumn{
				name:       colIdent,
				attributes: collectStdRecordValueJSONColumns(colIdent, m.Attributes...),
			})
		} else {
			cols = append(cols, NewSingleValueColumn(d, att))
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

func (t *GenericTableCodec) Encode(r crs.ValueGetter) (_ []any, err error) {
	enc := make([]any, len(t.columns))

	for c := range t.columns {
		enc[c], err = t.columns[c].Encode(r)
		if err != nil {
			return
		}
	}

	return enc, nil
}

func (t *GenericTableCodec) Decode(buf []any, r crs.ValueSetter) (err error) {
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
	case *data.StoreCodecAlias:
		// using column directly
		return exp.NewLiteralExpression("?", exp.NewIdentifierExpression("", t.model.Ident, s.Ident)), nil

	case *data.StoreCodecStdRecordValueJSON:
		// using JSON to handle embedded values
		lit, err := t.dialect.DeepIdentJSON(exp.NewIdentifierExpression("", t.model.Ident, s.Ident), attr.Ident, 0)
		if err != nil {
			return nil, err
		}

		return t.dialect.AttributeCast(attr, lit)
	}

	return exp.NewLiteralExpression("?", exp.NewIdentifierExpression("", t.model.Ident, ident)), nil
}

func attrColumnIdent(att *data.Attribute) string {
	switch ss := att.Store.(type) {
	case *data.StoreCodecStdRecordValueJSON:
		return ss.Ident

	case *data.StoreCodecAlias:
		return ss.Ident

	default:
		return att.Ident
	}
}

func collectStdRecordValueJSONColumns(ident string, aa ...*data.Attribute) []*data.Attribute {
	filtered := make([]*data.Attribute, 0)
	for _, a := range aa {
		storeType, is := a.Store.(*data.StoreCodecStdRecordValueJSON)
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
