package resource

import (
	"encoding/json"
	"time"
)

type (
	Timestamp struct {
		// S is a stringified representation of the timestamp
		S string
		// T is a time representation of the timestamp
		T *time.Time
		// Tz is the timezone the timestamp will be encoded into
		Tz string
		// Tpl is the template the timestamp will be encoded into
		Tpl string
	}
	Timestamps struct {
		CreatedAt   *Timestamp
		UpdatedAt   *Timestamp
		DeletedAt   *Timestamp
		ArchivedAt  *Timestamp
		SuspendedAt *Timestamp
	}
)

func (ts *Timestamp) MarshalYAML() (interface{}, error) {
	return ts.S, nil
}

func (ts *Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(ts.S)
}

// MakeTimestamp initializes a Timestamp from the passed timestamp
func MakeTimestamp(ts string) *Timestamp {
	t := toTime(ts)
	if t == nil {
		return nil
	}

	return &Timestamp{
		S: ts,
		T: t,
	}
}

// MakeTimestampsCUDA initializes timestamps for createdAt, updatedAt, deletedAt and archivedAt
func MakeTimestampsCUDA(c, u, d, a *time.Time) *Timestamps {
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

// MakeTimestampsCUDAS initializes timestamps for createdAt, updatedAt, deletedAt, archivedAt and suspendedAt
func MakeTimestampsCUDAS(c, u, d, a, s *time.Time) *Timestamps {
	t := MakeTimestampsCUDA(c, u, d, a)

	if s != nil && !s.IsZero() {
		t.SuspendedAt = &Timestamp{T: s}
	}

	return t
}

// Model stringifies the timestamps based on the passed template and timezone
func (tt *Timestamps) Model(tpl string, tz string) (*Timestamps, error) {
	var err error
	if tt.CreatedAt != nil {
		tt.CreatedAt, err = tt.CreatedAt.Model(tpl, tz)
	}
	if tt.UpdatedAt != nil {
		tt.UpdatedAt, err = tt.UpdatedAt.Model(tpl, tz)
	}
	if tt.DeletedAt != nil {
		tt.DeletedAt, err = tt.DeletedAt.Model(tpl, tz)
	}
	if tt.ArchivedAt != nil {
		tt.ArchivedAt, err = tt.ArchivedAt.Model(tpl, tz)
	}
	if tt.SuspendedAt != nil {
		tt.SuspendedAt, err = tt.SuspendedAt.Model(tpl, tz)
	}

	return tt, err
}

// Model stringifies the timestamp based on the passed template and timezone
func (ts *Timestamp) Model(tpl string, tz string) (*Timestamp, error) {
	if tz != "" {
		tzL, err := time.LoadLocation(tz)
		if err != nil {
			return nil, err
		}
		t := ts.T.In(tzL)
		ts.T = &t
	}

	if tpl == "" {
		tpl = time.RFC3339
	}

	ts.S = ts.T.Format(tpl)
	return ts, nil
}
