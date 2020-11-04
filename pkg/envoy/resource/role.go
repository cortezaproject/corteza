package resource

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/system/types"
)

const (
	ROLE_RESOURCE_TYPE = "role"
)

type (
	// role represents a Role
	role struct {
		*base
		Res *types.Role
	}
)

func Role(rl *types.Role) *role {
	r := &role{base: &base{}}
	r.SetResourceType(ROLE_RESOURCE_TYPE)
	r.Res = rl

	r.AddIdentifier(identifiers(rl.Handle, rl.Name, rl.ID)...)

	return r
}

func (m *role) SearchQuery() types.RoleFilter {
	f := types.RoleFilter{
		Handle: m.Res.Handle,
		Name:   m.Res.Name,
	}

	if m.Res.ID > 0 {
		f.Query = fmt.Sprintf("roleID=%d", m.Res.ID)
	}

	return f
}
