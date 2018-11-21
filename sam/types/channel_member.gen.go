package types

// 	Hello! This file is auto-generated.

type (

	// ChannelMemberSet slice of ChannelMember
	//
	// This type is auto-generated.
	ChannelMemberSet []*ChannelMember
)

// Walk iterates through every slice item and calls w(ChannelMember) err
//
// This function is auto-generated.
func (set ChannelMemberSet) Walk(w func(*ChannelMember) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(ChannelMember) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ChannelMemberSet) Filter(f func(*ChannelMember) (bool, error)) (out ChannelMemberSet, err error) {
	var ok bool
	out = ChannelMemberSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}
