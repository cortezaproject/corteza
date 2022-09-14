package gvalfnc

import (
	"errors"
	"time"

	"github.com/lestrrat-go/strftime"
	"github.com/spf13/cast"
)

func Now() time.Time {
	return time.Now()
}

func Quarter(in any) (int, error) {
	t, _, err := PrepMod(in, 0)
	if err != nil {
		return 0, err
	}

	return int(t.Month() / 4), nil
}

func Year(in any) (int, error) {
	t, _, err := PrepMod(in, 0)
	if err != nil {
		return 0, err
	}

	return t.Year(), nil
}

func Month(in any) (int, error) {
	t, _, err := PrepMod(in, 0)
	if err != nil {
		return 0, err
	}

	return int(t.Month()), nil
}

func Date(in any) (int, error) {
	t, _, err := PrepMod(in, 0)
	if err != nil {
		return 0, err
	}

	return t.Day(), nil
}

func PrepMod(base interface{}, mod interface{}) (*time.Time, int, error) {
	var (
		t *time.Time
	)

	switch auxt := base.(type) {
	case time.Time:
		t = &auxt
	case *time.Time:
		t = auxt
	case string:
		tt, err := cast.ToTimeE(auxt)

		if err != nil {
			return nil, 0, err
		}

		t = &tt
	default:
		return nil, 0, errors.New("unexpected input type")
	}

	m, err := cast.ToIntE(mod)
	if err != nil {
		return nil, 0, err
	}

	return t, m, nil
}

// Strftime formats time with POSIX standard format
// More details here:
// https://github.com/lestrrat-go/strftime#supported-conversion-specifications
func StrfTime(base interface{}, f string) (string, error) {
	t, _, err := PrepMod(base, 0)

	if err != nil {
		return "", err
	}

	o, _ := strftime.Format(f, *t,
		strftime.WithMilliseconds('b'),
		strftime.WithUnixSeconds('L'))

	return o, nil
}
