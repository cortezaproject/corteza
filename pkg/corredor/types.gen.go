package corredor

// 	Hello! This file is auto-generated.

type (

	// ScriptSet slice of Script
	//
	// This type is auto-generated.
	ScriptSet []*Script
)

// Walk iterates through every slice item and calls w(Script) err
//
// This function is auto-generated.
func (set ScriptSet) Walk(w func(*Script) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Script) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ScriptSet) Filter(f func(*Script) (bool, error)) (out ScriptSet, err error) {
	var ok bool
	out = ScriptSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}
