package resource

import (
	"time"
)

func firstOkString(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}

// Taken (and modified) from compose/service/values/sanitizer.go
func toTime(v string) *time.Time {
	ff := []string{
		time.RFC3339,
		time.RFC1123Z,
		time.RFC1123,
		time.RFC850,
		time.RFC822Z,
		time.RFC822,
		time.RubyDate,
		time.UnixDate,
		time.ANSIC,
		"2006/_1/_2 15:04:05",
		"2006/_1/_2 15:04",
	}

	for _, f := range ff {
		parsed, err := time.Parse(f, v)
		if err == nil {
			return &parsed
		}
	}

	return nil
}
