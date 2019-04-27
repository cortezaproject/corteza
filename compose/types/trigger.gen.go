package types

// 	Hello! This file is auto-generated.

type (

	// TriggerSet slice of Trigger
	//
	// This type is auto-generated.
	TriggerSet []*Trigger
)

// Walk iterates through every slice item and calls w(Trigger) err
//
// This function is auto-generated.
func (set TriggerSet) Walk(w func(*Trigger) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Trigger) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set TriggerSet) Filter(f func(*Trigger) (bool, error)) (out TriggerSet, err error) {
	var ok bool
	out = TriggerSet{}
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
func (set TriggerSet) FindByID(ID uint64) *Trigger {
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
func (set TriggerSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
