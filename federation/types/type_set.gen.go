package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// federation/types/types.yaml

type (

	// ExposedModuleSet slice of ExposedModule
	//
	// This type is auto-generated.
	ExposedModuleSet []*ExposedModule

	// NodeSet slice of Node
	//
	// This type is auto-generated.
	NodeSet []*Node

	// SharedModuleSet slice of SharedModule
	//
	// This type is auto-generated.
	SharedModuleSet []*SharedModule
)

// Walk iterates through every slice item and calls w(ExposedModule) err
//
// This function is auto-generated.
func (set ExposedModuleSet) Walk(w func(*ExposedModule) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(ExposedModule) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ExposedModuleSet) Filter(f func(*ExposedModule) (bool, error)) (out ExposedModuleSet, err error) {
	var ok bool
	out = ExposedModuleSet{}
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
func (set ExposedModuleSet) FindByID(ID uint64) *ExposedModule {
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
func (set ExposedModuleSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(Node) err
//
// This function is auto-generated.
func (set NodeSet) Walk(w func(*Node) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Node) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set NodeSet) Filter(f func(*Node) (bool, error)) (out NodeSet, err error) {
	var ok bool
	out = NodeSet{}
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
func (set NodeSet) FindByID(ID uint64) *Node {
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
func (set NodeSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(SharedModule) err
//
// This function is auto-generated.
func (set SharedModuleSet) Walk(w func(*SharedModule) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(SharedModule) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set SharedModuleSet) Filter(f func(*SharedModule) (bool, error)) (out SharedModuleSet, err error) {
	var ok bool
	out = SharedModuleSet{}
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
func (set SharedModuleSet) FindByID(ID uint64) *SharedModule {
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
func (set SharedModuleSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
