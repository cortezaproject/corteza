package drivers

import (
	"database/sql/driver"
	"strconv"

	"github.com/spf13/cast"
)

// nullUint64 represents an uint64 that may be null.
// nullUint64 implements the Scanner interface so
// it can be used as a scan dåestination, similar to NullString.
type ID struct {
	ID    uint64
	Valid bool // Valid is true if Uint64 is not NULL
}

// Scan implements the Scanner interface.
func (n *ID) Scan(value any) (err error) {
	if value == nil {
		n.ID, n.Valid = 0, false
		return nil
	}

	// NOT NULL => Valid!
	n.Valid = true
	if b, is := value.([]byte); is {
		n.ID, err = strconv.ParseUint(string(b), 10, 64)
	} else {
		n.ID, err = cast.ToUint64E(value)
	}

	return
}

// Value implements the driver Valuer interface.
func (n ID) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.ID, nil
}

func (n ID) String() string {
	return strconv.FormatUint(n.ID, 10)
}
