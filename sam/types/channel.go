package types

import (
	"encoding/json"
	"time"
)

type (
	Channel struct {
		ID    uint64          `db:"id"`
		Name  string          `db:"name"`
		Topic string          `db:"topic"`
		Type  ChannelType     `db:"type"`
		Meta  json.RawMessage `db:"meta"`

		CreatorID      uint64 `db:"rel_creator"`
		OrganisationID uint64 `db:"rel_organisation"`

		CreatedAt  time.Time  `json:"created_at,omitempty" db:"created_at"`
		UpdatedAt  *time.Time `json:"updated_at,omitempty" db:"updated_at"`
		ArchivedAt *time.Time `json:"archived_at,omitempty" db:"archived_at"`
		DeletedAt  *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`

		LastMessageID uint64 `json:",omitempty" db:"rel_last_message"`
	}

	ChannelMember struct {
		ChannelID uint64 `db:"rel_channel"`
		UserID    uint64 `db:"rel_user"`

		Type ChannelMembershipType `db:"type"`

		CreatedAt time.Time  `json:"created_at,omitempty" db:"created_at"`
		UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	}

	ChannelFilter struct {
		Query string
	}

	ChannelMembershipType string
	ChannelType           string
)

const (
	ChannelMembershipTypeOwner  ChannelMembershipType = "owner"
	ChannelMembershipTypeMember                       = "member"

	ChannelTypePublic  ChannelType = "public"
	ChannelTypePrivate             = "private"
	ChannelTypeDirect              = "direct"
)
