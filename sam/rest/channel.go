package rest

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/sam/rest/server"
	"github.com/crusttech/crust/sam/types"
)

var _ = errors.Wrap

const (
	sqlChannelScope  = "deleted_at IS NULL AND archived_at IS NULL"
	sqlChannelSelect = "SELECT * FROM channels WHERE " + sqlChannelScope
)

type Channel struct{}

func (Channel) New() *Channel {
	return &Channel{}
}

func (*Channel) Create(ctx context.Context, r *server.ChannelCreateRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: topic message/log entry
	// @todo: channel name cmessage/log entry
	// @todo: permission check if user can add channel

	c := types.Channel{}.
		New().
		SetName(r.Name).
		SetTopic(r.Topic).
		SetMeta([]byte("{}")).
		SetID(factory.Sonyflake.NextID())

	return c, db.Insert("channels", c)
}

func (*Channel) Edit(ctx context.Context, r *server.ChannelEditRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	/*var c *types.Channel
	if c, err = c.load(r.ID); err != nil {
		return nil, err
	}*/

	// @todo: topic change message/log entry
	// @todo: channel name change message/log entry
	// @todo: permission check if user can edit channel
	// @todo: make sure archived & deleted entries can not be edited
	// @todo: handle channel moving
	// @todo: handle channel archiving

	c := types.Channel{}.
		New().
		SetName(r.Name).
		SetTopic(r.Topic)

	return c, db.Replace("channels", c)

}

func (*Channel) Delete(ctx context.Context, r *server.ChannelDeleteRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	/*var c *types.Channel
	if c, err = c.load(r.ID); err != nil {
		return nil, err
	}*/

	// @todo: make history unavailable
	// @todo: notify users that channel has been removed (remove from web UI)
	// @todo: permissions check if user cah remove channel

	stmt := "UPDATE channels SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL"
	spew.Dump(r.ID)
	return nil, func() error {
		_, err := db.Exec(stmt, r.ID)
		return err
	}()
}

func (s *Channel) Read(ctx context.Context, r *server.ChannelReadRequest) (interface{}, error) {
	return s.load(r.ID)
}

func (*Channel) List(ctx context.Context, r *server.ChannelListRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permission check to return only channels that user has access to
	// @todo: actual searching not just a full select

	res := make([]types.Channel, 0)
	err = db.Select(&res, sqlChannelSelect+" ORDER BY name ASC")
	return res, err
}

func (*Channel) load(id uint64) (*types.Channel, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	c := types.Channel{}.New()

	if id == 0 {
		return nil, errors.New("Provide channel ID")
	} else if err := db.Get(c, sqlChannelSelect+" AND id = ?", id); err != nil {
		return nil, err
	} else if c.ID != id {
		spew.Dump(c)
		return nil, errors.New("Unexisting channel")
	}

	// @todo: permission check if user can read channel

	return c, nil
}
