package drivers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/modern-go/reflect2"
	"github.com/spf13/cast"
)

type (
	Type interface {
		// MakeScanBuffer() any
		Decode(any) (any, bool, error)
		Encode(any) (any, error)
	}

	// @todo makes sense to rethink this strategy
	//       we do not need to have 1:1 type pairings with the pkg/dal
	//       it makes sense to define a few most common ones here (in the rdbms/drivers)
	//       and introduce per/driver exceptions to handle implementation specific things
	//       like MySQLs lack of BOOLEAN that is replaced with TINYINT(1)

	TypeID        struct{ *dal.TypeID }
	TypeRef       struct{ *dal.TypeRef }
	TypeTimestamp struct{ *dal.TypeTimestamp }
	TypeTime      struct{ *dal.TypeTime }
	TypeDate      struct{ *dal.TypeDate }
	TypeNumber    struct{ *dal.TypeNumber }
	TypeText      struct{ *dal.TypeText }
	TypeBoolean   struct{ *dal.TypeBoolean }
	TypeEnum      struct{ *dal.TypeEnum }
	TypeGeometry  struct{ *dal.TypeGeometry }
	TypeJSON      struct{ *dal.TypeJSON }
)

// TypeWrap wraps type from data package
func TypeWrap(dt dal.Type) Type {
	switch c := dt.(type) {
	case *dal.TypeID:
		return &TypeID{c}
	case *dal.TypeRef:
		return &TypeRef{c}
	case *dal.TypeTimestamp:
		return &TypeTimestamp{c}
	case *dal.TypeTime:
		return &TypeTime{c}
	case *dal.TypeDate:
		return &TypeDate{c}
	case *dal.TypeNumber:
		return &TypeNumber{c}
	case *dal.TypeText:
		return &TypeText{c}
	case *dal.TypeBoolean:
		return &TypeBoolean{c}
	case *dal.TypeEnum:
		return &TypeEnum{c}
	case *dal.TypeGeometry:
		return &TypeGeometry{c}
	case *dal.TypeJSON:
		return &TypeJSON{c}
	}

	if dt == nil {
		panic(fmt.Sprintf("attribute type not set (nil)"))
	}

	panic(fmt.Sprintf("type implementation missing: %s", dt.Type()))
}

// func (*TypeID) MakeScanBuffer() any        { return new(ID) }
// func (*TypeRef) MakeScanBuffer() any       { return new(ID) }
// func (*TypeTimestamp) MakeScanBuffer() any { return new(sql.NullTime) }
// func (*TypeTime) MakeScanBuffer() any      { return new(sql.NullString) }
// func (*TypeDate) MakeScanBuffer() any      { return new(sql.NullTime) }
// func (*TypeNumber) MakeScanBuffer() any    { return new(sql.NullString) }
// func (*TypeText) MakeScanBuffer() any      { return new(sql.NullString) }
// func (*TypeBoolean) MakeScanBuffer() any   { return new(sql.NullBool) }
// func (*TypeEnum) MakeScanBuffer() any      { return new(sql.NullString) }
// func (*TypeGeometry) MakeScanBuffer() any  { return new(sql.RawBytes) }
// func (*TypeJSON) MakeScanBuffer() any      { return new(sql.RawBytes) }

func (t *TypeID) Decode(raw any) (any, bool, error) {
	dec, err := cast.ToUint64E(raw)
	if err != nil {
		return 0, false, err
	}

	return dec, true, nil

}

func (t *TypeID) Encode(val any) (any, error) {
	// @todo :)
	return val, nil
}

func (t *TypeRef) Decode(raw any) (any, bool, error) {
	dec, is := raw.(*ID)
	if !is {
		return nil, false, fmt.Errorf("unexpected raw type %T for Ref", raw)
	}

	return dec.ID, dec.Valid, nil
}

func (t *TypeRef) Encode(val any) (any, error) {
	return val, nil
}

func (t *TypeTimestamp) Decode(raw any) (any, bool, error) {
	if reflect2.IsNil(raw) {
		return nil, false, nil
	}

	dec, err := cast.ToTimeE(raw)
	if err != nil {
		return nil, false, err
	}

	return dec.UTC().Format(TimestampLayout(t.Timezone, t.Precision)), true, nil
}

