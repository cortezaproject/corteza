package types

import (
	"time"

	"github.com/jmoiron/sqlx/types"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

type (
	Channel struct {
		ID    uint64         `json:"channelID" db:"id"`
		Name  string         `json:"name" db:"name"`
		Topic string         `json:"topic" db:"topic"`
		Type  ChannelType    `json:"type" db:"type"`
		Meta  types.JSONText `json:"-" db:"meta"`

		MembershipPolicy ChannelMembershipPolicy `json:"membershipPolicy" db:"membership_policy""`

		CreatorID      uint64 `json:"creatorId" db:"rel_creator"`
		OrganisationID uint64 `json:"organisationId" db:"rel_organisation"`

		CreatedAt time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty" db:"updated_at"`

		ArchivedAt *time.Time `json:"archivedAt,omitempty" db:"archived_at"`
		DeletedAt  *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`

		LastMessageID uint64 `json:",omitempty" db:"rel_last_message"`

		CanJoin                   bool `json:"-" db:"-"`
		CanPart                   bool `json:"-" db:"-"`
		CanObserve                bool `json:"-" db:"-"`
		CanSendMessages           bool `json:"-" db:"-"`
		CanDeleteMessages         bool `json:"-" db:"-"`
		CanDeleteOwnMessages      bool `json:"-" db:"-"`
		CanUpdateMessages         bool `json:"-" db:"-"`
		CanUpdateOwnMessages      bool `json:"-" db:"-"`
		CanChangeMembers          bool `json:"-" db:"-"`
		CanChangeMembershipPolicy bool `json:"-" db:"-"`
		CanUpdate                 bool `json:"-" db:"-"`
		CanArchive                bool `json:"-" db:"-"`
		CanUnarchive              bool `json:"-" db:"-"`
		CanDelete                 bool `json:"-" db:"-"`
		CanUndelete               bool `json:"-" db:"-"`

		Member  *ChannelMember `json:"-" db:"-"`
		Members []uint64       `json:"-" db:"-"`
		Unread  *Unread        `json:"-" db:"-"`
	}

	ChannelFilter struct {
		Query string

		ChannelID []uint64

		// Only return channels accessible by this user
		CurrentUserID uint64

		// Do not filter out deleted channels
		IncludeDeleted bool

		Sort string `json:"sort"`
	}

	ChannelMembershipPolicy string
	ChannelType             string
)

// Resource returns a system resource ID for this type
func (c Channel) PermissionResource() permissions.Resource {
	return ChannelPermissionResource.AppendID(c.ID)
}

func (c Channel) DynamicRoles(userID uint64) []uint64 {
	return nil
}

func (c *Channel) IsValid() bool {
	return c.ArchivedAt == nil && c.DeletedAt == nil
}

const (
	ChannelTypePublic  ChannelType = "public"
	ChannelTypePrivate ChannelType = "private"
	ChannelTypeGroup   ChannelType = "group"

	ChannelMembershipPolicyFeatured ChannelMembershipPolicy = "featured"
	ChannelMembershipPolicyForced   ChannelMembershipPolicy = "forced"
	ChannelMembershipPolicyDefault  ChannelMembershipPolicy = ""
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

func (cm ChannelMembershipPolicy) String() string {
	return string(cm)
}

func (cm ChannelMembershipPolicy) IsValid() bool {
	switch cm {
	case ChannelMembershipPolicyFeatured,
		ChannelMembershipPolicyForced,
		ChannelMembershipPolicyDefault:
		return true
	}

	return false
}

// FindByName finds items from slice by its name
func (set ChannelSet) FindByName(name string) *Channel {
	for i := range set {
		if set[i].Name == name {
			return set[i]
		}
	}

	return nil
}
