package rbac

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/rbac/types.yaml

type (

	// ResourceSet slice of Resource
	//
	// This type is auto-generated.
	ResourceSet []*Resource

	// RuleSet slice of Rule
	//
	// This type is auto-generated.
	RuleSet []*Rule
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

// Walk iterates through every slice item and calls w(Rule) err
//
// This function is auto-generated.
func (set RuleSet) Walk(w func(*Rule) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Rule) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set RuleSet) Filter(f func(*Rule) (bool, error)) (out RuleSet, err error) {
	var ok bool
	out = RuleSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}
