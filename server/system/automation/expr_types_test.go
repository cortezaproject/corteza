package automation

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	var (
		req    = require.New(t)
		u, err = NewUser(&types.User{Handle: "handle"})
	)

	req.NoError(err)
	req.Equal("handle", u.value.Handle)
	req.Error(u.AssignFieldValue("some-unexisting-field", nil))
	req.NoError(u.AssignFieldValue("email", expr.Must(expr.NewString("dummy@domain.tpl"))))
	req.Equal("dummy@domain.tpl", u.value.Email)
}

func TestUser_Expr(t *testing.T) {
	var (
		req   = require.New(t)
		u, _  = NewUser(&types.User{Handle: "hendl"})
		scope = &expr.Vars{}
	)

	req.NoError(scope.Set("user", u))

	eval, err := expr.NewParser().Parse("user.handle")
	req.NoError(err)

	res, err := eval.Eval(context.Background(), scope)
	req.NoError(err)

	req.Equal("hendl", res.(string))
}

func TestCastToRbacResource(t *testing.T) {
	var (
		req = require.New(t)
	)

	t.Run("string to RbacResource", func(t *testing.T) {
		resourceString := "corteza::system:user/461133310995726337"
		resource, err := CastToRbacResource(resourceString)

		req.NoError(err)
		req.NotNil(resource)
		req.Equal(resourceString, resource.RbacResource())
	})

	t.Run("existing RbacResource", func(t *testing.T) {
		originalResource := rbac.NewResource("corteza::system:role/123")
		resource, err := CastToRbacResource(originalResource)

		req.NoError(err)
		req.Equal(originalResource, resource)
	})

	t.Run("RbacResource expression type", func(t *testing.T) {
		originalResource := rbac.NewResource("corteza::system:role/123")
		rbacRes, err := NewRbacResource(originalResource)
		req.NoError(err)

		resource, err := CastToRbacResource(rbacRes)
		req.NoError(err)
		req.Equal(originalResource, resource)
	})

	t.Run("invalid type", func(t *testing.T) {
		resource, err := CastToRbacResource(123)
		req.Error(err)
		req.Nil(resource)
		req.Contains(err.Error(), "unable to cast type int to")
	})
}
