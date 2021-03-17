package expr

import (
	"time"

	"github.com/PaesslerAG/gval"
	"github.com/lestrrat-go/strftime"
)

func TimeFunctions() []gval.Language {
	return []gval.Language{
		gval.Function("earliest", earliest),
		gval.Function("latest", latest),
		gval.Function("parseISOTime", parseISOTime),
		gval.Function("modTime", modTime),
		gval.Function("parseDuration", time.ParseDuration),
		gval.Function("strftime", strfTime),
		gval.Function("isLeapYear", isLeapYear),
		gval.Function("now", now),
		gval.Function("isWeekDay", isWeekDay),
	}
}

func now() time.Time {
	return time.Now()
}

func isLeapYear(f time.Time) bool {
	return time.Date(f.Year(), time.December, 31, 0, 0, 0, 0, time.Local).YearDay() == 366
}

func isLeapDay(f time.Time) bool {
	return f.Day() == 29 && f.Month() == 2
}

func isWeekDay(f time.Time) bool {
	return time.Sunday < f.Weekday() && f.Weekday() < time.Saturday
}

func earliest(f time.Time, aa ...time.Time) time.Time {
	for _, s := range aa {
		if s.Before(f) {
			f = s
		}
	}

	return f
}

func latest(f time.Time, aa ...time.Time) time.Time {
	for _, s := range aa {
		if s.After(f) {
			f = s
		}
	}

	return f
}

func parseISOTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func modTime(t time.Time, mod interface{}) (time.Time, error) {
	var (
		err error
		d   time.Duration
	)
	switch c := mod.(type) {
	case time.Duration:
		d = c
	case string:
		d, err = time.ParseDuration(c)
	}

	if err != nil {
		return t, err
	}

	return t.Add(d), nil
}

// Strftime formats time with POSIX standard format
// More details here:
// https://github.com/lestrrat-go/strftime#supported-conversion-specifications
func strfTime(t time.Time, f string) string {
	o, _ := strftime.Format(f, t, strftime.WithMilliseconds('b'))
	return o
}
