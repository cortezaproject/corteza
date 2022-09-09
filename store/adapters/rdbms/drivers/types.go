package drivers

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/google/uuid"
	"github.com/spf13/cast"
)

type (
	Type interface {
		MakeScanBuffer() any
		Decode(any) (any, bool, error)
		Encode(any) (driver.Value, error)
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
	TypeBlob      struct{ *dal.TypeBlob }
	TypeUUID      struct{ *dal.TypeUUID }
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
	case *dal.TypeBlob:
		return &TypeBlob{c}
	case *dal.TypeUUID:
		return &TypeUUID{c}
	}

	panic(fmt.Sprintf("type implementation missing: %s", dt.Type()))
}

func (*TypeID) MakeScanBuffer() any        { return new(ID) }
func (*TypeRef) MakeScanBuffer() any       { return new(ID) }
func (*TypeTimestamp) MakeScanBuffer() any { return new(sql.NullTime) }
func (*TypeTime) MakeScanBuffer() any      { return new(sql.NullString) }
func (*TypeDate) MakeScanBuffer() any      { return new(sql.NullTime) }
func (*TypeNumber) MakeScanBuffer() any    { return new(sql.NullString) }
func (*TypeText) MakeScanBuffer() any      { return new(sql.NullString) }
func (*TypeBoolean) MakeScanBuffer() any   { return new(sql.NullBool) }
func (*TypeEnum) MakeScanBuffer() any      { return new(sql.NullString) }
func (*TypeGeometry) MakeScanBuffer() any  { return new(sql.RawBytes) }
func (*TypeJSON) MakeScanBuffer() any      { return new(sql.RawBytes) }
func (*TypeBlob) MakeScanBuffer() any      { return new(sql.RawBytes) }

func (t *TypeID) Decode(raw any) (any, bool, error) {
	dec, is := raw.(*ID)
	if !is {
		return nil, false, fmt.Errorf("unexpected raw type %T for ID", raw)
	}

	return dec.ID, dec.Valid, nil
}

func (t *TypeID) Encode(val any) (driver.Value, error) {
	return val, nil
}

func (t *TypeRef) Decode(raw any) (any, bool, error) {
	dec, is := raw.(*ID)
	if !is {
		return nil, false, fmt.Errorf("unexpected raw type %T for Ref", raw)
	}

	return dec.ID, dec.Valid, nil
}

func (t *TypeRef) Encode(val any) (driver.Value, error) {
	return val, nil
}

func (t *TypeTimestamp) Decode(raw any) (any, bool, error) {
	dec, is := raw.(*sql.NullTime)
	if !is {
		return nil, false, fmt.Errorf("unexpected raw type %T for Timestamp", raw)
	}

	if dec.Valid {
		return dec.Time.Format(TimestampLayout(t.Timezone, t.Precision)), dec.Valid, nil
	}

	return nil, false, nil
}

func (t *TypeTimestamp) Encode(val any) (driver.Value, error) {
	return val, nil
}

func (t *TypeTime) Decode(raw any) (any, bool, error) {
	dec, is := raw.(*sql.NullString)
	if !is {
		return nil, false, fmt.Errorf("unexpected raw type %T for Time", raw)
	}

	if !dec.Valid {
		return time.Time{}, false, nil
	}

	parsed, err := time.Parse(TimeLayout(t.Timezone, t.Precision), dec.String)
	if err != nil {
		return time.Time{}, false, err
	}

	// @todo should we gracefully handle other combinations of time&precision?
	//       maybe with Strict flag?

	return parsed.Format(TimeLayout(t.Timezone, t.Precision)), true, nil
}

func (t *TypeTime) Encode(val any) (driver.Value, error) {
	return val, nil
}

func (t *TypeDate) Decode(raw any) (any, bool, error) {
	dec, is := raw.(*sql.NullTime)
	if !is {
		return nil, false, fmt.Errorf("unexpected raw type %T for Date", raw)
	}

	return dec.Time.Format(DateLayout), dec.Valid, nil
}

func (t *TypeDate) Encode(val any) (driver.Value, error) {
	return val, nil
}

func (t *TypeNumber) Decode(raw any) (any, bool, error) {
	dec, is := raw.(*sql.NullString)
	if !is {
		return nil, false, fmt.Errorf("unexpected raw type %T for Number", raw)
	}

	return dec.String, dec.Valid, nil
}

func (t *TypeNumber) Encode(val any) (driver.Value, error) {
	return val, nil
}

func (t *TypeText) Decode(raw any) (any, bool, error) {
	dec, is := raw.(*sql.NullString)
	if !is {
		return nil, false, fmt.Errorf("unexpected raw type %T for Text", raw)
	}

	return dec.String, dec.Valid, nil
}

func (t *TypeText) Encode(val any) (driver.Value, error) {
	return val, nil
}

func (t *TypeBoolean) Decode(raw any) (any, bool, error) {
	dec, is := raw.(*sql.NullBool)
	if !is {
		return nil, false, fmt.Errorf("unexpected raw type %T for Boolean", raw)
	}

	return dec.Bool, dec.Valid, nil
}

func (t *TypeBoolean) Encode(val any) (driver.Value, error) {
	return cast.ToBool(val), nil
}

func (t *TypeEnum) Decode(raw any) (any, bool, error) {
	dec, is := raw.(*sql.NullString)
	if !is {
		return nil, false, fmt.Errorf("unexpected raw type %T for Enum", raw)
	}

	return dec.String, dec.Valid, nil
}

func (t *TypeEnum) Encode(val any) (driver.Value, error) {
	return val, nil
}

func (t *TypeGeometry) Decode(raw any) (any, bool, error) {
	bb, is := raw.(*sql.RawBytes)
	if !is {
		return nil, false, fmt.Errorf("unexpected raw type %T for Geometry", raw)
	}

	return []byte(*bb), bb != nil, nil
}

func (t *TypeGeometry) Encode(val any) (driver.Value, error) {
	return val, nil
}

func (t *TypeJSON) Decode(raw any) (any, bool, error) {
	bb, is := raw.(*sql.RawBytes)
	if !is {
		return nil, false, fmt.Errorf("unexpected raw type %T for JSON", raw)
	}

	return []byte(*bb), bb != nil, nil
}

func (t *TypeJSON) Encode(val any) (driver.Value, error) {
	switch c := val.(type) {
	case driver.Valuer:
		// does the value type know how to encode itself for the DB?
		return c.Value()

	// These types are native to driver.Value
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

func (t *TypeBlob) Decode(raw any) (any, bool, error) {
	bb, is := raw.(*sql.RawBytes)
	if !is {
		return nil, false, fmt.Errorf("unexpected raw type %T for Blob", raw)
	}

	return []byte(*bb), bb != nil, nil
}

func (t *TypeBlob) Encode(val any) (driver.Value, error) {
	return val, nil
}

func (*TypeUUID) MakeScanBuffer() any { return new(uuid.UUID) }

func (t *TypeUUID) Decode(raw any) (any, bool, error) {
	dec, is := raw.(*uuid.UUID)
	if !is {
		return nil, false, fmt.Errorf("unexpected raw type %T for UUID", raw)
	}

	return *dec, dec != nil, nil
}

func (t *TypeUUID) Encode(val any) (driver.Value, error) {
	return val, nil
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
