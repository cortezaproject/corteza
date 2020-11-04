package resource

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/system/types"
)

const (
	USER_RESOURCE_TYPE = "user"
)

type (
	// user represents a User
	user struct {
		*base
		Res *types.User
	}
)

func User(u *types.User) *user {
	r := &user{base: &base{}}
	r.SetResourceType(USER_RESOURCE_TYPE)
	r.Res = u

	r.AddIdentifier(identifiers(u.Handle, u.Email, u.Name, u.ID)...)

	return r
}

func (m *user) SearchQuery() types.UserFilter {
	f := types.UserFilter{
		Handle: m.Res.Handle,
		Email:  m.Res.Email,
	}

	if m.Res.ID > 0 {
		f.Query = fmt.Sprintf("userID=%d", m.Res.ID)
	}

	return f
}
