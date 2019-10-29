package repository

import (
	"context"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	ChannelRepository interface {
		With(ctx context.Context, db *factory.DB) ChannelRepository

		FindByID(id uint64) (*types.Channel, error)
		FindByMemberSet(memberID ...uint64) (*types.Channel, error)
		Find(types.ChannelFilter) (types.ChannelSet, types.ChannelFilter, error)

		Create(mod *types.Channel) (*types.Channel, error)
		Update(mod *types.Channel) (*types.Channel, error)

		ArchiveByID(id uint64) error
		UnarchiveByID(id uint64) error
		DeleteByID(id uint64) error
		UndeleteByID(id uint64) error
	}

	channel struct {
		*repository
	}
)

const (
	ErrChannelNotFound = repositoryError("ChannelNotFound")
)

func Channel(ctx context.Context, db *factory.DB) ChannelRepository {
	return (&channel{}).With(ctx, db)
}

func (r channel) With(ctx context.Context, db *factory.DB) ChannelRepository {
	return &channel{
		repository: r.repository.With(ctx, db),
	}
}

func (r channel) table() string {
	return "messaging_channel"
}

func (r channel) columns() []string {
	return []string{
		"c.id",
		"c.name",
		"c.meta",
		"c.membership_policy",
		"c.created_at",
		"c.updated_at",
		"c.archived_at",
		"c.deleted_at",
		"c.rel_organisation",
		"c.rel_creator",
		"c.type",
		"c.rel_last_message",
		"c.topic",
	}
}

func (r channel) query() squirrel.SelectBuilder {
	return squirrel.
		Select(r.columns()...).
		From(r.table() + " AS c")
}

func (r channel) FindByID(ID uint64) (*types.Channel, error) {
	return r.findOneBy(squirrel.Eq{"c.id": ID})
}

// FindByMemberSet searches for channel (group!) with exactly the same membership structure
func (r channel) FindByMemberSet(memberIDs ...uint64) (*types.Channel, error) {
	return r.findOneBy(
		squirrel.And{
			squirrel.Eq{"type": types.ChannelTypeGroup},
			squirrel.ConcatExpr("c.id IN (", (channelMember{}).queryExactMembers(memberIDs...), ")"),
		})
}

func (r channel) findOneBy(cnd squirrel.Sqlizer) (*types.Channel, error) {
	var (
		ch = &types.Channel{}

		q = r.query().
			Where(cnd)

		err = rh.FetchOne(r.db(), q, ch)
	)

	if err != nil {
		return nil, err
	} else if ch.ID == 0 {
		return nil, ErrChannelNotFound
	}

	return ch, nil
}

func (r channel) Find(filter types.ChannelFilter) (set types.ChannelSet, f types.ChannelFilter, err error) {
	f = filter

	if f.Sort == "" {
		f.Sort = "c.name ASC"
	}

	query := r.query()

	query = query.Where(squirrel.Eq{"c.archived_at": nil})

	if !f.IncludeDeleted {
		query = query.Where(squirrel.Eq{"c.deleted_at": nil})
	}

	if len(f.ChannelID) > 0 {
		query = query.Where(squirrel.Eq{"c.id": f.ChannelID})
	}

	if f.Query != "" {
		q := "%" + strings.ToLower(f.Query) + "%"
		query = query.Where(squirrel.Like{"LOWER(name)": q})
	}

	if f.CurrentUserID > 0 {
		query = query.Where(squirrel.Or{
			squirrel.Eq{"c.type": types.ChannelTypePublic},
			squirrel.ConcatExpr("c.id IN (", (channelMember{}).queryAnyMember(f.CurrentUserID), ")"),
		})
	}

	var orderBy []string

	if orderBy, err = rh.ParseOrder(f.Sort, r.columns()...); err != nil {
		return
	} else {
		query = query.OrderBy(orderBy...)
	}

	return set, f, rh.FetchAll(r.db(), query, &set)
}

func (r channel) Create(mod *types.Channel) (*types.Channel, error) {
	mod.ID = factory.Sonyflake.NextID()

	rh.SetCurrentTimeRounded(&mod.CreatedAt)
	mod.UpdatedAt = nil

	if mod.Type == "" {
		mod.Type = types.ChannelTypePublic
	}

	return mod, r.db().Insert("messaging_channel", mod)
}

func (r channel) Update(mod *types.Channel) (*types.Channel, error) {
	rh.SetCurrentTimeRounded(&mod.UpdatedAt)

	if mod.Type == "" {
		mod.Type = types.ChannelTypePublic
	}

	whitelist := []string{"id", "name", "type", "membership_policy", "topic", "meta", "updated_at"}

	return mod, r.db().UpdatePartial("messaging_channel", mod, whitelist, "id")
}

func (r channel) ArchiveByID(ID uint64) error {
	return rh.UpdateColumns(r.db(), r.table(), rh.Set{"archived_at": time.Now()}, squirrel.Eq{"id": ID})
}

func (r channel) UnarchiveByID(ID uint64) error {
	return rh.UpdateColumns(r.db(), r.table(), rh.Set{"archived_at": nil}, squirrel.Eq{"id": ID})
}

func (r channel) DeleteByID(ID uint64) error {
	return rh.UpdateColumns(r.db(), r.table(), rh.Set{"deleted_at": time.Now()}, squirrel.Eq{"id": ID})
}

func (r channel) UndeleteByID(ID uint64) error {
	return rh.UpdateColumns(r.db(), r.table(), rh.Set{"deleted_at": nil}, squirrel.Eq{"id": ID})
}
