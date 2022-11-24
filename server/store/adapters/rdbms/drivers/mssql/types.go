package mssql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers"
)

type (
	TypeTime struct{ *dal.TypeTime }
)

func (*TypeTime) MakeScanBuffer() any { return new(sql.NullTime) }

func (t *TypeTime) Decode(raw any) (any, bool, error) {
	dec, is := raw.(*sql.NullTime)
	if !is {
		return nil, false, fmt.Errorf("unexpected raw type %T for Time", raw)
	}

	if !dec.Valid {
		return time.Time{}, false, nil
	}

	//parsed, err := time.Parse(drivers.TimeLayout(t.Timezone, t.Precision), dec.String)
	//if err != nil {
	//	return time.Time{}, false, err
	//}

	// @todo should we gracefully handle other combinations of time&precision?
	//       maybe with Strict flag?

	return dec.Time.Format(drivers.TimeLayout(false, 0)), true, nil
}

func (t *TypeTime) Encode(val any) (driver.Value, error) {
	return val, nil
}
