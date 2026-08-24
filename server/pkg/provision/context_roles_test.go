package provision

import (
	"testing"

	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

// The whole point of this repair is that it is invisible when it works and
// invisible when it does not: a context role with no resource types imports,
// holds its rules and never activates. So assert on what it leaves alone as
// hard as on what it fills.
func Test_fillContextRoleTypes(t *testing.T) {
	const (
		record = "corteza::compose:record"
		chart  = "corteza::compose:chart"
	)

	declared := map[string][]string{
		"agent-creator": {record, chart},
	}

	contextual := func(expr string, rr ...string) *types.Role {
		return &types.Role{
			Handle: "agent-creator",
			Meta:   &types.RoleMeta{Context: &types.RoleContext{Expr: expr, Resource: rr}},
		}
	}

	t.Run("an expression with no types is what the regression left behind", func(t *testing.T) {
		r := contextual("agentID != 0")
		out := fillContextRoleTypes(declared, map[string]*types.Role{"agent-creator": r})

		require.Len(t, out, 1)
		require.Equal(t, []string{record, chart}, r.Meta.Context.Resource)
	})

	t.Run("a list a deployment narrowed is left alone", func(t *testing.T) {
		r := contextual("agentID != 0", record)
		out := fillContextRoleTypes(declared, map[string]*types.Role{"agent-creator": r})

		require.Empty(t, out)
		require.Equal(t, []string{record}, r.Meta.Context.Resource)
	})

	t.Run("a role with no expression is not made contextual", func(t *testing.T) {
		r := contextual("")
		out := fillContextRoleTypes(declared, map[string]*types.Role{"agent-creator": r})

		require.Empty(t, out)
		require.Empty(t, r.Meta.Context.Resource)
	})

	for _, tc := range []struct {
		name   string
		stored map[string]*types.Role
	}{
		{"the role is not in the store yet", map[string]*types.Role{}},
		{"the role carries no meta", map[string]*types.Role{"agent-creator": {Handle: "agent-creator"}}},
		{"the role carries no context", map[string]*types.Role{
			"agent-creator": {Handle: "agent-creator", Meta: &types.RoleMeta{}},
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Empty(t, fillContextRoleTypes(declared, tc.stored))
		})
	}
}

// The repair reads the same file provisioning imports, so a role renamed or a
// context dropped there silently narrows what it can fix.
func Test_declaredContextRoles(t *testing.T) {
	got, err := declaredContextRoles("../../provision/*")
	require.NoError(t, err)
	require.NotEmpty(t, got)

	require.Contains(t, got, "agent-creator")
	require.Contains(t, got["agent-creator"], "corteza::compose:chart")

	for handle, rr := range got {
		require.NotEmpty(t, rr, "context role %q declares no resource types", handle)
	}
}

func Test_declaredContextRoles_noBaseDir(t *testing.T) {
	got, err := declaredContextRoles("../../provision/001_settings")
	require.NoError(t, err)
	require.Empty(t, got)
}
