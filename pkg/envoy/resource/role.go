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
	r.SetResourceType(types.RoleResourceType)
	r.Res = rl

	r.AddIdentifier(identifiers(rl.Handle, rl.Name, rl.ID)...)

	// Initial timestamps
	r.SetTimestamps(MakeTimestampsCUDA(&rl.CreatedAt, rl.UpdatedAt, rl.DeletedAt, rl.ArchivedAt))

	return r
}

func (r *Role) SysID() uint64 {
	return r.Res.ID
}

// FindRole looks for the role in the resources
func FindRole(rr InterfaceSet, ii Identifiers) (rl *types.Role) {
	var rlRes *Role

	rr.Walk(func(r Interface) error {
		rr, ok := r.(*Role)
		if !ok {
			return nil
		}

		if rr.Identifiers().HasAny(ii) {
			rlRes = rr
		}
		return nil
	})

	// Found it
	if rlRes != nil {
		return rlRes.Res
	}

	return nil
}

func RoleErrUnresolved(ii Identifiers) error {
	return fmt.Errorf("role unresolved %v", ii.StringSlice())
}
