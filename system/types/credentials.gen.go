package types

// 	Hello! This file is auto-generated.

type (

	// CredentialsSet slice of Credentials
	//
	// This type is auto-generated.
	CredentialsSet []*Credentials
)

// Walk iterates through every slice item and calls w(Credentials) err
//
// This function is auto-generated.
func (set CredentialsSet) Walk(w func(*Credentials) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Credentials) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set CredentialsSet) Filter(f func(*Credentials) (bool, error)) (out CredentialsSet, err error) {
	var ok bool
	out = CredentialsSet{}
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
func (set CredentialsSet) FindByID(ID uint64) *Credentials {
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
func (set CredentialsSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
