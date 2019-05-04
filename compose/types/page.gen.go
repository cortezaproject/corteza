package types

// 	Hello! This file is auto-generated.

type (

	// PageSet slice of Page
	//
	// This type is auto-generated.
	PageSet []*Page
)

// Walk iterates through every slice item and calls w(Page) err
//
// This function is auto-generated.
func (set PageSet) Walk(w func(*Page) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Page) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set PageSet) Filter(f func(*Page) (bool, error)) (out PageSet, err error) {
	var ok bool
	out = PageSet{}
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
func (set PageSet) FindByID(ID uint64) *Page {
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
func (set PageSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
