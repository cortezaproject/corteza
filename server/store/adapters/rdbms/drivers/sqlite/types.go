package sqlite

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/spf13/cast"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers"
)

type (
	TypeDate      struct{ *dal.TypeDate }
	TypeTime      struct{ *dal.TypeTime }
	TypeTimestamp struct{ *dal.TypeTimestamp }
)

func (*TypeDate) MakeScanBuffer() any { return new(sql.NullString) }

func (t *TypeDate) Decode(raw any) (any, bool, error) {
	dec, is := raw.(*sql.NullString)
	if !is {
		return nil, false, fmt.Errorf("unexpected raw type %T for Time", raw)
	}

	if !dec.Valid {
		return time.Time{}, false, nil
	}

	parsed, err := cast.StringToDate(dec.String)
	if err != nil {
		return time.Time{}, false, err
	}

	return parsed.Format(drivers.DateLayout), true, nil
}

func (t *TypeDate) Encode(val any) (driver.Value, error) {
	if val == nil {
		return nil, nil
	}
	return val, nil
}

func (*TypeTime) MakeScanBuffer() any { return new(sql.NullTime) }

func (t *TypeTime) Decode(raw any) (any, bool, error) {
	dec, is := raw.(*sql.NullTime)
	if !is {
		return nil, false, fmt.Errorf("unexpected raw type %T for Time", raw)
	}

	if !dec.Valid {
		return time.Time{}, false, nil
	}

	return dec.Time.Format(drivers.TimeLayout(false, t.Precision)), dec.Valid, nil
}

func (t *TypeTime) Encode(val any) (driver.Value, error) {
	if val == nil {
		return nil, nil
	}
	return val, nil
}

func (*TypeTimestamp) MakeScanBuffer() any { return new(sql.NullTime) }

func (t *TypeTimestamp) Decode(raw any) (any, bool, error) {
	dec, is := raw.(*sql.NullTime)
	if !is {
		return nil, false, fmt.Errorf("unexpected raw type %T for Timestamp", raw)
	}

	if dec.Valid {
		val := dec.Time.Format(drivers.TimestampLayout(false, t.Precision))
		return val, dec.Valid, nil
	}

	return nil, false, nil
}

func (t *TypeTimestamp) Encode(val any) (driver.Value, error) {
	if val == nil {
		return val, nil
	}

	switch c := val.(type) {
	case string:
		parsed, err := cast.StringToDate(c)
		if err != nil {
			return val, nil
		}

		val = parsed.Format(drivers.TimestampLayout(false, t.Precision))

	case time.Time:
		val = c.Format(drivers.TimestampLayout(false, t.Precision))
	}

	return val, nil
}
