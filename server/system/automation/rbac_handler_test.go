package automation

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// A workflow grants a permission when its runner may grant on the resource's
// component, written as component resources are: with a trailing /
func TestRbacHandlerGrantChecksComponent(t *testing.T) {
	const (
		granter = 1
		other   = 2
	)

	var (
		svc  = rbac.NewService(zap.NewNop(), nil)
		h    = rbacHandler{rbac: svc}
		rule = rbac.AllowRule(other, "corteza::compose:record/1/2/3", "read")

		as = func(role uint64) context.Context {
			return auth.SetIdentityToContext(context.Background(), auth.Authenticated(5, role))
		}
	)

	svc.UpdateRoles(rbac.CommonRole.Make(granter, "granter"), rbac.CommonRole.Make(other, "other"))
	require.NoError(t, svc.Grant(context.Background(), rbac.AllowRule(granter, "corteza::compose/", "grant")))

	require.EqualError(t, h.grant(as(other), rule), "not allowed to grant "+rule.String())

	require.NoError(t, h.grant(as(granter), rule))
	require.Equal(t, rbac.Allow, svc.Check(rbac.ParamsToSession(context.Background(), 6, other), "read", rbac.NewResource("corteza::compose:record/1/2/3")))
}
