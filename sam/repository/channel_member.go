package repository

import (
	"context"
	"time"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/sam/types"
)

type (
	// ChannelMemberRepository interface to channel member repository
	ChannelMemberRepository interface {
		With(ctx context.Context, db *factory.DB) ChannelMemberRepository

		Find(filter *types.ChannelMemberFilter) (types.ChannelMemberSet, error)

		Create(mod *types.ChannelMember) (*types.ChannelMember, error)
		Update(mod *types.ChannelMember) (*types.ChannelMember, error)
		Delete(channelID, userID uint64) error
	}

	channelMember struct {
		*repository
	}
)

const (
	// Copy definitions to make it more obvious that we're reusing channel-scope sql
	sqlChannelMemberChannelValidOnly = sqlChannelValidOnly

	// Copy definitions to make it more obvious that we're reusing channel-scope sql
	sqlChannelMemberChannelAccess = sqlChannelAccess

	// Fetching channel members of all channels a specific user has access to
	sqlChannelMemberSelect = `SELECT m.*
        FROM channel_members AS m
             INNER JOIN channels AS c ON (m.rel_channel = c.id)
       WHERE ` + sqlChannelMemberChannelValidOnly

	// Selects all user's memberships
	sqlChannelMemberships = `SELECT *
        FROM channel_members AS cm
       WHERE true`
)

// ChannelMember creates new instance of channel member repository
func ChannelMember(ctx context.Context, db *factory.DB) ChannelMemberRepository {
	return (&channelMember{}).With(ctx, db)
}

// With context...
func (r *channelMember) With(ctx context.Context, db *factory.DB) ChannelMemberRepository {
	return &channelMember{
		repository: r.repository.With(ctx, db),
	}
}

// FindMembers fetches membership info
//
// If channelID > 0 it returns members of a specific channel
// If userID    > 0 it returns members of all channels this user is member of
func (r *channelMember) Find(filter *types.ChannelMemberFilter) (types.ChannelMemberSet, error) {
	params := make([]interface{}, 0)
	mm := types.ChannelMemberSet{}

	sql := sqlChannelMemberSelect

	if filter != nil {
		if filter.ComembersOf > 0 {
			// scope: only channel we have access to
			sql += " AND m.rel_channel IN " + sqlChannelMemberChannelAccess
			params = append(params, filter.ComembersOf, types.ChannelTypePublic)
		}

		if filter.MemberID > 0 {
			sql += " AND m.rel_user = ?"
			params = append(params, filter.MemberID)
		}

		if filter.ChannelID > 0 {
			sql += " AND m.rel_channel = ?"
			params = append(params, filter.ChannelID)
		}
	}

	return mm, r.db().Select(&mm, sql, params...)
}

// Create adds channel membership record
func (r *channelMember) Create(mod *types.ChannelMember) (*types.ChannelMember, error) {
	mod.CreatedAt = time.Now()
	mod.UpdatedAt = nil

	return mod, r.db().Insert("channel_members", mod)
}

// Update modifies existing channel membership record
func (r *channelMember) Update(mod *types.ChannelMember) (*types.ChannelMember, error) {
	mod.UpdatedAt = timeNowPtr()

	whitelist := []string{"type", "updated_at", "rel_channel", "rel_user"}

	return mod, r.db().UpdatePartial("channel_members", mod, whitelist, "rel_channel", "rel_user")
}

// Delete removes existing channel membership record
func (r *channelMember) Delete(channelID, userID uint64) error {
	sql := `DELETE FROM channel_members WHERE rel_channel = ? AND rel_user = ?`
	return exec(r.db().Exec(sql, channelID, userID))
}
