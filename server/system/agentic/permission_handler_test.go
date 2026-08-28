package agentic

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/stretchr/testify/require"
)

// stubAccessControl stands in for a component's generated access controller.
// buildRules only reaches it to record which component a rule belongs to.
type stubAccessControl struct{}

func (stubAccessControl) CanGrant(context.Context) bool { return true }
func (stubAccessControl) List() []map[string]string {
	return []map[string]string{
		{"type": "corteza::system:application", "any": "corteza::system:application/*", "op": "read"},
		{"type": "corteza::system:application", "any": "corteza::system:application/*", "op": "delete"},
		{"type": "corteza::system:role", "any": "corteza::system:role/*", "op": "read"},
	}
}
func (stubAccessControl) FindRulesByRoleID(context.Context, uint64) (rbac.RuleSet, error) {
	return nil, nil
}
func (stubAccessControl) Grant(context.Context, ...*rbac.Rule) error { return nil }

// holdsOnly builds the predicate buildRules gates on: the caller holds exactly
// the resource/operation pairs listed and nothing else.
func holdsOnly(pairs ...[2]string) func(string, string) bool {
	return func(resource, operation string) bool {
		for _, p := range pairs {
			if p[0] == resource && p[1] == operation {
				return true
			}
		}
		return false
	}
}

const app = "corteza::system:application/511"

// stubComponents is the component list buildRules resolves resources against.
// The real one is empty until the services boot, which is the safe direction —
// but it means a test has to supply its own.
func stubComponents() []component {
	return []component{{name: "system", prefix: "corteza::system", ac: stubAccessControl{}}}
}

// The whole point of the tool: an agent may pass on only what it already has.
func TestBuildRulesRefusesWhatTheCallerDoesNotHold(t *testing.T) {
	_, _, err := buildRules(1, []permissionRuleArg{
		{Resource: app, Operation: "delete", Access: "allow"},
	}, false, "permission grant", stubComponents(), holdsOnly([2]string{app, "read"}))

	require.Error(t, err)
	require.Contains(t, err.Error(), `you do not hold "delete"`)
}

func TestBuildRulesAllowsWhatTheCallerHolds(t *testing.T) {
	byComponent, _, err := buildRules(1, []permissionRuleArg{
		{Resource: app, Operation: "read", Access: "allow"},
	}, false, "permission grant", stubComponents(), holdsOnly([2]string{app, "read"}))

	require.NoError(t, err)
	require.Len(t, byComponent["system"], 1)
	require.Equal(t, rbac.Allow, byComponent["system"][0].Access)
}

// Revoking is a privilege change too, so it carries the same ceiling.
func TestBuildRulesRefusesARevokeTheCallerDoesNotHold(t *testing.T) {
	_, _, err := buildRules(1, []permissionRuleArg{
		{Resource: app, Operation: "update"},
	}, true, "permission revoke", stubComponents(), holdsOnly([2]string{app, "read"}))

	require.Error(t, err)
	require.Contains(t, err.Error(), `you do not hold "update"`)
}

// One refused rule refuses the call. A permission change that half-lands leaves
// a role in a state nobody asked for.
func TestBuildRulesIsAllOrNothing(t *testing.T) {
	byComponent, _, err := buildRules(1, []permissionRuleArg{
		{Resource: app, Operation: "read", Access: "allow"},
		{Resource: app, Operation: "delete", Access: "allow"},
	}, false, "permission grant", stubComponents(), holdsOnly([2]string{app, "read"}))

	require.Error(t, err)
	require.Nil(t, byComponent)
}

func TestBuildRulesRefusesAResourceOutsideEveryComponent(t *testing.T) {
	// The short form the server itself rejects, caught here with the reason.
	_, _, err := buildRules(1, []permissionRuleArg{
		{Resource: "application/*", Operation: "read"},
	}, false, "permission grant", stubComponents(), func(string, string) bool { return true })

	require.Error(t, err)
	require.Contains(t, err.Error(), "belongs to no component")
}

func TestRuleAccess(t *testing.T) {
	for _, tc := range []struct {
		raw      string
		revoking bool
		want     rbac.Access
		wantErr  string
	}{
		{raw: "", want: rbac.Allow},
		{raw: "allow", want: rbac.Allow},
		{raw: "DENY", want: rbac.Deny},
		{raw: "allow", revoking: true, want: rbac.Inherit},
		{raw: "inherit", wantErr: "system_permission_revoke"},
		{raw: "maybe", wantErr: `access "maybe"`},
	} {
		got, err := ruleAccess(tc.raw, tc.revoking)
		if tc.wantErr != "" {
			require.Errorf(t, err, "access %q", tc.raw)
			require.Contains(t, err.Error(), tc.wantErr)
			continue
		}
		require.NoErrorf(t, err, "access %q", tc.raw)
		require.Equalf(t, tc.want, got, "access %q", tc.raw)
	}
}

// The schema collapses one row per resource-type-and-operation pair into one
// entry per type; a caller reads the operations off that entry.
func TestPermissionTypesGroupsOperationsByResourceType(t *testing.T) {
	got := permissionTypes(component{name: "system", prefix: "corteza::system", ac: stubAccessControl{}})

	require.Len(t, got, 2)
	require.Equal(t, "corteza::system:application", got[0].ResourceType)
	require.Equal(t, []string{"read", "delete"}, got[0].Operations)
	require.Equal(t, "corteza::system:application/*", got[0].Resource)
}

func TestMatchesPermissionTypeAcceptsTheBareTail(t *testing.T) {
	require.True(t, matchesPermissionType("corteza::compose:module", "module"))
	require.True(t, matchesPermissionType("corteza::compose:module", "corteza::compose:module"))
	require.False(t, matchesPermissionType("corteza::compose:module", "modul"))
	require.False(t, matchesPermissionType("corteza::compose:module-field", "module"))
}
