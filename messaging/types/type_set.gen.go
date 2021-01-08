package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// messaging/types/types.yaml

type (

	// AttachmentSet slice of Attachment
	//
	// This type is auto-generated.
	AttachmentSet []*Attachment

	// ChannelSet slice of Channel
	//
	// This type is auto-generated.
	ChannelSet []*Channel

	// ChannelMemberSet slice of ChannelMember
	//
	// This type is auto-generated.
	ChannelMemberSet []*ChannelMember

	// CommandSet slice of Command
	//
	// This type is auto-generated.
	CommandSet []*Command

	// CommandParamSet slice of CommandParam
	//
	// This type is auto-generated.
	CommandParamSet []*CommandParam

	// MentionSet slice of Mention
	//
	// This type is auto-generated.
	MentionSet []*Mention

	// MessageSet slice of Message
	//
	// This type is auto-generated.
	MessageSet []*Message

	// MessageAttachmentSet slice of MessageAttachment
	//
	// This type is auto-generated.
	MessageAttachmentSet []*MessageAttachment

	// MessageFlagSet slice of MessageFlag
	//
	// This type is auto-generated.
	MessageFlagSet []*MessageFlag

	// UnreadSet slice of Unread
	//
	// This type is auto-generated.
	UnreadSet []*Unread
)

// Walk iterates through every slice item and calls w(Attachment) err
//
// This function is auto-generated.
func (set AttachmentSet) Walk(w func(*Attachment) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Attachment) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set AttachmentSet) Filter(f func(*Attachment) (bool, error)) (out AttachmentSet, err error) {
	var ok bool
	out = AttachmentSet{}
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
func (set AttachmentSet) FindByID(ID uint64) *Attachment {
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
func (set AttachmentSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(Channel) err
//
// This function is auto-generated.
func (set ChannelSet) Walk(w func(*Channel) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Channel) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set ChannelSet) Filter(f func(*Channel) (bool, error)) (out ChannelSet, err error) {
	var ok bool
	out = ChannelSet{}
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
func (set ChannelSet) FindByID(ID uint64) *Channel {
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
func (set ChannelSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

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

// Walk iterates through every slice item and calls w(CommandParam) err
//
// This function is auto-generated.
func (set CommandParamSet) Walk(w func(*CommandParam) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(CommandParam) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set CommandParamSet) Filter(f func(*CommandParam) (bool, error)) (out CommandParamSet, err error) {
	var ok bool
	out = CommandParamSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// Walk iterates through every slice item and calls w(Mention) err
//
// This function is auto-generated.
func (set MentionSet) Walk(w func(*Mention) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Mention) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set MentionSet) Filter(f func(*Mention) (bool, error)) (out MentionSet, err error) {
	var ok bool
	out = MentionSet{}
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
func (set MentionSet) FindByID(ID uint64) *Mention {
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
func (set MentionSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(Message) err
//
// This function is auto-generated.
func (set MessageSet) Walk(w func(*Message) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(Message) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set MessageSet) Filter(f func(*Message) (bool, error)) (out MessageSet, err error) {
	var ok bool
	out = MessageSet{}
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
func (set MessageSet) FindByID(ID uint64) *Message {
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
func (set MessageSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

// Walk iterates through every slice item and calls w(MessageAttachment) err
//
// This function is auto-generated.
func (set MessageAttachmentSet) Walk(w func(*MessageAttachment) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(MessageAttachment) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set MessageAttachmentSet) Filter(f func(*MessageAttachment) (bool, error)) (out MessageAttachmentSet, err error) {
	var ok bool
	out = MessageAttachmentSet{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// Walk iterates through every slice item and calls w(MessageFlag) err
//
// This function is auto-generated.
func (set MessageFlagSet) Walk(w func(*MessageFlag) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f(MessageFlag) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set MessageFlagSet) Filter(f func(*MessageFlag) (bool, error)) (out MessageFlagSet, err error) {
	var ok bool
	out = MessageFlagSet{}
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
func (set MessageFlagSet) FindByID(ID uint64) *MessageFlag {
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
func (set MessageFlagSet) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}

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
