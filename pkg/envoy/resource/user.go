package resource

import (
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	// User represents a User
	User struct {
		*base
		Res *types.User
	}
)

func NewUser(u *types.User) *User {
	r := &User{base: &base{}}
	r.SetResourceType(USER_RESOURCE_TYPE)
	r.Res = u

	r.AddIdentifier(identifiers(u.Handle, u.Email, u.Name, u.ID)...)

	return r
}

func (r *User) SysID() uint64 {
	return r.Res.ID
}
