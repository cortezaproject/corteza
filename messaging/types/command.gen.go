package types

// 	Hello! This file is auto-generated.

type (

	// CommandSet slice of Command
	//
	// This type is auto-generated.
	CommandSet []*Command
)

// Walk iterates through every slice item and calls w(Command) err
//
// This function is auto-generated.
func (set CommandSet) Walk(w func(*Command) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Command) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set CommandSet) Filter(f func(*Command) (bool, error)) (out CommandSet, err error) {
	var ok bool
	out = CommandSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}
