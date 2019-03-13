package types

// 	Hello! This file is auto-generated.

type (

	// ChartSet slice of Chart
	//
	// This type is auto-generated.
	ChartSet []*Chart
)

// Walk iterates through every slice item and calls w(Chart) err
//
// This function is auto-generated.
func (set ChartSet) Walk(w func(*Chart) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Chart) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ChartSet) Filter(f func(*Chart) (bool, error)) (out ChartSet, err error) {
	var ok bool
	out = ChartSet{}
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
func (set ChartSet) FindByID(ID uint64) *Chart {
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
func (set ChartSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
