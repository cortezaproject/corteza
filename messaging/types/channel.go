package types

import (
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"time"

	"github.com/jmoiron/sqlx/types"

	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

type (
	Channel struct {
		ID    uint64         `json:"channelID"`
		Name  string         `json:"name"`
		Topic string         `json:"topic"`
		Type  ChannelType    `json:"type"`
		Meta  types.JSONText `json:"-"`

		MembershipPolicy ChannelMembershipPolicy `json:"membershipPolicy"`

		CreatorID uint64 `json:"creatorId"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`

		ArchivedAt *time.Time `json:"archivedAt,omitempty"`
		DeletedAt  *time.Time `json:"deletedAt,omitempty"`

		LastMessageID uint64 `json:",omitempty"`

		CanJoin                   bool `json:"-"`
		CanPart                   bool `json:"-"`
		CanObserve                bool `json:"-"`
		CanSendMessages           bool `json:"-"`
		CanDeleteMessages         bool `json:"-"`
		CanDeleteOwnMessages      bool `json:"-"`
		CanUpdateMessages         bool `json:"-"`
		CanUpdateOwnMessages      bool `json:"-"`
		CanChangeMembers          bool `json:"-"`
		CanChangeMembershipPolicy bool `json:"-"`
		CanUpdate                 bool `json:"-"`
		CanArchive                bool `json:"-"`
		CanUnarchive              bool `json:"-"`
		CanDelete                 bool `json:"-"`
		CanUndelete               bool `json:"-"`

		Member  *ChannelMember `json:"-"`
		Members []uint64       `json:"-"`
		Unread  *Unread        `json:"-"`
	}

	ChannelFilter struct {
		Query string

		ChannelID []uint64

		// Only return channels accessible by this user
		CurrentUserID uint64

		// Do not filter out deleted channels
		// @deprecated
		IncludeDeleted bool

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(channel *Channel) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	ChannelMembershipPolicy string
	ChannelType             string
)

// Resource returns a system resource ID for this type
func (c Channel) RBACResource() rbac.Resource {
	return ChannelRBACResource.AppendID(c.ID)
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
