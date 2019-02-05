package types

// 	Hello! This file is auto-generated.

type (

	// TeamSet slice of Team
	//
	// This type is auto-generated.
	TeamSet []*Team
)

// Walk iterates through every slice item and calls w(Team) err
//
// This function is auto-generated.
func (set TeamSet) Walk(w func(*Team) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Team) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set TeamSet) Filter(f func(*Team) (bool, error)) (out TeamSet, err error) {
	var ok bool
	out = TeamSet{}
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
func (set TeamSet) FindByID(ID uint64) *Team {
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
func (set TeamSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Resources returns a slice of types.Resource from all items in the set
//
// This function is auto-generated.
func (set TeamSet) Resources() (Resources []Resource) {
	Resources = make([]Resource, len(set))

	for i := range set {
		Resources[i] = set[i].Resource()
	}

	return
}
