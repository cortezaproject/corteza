package yaml

import (
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	messagingChannel struct {
		// when messagingChannel is at least partially defined
		res *types.Channel
		ts  *resource.Timestamps
		us  *resource.Userstamps

		envoyConfig   *resource.EnvoyConfig
		encoderConfig *EncoderConfig

		// module's RBAC rules
		rbac rbacRuleSet
	}

	messagingChannelSet []*messagingChannel
)

func (nn messagingChannelSet) ConfigureEncoder(cfg *EncoderConfig) {
	for _, n := range nn {
		n.encoderConfig = cfg
	}
}
