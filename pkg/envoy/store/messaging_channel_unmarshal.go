package store

import (
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

func newMessagingChannel(ch *types.Channel) *messagingChannel {
	return &messagingChannel{
		ch: ch,
	}
}

// MarshalEnvoy converts the messaging channel struct to a resource
func (ch *messagingChannel) MarshalEnvoy() ([]resource.Interface, error) {
	return envoy.CollectNodes(
		resource.NewMessagingChannel(ch.ch),
	)
}
