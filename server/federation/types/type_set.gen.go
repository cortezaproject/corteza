package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

type (
	NodeSet          []*Node
	NodeSyncSet      []*NodeSync
	ExposedModuleSet []*ExposedModule
	SharedModuleSet  []*SharedModule
	ModuleMappingSet []*ModuleMapping
)

func (set NodeSet) Walk(w func(*Node) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set NodeSet) FindByID(ID uint64) *Node {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set NodeSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set NodeSyncSet) Walk(w func(*NodeSync) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set NodeSyncSet) Filter(f func(*NodeSync) (bool, error)) (out NodeSyncSet, err error) {
	var ok bool
	out = NodeSyncSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

func (set ExposedModuleSet) Walk(w func(*ExposedModule) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set ExposedModuleSet) FindByID(ID uint64) *ExposedModule {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set ExposedModuleSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set SharedModuleSet) Walk(w func(*SharedModule) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

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

func (set SharedModuleSet) FindByID(ID uint64) *SharedModule {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set SharedModuleSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

func (set ModuleMappingSet) Walk(w func(*ModuleMapping) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set ModuleMappingSet) Filter(f func(*ModuleMapping) (bool, error)) (out ModuleMappingSet, err error) {
	var ok bool
	out = ModuleMappingSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}
