package store

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/system/types"
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
