package types

// 	Hello! This file is auto-generated.

type (

	// RecordValueSet slice of RecordValue
	//
	// This type is auto-generated.
	RecordValueSet []*RecordValue
)

// Walk iterates through every slice item and calls w(RecordValue) err
//
// This function is auto-generated.
func (set RecordValueSet) Walk(w func(*RecordValue) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(RecordValue) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set RecordValueSet) Filter(f func(*RecordValue) (bool, error)) (out RecordValueSet, err error) {
	var ok bool
	out = RecordValueSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}
