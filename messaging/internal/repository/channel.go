package repository

import (
	"context"
	"sort"
	"strconv"
	"time"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	ChannelRepository interface {
		With(ctx context.Context, db *factory.DB) ChannelRepository

		FindByID(id uint64) (*types.Channel, error)
		FindByMemberSet(memberID ...uint64) (*types.Channel, error)
		Find(filter *types.ChannelFilter) (types.ChannelSet, error)

		Create(mod *types.Channel) (*types.Channel, error)
		Update(mod *types.Channel) (*types.Channel, error)

		ArchiveByID(id uint64) error
		UnarchiveByID(id uint64) error
		DeleteByID(id uint64) error
		UndeleteByID(id uint64) error

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
        FROM messaging_channel AS c
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
                  FROM messaging_channel c
                       LEFT OUTER JOIN messaging_channel_member AS m ON (c.id = m.rel_channel)
                 WHERE rel_user = ?
              UNION
                SELECT id
                  FROM messaging_channel c
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

func (r *channel) FindByID(id uint64) (*types.Channel, error) {
	mod := &types.Channel{}
	sql := sqlChannelSelect + " AND id = ?"

	return mod, isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrChannelNotFound)
}

// FindChannelByMemberSet searches for channel (group!) with exactly the same membership structure
func (r *channel) FindByMemberSet(memberIDs ...uint64) (*types.Channel, error) {
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

func (r *channel) Find(filter *types.ChannelFilter) (types.ChannelSet, error) {
	// @todo: actual searching (filter.Query) not just a full select

	params := make([]interface{}, 0)
	rval := types.ChannelSet{}

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

		if !filter.IncludeDeleted {
			sql += " AND deleted_at IS NULL"
		}
	}

	sql += " ORDER BY c.name ASC"

	return rval, r.db().Select(&rval, sql, params...)
}

func (r *channel) Create(mod *types.Channel) (*types.Channel, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()
	mod.UpdatedAt = nil

	if mod.Type == "" {
		mod.Type = types.ChannelTypePublic
	}

	return mod, r.db().Insert("messaging_channel", mod)
}

func (r *channel) Update(mod *types.Channel) (*types.Channel, error) {
	mod.UpdatedAt = timeNowPtr()
	if mod.Type == "" {
		mod.Type = types.ChannelTypePublic
	}

	whitelist := []string{"id", "name", "type", "topic", "meta", "updated_at"}

	return mod, r.db().UpdatePartial("messaging_channel", mod, whitelist, "id")
}

func (r *channel) ArchiveByID(id uint64) error {
	return r.updateColumnByID("messaging_channel", "archived_at", time.Now(), id)
}

func (r *channel) UnarchiveByID(id uint64) error {
	return r.updateColumnByID("messaging_channel", "archived_at", nil, id)
}

func (r *channel) DeleteByID(id uint64) error {
	return r.updateColumnByID("messaging_channel", "deleted_at", time.Now(), id)
}

func (r *channel) UndeleteByID(id uint64) error {
	return r.updateColumnByID("messaging_channel", "deleted_at", nil, id)
}

func (r *channel) CountCreated(userID uint64) (c int, err error) {
	return c, r.db().Get(&c, "SELECT COUNT(*) FROM messaging_channel WHERE rel_creator = ?", userID)
}

func (r *channel) ChangeCreator(userID, target uint64) error {
	_, err := r.db().Exec("UPDATE messaging_channel SET rel_creator = ? WHERE rel_creator = ?", target, userID)
	return err
}
