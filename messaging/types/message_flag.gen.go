package types

// 	Hello! This file is auto-generated.

type (

	// MessageFlagSet slice of MessageFlag
	//
	// This type is auto-generated.
	MessageFlagSet []*MessageFlag
)

// Walk iterates through every slice item and calls w(MessageFlag) err
//
// This function is auto-generated.
func (set MessageFlagSet) Walk(w func(*MessageFlag) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(MessageFlag) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set MessageFlagSet) Filter(f func(*MessageFlag) (bool, error)) (out MessageFlagSet, err error) {
	var ok bool
	out = MessageFlagSet{}
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
func (set MessageFlagSet) FindByID(ID uint64) *MessageFlag {
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
func (set MessageFlagSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
