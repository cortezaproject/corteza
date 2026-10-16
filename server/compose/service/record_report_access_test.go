package service

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// reportReadAC reads every record and lists the given rules.
type reportReadAC struct {
	recordAccessController
	rules rbac.RuleSet
}

func (reportReadAC) CanReadRecord(context.Context, *types.Record) bool { return true }
func (a reportReadAC) ruleSet() (rbac.RuleSet, bool)                   { return a.rules, true }
func (reportReadAC) readsWithoutContextRoles(context.Context, *types.Record) (bool, bool) {
	return true, true
}

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

// A contextual role whose expression matches the ownerless probe, and which
// alone grants reading, says nothing about the owned records: the report reads
// them one by one. A module-wide allow keeps reports on the database, an
// owner's allow next to it included; a contextual deny does not.
func TestReportReadsAllRecordsContextRoles(t *testing.T) {
	const (
		userRole    = 10
		unownedRole = 20
		ownerRole   = 30
	)

	var (
		m       = &types.Module{ID: 2, NamespaceID: 1}
		records = types.RecordRbacResource(1, 2, 0)

		setup = func(rr ...*rbac.Rule) (context.Context, accessControl) {
			svc := rbac.NewService(zap.NewNop(), nil)
			svc.UpdateRoles(
				rbac.CommonRole.Make(userRole, "user"),
				rbac.MakeContextRole(unownedRole, "unowned", func(scope map[string]interface{}) bool {
					res, _ := scope["resource"].(map[string]interface{})
					return res["ownedBy"] == uint64(0)
				}, types.RecordResourceType),
				rbac.MakeContextRole(ownerRole, "owner", func(scope map[string]interface{}) bool {
					res, _ := scope["resource"].(map[string]interface{})
					return res["ownedBy"] == scope["userID"]
				}, types.RecordResourceType),
			)
			require.NoError(t, svc.Grant(context.Background(), rr...))

			return auth.SetIdentityToContext(context.Background(), auth.Authenticated(5, userRole)), accessControl{rbac: svc}
		}
	)

	ctx, ac := setup(rbac.AllowRule(unownedRole, records, "read"))
	require.False(t, readsAllRecords(ctx, ac, m), "only a contextual role matching the probe allows reading")

	ctx, ac = setup(
		rbac.AllowRule(userRole, records, "read"),
		rbac.AllowRule(ownerRole, records, "read"),
	)
	require.True(t, readsAllRecords(ctx, ac, m), "a module-wide allow, an owner's allow next to it")

	ctx, ac = setup(
		rbac.AllowRule(userRole, records, "read"),
		rbac.DenyRule(ownerRole, records, "read"),
	)
	require.False(t, readsAllRecords(ctx, ac, m), "a contextual deny")
}
