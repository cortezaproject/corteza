package store

import (
	"github.com/crusttech/human/server/pkg/envoy"
	"github.com/crusttech/human/server/pkg/envoy/resource"
	"github.com/crusttech/human/server/system/types"
)

func newRole(rl *types.Role) *role {
	return &role{
		rl: rl,
	}
}

func (rl *role) MarshalEnvoy() ([]resource.Interface, error) {
	return envoy.CollectNodes(
		resource.NewRole(rl.rl),
	)
}
