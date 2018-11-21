package types

type (

	// ModuleSet slice of Module
	ModuleSet []*Module
	// PageSet slice of Page
	PageSet []*Page
)

// Walk iterates through every slice item and calls w(Module) err
func (set ModuleSet) Walk(w func(*Module) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Module) (bool, err) and return filtered slice
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

// Finds slice item by its ID property
func (set ModuleSet) FindByID(ID uint64) *Module {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// Returns a slice of uint64s from all items in the set
func (set ModuleSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(Page) err
func (set PageSet) Walk(w func(*Page) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Page) (bool, err) and return filtered slice
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

// Finds slice item by its ID property
func (set PageSet) FindByID(ID uint64) *Page {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// Returns a slice of uint64s from all items in the set
func (set PageSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
