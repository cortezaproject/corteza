package types

import (
	"time"

	systemTypes "github.com/crusttech/crust/system/types"
)

type (
	ChannelMember struct {
		ChannelID uint64 `db:"rel_channel"`

		UserID uint64            `db:"rel_user"`
		User   *systemTypes.User `db:"-"`

		Type ChannelMembershipType `db:"type"`

		CreatedAt time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
	}

	ChannelMemberFilter struct {
		ComembersOf uint64
		ChannelID   uint64
		MemberID    uint64
	}

	ChannelMembershipType string

	ChannelMemberSet []*ChannelMember
)

func (mm ChannelMemberSet) Walk(w func(*ChannelMember) error) (err error) {
	for i := range mm {
		if err = w(mm[i]); err != nil {
			return
		}
	}

	return
}

// MembersOf extracts member IDs from channel member set
//
// It filters out only members that match a particular channel
func (mm ChannelMemberSet) MembersOf(channelID uint64) []uint64 {
	var mmof = make([]uint64, 0)

	for i := range mm {
		if mm[i].ChannelID == channelID {
			mmof = append(mmof, mm[i].UserID)
		}
	}

	return mmof
}

// AllMemberIDs returns IDs of all members
func (mm ChannelMemberSet) AllMemberIDs() []uint64 {
	var mmof = make([]uint64, 0)

	for i := range mm {
		mmof = append(mmof, mm[i].UserID)
	}

	return mmof
}

func (mm ChannelMemberSet) FindByUserID(userID uint64) *ChannelMember {
	for i := range mm {
		if mm[i].UserID == userID {
			return mm[i]
		}
	}

	return nil
}

func (mm ChannelMemberSet) FindByChannelID(channelID uint64) (out ChannelMemberSet) {
	out = ChannelMemberSet{}

	for i := range mm {
		if mm[i].ChannelID == channelID {
			out = append(out, mm[i])
		}
	}

	return
}

const (
	ChannelMembershipTypeOwner   ChannelMembershipType = "owner"
	ChannelMembershipTypeMember                        = "member"
	ChannelMembershipTypeInvitee                       = "invitee"
)
