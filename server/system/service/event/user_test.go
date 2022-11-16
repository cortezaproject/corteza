package event

import (
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/system/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserMatching(t *testing.T) {
	var (
		a   = assert.New(t)
		res = &userBase{
			user: &types.User{Email: "user@example.tld"},
		}

		cUsr = eventbus.MustMakeConstraint("user.email", "eq", "user@example.tld")
	)

	a.True(res.Match(cUsr))
}
