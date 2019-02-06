package repository

import (
	"context"
	"sort"
	"strconv"
	"time"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/messaging/types"
)

type (
	ChannelRepository interface {
		With(ctx context.Context, db *factory.DB) ChannelRepository

		FindChannelByID(id uint64) (*types.Channel, error)
		FindChannelByMemberSet(memberID ...uint64) (*types.Channel, error)
		FindChannels(filter *types.ChannelFilter) ([]*types.Channel, error)
		CreateChannel(mod *types.Channel) (*types.Channel, error)
		UpdateChannel(mod *types.Channel) (*types.Channel, error)

		ArchiveChannelByID(id uint64) error
		UnarchiveChannelByID(id uint64) error
		DeleteChannelByID(id uint64) error
		UndeleteChannelByID(id uint64) error

		CountCreated(userID uint64) (c int, err error)
		ChangeCreator(userID, target uint64) error
	}

	channel struct {
		*repository
	}
)

const (
	sqlChannelColumns = " id," +
		"name, " +
		"meta, " +
		"created_at, " +
		"updated_at, " +
		"archived_at, " +
		"deleted_at, " +
		"rel_organisation, " +
		"rel_creator, " +
		"type  , " +
		"rel_last_message, " +
		"topic"

	sqlChannelSelect = `SELECT ` + sqlChannelColumns + `
        FROM channels AS c
       WHERE true `

	sqlChannelGroupByMemberSet = sqlChannelSelect + ` AND c.type = ? AND c.id IN (
            SELECT rel_channel 
              FROM messaging_channel_member 
             GROUP BY rel_channel 
            HAVING COUNT(*) = ? 
               AND CONCAT(GROUP_CONCAT(rel_user ORDER BY 1 ASC SEPARATOR ','),',') = ?
        )`

	// subquery that filters out all channels that current user has access to as a member
	// or via channel type (public channels)
	sqlChannelAccess = ` (
				SELECT id
                  FROM channels c
                       LEFT OUTER JOIN messaging_channel_member AS m ON (c.id = m.rel_channel)
                 WHERE rel_user = ?
              UNION
                SELECT id
                  FROM channels c
                 WHERE c.type = ?
			)`

	ErrChannelNotFound = repositoryError("ChannelNotFound")
)

func Channel(ctx context.Context, db *factory.DB) ChannelRepository {
	return (&channel{}).With(ctx, db)
}

func (r *channel) With(ctx context.Context, db *factory.DB) ChannelRepository {
	return &channel{
		repository: r.repository.With(ctx, db),
	}
}

func (r *channel) FindChannelByID(id uint64) (*types.Channel, error) {
	mod := &types.Channel{}
	sql := sqlChannelSelect + " AND id = ?"

	return mod, isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrChannelNotFound)
}

// FindChannelByMemberSet searches for channel (group!) with exactly the same membership structure
func (r *channel) FindChannelByMemberSet(memberIDs ...uint64) (*types.Channel, error) {
	mod := &types.Channel{}

	sort.Slice(memberIDs, func(i, j int) bool {
		return memberIDs[i] < memberIDs[j]
	})

	membersConcat := ""
	for i := range memberIDs {
		// Don't panic, we're adding , in the SQL as well
		membersConcat += strconv.FormatUint(memberIDs[i], 10) + ","
	}

	return mod, isFound(r.db().Get(mod, sqlChannelGroupByMemberSet, types.ChannelTypeGroup, len(memberIDs), membersConcat), mod.ID > 0, ErrChannelNotFound)
}

func (r *channel) FindChannels(filter *types.ChannelFilter) ([]*types.Channel, error) {
	// @todo: actual searching (filter.Query) not just a full select

	params := make([]interface{}, 0)
	rval := make([]*types.Channel, 0)

	sql := sqlChannelSelect

	if filter != nil {
		if filter.Query != "" {
			sql += " AND c.name LIKE ?"
			params = append(params, filter.Query+"%")
		}

		if filter.CurrentUserID > 0 {
			sql += " AND c.id IN " + sqlChannelAccess
			params = append(params, filter.CurrentUserID, types.ChannelTypePublic)
		}
	}

	sql += " ORDER BY c.name ASC"

	return rval, r.db().Select(&rval, sql, params...)
}

func (r *channel) CreateChannel(mod *types.Channel) (*types.Channel, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()
	mod.UpdatedAt = nil
	mod.Meta = coalesceJson(mod.Meta, []byte("{}"))

	if mod.Type == "" {
		mod.Type = types.ChannelTypePublic
	}

	return mod, r.db().Insert("channels", mod)
}

func (r *channel) UpdateChannel(mod *types.Channel) (*types.Channel, error) {
	mod.UpdatedAt = timeNowPtr()
	mod.Meta = coalesceJson(mod.Meta, []byte("{}"))
	if mod.Type == "" {
		mod.Type = types.ChannelTypePublic
	}

	whitelist := []string{"id", "name", "type", "topic", "meta", "updated_at"}

	return mod, r.db().
		UpdatePartial("channels", mod, whitelist, "id")
}

func (r *channel) ArchiveChannelByID(id uint64) error {
	return r.updateColumnByID("channels", "archived_at", time.Now(), id)
}

func (r *channel) UnarchiveChannelByID(id uint64) error {
	return r.updateColumnByID("channels", "archived_at", nil, id)
}

func (r *channel) DeleteChannelByID(id uint64) error {
	return r.updateColumnByID("channels", "deleted_at", time.Now(), id)
}

func (r *channel) UndeleteChannelByID(id uint64) error {
	return r.updateColumnByID("channels", "deleted_at", nil, id)
}

func (r *channel) CountCreated(userID uint64) (c int, err error) {
	return c, r.db().Get(&c,
		"SELECT COUNT(*) FROM channels WHERE rel_creator = ?",
		userID)
}

func (r *channel) ChangeCreator(userID, target uint64) error {
	_, err := r.db().Exec("UPDATE channels SET rel_creator = ? WHERE rel_creator = ?", target, userID)
	return err
}
