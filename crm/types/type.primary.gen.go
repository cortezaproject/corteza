package types

// 	Hello! This file is auto-generated.

type (

	// ModuleSet slice of Module
	//
	// This type is auto-generated.
	ModuleSet []*Module

	// PageSet slice of Page
	//
	// This type is auto-generated.
	PageSet []*Page

	// ChartSet slice of Chart
	//
	// This type is auto-generated.
	ChartSet []*Chart

	// TriggerSet slice of Trigger
	//
	// This type is auto-generated.
	TriggerSet []*Trigger

	// RecordSet slice of Record
	//
	// This type is auto-generated.
	RecordSet []*Record
)

// Walk iterates through every slice item and calls w(Module) err
//
// This function is auto-generated.
func (set ModuleSet) Walk(w func(*Module) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Module) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ModuleSet) Filter(f func(*Module) (bool, error)) (out ModuleSet, err error) {
	var ok bool
	out = ModuleSet{}
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
func (set ModuleSet) FindByID(ID uint64) *Module {
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
func (set ModuleSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(Page) err
//
// This function is auto-generated.
func (set PageSet) Walk(w func(*Page) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Page) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set PageSet) Filter(f func(*Page) (bool, error)) (out PageSet, err error) {
	var ok bool
	out = PageSet{}
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
func (set PageSet) FindByID(ID uint64) *Page {
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
func (set PageSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(Chart) err
//
// This function is auto-generated.
func (set ChartSet) Walk(w func(*Chart) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Chart) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ChartSet) Filter(f func(*Chart) (bool, error)) (out ChartSet, err error) {
	var ok bool
	out = ChartSet{}
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
func (set ChartSet) FindByID(ID uint64) *Chart {
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
func (set ChartSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

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

// Walk iterates through every slice item and calls w(Record) err
//
// This function is auto-generated.
func (set RecordSet) Walk(w func(*Record) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Record) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set RecordSet) Filter(f func(*Record) (bool, error)) (out RecordSet, err error) {
	var ok bool
	out = RecordSet{}
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
func (set RecordSet) FindByID(ID uint64) *Record {
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
func (set RecordSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
