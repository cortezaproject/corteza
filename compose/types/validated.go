package types

import (
	"encoding/json"
	"errors"
	"fmt"
)

type (
	RecordValueError struct {
		Kind    string                 `json:"kind"`
		Message string                 `json:"message"`
		Meta    map[string]interface{} `json:"meta"`
	}

	RecordValueErrorSet struct {
		Set []RecordValueError `json:"set"`
	}
)

// safe to show details of this error
func (RecordValueErrorSet) Safe() bool { return true }

func (v *RecordValueErrorSet) Push(err ...RecordValueError) {
	v.Set = append(v.Set, err...)
}

func (v *RecordValueErrorSet) IsValid() bool {
	return v == nil || len(v.Set) == 0
}

func (v *RecordValueErrorSet) Error() string {
	var no = 0
	if v != nil {
		no = len(v.Set)
	}

	return fmt.Sprintf("%d issue(s) found", no)
}

func (v RecordValueErrorSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message string             `json:"message"`
		Set     []RecordValueError `json:"set,omitempty"`
	}{
		Message: v.Error(),
		Set:     v.Set,
	})
}

func (v *RecordValueErrorSet) HasKind(kind string) bool {
	if v == nil || v.IsValid() {
		return false
	}

	for _, e := range v.Set {
		if e.Kind == kind {
			return true
		}
	}

	return false
}

// IsRecordValueErrorSet tests if given error is RecordValueErrorSet (or it wraps it) and it has errors
// If not is not (or !IsValid), it return nil!
func IsRecordValueErrorSet(err error) *RecordValueErrorSet {
	for {
		if err == nil {
			return nil
		}

		if rve, ok := err.(*RecordValueErrorSet); ok {
			if rve.IsValid() {
				return nil
			}

			return rve
		}

		err = errors.Unwrap(err)
	}
}
