package workflows

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/store"
	"github.com/stretchr/testify/require"
)

// Roles imported without timestamps get a creation time instead of the zero
// time (which SQL Server can not even store)
func Test_import_without_timestamps(t *testing.T) {
	ctx := bypassRBAC(context.Background())
	loadScenarioWithName(ctx, t, "S0007_iterator_roles")

	r, err := store.LookupRoleByHandle(ctx, defStore, "r1")
	require.NoError(t, err)
	require.False(t, r.CreatedAt.IsZero(), "created at %v", r.CreatedAt)
}
