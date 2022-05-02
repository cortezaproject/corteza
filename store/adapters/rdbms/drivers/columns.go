package drivers

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/crs"
	"github.com/cortezaproject/corteza-server/pkg/data"
	"github.com/spf13/cast"
)

type (
	Column interface {
		Name() string
		Attribute() *data.Attribute
		IsPrimaryKey() bool
		Encode(crs.ValueGetter) (any, error)
		Decode(any, crs.ValueSetter) error
		Type() Type
	}

	SingleValueColumn struct {
		typ  Type
		name string
		attr *data.Attribute
	}

	SimpleJsonDocColumn struct {
		name       string
		attributes []*data.Attribute
	}
)

func NewSingleValueColumn(d Dialect, a *data.Attribute) *SingleValueColumn {
	return &SingleValueColumn{
		typ:  d.TypeWrap(a.Type),
		attr: a,
		name: attrColumnIdent(a),
	}
}

func (c *SingleValueColumn) Name() string {
	return c.name
}

func (c *SingleValueColumn) IsPrimaryKey() bool {
	return c.attr.PrimaryKey
}

func (c *SingleValueColumn) Attribute() *data.Attribute {
	return c.attr
}

func (c *SingleValueColumn) Type() Type {
	return c.typ
}

func (c *SingleValueColumn) Encode(r crs.ValueGetter) (any, error) {
	val, err := r.GetValue(c.attr.Ident, 0)
	if err != nil {
		return nil, err
	}

	return c.typ.Encode(val)
}

func (c *SingleValueColumn) Decode(raw any, r crs.ValueSetter) error {
	value, valid, err := c.typ.Decode(raw)
	if err != nil {
		return err
	}

	ident := c.attr.Ident
	if !valid {
		return r.SetValue(ident, 0, nil)
	}

	return r.SetValue(ident, 0, value)
}

func (c *SimpleJsonDocColumn) Name() string {
	return c.name
}

func (c *SimpleJsonDocColumn) IsPrimaryKey() bool {
	return false
}

func (c *SimpleJsonDocColumn) Attribute() *data.Attribute {
	return c.attributes[0]
}

func (c *SimpleJsonDocColumn) Type() Type {
	return &TypeJSON{}
}

func (c *SimpleJsonDocColumn) Encode(r crs.ValueGetter) (_ any, err error) {
	var (
		aux   = make(map[string][]any)
		value any
		place uint
		size  uint

		count = r.CountValues()
	)

	for _, attr := range c.attributes {
		size = count[attr.Ident]

		aux[attr.Ident] = make([]any, size)

		for place = 0; place < size; place++ {
			value, err = r.GetValue(attr.Ident, 0)
			if err != nil {
				return nil, err
			}

			// now, encode the value according to JSON format constraints
			switch attr.Type.(type) {
			case *data.TypeBoolean:
				// we want booleans stored as booleans
				aux[attr.Ident][place] = cast.ToBool(value)

			default:
				// every other value is
				aux[attr.Ident][place] = value
			}

			if !attr.MultiValue {
				// model attribute supports storing of single values only.
				break
			}
		}
	}

	return json.Marshal(aux)
}

func (c *SimpleJsonDocColumn) Decode(raw any, r crs.ValueSetter) (err error) {
	rawJson, is := raw.(*sql.RawBytes)
	if !is {
		return fmt.Errorf("incompatible input value type (%T), expecting *sql.RawBytes", raw)
	}

	buf := make(map[string][]any)
	if err = json.Unmarshal(*rawJson, &buf); err != nil {
		return
	}

	// @todo this is too naive
	for name, vv := range buf {
		var attr *data.Attribute
		for a := range c.attributes {
			if c.attributes[a].Ident != name {
				continue
			}

			attr = c.attributes[a]
		}

		if attr == nil {
			// unrecognized value in the json doc
			continue
		}

		for pos, v := range vv {
			if err = r.SetValue(name, uint(pos), v); err != nil {
				return
			}

			if !attr.MultiValue {
				// model attribute supports storing of single values only.
				break
			}
		}

		attr = nil
	}

	return
}
