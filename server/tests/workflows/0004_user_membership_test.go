package workflows

import (
	"context"
	"testing"

	autTypes "github.com/cortezaproject/corteza/server/automation/types"
	sysTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

func Test0004_user_membership(t *testing.T) {
	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
	)

	req.NoError(defStore.TruncateRoleMembers(ctx))
	req.NoError(defStore.TruncateUsers(ctx))
	req.NoError(defStore.TruncateRoles(ctx))

	loadScenario(ctx, t)

	var (
		aux = struct {
			RolesPre_u1  sysTypes.RoleSet
			TotalPre_u1  uint64
			RolesPre_u2  sysTypes.RoleSet
			TotalPre_u2  uint64
			MemberPre_u1 bool
			MemberPre_u2 bool

			RolesPost_u1  sysTypes.RoleSet
			TotalPost_u1  uint64
			RolesPost_u2  sysTypes.RoleSet
			TotalPost_u2  uint64
			MemberPost_u1 bool
			MemberPost_u2 bool
		}{}
		vars, _ = mustExecWorkflow(ctx, t, "user-membership", autTypes.WorkflowExecParams{})
	)

	req.NoError(vars.Decode(&aux))

	req.Len(aux.RolesPre_u1, 2)
	req.Equal(uint64(2), aux.TotalPre_u1)
	req.True(aux.MemberPre_u1)
	req.Len(aux.RolesPre_u2, 0)
	req.Equal(uint64(0), aux.TotalPre_u2)
	req.False(aux.MemberPre_u2)

	req.Len(aux.RolesPost_u1, 1)
	req.Equal(uint64(1), aux.TotalPost_u1)
	req.False(aux.MemberPost_u1)
	req.Len(aux.RolesPost_u2, 0)
	req.Equal(uint64(0), aux.TotalPost_u2)
	req.False(aux.MemberPost_u2)
}
