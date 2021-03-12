package automation

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUser(t *testing.T) {
	var (
		req    = require.New(t)
		u, err = NewUser(&types.User{Handle: "handle"})
	)

	req.NoError(err)
	req.Equal("handle", u.value.Handle)
	req.Error(u.AssignFieldValue("some-unexisting-field", nil))
	req.NoError(u.AssignFieldValue("email", "dummy@domain.tpl"))
	req.Equal("dummy@domain.tpl", u.value.Email)
}

func TestUser_Expr(t *testing.T) {
	var (
		req    = require.New(t)
		u, err = NewUser(&types.User{Handle: "hendl"})
	)

	req.NoError(err)

	eval, err := expr.NewParser().Parse("user.handle")
	req.NoError(err)

	res, err := eval.Eval(context.Background(), expr.RVars{"user": u}.Vars())
	req.NoError(err)

	req.Equal("hendl", res.(string))
}
