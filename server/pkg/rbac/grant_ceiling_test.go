package rbac

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/pkg/auth"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// ceilingService is a real service with one rule in it, so the ceiling is read
// through the same evaluation a request would take.
func ceilingService(t *testing.T, callerRoleID uint64, rules ...*Rule) *service {
	t.Helper()

	svc := NewService(zap.NewNop(), nil)
	svc.rules = rules
	svc.indexed = buildRuleIndex(svc.rules)
	svc.orgTree = &orgTree{}
	svc.roles = []*Role{CommonRole.Make(callerRoleID, "caller")}

	prev := Global()
	SetGlobal(svc)
	t.Cleanup(func() { SetGlobal(prev) })

	return svc
}

// ceilingContext is what a request carries: an identity in the context, which
// is where CanPassOn reads the caller's roles from.
func ceilingContext(roleID uint64) context.Context {
	return auth.SetIdentityToContext(context.Background(), auth.Authenticated(1, roleID))
}

func TestEnforceGrantCeiling(t *testing.T) {
	const (
		callerRole = uint64(10)
		targetRole = uint64(20)
		resource   = "corteza::system:role/*"
	)

	t.Run("a rule for an operation the caller holds passes", func(t *testing.T) {
		ceilingService(t, callerRole, AllowRule(callerRole, resource, "read"))
		ctx := ceilingContext(callerRole)

		require.NoError(t, EnforceGrantCeiling(ctx, AllowRule(targetRole, resource, "read")))
	})

	// The hole this closes: holding a component's grant let a caller hand out
	// every operation the component declares, including ones they could not
	// perform themselves.
	t.Run("a rule for an operation the caller lacks is refused", func(t *testing.T) {
		ceilingService(t, callerRole, AllowRule(callerRole, resource, "read"))
		ctx := ceilingContext(callerRole)

		err := EnforceGrantCeiling(ctx, AllowRule(targetRole, resource, "delete"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "delete")
		require.Contains(t, err.Error(), "Nothing was changed")
	})

	t.Run("one bad rule refuses the whole set", func(t *testing.T) {
		ceilingService(t, callerRole, AllowRule(callerRole, resource, "read"))
		ctx := ceilingContext(callerRole)

		require.Error(t, EnforceGrantCeiling(ctx,
			AllowRule(targetRole, resource, "read"),
			AllowRule(targetRole, resource, "delete"),
		))
	})

	// An inherit deletes a rule, and the rule it deletes may be the deny that
	// was holding an operation shut.
	t.Run("inherit is judged like the rest", func(t *testing.T) {
		ceilingService(t, callerRole, AllowRule(callerRole, resource, "read"))
		ctx := ceilingContext(callerRole)

		require.Error(t, EnforceGrantCeiling(ctx, InheritRule(targetRole, resource, "delete")))
		require.NoError(t, EnforceGrantCeiling(ctx, InheritRule(targetRole, resource, "read")))
	})

	t.Run("no rules, nothing to refuse", func(t *testing.T) {
		ceilingService(t, callerRole, AllowRule(callerRole, resource, "read"))
		require.NoError(t, EnforceGrantCeiling(context.Background()))
	})

	// Refusing when nothing is evaluating access would break the paths that
	// write rules before the service exists, not close a hole.
	t.Run("no RBAC service means no ceiling to read", func(t *testing.T) {
		prev := Global()
		SetGlobal(nil)
		t.Cleanup(func() { SetGlobal(prev) })

		require.NoError(t, EnforceGrantCeiling(context.Background(), AllowRule(targetRole, resource, "delete")))
	})
}
