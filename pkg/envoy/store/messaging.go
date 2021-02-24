package store

import (
	"context"

	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	messagingChannelsFilter types.ChannelFilter

	messagingStore interface {
		store.MessagingChannels
	}

	messagingDecoder struct{}
)

func newMessagingDecoder() *messagingDecoder {
	return &messagingDecoder{}
}

func (d *messagingDecoder) decodeMessagingChannel(ctx context.Context, s messagingStore, ff []*messagingChannelsFilter) *auxRsp {
	mm := make([]envoy.Marshaller, 0, 100)
	if ff == nil {
		return &auxRsp{
			mm: mm,
		}
	}

	var nn types.ChannelSet
	var fn types.ChannelFilter
	var err error
	for _, f := range ff {
		aux := *f

		if aux.Limit == 0 {
			aux.Limit = 1000
		}

		for {
			nn, fn, err = s.SearchMessagingChannels(ctx, types.ChannelFilter(aux))
			if err != nil {
				return &auxRsp{
					err: err,
				}
			}

			for _, n := range nn {
				mm = append(mm, newMessagingChannel(n))
			}

			if fn.NextPage != nil {
				aux.PageCursor = fn.NextPage
			} else {
				break
			}
		}
	}

	return &auxRsp{
		mm: mm,
	}
}

// MessagingChannels adds a new messaging ChannelFilter
func (df *DecodeFilter) MessagingChannels(f *types.ChannelFilter) *DecodeFilter {
	if df.messagingChannels == nil {
		df.messagingChannels = make([]*messagingChannelsFilter, 0, 1)
	}
	df.messagingChannels = append(df.messagingChannels, (*messagingChannelsFilter)(f))
	return df
}
