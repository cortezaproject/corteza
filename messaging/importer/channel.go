package importer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
	"github.com/cortezaproject/corteza-server/pkg/importer"
)

type (
	Channel struct {
		set         types.ChannelSet
		dirty       map[uint64]bool
		permissions importer.PermissionImporter
	}

	channelKeeper interface {
		Update(*types.Channel) (*types.Channel, error)
		Create(*types.Channel) (*types.Channel, error)
	}
)

func NewChannelImport(permissions importer.PermissionImporter, set types.ChannelSet) *Channel {
	if set == nil {
		set = types.ChannelSet{}
	}

	out := &Channel{
		set:         set,
		dirty:       make(map[uint64]bool),
		permissions: permissions,
	}

	return out
}

func (cImp *Channel) CastSet(set interface{}) error {
	var name string
	return deinterfacer.Each(set, func(index int, _ string, def interface{}) error {
		if index > -1 {
			// Channels defined as collection
			deinterfacer.KVsetString(&name, "name", def)
		}

		return cImp.Cast(name, def)
	})
}

func (cImp *Channel) Cast(name string, def interface{}) (err error) {
	var channel *types.Channel

	// if !importer.IsValidHandle(handle) {
	// 	return errors.New("invalid channel handle")
	// }
	//
	// handle = importer.NormalizeHandle(handle)
	if channel, err = cImp.Get(name); err != nil {
		return err
	} else if channel == nil {
		channel = &types.Channel{
			Name: name,
		}

		cImp.set = append(cImp.set, channel)
	} else if channel.ID == 0 {
		return errors.Errorf("channel name %q already defined in this import session", channel.Name)
	} else {
		cImp.dirty[channel.ID] = true
	}

	if name, ok := def.(string); ok && name != "" {
		channel.Name = name
		return nil
	}

	return deinterfacer.Each(def, func(_ int, key string, val interface{}) (err error) {
		switch key {
		case "name":
			// already handled
		case "type":
			channel.Type = types.ChannelType(deinterfacer.ToString(val))
			if !channel.Type.IsValid() {
				return fmt.Errorf("invalid channel type %q for channel %q", channel.Type, channel.Name)

			}

		case "topic":
			channel.Topic = deinterfacer.ToString(val)

		case "allow", "deny":
			return cImp.permissions.CastSet(types.ChannelPermissionResource.String()+channel.Name, key, val)

		default:
			return fmt.Errorf("unexpected key %q for channel %q", key, channel.Name)
		}

		return err
	})
}

func (cImp *Channel) Get(name string) (*types.Channel, error) {
	// name = importer.NormalizeHandle(name)
	//
	// if !importer.IsValidHandle(name) {
	// 	return nil, errors.New("invalid channel name")
	// }

	return cImp.set.FindByName(name), nil
}

func (cImp *Channel) Store(ctx context.Context, k channelKeeper) error {
	return cImp.set.Walk(func(channel *types.Channel) (err error) {
		var handle = channel.Name

		if channel.ID == 0 {
			channel, err = k.Create(channel)
		} else if cImp.dirty[channel.ID] {
			channel, err = k.Update(channel)
		}

		if err != nil {
			return
		}

		cImp.permissions.UpdateResources(types.ChannelPermissionResource.String(), handle, channel.ID)
		cImp.permissions.UpdateRoles(channel.Name, channel.ID)

		return
	})
}
