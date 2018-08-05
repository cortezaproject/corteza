package repository

import (
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
	"time"
)

type (
	Channel interface {
		FindChannelByID(id uint64) (*types.Channel, error)
		FindDirectChannelByUserID(fromUserID, toUserID uint64) (*types.Channel, error)
		FindChannels(filter *types.ChannelFilter) ([]*types.Channel, error)
		CreateChannel(mod *types.Channel) (*types.Channel, error)
		UpdateChannel(mod *types.Channel) (*types.Channel, error)

		FindChannelsMembershipsByMemberId(memberId uint64) ([]*types.ChannelMember, error)
		AddChannelMember(mod *types.ChannelMember) (*types.ChannelMember, error)
		RemoveChannelMember(channelID, userID uint64) error
		ArchiveChannelByID(id uint64) error
		UnarchiveChannelByID(id uint64) error
		DeleteChannelByID(id uint64) error
	}
)

const (
	sqlChannelValidOnly = ` true
         AND c.archived_at IS NULL         
         AND c.deleted_at IS NULL`

	sqlChannelSelect = `SELECT *
        FROM channels AS c
       WHERE ` + sqlChannelValidOnly

	sqlChannelDirect = `SELECT *
        FROM channels AS c
       WHERE c.type = 'group' 
         AND c.id IN (SELECT rel_channel 
                        FROM channel_members 
                       GROUP BY rel_channel
                      HAVING COUNT(*) = 2
                         AND MIN(rel_user) = ?
                         AND MAX(rel_user) = ?)`

	sqlChannelMemberships = `SELECT *
        FROM channel_members AS cm
       WHERE true`

	ErrChannelNotFound = repositoryError("ChannelNotFound")
)

func (r *repository) FindChannelByID(id uint64) (*types.Channel, error) {
	mod := &types.Channel{}
	sql := sqlChannelSelect + " AND id = ?"

	return mod, isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrChannelNotFound)
}

func (r *repository) FindDirectChannelByUserID(fromUserID, toUserID uint64) (*types.Channel, error) {
	mod := &types.Channel{}

	if fromUserID == toUserID {
		// do not waste any cpu cycles for his...
		return nil, ErrChannelNotFound
	}

	// We're grouping and aggregating values by min/max value of the user ID
	// so we need to swap valuess
	if fromUserID > toUserID {
		// Order by user idso we can simplifiy the search
		toUserID, fromUserID = fromUserID, toUserID
	}

	return mod, isFound(r.db().Get(mod, sqlChannelDirect, fromUserID, toUserID), mod.ID > 0, ErrChannelNotFound)
}

func (r *repository) FindChannels(filter *types.ChannelFilter) ([]*types.Channel, error) {
	params := make([]interface{}, 0)
	rval := make([]*types.Channel, 0)

	sql := sqlChannelSelect

	if filter != nil {
		if filter.Query != "" {
			sql += " AND c.name LIKE ?"
			params = append(params, filter.Query+"%")
		}
	}

	sql += " ORDER BY c.name ASC"

	return rval, r.db().Select(&rval, sql, params...)
}

func (r *repository) CreateChannel(mod *types.Channel) (*types.Channel, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()
	mod.Meta = coalesceJson(mod.Meta, []byte("{}"))

	if mod.Type == "" {
		mod.Type = types.ChannelTypePublic
	}

	return mod, r.db().Insert("channels", mod)
}

func (r *repository) UpdateChannel(mod *types.Channel) (*types.Channel, error) {
	mod.UpdatedAt = timeNowPtr()
	mod.Meta = coalesceJson(mod.Meta, []byte("{}"))
	if mod.Type == "" {
		mod.Type = types.ChannelTypePublic
	}

	whitelist := []string{"id", "name", "type", "topic", "meta", "updated_at"}

	return mod, r.db().
		UpdatePartial("channels", mod, whitelist, "id")
}

func (r *repository) FindChannelsMembershipsByMemberId(memberId uint64) ([]*types.ChannelMember, error) {
	var rval = make([]*types.ChannelMember, 0)

	return rval, r.db().Select(&rval, sqlChannelMemberships+" AND cm.rel_user = ? ", memberId)
}

func (r *repository) AddChannelMember(mod *types.ChannelMember) (*types.ChannelMember, error) {
	sql := `INSERT INTO channel_members (rel_channel, rel_user) VALUES (?, ?)`
	mod.CreatedAt = time.Now()

	return mod, exec(r.db().Exec(sql, mod.ChannelID, mod.UserID))
}

func (r *repository) RemoveChannelMember(channelID, userID uint64) error {
	sql := `DELETE FROM channel_members WHERE rel_channel = ? AND rel_user = ?`
	return exec(r.db().Exec(sql, channelID, userID))
}

func (r *repository) ArchiveChannelByID(id uint64) error {
	return r.updateColumnByID("channels", "archived_at", time.Now(), id)
}

func (r *repository) UnarchiveChannelByID(id uint64) error {
	return r.updateColumnByID("channels", "archived_at", nil, id)
}

func (r *repository) DeleteChannelByID(id uint64) error {
	return r.updateColumnByID("channels", "deleted_at", time.Now(), id)
}

func (r *repository) RecoverChannelByID(id uint64) error {
	return r.updateColumnByID("channels", "deleted_at", nil, id)
}
