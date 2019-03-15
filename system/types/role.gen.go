package types

// 	Hello! This file is auto-generated.

type (

	// RoleSet slice of Role
	//
	// This type is auto-generated.
	RoleSet []*Role
)

// Walk iterates through every slice item and calls w(Role) err
//
// This function is auto-generated.
func (set RoleSet) Walk(w func(*Role) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Role) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set RoleSet) Filter(f func(*Role) (bool, error)) (out RoleSet, err error) {
	var ok bool
	out = RoleSet{}
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
func (set RoleSet) FindByID(ID uint64) *Role {
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
func (set RoleSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
