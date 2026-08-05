package automation

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/system/types"
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

func TestTemplate_metaAssign(t *testing.T) {
	var (
		req     = require.New(t)
		t1, err = NewTemplate(&types.Template{
			Handle: "invoice",
			Meta:   types.TemplateMeta{Short: "before"},
		})
	)

	req.NoError(err)

	// fields of the meta struct
	req.NoError(expr.Assign(t1, "meta.short", expr.Must(expr.NewString("after"))))
	req.Equal("after", t1.value.Meta.Short)

	req.NoError(expr.Assign(t1, "meta.headerTemplateID", expr.Must(expr.NewID(42))))
	req.Equal(uint64(42), t1.value.Meta.HeaderTemplateID)

	req.NoError(expr.Assign(t1, "meta.footerTemplateID", expr.Must(expr.NewID(43))))
	req.Equal(uint64(43), t1.value.Meta.FooterTemplateID)

	req.Error(expr.Assign(t1, "meta.unknown-field", expr.Must(expr.NewString("nope"))))

	// whole meta struct
	meta, err := NewTemplateMeta(types.TemplateMeta{Short: "whole", HeaderTemplateID: 99})
	req.NoError(err)
	req.NoError(expr.Assign(t1, "meta", meta))
	req.Equal("whole", t1.value.Meta.Short)
	req.Equal(uint64(99), t1.value.Meta.HeaderTemplateID)

	// meta from a map, the way it comes in from a workflow
	fromMap, err := NewTemplateMeta(map[string]interface{}{"short": "mapped", "footerTemplateID": "13"})
	req.NoError(err)
	req.NoError(expr.Assign(t1, "meta", fromMap))
	req.Equal("mapped", t1.value.Meta.Short)
	req.Equal(uint64(13), t1.value.Meta.FooterTemplateID)

	// other fields still assign
	req.NoError(expr.Assign(t1, "handle", expr.Must(expr.NewString("receipt"))))
	req.Equal("receipt", t1.value.Handle)
	req.Error(expr.Assign(t1, "some-unexisting-field", expr.Must(expr.NewString("nope"))))
}

// Assigning through a scope, the way workflow expression steps do it
func TestTemplate_metaAssignInScope(t *testing.T) {
	var (
		req     = require.New(t)
		t1, err = NewTemplate(&types.Template{Handle: "invoice"})
		scope   = &expr.Vars{}
	)

	req.NoError(err)
	req.NoError(scope.Set("template", t1))

	req.NoError(expr.Assign(scope, "template.meta.headerTemplateID", expr.Must(expr.NewID(42))))
	req.NoError(expr.Assign(scope, "template.meta.short", expr.Must(expr.NewString("scoped"))))

	selected, err := scope.Select("template")
	req.NoError(err)

	tpl, err := CastToTemplate(selected)
	req.NoError(err)
	req.Equal(uint64(42), tpl.Meta.HeaderTemplateID)
	req.Equal("scoped", tpl.Meta.Short)
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
