package types

// 	Hello! This file is auto-generated.

type (

	// ModuleFieldSet slice of ModuleField
	//
	// This type is auto-generated.
	ModuleFieldSet []*ModuleField
)

// Walk iterates through every slice item and calls w(ModuleField) err
//
// This function is auto-generated.
func (set ModuleFieldSet) Walk(w func(*ModuleField) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(ModuleField) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ModuleFieldSet) Filter(f func(*ModuleField) (bool, error)) (out ModuleFieldSet, err error) {
	var ok bool
	out = ModuleFieldSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}
