package actionlog

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/actionlog/types.yaml

type (

	// ActionSet slice of Action
	//
	// This type is auto-generated.
	ActionSet []*Action
)

// Walk iterates through every slice item and calls w(Action) err
//
// This function is auto-generated.
func (set ActionSet) Walk(w func(*Action) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Action) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ActionSet) Filter(f func(*Action) (bool, error)) (out ActionSet, err error) {
	var ok bool
	out = ActionSet{}
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
func (set ActionSet) FindByID(ID uint64) *Action {
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
func (set ActionSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
