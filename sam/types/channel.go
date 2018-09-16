package types

import (
	"encoding/json"
	"time"
)

type (
	Channel struct {
		ID    uint64          `json:"id" db:"id"`
		Name  string          `json:"name" db:"name"`
		Topic string          `json:"topic" db:"topic"`
		Type  ChannelType     `json:"type" db:"type"`
		Meta  json.RawMessage `json:"-" db:"meta"`

		CreatorID      uint64 `json:"creatorId" db:"rel_creator"`
		OrganisationID uint64 `json:"organisationId" db:"rel_organisation"`

		CreatedAt  time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt  *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
		ArchivedAt *time.Time `json:"archivedAt,omitempty" db:"archived_at"`
		DeletedAt  *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`

		LastMessageID uint64 `json:",omitempty" db:"rel_last_message"`

		Member  *ChannelMember `json:"-" db:"-"`
		Members []*uint64      `json:"-" db:"-"`
	}

	ChannelMember struct {
		ChannelID uint64 `db:"rel_channel"`
		UserID    uint64 `db:"rel_user"`

		Type ChannelMembershipType `db:"type"`

		CreatedAt time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
	}

	ChannelFilter struct {
		Query          string
		IncludeMembers bool
	}

	ChannelMembershipType string
	ChannelType           string

	ChannelSet []*Channel
)

func (cc ChannelSet) Walk(w func(*Channel) error) (err error) {
	for i := range cc {
		if err = w(cc[i]); err != nil {
			return
		}
	}

	return
}

const (
	ChannelMembershipTypeOwner  ChannelMembershipType = "owner"
	ChannelMembershipTypeMember                       = "member"

	ChannelTypePublic  ChannelType = "public"
	ChannelTypePrivate             = "private"
	ChannelTypeGroup               = "group"
	ChannelTypeDirect              = "direct"
)