func (t *TypeTimestamp) Encode(val any) (any, error) {
	if reflect2.IsNil(val) || val == "" {
		return nil, nil
	}

	return cast.ToTimeE(val)
}

func (t *TypeTime) Decode(raw any) (any, bool, error) {
	dec, err := cast.ToStringE(raw)
	if err != nil {
		return nil, false, err
	}

	if len(dec) == 0 {
		return time.Time{}, false, nil
	}

	parsed, err := time.Parse(TimeLayout(t.Timezone, t.Precision), dec)
	if err != nil {
		return time.Time{}, false, err
	}

	// @todo should we gracefully handle other combinations of time&precision?
	//       maybe with Strict flag?

	return parsed.Format(TimeLayout(t.Timezone, t.Precision)), true, nil
}

func (t *TypeTime) Encode(val any) (any, error) {
	if reflect2.IsNil(val) {
		return nil, nil
	}

	return cast.ToTimeE(val)
}

func (t *TypeDate) Decode(raw any) (any, bool, error) {
	dec, err := cast.ToTimeE(raw)
	if err != nil {
		return time.Time{}, false, err
	}

	return dec.Format(DateLayout), true, nil
}

func (t *TypeDate) Encode(val any) (any, error) {
	if reflect2.IsNil(val) {
		return nil, nil
	}

	return cast.ToTimeE(val)
}

func (t *TypeNumber) Decode(raw any) (any, bool, error) {
	dec, err := cast.ToInt64E(raw)
	if err != nil {
		return 0, false, err
	}

	return strconv.FormatInt(dec, 10), true, nil
}

func (t *TypeNumber) Encode(val any) (any, error) {
	num, err := cast.ToInt64E(val)
	if err != nil {
		return 0, err
	}

	return num, nil
}

func (t *TypeText) Decode(raw any) (any, bool, error) {
	dec, err := cast.ToStringE(raw)
	if err != nil {
		return "", false, err
	}

	return dec, true, nil
}

func (t *TypeText) Encode(val any) (any, error) {
	return cast.ToStringE(val)
}

func (t *TypeBoolean) Decode(raw any) (any, bool, error) {
	dec, err := cast.ToBoolE(raw)
	if err != nil {
		return false, false, err
	}

	return fmt.Sprintf("%v", dec), true, nil
}

func (t *TypeBoolean) Encode(val any) (any, error) {
	if reflect2.IsNil(val) {
		return nil, nil
	}

	return cast.ToBool(val), nil
}

func (t *TypeEnum) Decode(raw any) (any, bool, error) {
	dec, err := cast.ToStringE(raw)
	if err != nil {
		return "", false, err
	}

	return dec, true, nil
}

func (t *TypeEnum) Encode(val any) (any, error) {
	return cast.ToStringE(val)
}

func (t *TypeGeometry) Decode(raw any) (any, bool, error) {
	dec, err := cast.ToStringE(raw)
	if err != nil {
		return nil, false, err
	}

	// @todo ??

	return []byte(dec), len(dec) > 0, nil
}

func (t *TypeGeometry) Encode(val any) (any, error) {
	return val, nil
}

func (t *TypeJSON) Decode(raw any) (any, bool, error) {
	bb, is := raw.(string)
	if !is {
		return nil, false, fmt.Errorf("unexpected raw type %T for JSON", raw)
	}

	return []byte(bb), len(bb) > 0, nil
}

func (t *TypeJSON) Encode(val any) (any, error) {
	switch c := val.(type) {
	case int64, float64, bool, []byte, string, time.Time:
		return c, nil

	case json.Marshaler:
		// does the value type know how to encode itself as JSON?
		return c.MarshalJSON()

	default:
		// Last resort - just encode with JSON pkg
		return json.Marshal(val)
	}
}

const (
	DateLayout = "2006-01-02"
)

func TimestampLayout(tz bool, precision int) string {
	return DateLayout + "T" + TimeLayout(tz, precision)
}

func TimeLayout(tz bool, precision int) string {
	var layout = "15:04:05"
	if precision > 0 {
		layout += "." + strings.Repeat("9", precision)
	}

	if tz {
		layout += "Z07:00"
	}

	return layout
}
