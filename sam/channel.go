package sam

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
)

var _ = errors.Wrap

const (
	sqlChannelScope  = "deleted_at IS NULL AND archived_at IS NULL"
	sqlChannelSelect = "SELECT * FROM channels WHERE " + sqlChannelScope
)

func (*Channel) Create(r *channelCreateRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: topic message/log entry
	// @todo: channel name cmessage/log entry
	// @todo: permission check if user can add channel

	c := Channel{}.
		New().
		SetName(r.name).
		SetTopic(r.topic).
		SetMeta([]byte("{}")).
		SetID(factory.Sonyflake.NextID())

	return c, db.Insert("channels", c)
}

func (*Channel) Edit(r *channelEditRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	var c *Channel
	if c, err = c.load(r.id); err != nil {
		return nil, err
	}

	// @todo: topic change message/log entry
	// @todo: channel name change message/log entry
	// @todo: permission check if user can edit channel
	// @todo: make sure archived & deleted entries can not be edited
	// @todo: handle channel moving
	// @todo: handle channel archiving

	c.SetName(r.name).SetTopic(r.topic)

	return c, db.Replace("channels", c)

}

func (*Channel) Delete(r *channelDeleteRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	var c *Channel
	if c, err = c.load(r.id); err != nil {
		return nil, err
	}

	// @todo: make history unavailable
	// @todo: notify users that channel has been removed (remove from web UI)
	// @todo: permissions check if user cah remove channel

	stmt := "UPDATE channels SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL"
	spew.Dump(r.id)
	return nil, func() error {
		_, err := db.Exec(stmt, r.id)
		return err
	}()
}

func (*Channel) Read(r *channelReadRequest) (interface{}, error) {
	return (&Channel{}).load(r.id)
}

func (*Channel) List(r *channelListRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permission check to return only channels that user has access to
	// @todo: actual searching not just a full select

	res := make([]Channel, 0)
	err = db.Select(&res, sqlChannelSelect+" ORDER BY name ASC")
	return res, err
}

func (*Channel) load(id uint64) (*Channel, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	c := Channel{}.New()

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
