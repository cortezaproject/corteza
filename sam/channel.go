package sam

import (
	"fmt"

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

	c := Channel{}.New().SetName(r.name).SetTopic(r.topic)
	if c.GetID() > 0 {
		if is("topic", c.changed...) {
			fmt.Println("Topic for channel was set:", c.GetTopic())
		}
		return c, db.Replace("channel", c)
	}
	c.SetID(factory.Sonyflake.NextID())
	return c, db.Insert("channel", c)
}

func (*Channel) Edit(r *channelEditRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: topic change message/log entry
	// @todo: channel name change message/log entry
	// @todo: permission check if user can edit channel
	// @todo: make sure archived & deleted entries can not be edited
	// @todo: handle channel moving
	// @todo: handle channel archiving

	c := Channel{}.New().SetID(r.id).SetName(r.name).SetTopic(r.topic)
	if c.GetID() > 0 {
		if is("topic", c.changed...) {
			fmt.Println("Topic for channel was changed:", c.GetTopic())
		}
		return c, db.Replace("channel", c)
	}

	return c, db.Insert("channel", c)
}

func (*Channel) Remove(r *channelRemoveRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: make history unavailable
	// @todo: notify users that channel has been removed (remove from web UI)
	// @todo: permissions check if user cah remove channel

	stmt := "UPDATE channels SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL"

	return nil, func() error {
		_, err := db.Exec(stmt, r.id)
		return err
	}()
}

func (*Channel) Read(r *channelReadRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permission check if user can read channel

	c := Channel{}.New()
	return c, db.Get(c, sqlChannelSelect+" AND id = ?", r.id)
}

func (*Channel) Search(r *channelSearchRequest) (interface{}, error) {
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
