package types

// 	Hello! This file is auto-generated.

type (

	// ChannelSet slice of Channel
	//
	// This type is auto-generated.
	ChannelSet []*Channel
)

// Walk iterates through every slice item and calls w(Channel) err
//
// This function is auto-generated.
func (set ChannelSet) Walk(w func(*Channel) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Channel) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ChannelSet) Filter(f func(*Channel) (bool, error)) (out ChannelSet, err error) {
	var ok bool
	out = ChannelSet{}
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
func (set ChannelSet) FindByID(ID uint64) *Channel {
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
func (set ChannelSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
