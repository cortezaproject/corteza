package store

import (
	"github.com/crusttech/human/server/pkg/envoy"
	"github.com/crusttech/human/server/pkg/envoy/resource"
	"github.com/crusttech/human/server/system/types"
)

func newUser(u *types.User) *user {
	return &user{
		u: u,
	}
}

func (u *user) MarshalEnvoy() ([]resource.Interface, error) {
	return envoy.CollectNodes(
		resource.NewUser(u.u),
	)
}
