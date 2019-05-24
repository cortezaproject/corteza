package repository

import (
	"context"
	"time"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	// ChannelMemberRepository interface to channel member repository
	ChannelMemberRepository interface {
		With(ctx context.Context, db *factory.DB) ChannelMemberRepository

		Find(filter *types.ChannelMemberFilter) (types.ChannelMemberSet, error)

		Create(mod *types.ChannelMember) (*types.ChannelMember, error)
		Update(mod *types.ChannelMember) (*types.ChannelMember, error)
		Delete(channelID, userID uint64) error

		CountMemberships(userID uint64) (c int, err error)
		ChangeMembership(userID, target uint64) error
	}

	channelMember struct {
		*repository
	}
)

const (
	// Copy definitions to make it more obvious that we're reusing channel-scope sql
	sqlChannelMemberChannelAccess = sqlChannelAccess

	// Fetching channel members of all channels a specific user has access to
	sqlChannelMemberSelect = `SELECT m.*
        FROM messaging_channel_member AS m
             INNER JOIN messaging_channel AS c ON (m.rel_channel = c.id)
       WHERE c.archived_at IS NULL         
         AND c.deleted_at IS NULL`
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

	return mod, r.db().Insert("messaging_channel_member", mod)
}

// Update modifies existing channel membership record
func (r *channelMember) Update(mod *types.ChannelMember) (*types.ChannelMember, error) {
	mod.UpdatedAt = timeNowPtr()

	whitelist := []string{"type", "flag", "updated_at", "rel_channel", "rel_user"}

	return mod, r.db().UpdatePartial("messaging_channel_member", mod, whitelist, "rel_channel", "rel_user")
}

// Delete removes existing channel membership record
func (r *channelMember) Delete(channelID, userID uint64) error {
	sql := `DELETE FROM messaging_channel_member WHERE rel_channel = ? AND rel_user = ?`
	return exec(r.db().Exec(sql, channelID, userID))
}

func (r *channelMember) CountMemberships(userID uint64) (c int, err error) {
	return c, r.db().Get(&c,
		"SELECT COUNT(*) FROM messaging_channel_member WHERE rel_user = ?",
		userID)
}

func (r *channelMember) ChangeMembership(userID, target uint64) (err error) {
	// Remove dups
	// with an ugly mysql workaround
	_, err = r.db().Exec(
		"DELETE FROM messaging_channel_member WHERE rel_user = ? "+
			"AND rel_channel IN (SELECT rel_channel FROM (SELECT * FROM messaging_channel_member) AS workaround WHERE rel_user = ?)",
		userID,
		target)

	if err != nil {
		return err
	}

	_, err = r.db().Exec("UPDATE messaging_channel_member SET rel_user = ? WHERE rel_user = ?", target, userID)
	return err
}
