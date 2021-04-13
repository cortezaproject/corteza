package messagebus

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/messagebus/types.yaml

type (

	// QueueMessageSet slice of QueueMessage
	//
	// This type is auto-generated.
	QueueMessageSet []*QueueMessage

	// QueueSettingsSet slice of QueueSettings
	//
	// This type is auto-generated.
	QueueSettingsSet []*QueueSettings
)

// Walk iterates through every slice item and calls w(QueueMessage) err
//
// This function is auto-generated.
func (set QueueMessageSet) Walk(w func(*QueueMessage) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(QueueMessage) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set QueueMessageSet) Filter(f func(*QueueMessage) (bool, error)) (out QueueMessageSet, err error) {
	var ok bool
	out = QueueMessageSet{}
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
func (set QueueMessageSet) FindByID(ID uint64) *QueueMessage {
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
func (set QueueMessageSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(QueueSettings) err
//
// This function is auto-generated.
func (set QueueSettingsSet) Walk(w func(*QueueSettings) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(QueueSettings) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set QueueSettingsSet) Filter(f func(*QueueSettings) (bool, error)) (out QueueSettingsSet, err error) {
	var ok bool
	out = QueueSettingsSet{}
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
func (set QueueSettingsSet) FindByID(ID uint64) *QueueSettings {
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
func (set QueueSettingsSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
