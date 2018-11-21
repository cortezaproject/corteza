package types

// 	Hello! This file is auto-generated.

type (

	// UnreadSet slice of Unread
	//
	// This type is auto-generated.
	UnreadSet []*Unread
)

// Walk iterates through every slice item and calls w(Unread) err
//
// This function is auto-generated.
func (set UnreadSet) Walk(w func(*Unread) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Unread) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set UnreadSet) Filter(f func(*Unread) (bool, error)) (out UnreadSet, err error) {
	var ok bool
	out = UnreadSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}
