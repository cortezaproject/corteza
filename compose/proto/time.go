package proto

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

// Converts time.Time (ptr AND value) to *timestamp.Timestamp
//
// Intentionally ignoring
func fromTime(i interface{}) *timestamp.Timestamp {
	switch t := i.(type) {
	case *time.Time:
		if t == nil {
			return nil
		}
		return &timestamp.Timestamp{Seconds: t.Unix(), Nanos: int32(t.Nanosecond())}
	case time.Time:
		return &timestamp.Timestamp{Seconds: t.Unix(), Nanos: int32(t.Nanosecond())}
	default:
		return nil
	}
}
