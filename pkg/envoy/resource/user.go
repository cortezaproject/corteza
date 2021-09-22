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
	r.SetResourceType(types.UserResourceType)
	r.Res = u

	r.AddIdentifier(identifiers(u.Handle, u.Email, u.Name, u.ID)...)

	// Initial timestamps
	r.SetTimestamps(MakeTimestampsCUDAS(&u.CreatedAt, u.UpdatedAt, u.DeletedAt, nil, u.SuspendedAt))

	return r
}

func (r *User) SysID() uint64 {
	return r.Res.ID
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
