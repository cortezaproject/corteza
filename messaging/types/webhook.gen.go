package types

// 	Hello! This file is auto-generated.

type (

	// WebhookSet slice of Webhook
	//
	// This type is auto-generated.
	WebhookSet []*Webhook
)

// Walk iterates through every slice item and calls w(Webhook) err
//
// This function is auto-generated.
func (set WebhookSet) Walk(w func(*Webhook) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Webhook) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set WebhookSet) Filter(f func(*Webhook) (bool, error)) (out WebhookSet, err error) {
	var ok bool
	out = WebhookSet{}
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
func (set WebhookSet) FindByID(ID uint64) *Webhook {
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
func (set WebhookSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
