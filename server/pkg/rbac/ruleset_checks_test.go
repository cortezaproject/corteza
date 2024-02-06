package rbac

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_check(t *testing.T) {
	var (
		cc = []struct {
			name string
			exp  Access
			res  string
			op   string
			rr   []*Role
			set  RuleSet
		}{
			{"inherit when no roles or rules",
				Inherit, "", "", nil, nil},
			{
				"allow when checking with bypass roles",
				Allow,
				"",
				"",
				[]*Role{
					{id: 1, kind: BypassRole},
				},
				nil,
			},
			{
				"inherit when no matching roles",
				Inherit,
				"",
				"",
				[]*Role{
					{id: 1, kind: CommonRole},
				},
				[]*Rule{
					{RoleID: 2, Access: Deny},
				},
			},
			{
				"allow when matching rule",
				Allow,
				"",
				"",
				[]*Role{
					{id: 1, kind: CommonRole},
				},
				[]*Rule{
					{RoleID: 1, Access: Allow},
					{RoleID: 2, Access: Deny},
				},
			},
			{
				"multiple matching roles of same kind with deny",
				Deny,
				"",
				"",
				[]*Role{
					{id: 1, kind: CommonRole},
					{id: 2, kind: CommonRole},
				},
				[]*Rule{
					{RoleID: 1, Access: Allow},
					{RoleID: 2, Access: Deny},
				},
			},
			{
				"multiple matching matching roles of different with deny last",
				Allow,
				"",
				"",
				[]*Role{
					{id: 1, kind: CommonRole},
					{id: 2, kind: AuthenticatedRole},
				},
				[]*Rule{
					{RoleID: 1, Access: Allow},
					{RoleID: 2, Access: Deny},
				},
			},
			{
				"complex inheritance",
				Deny,
				"test::test:test/1/2/3",
				"",
				[]*Role{
					{id: 1, kind: CommonRole},
					{id: 2, kind: CommonRole},
				},
				[]*Rule{
					{RoleID: 1, Operation: "", Resource: "test::test:test/1/*/*", Access: Allow},
					{RoleID: 2, Operation: "", Resource: "test::test:test/*/*/3", Access: Allow},
					{RoleID: 2, Operation: "", Resource: "test::test:test/1/2/3", Access: Deny},
					{RoleID: 1, Operation: "", Resource: "test::test:test/*/2/3", Access: Allow},
				},
			},
		}
	)

	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.exp.String(), check(buildRuleIndex(c.set), partitionRoles(c.rr...), c.op, c.res, nil).String())
		})
	}
}

func Test_checkWithTrace(t *testing.T) {
	var (
		trace *Trace

		cc = []struct {
			name  string
			res   string
			exp   Access
			rr    []*Role
			set   RuleSet
			trace *Trace
		}{
			{
				"fail on integrity check (multiple anonymous roles)",
				"res-trace",
				Deny,
				[]*Role{
					{id: 1, kind: AnonymousRole},
					{id: 2, kind: CommonRole},
					{id: 3, kind: CommonRole},
				},
				nil,
				&Trace{
					Resource:   "res-trace",
					Operation:  "op-trace",
					Access:     Deny,
					Roles:      []uint64{1, 2, 3},
					Rules:      nil,
					Resolution: failedIntegrityCheck,
				},
			},
			{
				"allow when checking with bypass roles",
				"res-trace",
				Allow,
				[]*Role{
					{id: 1, kind: BypassRole},
				},
				nil,
				&Trace{
					Resource:   "res-trace",
					Operation:  "op-trace",
					Access:     Allow,
					Roles:      []uint64{1},
					Rules:      nil,
					Resolution: bypassRoleMembership,
				},
			},
			{
				"no rules",
				"res-trace",
				Allow,
				[]*Role{
					{id: 1, kind: CommonRole},
				},
				nil,
				&Trace{
					Resource:   "res-trace",
					Operation:  "op-trace",
					Access:     Inherit,
					Roles:      []uint64{1},
					Rules:      nil,
					Resolution: noRules,
				},
			},
			{
				"multi-role",
				"res-trace",
				Allow,
				[]*Role{
					{id: 1, kind: CommonRole},
					{id: 2, kind: CommonRole},
					{id: 3, kind: AuthenticatedRole},
				},
				RuleSet{
					AllowRule(1, "res-trace", "op-trace"),
					AllowRule(2, "res-trace", "op-trace"),
					AllowRule(3, "res-trace", "op-trace"),
					AllowRule(1, "res-trace-2", "op-trace"),
					AllowRule(2, "res-trace", "op-trace-2"),
				},
				&Trace{
					Resource:  "res-trace",
					Operation: "op-trace",
					Access:    Allow,
					Roles:     []uint64{1, 2, 3},
					Rules: RuleSet{
						AllowRule(1, "res-trace", "op-trace"),
						AllowRule(2, "res-trace", "op-trace"),
					},
				},
			},
			{
				"nested resource",
				"res-trace/2",
				Allow,
				[]*Role{
					{id: 1, kind: CommonRole},
					{id: 2, kind: CommonRole},
					{id: 3, kind: AuthenticatedRole},
				},
				RuleSet{
					AllowRule(1, "res-trace/*", "op-trace"),
					AllowRule(2, "res-trace/*", "op-trace"),
					AllowRule(2, "res-trace/1", "op-trace"),
				},
				&Trace{
					Resource:  "res-trace/2",
					Operation: "op-trace",
					Access:    Allow,
					Roles:     []uint64{1, 2, 3},
					Rules: RuleSet{
						AllowRule(1, "res-trace/*", "op-trace"),
						AllowRule(2, "res-trace/*", "op-trace"),
					},
				},
			},
		}
	)

	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			trace = new(Trace)
			check(buildRuleIndex(c.set), partitionRoles(c.rr...), "op-trace", c.res, trace)
			require.Equal(t, c.trace, trace)

		})
	}
}
