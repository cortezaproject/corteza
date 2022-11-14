package expr

import (
	"fmt"
)

func compareToBoolean(a Boolean, b TypedValue) (int, error) {
	c, err := NewBoolean(b)
	if err != nil {
		return 0, fmt.Errorf("cannot compare %s and %s: %s", a.Type(), b.Type(), err.Error())
	}

	switch {
	case a.value == c.value:
		return 0, nil
	case !a.value && c.value:
		return -1, nil
	default:
		return 1, nil
	}
}

func compareToDateTime(a DateTime, b TypedValue) (int, error) {
	c, err := NewDateTime(b)
	if err != nil {
		return 0, fmt.Errorf("cannot compare %s and %s: %s", a.Type(), b.Type(), err.Error())
	}

	ar := a.value.Nanosecond()
	cr := c.value.Nanosecond()

	switch {
	case ar == cr:
		return 0, nil
	case ar < cr:
		return -1, nil
	default:
		return 1, nil
	}
}

func compareToDuration(a Duration, b TypedValue) (int, error) {
	c, err := NewDuration(b)
	if err != nil {
		return 0, fmt.Errorf("cannot compare %s and %s: %s", a.Type(), b.Type(), err.Error())
	}

	switch {
	case a.value == c.value:
		return 0, nil
	case a.value < c.value:
		return -1, nil
	default:
		return 1, nil
	}
}
