package drivers

import (
	"encoding/json"

	"github.com/cortezaproject/corteza/server/pkg/dal"
)

type (
	// TableCodec is an RDBMS representation of data.Model structure and its arguments
	TableCodec interface {
		Columns() []Column
		Ident() string
		Encode(r dal.ValueGetter) (_ []byte, err error)
		Decode(buf []byte, r dal.ValueSetter) (err error)

		AttributeExpression(col string) (out string, err error)
	}

	// GenericTableCodec is a generic implementation of TableCodec
	GenericTableCodec struct {
		// table identifier (name)
		ident string

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
		ident:   m.Ident,
	}

	var (
		colIdent string
		attr     *dal.Attribute
		done     = make(map[string]bool)
		cols     = make([]Column, 0, len(m.Attributes))
	)

	for a := range m.Attributes {
		attr = m.Attributes[a]
		colIdent = attr.StoreIdent()

		switch attr.Store.(type) {
		case *dal.CodecRecordValueSetJSON:
			if done[colIdent] {
				continue
			}

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

func (t *GenericTableCodec) Ident() string {
	return t.ident
}

func (t *GenericTableCodec) Columns() []Column {
	return t.columns
}

func (t *GenericTableCodec) AttributeExpression(col string) (out string, err error) {
	return col, nil
}

func (t *GenericTableCodec) GetColumn(name string) (c Column, ok bool) {
	for _, c = range t.Columns() {
		if c.Name() == name {
			return c, true
		}
	}

	return nil, false
}

func (t *GenericTableCodec) Encode(r dal.ValueGetter) (_ []byte, err error) {
	aux := make(map[string]any)

	for _, c := range t.columns {
		// @todo this won't fly
		aux[c.Name()], err = c.Encode(r)
		if err != nil {
			return
		}
	}

	bb, err := json.Marshal(aux)
	if err != nil {
		return
	}

	return bb, err
}

func (t *GenericTableCodec) Decode(buf []byte, r dal.ValueSetter) (err error) {
	payload := make(map[string]any)
	err = json.Unmarshal(buf, &payload)
	if err != nil {
		return
	}

	for _, c := range t.columns {
		if err = c.Decode(payload[c.Name()], r); err != nil {
			return err
		}
	}

	return
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
