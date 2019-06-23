package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

// MembersOf extracts member IDs from channel member set
//
// It filters out only members that match a particular channel
func (set ChannelMemberSet) MembersOf(channelID uint64) []uint64 {
	var mmof = make([]uint64, 0)

	for i := range set {
		if set[i].ChannelID == channelID {
			mmof = append(mmof, set[i].UserID)
		}
	}

	return mmof
}

// AllMemberIDs returns IDs of all members
func (set ChannelMemberSet) AllMemberIDs() []uint64 {
	var mmof = make([]uint64, 0)

	for i := range set {
		mmof = append(mmof, set[i].UserID)
	}

	return mmof
}

func (set ChannelMemberSet) FindByUserID(userID uint64) *ChannelMember {
	for i := range set {
		if set[i].UserID == userID {
			return set[i]
		}
	}

	return nil
}

func (set ChannelMemberSet) FindByChannelID(channelID uint64) (out ChannelMemberSet) {
	out = ChannelMemberSet{}

	for i := range set {
		if set[i].ChannelID == channelID {
			out = append(out, set[i])
		}
	}

	return
}

func (set *CommandParamSet) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*set = CommandParamSet{}
	case []uint8:
		if err := json.Unmarshal(value.([]byte), set); err != nil {
			return errors.Wrapf(err, "Can not scan '%v' into CommandParamSet", value)
		}
	}

	return nil
}

func (set CommandParamSet) Value() (driver.Value, error) {
	return json.Marshal(set)
}

func (set MessageFlagSet) IsBookmarked(UserID uint64) bool {
	for i := range set {
		if set[i].UserID == UserID && set[i].IsBookmark() {
			return true
		}
	}

	return false
}

func (set MessageFlagSet) IsPinned() bool {
	for i := range set {
		if set[i].IsPin() {
			return true
		}
	}

	return false
}

func (set MentionSet) FindByUserID(ID uint64) (out MentionSet) {
	out = MentionSet{}

	for i := range set {
		if set[i].UserID == ID {
			out = append(out, set[i])
		}
	}

	return
}

func (set MentionSet) FindByMessageID(ID uint64) (out MentionSet) {
	out = MentionSet{}

	for i := range set {
		if set[i].MessageID == ID {
			out = append(out, set[i])
		}
	}

	return
}

func (set MentionSet) UserIDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].UserID
	}

	return
}

func (set MentionSet) Diff(in MentionSet) (add, upd, del MentionSet) {
	add, upd, del = MentionSet{}, MentionSet{}, MentionSet{}

	for _, m := range in {
		if m.ID == 0 {
			// Mark for adding all new
			add = append(add, m)
		}
	}

	for _, m := range set {
		if m.ID == 0 {
			// Ignore all unsaved
			continue
		}

		if in.FindByID(m.ID) == nil {
			// Mark for removal all that are not added
			del = append(del, m)
		} else {
			// Mark for update all that are still there
			upd = append(upd, m)
		}
	}

	return
}

func (set UnreadSet) FindByChannelId(channelID uint64) *Unread {
	for i := range set {
		if set[i].ChannelID == channelID {
			return set[i]
		}
	}

	return nil
}

func (set UnreadSet) FindByThreadId(threadID uint64) *Unread {
	for i := range set {
		if set[i].ReplyTo == threadID {
			return set[i]
		}
	}

	return nil
}
