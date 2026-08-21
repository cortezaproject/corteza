package provision

import (
	"testing"

	"github.com/crusttech/human/server/pkg/rbac"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// Whether a base grant ever reaches an existing install comes down to this
// check, and getting it wrong is silent: the grant lands for new installs and
// nobody else, which is exactly the failure the markers exist to prevent.
func Test_baseMarkerMissing(t *testing.T) {
	const (
		adminID = 10
		probeID = 20
	)

	roles := types.RoleSet{
		{ID: adminID, Handle: "admin"},
		// The kind of role a deployment or an earlier probe leaves lying about.
		{ID: probeID, Handle: "some-ad-hoc-role"},
	}

	allow := func(roleID uint64, resource, op string) *rbac.Rule {
		return &rbac.Rule{RoleID: roleID, Resource: resource, Operation: op, Access: rbac.Allow}
	}

	// Everything the markers ask for, granted to the role they name.
	complete := rbac.RuleSet{}
	for _, m := range baseMarkers {
		complete = append(complete, allow(adminID, m.resource, m.operation))
	}

	t.Run("an install holding every marker is left alone", func(t *testing.T) {
		require.False(t, baseMarkerMissing(roles, complete, zap.NewNop()))
	})

	t.Run("a missing marker triggers the re-import", func(t *testing.T) {
		for i := range baseMarkers {
			short := rbac.RuleSet{}
			short = append(short, complete[:i]...)
			short = append(short, complete[i+1:]...)
			require.True(t, baseMarkerMissing(roles, short, zap.NewNop()),
				"dropping %q should have triggered a re-import", baseMarkers[i].what)
		}
	})

	t.Run("another role holding the grant does not count", func(t *testing.T) {
		// The regression: a probe role carrying one user-group rule convinced
		// the old check that the install had permissions admin never got.
		foreign := rbac.RuleSet{}
		for _, m := range baseMarkers {
			foreign = append(foreign, allow(probeID, m.resource, m.operation))
		}
		require.True(t, baseMarkerMissing(roles, foreign, zap.NewNop()))
	})

	t.Run("a deny does not count as holding the grant", func(t *testing.T) {
		denied := rbac.RuleSet{}
		for _, m := range baseMarkers {
			denied = append(denied, &rbac.Rule{
				RoleID: adminID, Resource: m.resource, Operation: m.operation, Access: rbac.Deny,
			})
		}
		require.True(t, baseMarkerMissing(roles, denied, zap.NewNop()))
	})

	t.Run("a rule whose role no longer exists is ignored", func(t *testing.T) {
		orphan := rbac.RuleSet{}
		for _, m := range baseMarkers {
			orphan = append(orphan, allow(999, m.resource, m.operation))
		}
		require.True(t, baseMarkerMissing(roles, orphan, zap.NewNop()))
	})
}
