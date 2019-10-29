package repository

import (
	"context"
	"sort"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	// ChannelMemberRepository interface to channel member repository
	ChannelMemberRepository interface {
		With(ctx context.Context, db *factory.DB) ChannelMemberRepository

		Find(filter types.ChannelMemberFilter) (types.ChannelMemberSet, error)

		Create(mod *types.ChannelMember) (*types.ChannelMember, error)
		Update(mod *types.ChannelMember) (*types.ChannelMember, error)
		Delete(channelID, userID uint64) error
	}

	channelMember struct {
		*repository
	}
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

func (r channelMember) table() string {
	return "messaging_channel_member"
}

func (r channelMember) columns() []string {
	return []string{
		"cm.rel_channel",
		"cm.rel_user",
		"cm.type",
		"cm.flag",
		"cm.created_at",
		"cm.updated_at",
	}
}

func (r channelMember) query() squirrel.SelectBuilder {
	return squirrel.
		Select(r.columns()...).
		From(r.table() + " AS cm")
}

// Finds channel ID(s) with any of the members
//
// Builds a (sub)query that returns list of channel IDs at least one of the members
//
func (r channelMember) queryAnyMember(memberIDs ...uint64) squirrel.SelectBuilder {
	return squirrel.
		Select("cm.rel_channel").
		From(r.table() + " AS cm").
		Where(squirrel.Eq{"cm.rel_user": memberIDs})
}

// Finds channel ID(s) with exact membership
//
// Builds a (sub)query that returns list of channel IDs that have this exact membership
//
func (r channelMember) queryExactMembers(memberIDs ...uint64) squirrel.SelectBuilder {
	if len(memberIDs) == 0 {
		return squirrel.
			Select("null")
	}

	// Make sure members are sorted
	sort.Slice(memberIDs, func(i, j int) bool {
		return memberIDs[i] < memberIDs[j]
	})

	// Concatentating members fore
	membersConcat := ""
	for i := range memberIDs {
		// Don't panic, we're adding , in the SQL as well
		membersConcat += strconv.FormatUint(memberIDs[i], 10) + ","
	}

	return r.queryAnyMember(memberIDs...).
		GroupBy("cm.rel_channel").
		Having(squirrel.Eq{
			"COUNT(*)": len(memberIDs),
			"CONCAT(GROUP_CONCAT(cm.rel_user ORDER BY 1 ASC SEPARATOR ','),',')": membersConcat,
		})
}

// Find fetches membership info
func (r *channelMember) Find(filter types.ChannelMemberFilter) (set types.ChannelMemberSet, err error) {
	query := r.query()

	if len(filter.MemberID) > 0 {
		query = query.Where(squirrel.Eq{"cm.rel_user": filter.MemberID})
	}

	if len(filter.ChannelID) > 0 {
		query = query.Where(squirrel.Eq{"cm.rel_channel": filter.ChannelID})
	}

	return set, rh.FetchAll(r.db(), query, &set)
}

// Create adds channel membership record
func (r *channelMember) Create(mod *types.ChannelMember) (*types.ChannelMember, error) {
	rh.SetCurrentTimeRounded(&mod.CreatedAt)
	mod.UpdatedAt = nil

	return mod, r.db().Insert("messaging_channel_member", mod)
}

// Update modifies existing channel membership record
func (r *channelMember) Update(mod *types.ChannelMember) (*types.ChannelMember, error) {
	rh.SetCurrentTimeRounded(&mod.UpdatedAt)

	whitelist := []string{"type", "flag", "updated_at", "rel_channel", "rel_user"}

	return mod, r.db().UpdatePartial("messaging_channel_member", mod, whitelist, "rel_channel", "rel_user")
}

// Delete removes existing channel membership record
func (r *channelMember) Delete(channelID, userID uint64) error {
	return rh.Delete(r.db(), r.table(), squirrel.Eq{
		"rel_channel": channelID,
		"rel_user":    userID,
	})
}
