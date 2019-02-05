package types

// 	Hello! This file is auto-generated.

type (

	// MentionSet slice of Mention
	//
	// This type is auto-generated.
	MentionSet []*Mention
)

// Walk iterates through every slice item and calls w(Mention) err
//
// This function is auto-generated.
func (set MentionSet) Walk(w func(*Mention) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Mention) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set MentionSet) Filter(f func(*Mention) (bool, error)) (out MentionSet, err error) {
	var ok bool
	out = MentionSet{}
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
func (set MentionSet) FindByID(ID uint64) *Mention {
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
func (set MentionSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
