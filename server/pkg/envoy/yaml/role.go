package yaml

import (
	"github.com/crusttech/human/server/pkg/envoy/resource"
	"github.com/crusttech/human/server/system/types"
)

type (
	role struct {
		res *types.Role
		ts  *resource.Timestamps

		envoyConfig   *resource.EnvoyConfig
		encoderConfig *EncoderConfig

		// module's RBAC rules
		rbac rbacRuleSet
	}
	roleSet []*role
)

func (nn roleSet) configureEncoder(cfg *EncoderConfig) {
	for _, n := range nn {
		n.encoderConfig = cfg
	}
}
