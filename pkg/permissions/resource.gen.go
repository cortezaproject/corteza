package permissions

// 	Hello! This file is auto-generated.

type (

	// ResourceSet slice of Resource
	//
	// This type is auto-generated.
	ResourceSet []*Resource
)

// Walk iterates through every slice item and calls w(Resource) err
//
// This function is auto-generated.
func (set ResourceSet) Walk(w func(*Resource) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Resource) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ResourceSet) Filter(f func(*Resource) (bool, error)) (out ResourceSet, err error) {
	var ok bool
	out = ResourceSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}
