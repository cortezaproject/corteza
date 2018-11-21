package types

// 	Hello! This file is auto-generated.

type (

	// MessageAttachmentSet slice of MessageAttachment
	//
	// This type is auto-generated.
	MessageAttachmentSet []*MessageAttachment
)

// Walk iterates through every slice item and calls w(MessageAttachment) err
//
// This function is auto-generated.
func (set MessageAttachmentSet) Walk(w func(*MessageAttachment) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(MessageAttachment) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set MessageAttachmentSet) Filter(f func(*MessageAttachment) (bool, error)) (out MessageAttachmentSet, err error) {
	var ok bool
	out = MessageAttachmentSet{}
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
func (set MessageAttachmentSet) FindByID(ID uint64) *MessageAttachment {
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
func (set MessageAttachmentSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
