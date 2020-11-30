package resource

import (
	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	MessagingChannel struct {
		*base
		Res *types.Channel
	}
)

func NewMessagingChannel(res *types.Channel) *MessagingChannel {
	r := &MessagingChannel{
		base: &base{},
	}
	r.SetResourceType(MESSAGING_CHANNEL_RESOURCE_TYPE)
	r.Res = res

	r.AddIdentifier(identifiers(res.ID, res.Name)...)

	return r
}

func (r *MessagingChannel) SysID() uint64 {
	return r.Res.ID
}
