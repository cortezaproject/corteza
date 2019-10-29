package rh

import (
	"time"
)

var (
	now = func() time.Time {
		return time.Now()
	}
)

// SetCurrentTimeRounded sets current time (rounded to seconds) to a given ptr
func SetCurrentTimeRounded(v interface{}) {
	n := now().Truncate(time.Second)

	switch t := v.(type) {
	case *time.Time:
		*t = n
	case **time.Time:
		_ = t
		*t = &n
	}
}
