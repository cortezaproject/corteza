package resource

import (
	"fmt"

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

func (m *Role) SearchQuery() types.RoleFilter {
	f := types.RoleFilter{
		Handle: m.Res.Handle,
		Name:   m.Res.Name,
	}

	if m.Res.ID > 0 {
		f.Query = fmt.Sprintf("roleID=%d", m.Res.ID)
	}

	return f
}
