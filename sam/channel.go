package sam

import (
	"github.com/pkg/errors"

	"github.com/titpetric/factory"
)

func (*Channel) Edit(r *channelEditRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: topic change message/log entry
	// @todo: channel name change message/log entry
	// @todo: permission check if user can edit channel
	// @todo: permission check if user can add channel

	c := Channel{}.new().SetID(r.id).SetName(r.name).SetTopic(r.topic)
	if c.GetID() > 0 {
		return c, db.Replace("channel", c)
	}
	c.SetID(factory.Sonyflake.NextID())
	return c, db.Insert("channel", c)
}

func (*Channel) Remove(r *channelRemoveRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: don't actually delete the channel (mark as deleted, history becomes unavailable)
	// @todo: notify users that channel has been removed (remove from web UI)
	// @todo: permissions check if user cah remove channel

	return nil, func() error {
		_, err := db.Exec("delete from channel where id=?", r.id)
		return err
	}()
}

func (*Channel) Read(r *channelReadRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permission check if user can read channel

	c := Channel{}.new()
	return c, db.Get(c, "select * from channel where id=?", r.id)
}

func (*Channel) Search(r *channelSearchRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permission check to return only channels that user has access to
	// @todo: actual searching not just a full select

	res := make([]Channel, 0)
	err = db.Select(&res, "select * from channel order by name asc")
	return res, err
}

func (*Channel) Archive(r *channelArchiveRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: don't actually delete the channel (mark as archived, history stays available)
	// @todo: notify users that channel has been archived (last message - archival, disable new messages)
	// @todo: permissions check if user cah archive channel

	return nil, func() error {
		_, err = db.Exec("delete from channel where id=?", r.id)
		return err
	}()
}

func (*Channel) Move(r *channelMoveRequest) (interface{}, error) {
	// @todo: move channel from r.source to r.destination (organisation)
	return nil, errors.New("Not implemented: Channel.move")
}
