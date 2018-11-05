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

func (uu ChannelMemberSet) FindByUserId(userID uint64) *ChannelMember {
	for i := range uu {
		if uu[i].UserID == userID {
			return uu[i]
		}
	}

	return nil
}

const (
	ChannelMembershipTypeOwner   ChannelMembershipType = "owner"
	ChannelMembershipTypeMember                        = "member"
	ChannelMembershipTypeInvitee                       = "invitee"
)
