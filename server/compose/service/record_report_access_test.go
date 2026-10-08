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

// Whoever reads every record of a module reports straight from the database;
// a rule denying one of its records sends the report through the records one
// by one.
func TestReportReadableFastPath(t *testing.T) {
	m := &types.Module{ID: 2, NamespaceID: 1}

	ids, err := record{ac: reportReadAC{}}.reportReadable(context.Background(), m, "")
	require.NoError(t, err)
	require.Nil(t, ids, "no record is read one by one")

	require.False(t, recordReadDenied(reportReadAC{rules: rbac.RuleSet{
		rbac.DenyRule(9, types.RecordRbacResource(1, 3, 7), "read"),
		rbac.DenyRule(9, types.RecordRbacResource(1, 2, 7), "update"),
	}}, m), "a deny on another module's record or another operation")
	require.True(t, recordReadDenied(reportReadAC{rules: rbac.RuleSet{
		rbac.DenyRule(9, types.RecordRbacResource(1, 2, 7), "read"),
	}}, m), "a deny on one of the module's records")
	require.True(t, recordReadDenied(reportReadAC{rules: rbac.RuleSet{
		rbac.DenyRule(9, types.RecordRbacResource(0, 0, 0), "read"),
	}}, m), "a deny on every record")
	require.True(t, recordReadDenied(struct{}{}, m), "rules that cannot be listed may hold a deny")
}
