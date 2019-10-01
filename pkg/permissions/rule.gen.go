package permissions

// 	Hello! This file is auto-generated.

type (

	// RuleSet slice of Rule
	//
	// This type is auto-generated.
	RuleSet []*Rule
)

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
