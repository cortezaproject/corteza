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
		}
	)

	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.exp.String(), check(indexRules(c.set), partitionRoles(c.rr...), c.op, c.res).String())
		})
	}
}

//cpu: Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz
//Benchmark_Check100-16                	15626196	        88.85 ns/op
//Benchmark_Check1000-16               	15976252	        74.09 ns/op
//Benchmark_Check10000-16              	15025586	        78.12 ns/op
//Benchmark_Check100000-16             	13760616	        84.70 ns/op
//Benchmark_Check1000000-16            	 2602420	       415.8 ns/op
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
		check(iRules, pr, "res-0", "op-0")
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
			require.Equal(t, c.exp.String(), checkRulesByResource(c.set, c.op, c.res).String())
		})
	}
}
