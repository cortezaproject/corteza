package types

import (
	"time"
)

type (
	ChannelMember struct {
		ChannelID uint64

		UserID uint64

		Type ChannelMembershipType
		Flag ChannelMembershipFlag

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	}

	ChannelMemberFilter struct {
		ChannelID []uint64
		MemberID  []uint64
	}

	ChannelMembershipType string
	ChannelMembershipFlag string
)

const (
	ChannelMembershipTypeOwner   ChannelMembershipType = "owner"
	ChannelMembershipTypeMember  ChannelMembershipType = "member"
	ChannelMembershipTypeInvitee ChannelMembershipType = "invitee"

	ChannelMembershipFlagPinned  ChannelMembershipFlag = "pinned"
	ChannelMembershipFlagHidden  ChannelMembershipFlag = "hidden"
	ChannelMembershipFlagIgnored ChannelMembershipFlag = "ignored"
	ChannelMembershipFlagNone    ChannelMembershipFlag = ""
)

// ChannelMemberFilterChannels helper func for building channel member filter with list of channels
func ChannelMemberFilterChannels(ID ...uint64) ChannelMemberFilter {
	return ChannelMemberFilter{ChannelID: ID}
}
