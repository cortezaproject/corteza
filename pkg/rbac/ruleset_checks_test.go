package rbac

import (
	"fmt"
	"math/rand"
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
			require.Equal(t, c.exp.String(), check(indexRules(c.set), partitionRoles(c.rr...), c.op, c.res, nil).String())
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
			check(indexRules(c.set), partitionRoles(c.rr...), "op-trace", c.res, trace)
			require.Equal(t, c.trace, trace)

		})
	}
}

//cpu: Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz
//Benchmark_Check100-16        	12395438	        95.24 ns/op
//Benchmark_Check1000-16       	12507883	        96.34 ns/op
//Benchmark_Check10000-16      	11788594	        96.85 ns/op
//Benchmark_Check100000-16     	11679951	       100.1 ns/op
//Benchmark_Check1000000-16    	 4670353	       287.3 ns/op
func benchmarkCheck(b *testing.B, c int) {
	var (
		// resting with 50 roles
		rules = make(RuleSet, 0, c)

		pr = partitionRoles(
			&Role{id: 1, kind: CommonRole},
			&Role{id: 2, kind: CommonRole},
			&Role{id: 3, kind: CommonRole},
			&Role{id: 4, kind: ContextRole},
			&Role{id: 5, kind: ContextRole},
			&Role{id: 6, kind: AuthenticatedRole},
		)
	)

	for i := 0; i < cap(rules); i++ {
		rules = append(rules, &Rule{
			RoleID:    uint64(rand.Int31n(50)),
			Resource:  fmt.Sprintf("res-%d", rand.Int31n(1000)),
			Operation: fmt.Sprintf("op-%d", rand.Int31n(100)),
			Access:    Access(rand.Int31n(2)),
		})
	}

	iRules := indexRules(rules)

	b.StartTimer()

	for n := 0; n < b.N; n++ {
		check(iRules, pr, "res-0", "op-0", nil)
	}

	b.StopTimer()
}

func Benchmark_Check100(b *testing.B)     { benchmarkCheck(b, 100) }
func Benchmark_Check1000(b *testing.B)    { benchmarkCheck(b, 1000) }
func Benchmark_Check10000(b *testing.B)   { benchmarkCheck(b, 10000) }
func Benchmark_Check100000(b *testing.B)  { benchmarkCheck(b, 100000) }
func Benchmark_Check1000000(b *testing.B) { benchmarkCheck(b, 1000000) }

func Test_checkRulesByResource(t *testing.T) {
	var (
		cc = []struct {
			exp Access
			res string
			op  string
			set []*Rule
		}{
			{Inherit, "", "", nil},
			{Inherit, "res", "op", nil},
			{Allow, "res/1", "op", []*Rule{
				{Resource: "---", Operation: "--", Access: Deny},
				{Resource: "res/1", Operation: "op", Access: Allow},
			}},
			{Allow, "res/2", "op", []*Rule{
				{Resource: "res/*", Operation: "op", Access: Allow},
			}},
			{Allow, "res/3", "op", []*Rule{
				{Resource: "res/3", Operation: "op", Access: Allow},
				{Resource: "res/*", Operation: "op", Access: Deny},
			}},
		}
	)

	for _, c := range cc {
		t.Run(c.res, func(t *testing.T) {
			a := findRuleByResOp(c.set, c.op, c.res)

			if a == nil {
				a = InheritRule(0, "", "")
			}

			require.Equal(t, c.exp.String(), a.Access.String())
		})
	}
}
