package sqlite

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/spf13/cast"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers"
)

type (
	TypeDate struct{ *dal.TypeDate }
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
	return val, nil
}
