package automation

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

// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set ScriptSet) FindByID(ID uint64) *Script {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// IDs returns a slice of uint64s from all items in the set
//
// This function is auto-generated.
func (set ScriptSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
