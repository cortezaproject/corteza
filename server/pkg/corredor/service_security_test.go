package corredor

import (
	"context"
	"fmt"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type (
	// resolves roles by handle
	mockRoleByHandleSvc map[string]uint64
)

func (rr mockRoleByHandleSvc) FindByAny(_ context.Context, identifier interface{}) (*types.Role, error) {
	if id, ok := rr[fmt.Sprintf("%v", identifier)]; ok {
		return &types.Role{ID: id, Handle: fmt.Sprintf("%v", identifier)}, nil
	}

	return nil, fmt.Errorf("role not found")
}

func TestService_canExecAllow(t *testing.T) {
	const (
		adminRoleID   uint64 = 100
		managerRoleID uint64 = 200
		clientRoleID  uint64 = 300
		bypassRoleID  uint64 = 900
	)

	var (
		svc = &service{
			log:       zap.NewNop(),
			users:     &mockUserSvc{user: &types.User{ID: 42, Handle: "dummy"}},
			roles:     mockRoleByHandleSvc{"admin": adminRoleID, "manager": managerRoleID, "client": clientRoleID},
			denyExec:  make(map[string]map[uint64]bool),
			allowExec: make(map[string]map[uint64]bool),
		}

		manual = []*Trigger{{EventTypes: []string{onManualEventType}, ResourceTypes: []string{"res"}}}

		allowOnly    = &ServerScript{Name: "allow-only", Triggers: manual, Security: &Security{Allow: []string{"admin", "manager"}}}
		allowAndDeny = &ServerScript{Name: "allow-and-deny", Triggers: manual, Security: &Security{Allow: []string{"admin", "manager"}, Deny: []string{"manager"}}}
		denyOnly     = &ServerScript{Name: "deny-only", Triggers: manual, Security: &Security{Deny: []string{"client"}}}
		runAsOnly    = &ServerScript{Name: "run-as-only", Triggers: manual, Security: &Security{RunAs: "dummy"}}
		unrestricted = &ServerScript{Name: "unrestricted", Triggers: manual}

		as = func(roles ...uint64) context.Context {
			return auth.SetIdentityToContext(context.Background(), auth.Authenticated(42, roles...))
		}
	)

	auth.SetSystemRoles(types.RoleSet{{ID: bypassRoleID, Handle: auth.BypassRoleHandle}})

	svc.registerServerScripts(context.Background(), allowOnly, allowAndDeny, denyOnly, runAsOnly, unrestricted)
	require.Len(t, svc.sScripts, 5)

	t.Run("allow only", func(t *testing.T) {
		require.True(t, svc.canExec(as(adminRoleID), allowOnly.Name))
		require.True(t, svc.canExec(as(clientRoleID, managerRoleID), allowOnly.Name))
		require.False(t, svc.canExec(as(clientRoleID), allowOnly.Name), "role is not on the allow list")
		require.False(t, svc.canExec(as(), allowOnly.Name), "user without roles")
	})

	t.Run("deny wins over allow", func(t *testing.T) {
		require.True(t, svc.canExec(as(adminRoleID), allowAndDeny.Name))
		require.False(t, svc.canExec(as(managerRoleID), allowAndDeny.Name))
		require.False(t, svc.canExec(as(adminRoleID, managerRoleID), allowAndDeny.Name))
		require.False(t, svc.canExec(as(clientRoleID), allowAndDeny.Name))
	})

	t.Run("deny only", func(t *testing.T) {
		require.True(t, svc.canExec(as(adminRoleID), denyOnly.Name))
		require.True(t, svc.canExec(as(), denyOnly.Name))
		require.False(t, svc.canExec(as(clientRoleID), denyOnly.Name))
	})

	t.Run("without allow and deny", func(t *testing.T) {
		require.True(t, svc.canExec(as(clientRoleID), runAsOnly.Name))
		require.True(t, svc.canExec(as(clientRoleID), unrestricted.Name))
	})

	t.Run("bypass roles skip the allow list", func(t *testing.T) {
		require.True(t, svc.canExec(as(bypassRoleID), allowOnly.Name))
		require.True(t, svc.canExec(as(bypassRoleID), allowAndDeny.Name))
	})

	t.Run("restricted scripts are not listed", func(t *testing.T) {
		var (
			filter = Filter{ResourceTypes: []string{"res"}, EventTypes: []string{onManualEventType}}

			names = func(ctx context.Context) (nn []string) {
				set, _, err := svc.Find(ctx, filter)
				require.NoError(t, err)
				for _, s := range set {
					nn = append(nn, s.Name)
				}
				return
			}
		)

		require.NotContains(t, names(as(clientRoleID)), allowOnly.Name)
		require.Contains(t, names(as(adminRoleID)), allowOnly.Name)
	})

	t.Run("unknown role on the allow list", func(t *testing.T) {
		broken := &ServerScript{Name: "broken", Triggers: manual, Security: &Security{Allow: []string{"missing"}}}
		svc.registerServerScripts(context.Background(), broken)

		// script with invalid security must not end up as unrestricted
		require.False(t, svc.canExec(as(clientRoleID), broken.Name))
	})
}
