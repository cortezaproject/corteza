package drivers

import (
	"github.com/cortezaproject/corteza/server/pkg/dal"
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
		name: a.StoreIdent(),
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
	panic("not implemented")
}

func (c *SimpleJsonDocColumn) Decode(raw any, r dal.ValueSetter) (err error) {
	panic("not implemented")
}

func (c *SimpleJsonDocColumn) DecodeOld(raw any, r dal.ValueSetter) (err error) {
	panic("not implemented")
}
