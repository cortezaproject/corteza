package types

//go:generate go run ../../codegen/v2/type-set.go --types Channel --output channel.gen.go

import (
	"encoding/json"
	"fmt"
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

		CreatedAt time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty" db:"updated_at"`

		ArchivedAt *time.Time `json:"archivedAt,omitempty" db:"archived_at"`
		DeletedAt  *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`

		LastMessageID uint64 `json:",omitempty" db:"rel_last_message"`

		CanJoin           bool `json:"-" db:"-"`
		CanPart           bool `json:"-" db:"-"`
		CanObserve        bool `json:"-" db:"-"`
		CanSendMessages   bool `json:"-" db:"-"`
		CanDeleteMessages bool `json:"-" db:"-"`
		CanChangeMembers  bool `json:"-" db:"-"`
		CanUpdate         bool `json:"-" db:"-"`
		CanArchive        bool `json:"-" db:"-"`
		CanDelete         bool `json:"-" db:"-"`

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

// Scope returns permissions group that for this type
func (c *Channel) Scope() string {
	return "channel"
}

// Resource returns a RBAC resource ID for this type
func (c *Channel) Resource() string {
	return fmt.Sprintf("%s:%d", c.Scope(), c.ID)
}

// Operation returns a RBAC resource-scoped role name for an operation
func (c *Channel) Operation(name string) string {
	return fmt.Sprintf("%s/%s", c.Resource(), name)
}

func (c *Channel) IsValid() bool {
	return c.ArchivedAt == nil && c.DeletedAt == nil
}

const (
	ChannelTypePublic  ChannelType = "public"
	ChannelTypePrivate             = "private"
	ChannelTypeGroup               = "group"
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
