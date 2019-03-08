package types

// 	Hello! This file is auto-generated.

type (

	// ApplicationSet slice of Application
	//
	// This type is auto-generated.
	ApplicationSet []*Application
)

// Walk iterates through every slice item and calls w(Application) err
//
// This function is auto-generated.
func (set ApplicationSet) Walk(w func(*Application) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Application) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ApplicationSet) Filter(f func(*Application) (bool, error)) (out ApplicationSet, err error) {
	var ok bool
	out = ApplicationSet{}
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
func (set ApplicationSet) FindByID(ID uint64) *Application {
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
func (set ApplicationSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
