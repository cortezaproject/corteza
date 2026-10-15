package service

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/stretchr/testify/require"
)

// reportReadAC reads every record and lists the given rules.
type reportReadAC struct {
	recordAccessController
	rules rbac.RuleSet
}

func (reportReadAC) CanReadRecord(context.Context, *types.Record) bool { return true }
func (a reportReadAC) ruleSet() (rbac.RuleSet, bool)                   { return a.rules, true }

// reportRolesAC also lists the roles that can apply to the caller.
type reportRolesAC struct {
	reportReadAC
	roles map[uint64]bool
}

func (a reportRolesAC) applicableRoles(context.Context) (map[uint64]bool, bool) { return a.roles, true }

// Whoever reads every record of a module reports straight from the database;
// a rule denying one of its records sends the report through the records one
// by one.
func TestReportReadableFastPath(t *testing.T) {
	m := &types.Module{ID: 2, NamespaceID: 1}

	ids, err := record{ac: reportReadAC{}}.reportReadable(context.Background(), m, "")
	require.NoError(t, err)
	require.Nil(t, ids, "no record is read one by one")

	require.False(t, recordReadDenied(context.Background(), reportReadAC{rules: rbac.RuleSet{
		rbac.DenyRule(9, types.RecordRbacResource(1, 3, 7), "read"),
		rbac.DenyRule(9, types.RecordRbacResource(1, 2, 7), "update"),
	}}, m), "a deny on another module's record or another operation")
	require.True(t, recordReadDenied(context.Background(), reportReadAC{rules: rbac.RuleSet{
		rbac.DenyRule(9, types.RecordRbacResource(1, 2, 7), "read"),
	}}, m), "a deny on one of the module's records")
	require.True(t, recordReadDenied(context.Background(), reportReadAC{rules: rbac.RuleSet{
		rbac.DenyRule(9, types.RecordRbacResource(0, 0, 0), "read"),
	}}, m), "a deny on every record")
	require.True(t, recordReadDenied(context.Background(), struct{}{}, m), "rules that cannot be listed may hold a deny")
}

// A deny given to a role the caller can not hold leaves their reports on the
// database; one given to a role they hold, or may hold, does not.
func TestReportReadDeniedForCallerRoles(t *testing.T) {
	var (
		ctx  = context.Background()
		m    = &types.Module{ID: 2, NamespaceID: 1}
		deny = rbac.RuleSet{rbac.DenyRule(9, types.RecordRbacResource(1, 2, 0), "read")}
	)

	require.False(t, recordReadDenied(ctx, reportRolesAC{reportReadAC{rules: deny}, map[uint64]bool{5: true}}, m),
		"a deny for a role the caller can not hold")
	require.True(t, recordReadDenied(ctx, reportRolesAC{reportReadAC{rules: deny}, map[uint64]bool{5: true, 9: true}}, m),
		"a deny for a role the caller can hold")
	require.True(t, recordReadDenied(ctx, reportReadAC{rules: deny}, m),
		"roles that cannot be listed may include the denied one")
}
