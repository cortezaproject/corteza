package resource

import (
	"strconv"

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

	// Initial timestamps
	r.SetTimestamps(MakeCUDATimestamps(&res.CreatedAt, res.UpdatedAt, res.DeletedAt, res.ArchivedAt))
	// Initial userstamps
	r.SetUserstamps(&Userstamps{
		CreatedBy: &Userstamp{UserID: res.CreatorID},
	})

	return r
}

func (r *MessagingChannel) SysID() uint64 {
	return r.Res.ID
}

func (r *MessagingChannel) Ref() string {
	return FirstOkString(r.Res.Name, strconv.FormatUint(r.Res.ID, 10))
}

// FindMessagingChannel looks for the ch in the resource set
func FindMessagingChannel(rr InterfaceSet, ii Identifiers) (ap *types.Channel) {
	var chRes *MessagingChannel

	rr.Walk(func(r Interface) error {
		ar, ok := r.(*MessagingChannel)
		if !ok {
			return nil
		}

		if ar.Identifiers().HasAny(ii) {
			chRes = ar
		}

		return nil
	})

	// Found it
	if chRes != nil {
		return chRes.Res
	}

	return nil
}
