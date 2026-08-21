package rbac

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Revoking a grant an earlier baseline gave the authenticated role means
// denying it there — provisioning only ever adds rules, so an install that
// carried the grant keeps it otherwise.
//
// That is only a safe revocation because role kinds are resolved in priority
// order: a common role's allow returns before the authenticated tier is read
// at all. These cases pin that, since the whole migration rests on it.
func Test_checkAuthenticatedDenyDoesNotReachRoleHolders(t *testing.T) {
	const (
		authenticatedRoleID = 1
		adminRoleID         = 2
	)

	var (
		roleless = []*Role{
			{id: authenticatedRoleID, kind: AuthenticatedRole},
		}

		admin = []*Role{
			{id: authenticatedRoleID, kind: AuthenticatedRole},
			{id: adminRoleID, kind: CommonRole},
		}

		// The shape 000_base produces: the baseline denies the operation to
		// everyone logged in, and the admin role allows it.
		baseline = RuleSet{
			{RoleID: authenticatedRoleID, Resource: "corteza::system:application/1", Operation: "read", Access: Deny},
			{RoleID: adminRoleID, Resource: "corteza::system:application/1", Operation: "read", Access: Allow},
		}
	)

	cc := []struct {
		name string
		exp  Access
		rr   []*Role
		set  RuleSet
	}{
		{"role-less user is denied by the baseline", Deny, roleless, baseline},
		{"a role that allows it wins over the baseline deny", Allow, admin, baseline},
		{"without the deny the role-less user only ever gets inherit", Inherit, roleless, RuleSet{baseline[1]}},
	}

	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(
				t,
				c.exp.String(),
				check(buildRuleIndex(c.set), partitionRoles(c.rr...), "read", "corteza::system:application/1", nil).String(),
			)
		})
	}
}
