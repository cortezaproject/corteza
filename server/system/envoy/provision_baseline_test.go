package envoy

import (
	"testing"

	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/stretchr/testify/require"
)

// The base access-control config is the whole of the deny-by-default baseline:
// what a logged-in user holds before any role, and what an older install has
// taken away from it. A typo in a resource or an operation there is silent —
// the rule imports and simply never matches — so decode the real file and read
// back what it actually says.
func TestProvisionBaseAccessControl(t *testing.T) {
	type ruleKey = provisionRuleKey

	got := decodeProvisionedRules(t, "../../provision/000_base/system_access_control.yaml")

	// rbac.Deny is the zero value of rbac.Access, so a rule that is simply
	// absent reads as a denial. Every assertion below goes through this.
	access := func(k ruleKey) string {
		v, ok := got[k]
		if !ok {
			return "<no rule>"
		}
		return v.String()
	}

	// A plain logged-in user may fetch the app menu; every entry in it is then
	// filtered per application, so this reveals nothing on its own.
	require.Equal(t, "allow", access(ruleKey{"authenticated", "corteza::system/", "applications.search"}),
		"the shell cannot gate section entry on a list it may not fetch")

	// ...and pin one to their own menu, which reaches only apps they can
	// already see. The operation is component-level; a rule written against
	// corteza::system:application/* is never consulted.
	require.Equal(t, "allow", access(ruleKey{"authenticated", "corteza::system/", "application.flag.self"}))

	// ...but not see any application, user or role by default.
	for _, denied := range []ruleKey{
		{"authenticated", "corteza::system:application/*", "read"},
		{"authenticated", "corteza::system:user/*", "read"},
		{"authenticated", "corteza::system:role/*", "read"},
		{"authenticated", "corteza::system/", "users.search"},
		{"authenticated", "corteza::system/", "roles.search"},
		{"authenticated", "corteza::system/", "auth-clients.search"},
	} {
		require.Equal(t, "deny", access(denied), "expected a deny for %v", denied)
	}

	// Opening a webapp is its own operation, and administrators hold it.
	require.Equal(t, "allow", access(ruleKey{"admin", "corteza::system:application/*", "access"}))
	require.Equal(t, "allow", access(ruleKey{"security-admin", "corteza::system:application/*", "access"}))

	// Four admin screens were reachable but answered "not allowed" to the role
	// that owns them, because nothing here granted their operations.
	for _, held := range []ruleKey{
		{"admin", "corteza::system/", "user-groups.search"},
		{"admin", "corteza::system/", "user-group.create"},
		{"admin", "corteza::system/", "connections.search"},
		{"admin", "corteza::system/", "connection.create"},
		{"admin", "corteza::system/", "dal-connections.search"},
		{"admin", "corteza::system/", "dal-connection.create"},
		{"admin", "corteza::system:user-group/*", "members.manage"},
		{"admin", "corteza::system:connection/*", "install"},
		{"admin", "corteza::system:dal-connection/*", "dal-config.manage"},
		{"admin", "corteza::system/", "labels.search"},
		{"admin", "corteza::system/", "corredor-scripts.search"},
		{"admin", "corteza::system/", "chatbot-sessions.search"},
		{"admin", "corteza::system:chatbot/*", "sessions.view"},
		{"admin", "corteza::system:chatbot/*", "sessions.handoff.manage"},
		{"admin", "corteza::system:chatbot-session/*", "read"},
		{"admin", "corteza::system:configured-connection/*", "read"},
		{"admin", "corteza::system/", "dal-sensitivity-level.manage"},
		{"admin", "corteza::system/", "resource-translations.manage"},
		{"admin", "corteza::system:project-ai-system/*", "resources.manage"},
		{"low-code-admin", "corteza::system:application/*", "access"},
		{"security-admin", "corteza::system/", "user-groups.search"},
	} {
		require.Equal(t, "allow", access(held), "expected an allow for %v", held)
	}

	// Nothing else is granted to everyone: the baseline is login and the menu.
	var authenticatedAllows []ruleKey
	for k, v := range got {
		if k.role == "authenticated" && v == rbac.Allow {
			authenticatedAllows = append(authenticatedAllows, k)
		}
	}
	require.Len(t, authenticatedAllows, 3, "authenticated baseline grew: %v", authenticatedAllows)
}
