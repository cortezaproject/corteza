package types

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
	return "invalid values input"
}
