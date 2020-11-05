package resource

import (
	"fmt"

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

func (m *User) SearchQuery() types.UserFilter {
	f := types.UserFilter{
		Handle: m.Res.Handle,
		Email:  m.Res.Email,
	}

	if m.Res.ID > 0 {
		f.Query = fmt.Sprintf("userID=%d", m.Res.ID)
	}

	return f
}
