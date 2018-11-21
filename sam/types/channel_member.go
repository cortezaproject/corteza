package types

//go:generate go run ../../codegen/v2/type-set.go --no-pk-types ChannelMember --output channel_member.gen.go

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
)

const (
	ChannelMembershipTypeOwner   ChannelMembershipType = "owner"
	ChannelMembershipTypeMember                        = "member"
	ChannelMembershipTypeInvitee                       = "invitee"
)
