package types

import (
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

func (v *RecordValueErrorSet) Push(err ...RecordValueError) {
	v.Set = append(v.Set, err...)
}

func (v *RecordValueErrorSet) IsValid() bool {
	return v == nil || len(v.Set) == 0
}

func (v *RecordValueErrorSet) Error() string {
	return fmt.Sprintf("%d issue(s) found", len(v.Set))
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
