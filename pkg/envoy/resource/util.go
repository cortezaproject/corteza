package resource

import (
	"fmt"
	"strconv"
	"time"
)

// fn converts identifier values (string, fmt.Stringer, uint64) to string slice
//
// Each value is checked and should not be empty or zero
func identifiers(ii ...interface{}) []string {
	ss := make([]string, 0, len(ii))

	for _, i := range ii {
		switch c := i.(type) {
		case uint64:
			if c == 0 {
				continue
			}

			ss = append(ss, strconv.FormatUint(c, 10))

		case fmt.Stringer:
			if c.String() == "" {
				continue
			}

			ss = append(ss, c.String())

		case string:
			if c == "" {
				continue
			}

			ss = append(ss, c)
		}
	}

	return ss
}

// Check checks if the identifier c is in the set of identifiers ii
func Check(a string, ii ...interface{}) bool {
	for _, i := range ii {
		switch pi := i.(type) {
		case string:
			if a == pi {
				return true
			}
		case uint64:
			if pi > 0 && strconv.FormatUint(pi, 10) == i {
				return true
			}
		}
	}

	return false
}

func MakeCUDATimestamps(c, u, d, a *time.Time) *Timestamps {
	t := &Timestamps{}

	if c != nil && !c.IsZero() {
		t.CreatedAt = &Timestamp{T: c}
	}
	if u != nil && !u.IsZero() {
		t.UpdatedAt = &Timestamp{T: u}
	}
	if d != nil && !d.IsZero() {
		t.DeletedAt = &Timestamp{T: d}
	}
	if a != nil && !a.IsZero() {
		t.ArchivedAt = &Timestamp{T: a}
	}

	return t
}

func MakeCUDASTimestamps(c, u, d, a, s *time.Time) *Timestamps {
	t := MakeCUDATimestamps(c, u, d, a)

	if s != nil && !s.IsZero() {
		t.SuspendedAt = &Timestamp{T: s}
	}

	return t
}

func MakeCUDOUserstamps(c, u, d, o uint64) *Userstamps {
	us := &Userstamps{}

	if c > 0 {
		us.CreatedBy = &Userstamp{
			UserID: c,
			Ref:    strconv.FormatUint(c, 10),
		}
	}
	if u > 0 {
		us.UpdatedBy = &Userstamp{
			UserID: u,
			Ref:    strconv.FormatUint(u, 10),
		}
	}
	if d > 0 {
		us.DeletedBy = &Userstamp{
			UserID: d,
			Ref:    strconv.FormatUint(d, 10),
		}
	}
	if o > 0 {
		us.OwnedBy = &Userstamp{
			UserID: o,
			Ref:    strconv.FormatUint(o, 10),
		}
	}

	return us
}

func FirstOkString(ss ...string) string {
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
