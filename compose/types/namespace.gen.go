package types

// 	Hello! This file is auto-generated.

type (

	// NamespaceSet slice of Namespace
	//
	// This type is auto-generated.
	NamespaceSet []*Namespace
)

// Walk iterates through every slice item and calls w(Namespace) err
//
// This function is auto-generated.
func (set NamespaceSet) Walk(w func(*Namespace) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Namespace) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set NamespaceSet) Filter(f func(*Namespace) (bool, error)) (out NamespaceSet, err error) {
	var ok bool
	out = NamespaceSet{}
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
func (set NamespaceSet) FindByID(ID uint64) *Namespace {
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
func (set NamespaceSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
