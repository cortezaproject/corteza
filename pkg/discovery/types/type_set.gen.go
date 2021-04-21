package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/discovery/types.yaml

type (

	// ResourceActivitySet slice of ResourceActivity
	//
	// This type is auto-generated.
	ResourceActivitySet []*ResourceActivity
)

// Walk iterates through every slice item and calls w(ResourceActivity) err
//
// This function is auto-generated.
func (set ResourceActivitySet) Walk(w func(*ResourceActivity) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(ResourceActivity) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ResourceActivitySet) Filter(f func(*ResourceActivity) (bool, error)) (out ResourceActivitySet, err error) {
	var ok bool
	out = ResourceActivitySet{}
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
func (set ResourceActivitySet) FindByID(ID uint64) *ResourceActivity {
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
func (set ResourceActivitySet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
