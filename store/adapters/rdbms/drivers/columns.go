package drivers

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/spf13/cast"
)

type (
	Column interface {
		Name() string
		Attribute() *dal.Attribute
		IsPrimaryKey() bool
		Encode(dal.ValueGetter) (any, error)
		Decode(any, dal.ValueSetter) error
		Type() Type
	}

	SingleValueColumn struct {
		typ  Type
		name string
		attr *dal.Attribute
	}

	SimpleJsonDocColumn struct {
		name       string
		attributes []*dal.Attribute
	}
)

func NewSingleValueColumn(d Dialect, a *dal.Attribute) *SingleValueColumn {
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

func (c *SingleValueColumn) Attribute() *dal.Attribute {
	return c.attr
}

func (c *SingleValueColumn) Type() Type {
	return c.typ
}

func (c *SingleValueColumn) Encode(r dal.ValueGetter) (any, error) {
	val, err := r.GetValue(c.attr.Ident, 0)
	if err != nil {
		return nil, err
	}

	return c.typ.Encode(val)
}

func (c *SingleValueColumn) Decode(raw any, r dal.ValueSetter) error {
	value, valid, err := c.typ.Decode(raw)
	if err != nil {
		return err
	}

	// now, encode the value according to JSON format constraints
	switch c.attr.Type.(type) {
	case *dal.TypeBoolean:
		// for backward compatibility reasons
		// we need to cast true bool values to "1"
		// and use "" for other (false) values
		if cast.ToBool(value) {
			value = "1"
		} else {
			value = ""
		}
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

func (c *SimpleJsonDocColumn) Attribute() *dal.Attribute {
	return c.attributes[0]
}

func (c *SimpleJsonDocColumn) Type() Type {
	return &TypeJSON{}
}

func (c *SimpleJsonDocColumn) Encode(r dal.ValueGetter) (_ any, err error) {
	var (
		aux   = make(map[string][]any)
		value any
		place uint

		// determinate how many values can we expected to store in the attribute
		//
		// implementations can return nil for the value to signal
		// that each attribute holds exactly one value
		count = r.CountValues()

		// preset this to one just in case CountValues()
		// returns nil!
		size uint = 1
	)

	for _, attr := range c.attributes {
		if count != nil {
			size = count[attr.Ident]
		}

		aux[attr.Ident] = make([]any, size)

		for place = 0; place < size; place++ {
			value, err = r.GetValue(attr.Ident, 0)
			if err != nil {
				return nil, err
			}

			// now, encode the value according to JSON format constraints
			switch attr.Type.(type) {
			case *dal.TypeBoolean:
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

func (c *SimpleJsonDocColumn) Decode(raw any, r dal.ValueSetter) (err error) {
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
		var attr *dal.Attribute
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
			// now, encode the value according to JSON format constraints
			switch attr.Type.(type) {
			case *dal.TypeBoolean:
				// for backward compatibility reasons
				// we need to cast true bool values to "1"
				// and use "" for other (false) values
				if cast.ToBool(v) {
					v = "1"
				} else {
					v = ""
				}
			}

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
