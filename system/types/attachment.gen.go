package types

// 	Hello! This file is auto-generated.

type (

	// AttachmentSet slice of Attachment
	//
	// This type is auto-generated.
	AttachmentSet []*Attachment
)

// Walk iterates through every slice item and calls w(Attachment) err
//
// This function is auto-generated.
func (set AttachmentSet) Walk(w func(*Attachment) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Attachment) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set AttachmentSet) Filter(f func(*Attachment) (bool, error)) (out AttachmentSet, err error) {
	var ok bool
	out = AttachmentSet{}
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
func (set AttachmentSet) FindByID(ID uint64) *Attachment {
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
func (set AttachmentSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
