package event

import (
	"github.com/crusttech/human/server/pkg/eventbus"
	"github.com/crusttech/human/server/system/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleMemberMatching(t *testing.T) {
	var (
		a   = assert.New(t)
		res = &roleMemberBase{
			role: &types.Role{Handle: "admin"},
			user: &types.User{Handle: "user"},
		}

		cRol = eventbus.MustMakeConstraint("role", "eq", "admin")
		cUsr = eventbus.MustMakeConstraint("user", "eq", "user")
	)

	a.True(res.Match(cRol))
	a.True(res.Match(cUsr))
}
