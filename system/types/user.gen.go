package types

// 	Hello! This file is auto-generated.

type (

	// UserSet slice of User
	//
	// This type is auto-generated.
	UserSet []*User
)

// Walk iterates through every slice item and calls w(User) err
//
// This function is auto-generated.
func (set UserSet) Walk(w func(*User) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(User) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set UserSet) Filter(f func(*User) (bool, error)) (out UserSet, err error) {
	var ok bool
	out = UserSet{}
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
func (set UserSet) FindByID(ID uint64) *User {
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
func (set UserSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
