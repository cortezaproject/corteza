package resource

import (
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	// User represents a User
	User struct {
		*base
		Res *types.User

		RoleMembership []Identifiers
		RefRoles       RefSet
	}
)

func NewUser(u *types.User, roles ...string) *User {
	r := &User{base: &base{}}
	r.SetResourceType(types.UserResourceType)
	r.Res = u

	r.AddIdentifier(identifiers(u.Handle, u.Email, u.Name, u.ID)...)

	// Role membership
	for _, role := range roles {
		rid := MakeIdentifiers(role)
		r.RoleMembership = append(r.RoleMembership, rid)

		r.RefRoles = append(r.RefRoles, r.AddRef(types.RoleResourceType, role))
	}

	// Initial timestamps
	r.SetTimestamps(MakeTimestampsCUDAS(&u.CreatedAt, u.UpdatedAt, u.DeletedAt, nil, u.SuspendedAt))

	return r
}

func (u *User) AddRoles(roles ...*types.Role) *User {
	for _, r := range roles {
		idf := firstOkString(r.Handle, strconv.FormatUint(r.ID, 10))
		rid := MakeIdentifiers(idf)
		u.RoleMembership = append(u.RoleMembership, rid)

		u.RefRoles = append(u.RefRoles, u.AddRef(types.RoleResourceType, idf))
	}

	return u
}

func (r *User) Resource() interface{} {
	return r.Res
}

func (r *User) SysID() uint64 {
	return r.Res.ID
}

func (r *User) RBACParts() (resource string, ref *Ref, path []*Ref) {
	ref = r.Ref()
	path = nil
	resource = fmt.Sprintf(types.UserRbacResourceTpl(), types.UserResourceType, firstOkString(strconv.FormatUint(r.Res.ID, 10), r.Res.Handle, r.Res.Username))

	return
}

// FindUser looks for the user in the resources
func FindUser(rr InterfaceSet, ii Identifiers) (u *types.User) {
	var uRes *User

	rr.Walk(func(r Interface) error {
		ur, ok := r.(*User)
		if !ok {
			return nil
		}

		if ur.Identifiers().HasAny(ii) {
			uRes = ur
		}
		return nil
	})

	// Found it
	if uRes != nil {
		return uRes.Res
	}

	return nil
}

func UserErrUnresolved(ii Identifiers) error {
	return fmt.Errorf("user unresolved %v", ii.StringSlice())
}
