package expr

import (
	"errors"
	"time"

	"github.com/PaesslerAG/gval"
	"github.com/cortezaproject/corteza-server/pkg/gvalfnc"
	"github.com/lestrrat-go/strftime"
)

func TimeFunctions() []gval.Language {
	return []gval.Language{
		gval.Function("earliest", earliest),
		gval.Function("latest", latest),
		gval.Function("parseISOTime", parseISOTime),
		gval.Function("modTime", modTime),
		gval.Function("modDate", modDate),
		gval.Function("modWeek", modWeek),
		gval.Function("modMonth", modMonth),
		gval.Function("modYear", modYear),
		gval.Function("parseDuration", time.ParseDuration),
		gval.Function("strftime", strfTime),
		gval.Function("isLeapYear", isLeapYear),
		gval.Function("now", now),
		gval.Function("isWeekDay", isWeekDay),
		gval.Function("sub", sub),
	}
}

func now() time.Time {
	return gvalfnc.Now()
}

func isLeapYear(base interface{}) (bool, error) {
	t, _, err := gvalfnc.PrepMod(base, 0)
	if err != nil {
		return false, err
	}
	return time.Date(t.Year(), time.December, 31, 0, 0, 0, 0, time.Local).YearDay() == 366, nil
}

func isLeapDay(base interface{}) (bool, error) {
	t, _, err := gvalfnc.PrepMod(base, 0)
	if err != nil {
		return false, err
	}
	return t.Day() == 29 && t.Month() == 2, nil
}

func isWeekDay(base interface{}) (bool, error) {
	t, _, err := gvalfnc.PrepMod(base, 0)
	if err != nil {
		return false, err
	}
	return time.Sunday < t.Weekday() && t.Weekday() < time.Saturday, nil
}

func earliest(f interface{}, aa ...interface{}) (*time.Time, error) {
	t, _, err := gvalfnc.PrepMod(f, 0)
	if err != nil {
		return nil, err
	}
	for _, a := range aa {
		s, _, err := gvalfnc.PrepMod(a, 0)
		if err != nil {
			return nil, err
		}
		if (*s).Before(*t) {
			t = s
		}
	}

	return t, nil
}

func latest(f interface{}, aa ...interface{}) (*time.Time, error) {
	t, _, err := gvalfnc.PrepMod(f, 0)
	if err != nil {
		return nil, err
	}
	for _, a := range aa {
		s, _, err := gvalfnc.PrepMod(a, 0)
		if err != nil {
			return nil, err
		}
		if s.After(*t) {
			t = s
		}
	}

	return t, nil
}

func parseISOTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func modTime(base interface{}, mod interface{}) (*time.Time, error) {
	var (
		err error
		d   time.Duration
		t   *time.Time
	)

	t, _, err = gvalfnc.PrepMod(base, 0)
	if err != nil {
		return nil, err
	}

	switch c := mod.(type) {
	case time.Duration:
		d = c
	case string:
		d, err = time.ParseDuration(c)
	}

	if err != nil {
		return t, err
	}

	tmp := t.Add(d)
	return &tmp, nil
}

func modDate(base interface{}, mod interface{}) (*time.Time, error) {
	t, m, err := gvalfnc.PrepMod(base, mod)
	if err != nil {
		return nil, err
	}

	tmp := t.AddDate(0, 0, m)
	return &tmp, nil
}

func modWeek(base interface{}, mod interface{}) (*time.Time, error) {
	t, m, err := gvalfnc.PrepMod(base, mod)
	if err != nil {
		return nil, err
	}

	tmp := t.AddDate(0, 0, 7*m)
	return &tmp, nil
}

func modMonth(base interface{}, mod interface{}) (*time.Time, error) {
	t, m, err := gvalfnc.PrepMod(base, mod)
	if err != nil {
		return nil, err
	}

	tmp := t.AddDate(0, m, 0)
	return &tmp, nil
}

func modYear(base interface{}, mod interface{}) (*time.Time, error) {
	t, m, err := gvalfnc.PrepMod(base, mod)
	if err != nil {
		return nil, err
	}

	tmp := t.AddDate(m, 0, 0)
	return &tmp, nil
}

// Strftime formats time with POSIX standard format
// More details here:
// https://github.com/lestrrat-go/strftime#supported-conversion-specifications
func strfTime(base interface{}, f string) (string, error) {
	t, _, err := gvalfnc.PrepMod(base, 0)

	if err != nil {
		return "", err
	}

	o, _ := strftime.Format(f, *t,
		strftime.WithMilliseconds('b'),
		strftime.WithUnixSeconds('L'))

	return o, nil
}

// sub returns difference between two date into milliseconds
func sub(from interface{}, to interface{}) (out int64, err error) {
	t1, _, err := gvalfnc.PrepMod(from, 0)
	if err != nil {
		return
	}

	t2, _, err := gvalfnc.PrepMod(to, 0)
	if err != nil {
		return
	}

	var duration time.Duration

	if t1.After(*t2) {
		duration = t1.Sub(*t2)
	} else {
		return -1, errors.New("expecting 2nd input to be less than 1st input")
	}

	return duration.Milliseconds(), nil
}
