package actionlog

// 	Hello! This file is auto-generated.

type (

	// ActionSet slice of Action
	//
	// This type is auto-generated.
	ActionSet []*Action
)

// Walk iterates through every slice item and calls w(Action) err
//
// This function is auto-generated.
func (set ActionSet) Walk(w func(*Action) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Action) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ActionSet) Filter(f func(*Action) (bool, error)) (out ActionSet, err error) {
	var ok bool
	out = ActionSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}
