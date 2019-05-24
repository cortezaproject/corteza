package types

import (
	"time"

	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/jmoiron/sqlx/types"
)

type (
	Channel struct {
		ID    uint64         `json:"channelID" db:"id"`
		Name  string         `json:"name" db:"name"`
		Topic string         `json:"topic" db:"topic"`
		Type  ChannelType    `json:"type" db:"type"`
		Meta  types.JSONText `json:"-" db:"meta"`

		CreatorID      uint64 `json:"creatorId" db:"rel_creator"`
		OrganisationID uint64 `json:"organisationId" db:"rel_organisation"`

		CreatedAt time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty" db:"updated_at"`

		ArchivedAt *time.Time `json:"archivedAt,omitempty" db:"archived_at"`
		DeletedAt  *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`

		LastMessageID uint64 `json:",omitempty" db:"rel_last_message"`

		CanJoin              bool `json:"-" db:"-"`
		CanPart              bool `json:"-" db:"-"`
		CanObserve           bool `json:"-" db:"-"`
		CanSendMessages      bool `json:"-" db:"-"`
		CanDeleteMessages    bool `json:"-" db:"-"`
		CanDeleteOwnMessages bool `json:"-" db:"-"`
		CanUpdateMessages    bool `json:"-" db:"-"`
		CanUpdateOwnMessages bool `json:"-" db:"-"`
		CanChangeMembers     bool `json:"-" db:"-"`
		CanUpdate            bool `json:"-" db:"-"`
		CanArchive           bool `json:"-" db:"-"`
		CanUnarchive         bool `json:"-" db:"-"`
		CanDelete            bool `json:"-" db:"-"`
		CanUndelete          bool `json:"-" db:"-"`

		Member  *ChannelMember `json:"-" db:"-"`
		Members []uint64       `json:"-" db:"-"`
		Unread  *Unread        `json:"-" db:"-"`
	}

	ChannelFilter struct {
		Query string

		// Only return channels accessible by this user
		CurrentUserID uint64

		// Do not filter out deleted channels
		IncludeDeleted bool
	}

	ChannelType string
)

// Resource returns a system resource ID for this type
func (c Channel) PermissionResource() permissions.Resource {
	return ChannelPermissionResource.AppendID(c.ID)
}

func (c *Channel) IsValid() bool {
	return c.ArchivedAt == nil && c.DeletedAt == nil
}

const (
	ChannelTypePublic  ChannelType = "public"
	ChannelTypePrivate ChannelType = "private"
	ChannelTypeGroup   ChannelType = "group"
)

func (mtype ChannelType) String() string {
	return string(mtype)
}

func (mtype ChannelType) IsValid() bool {
	switch mtype {
	case ChannelTypePublic,
		ChannelTypePrivate,
		ChannelTypeGroup:
		return true
	}

	return false
}
