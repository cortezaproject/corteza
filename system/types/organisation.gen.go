package types

// 	Hello! This file is auto-generated.

type (

	// OrganisationSet slice of Organisation
	//
	// This type is auto-generated.
	OrganisationSet []*Organisation
)

// Walk iterates through every slice item and calls w(Organisation) err
//
// This function is auto-generated.
func (set OrganisationSet) Walk(w func(*Organisation) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Organisation) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set OrganisationSet) Filter(f func(*Organisation) (bool, error)) (out OrganisationSet, err error) {
	var ok bool
	out = OrganisationSet{}
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
func (set OrganisationSet) FindByID(ID uint64) *Organisation {
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
func (set OrganisationSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
