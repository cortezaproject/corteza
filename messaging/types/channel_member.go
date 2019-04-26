package types

import (
	"time"
)

type (
	ChannelMember struct {
		ChannelID uint64 `db:"rel_channel"`

		UserID uint64 `db:"rel_user"`

		Type ChannelMembershipType `db:"type"`
		Flag ChannelMembershipFlag `db:"flag"`

		CreatedAt time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
	}

	ChannelMemberFilter struct {
		ComembersOf uint64
		ChannelID   uint64
		MemberID    uint64
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
