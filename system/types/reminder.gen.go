package types

// 	Hello! This file is auto-generated.

type (

	// ReminderSet slice of Reminder
	//
	// This type is auto-generated.
	ReminderSet []*Reminder
)

// Walk iterates through every slice item and calls w(Reminder) err
//
// This function is auto-generated.
func (set ReminderSet) Walk(w func(*Reminder) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Reminder) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ReminderSet) Filter(f func(*Reminder) (bool, error)) (out ReminderSet, err error) {
	var ok bool
	out = ReminderSet{}
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
func (set ReminderSet) FindByID(ID uint64) *Reminder {
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
func (set ReminderSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
