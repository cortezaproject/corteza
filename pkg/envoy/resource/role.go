package resource

import (
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	// Role represents a Role
	Role struct {
		*base
		Res *types.Role
	}
)

func NewRole(rl *types.Role) *Role {
	r := &Role{base: &base{}}
	r.SetResourceType(ROLE_RESOURCE_TYPE)
	r.Res = rl

	r.AddIdentifier(identifiers(rl.Handle, rl.Name, rl.ID)...)

	return r
}

func (r *Role) SysID() uint64 {
	return r.Res.ID
}
