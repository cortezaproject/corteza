package types

// 	Hello! This file is auto-generated.

type (

	// MessageSet slice of Message
	//
	// This type is auto-generated.
	MessageSet []*Message
)

// Walk iterates through every slice item and calls w(Message) err
//
// This function is auto-generated.
func (set MessageSet) Walk(w func(*Message) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Message) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set MessageSet) Filter(f func(*Message) (bool, error)) (out MessageSet, err error) {
	var ok bool
	out = MessageSet{}
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
func (set MessageSet) FindByID(ID uint64) *Message {
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
func (set MessageSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
