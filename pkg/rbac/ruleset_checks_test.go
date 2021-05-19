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
			require.Equal(t, c.exp, c.set.Check(partitionRoles(c.rr...), c.res, c.op))
		})
	}
}

func benchmarkCheck(b *testing.B, c int) {
	var (
		// resting with 50 roles
		rules = make(RuleSet, 0, c)

		pr = partitionRoles(
			&Role{id: 1, kind: CommonRole},
			&Role{id: 2, kind: CommonRole},
			&Role{id: 3, kind: CommonRole},
			&Role{id: 4, kind: CommonRole},
			&Role{id: 5, kind: CommonRole},
			&Role{id: 6, kind: CommonRole},
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
		checkOptimised(iRules, pr, "res-0", "op-0")
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
			{Allow, "res", "op", []*Rule{
				{Resource: "---", Operation: "--", Access: Deny},
				{Resource: "res", Operation: "op", Access: Allow},
			}},
		}
	)

	for _, c := range cc {
		t.Run("", func(t *testing.T) {
			require.Equal(t, c.exp, checkRulesByResource(c.set, c.res, c.op))
		})
	}
}

//// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //
//// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //
//// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //// //
//
//const (
//	role1 uint64 = 10001
//	role2 uint64 = 10002
//
//	resService1 = Resource("service1")
//	resService2 = Resource("service2")
//
//	resThingWc = Resource("some:answer:*")
//	resThing13 = Resource("some:answer:13")
//	resThing42 = Resource("some:answer:42")
//
//	opAccess = "access"
//	opRead   = "read"
//	opWrite  = "write"
//)
//
//func TestRuleSet_check(t *testing.T) {
//	var (
//		req = require.New(t)
//
//		rr = RuleSet{
//			AllowRule(role1, resThing42, opRead),
//			DenyRule(role1, resThing13, opWrite),
//			AllowRule(role2, resThing13, opWrite),
//		}
//
//		sCases = []struct {
//			roles    []uint64
//			res      Resource
//			op       Operation
//			expected Access
//		}{
//			{[]uint64{role1}, resThing42, opRead, Allow},
//			{[]uint64{role1}, resThing42, opWrite, Inherit},
//			{[]uint64{role1}, resThing13, opWrite, Deny},
//			{[]uint64{role2}, resThing13, opWrite, Allow},
//			{[]uint64{role1, role2}, resThing13, opWrite, Deny},
//			{[]uint64{role1, role2}, resThing42, opRead, Allow},
//		}
//	)
//
//	for c, sc := range sCases {
//		v := rr.check(sc.res, sc.op, sc.roles...)
//		req.Equalf(sc.expected, v, "Check test #%d failed, expected %s, got %s", c, sc.expected, v)
//	}
//}
//
//// Test resource inheritance
//func TestRuleSet_checkResource(t *testing.T) {
//	const (
//		role1 uint64 = 10001
//
//		resService1 = Resource("service1")
//		resService2 = Resource("service2")
//
//		resThingWc = Resource("some:answer:*")
//		resThing13 = Resource("some:answer:13")
//		resThing42 = Resource("some:answer:42")
//
//		opAccess = "access"
//	)
//
//	var (
//		r = require.New(t)
//
//		sCases = []struct {
//			rr       RuleSet
//			roles    []uint64
//			res      Resource
//			op       Operation
//			expected Access
//		}{
//			{
//				RuleSet{
//					AllowRule(role1, resService1, opAccess),
//				},
//				[]uint64{role1},
//				resService1,
//				opAccess,
//				Allow,
//			},
//			{
//				RuleSet{
//					AllowRule(role1, resThingWc, opAccess),
//				},
//				[]uint64{role1},
//				resThing42,
//				opAccess,
//				Allow,
//			},
//			{ // deny wc and explictly allow 42
//				RuleSet{
//					DenyRule(role1, resThingWc, opAccess),
//					AllowRule(role1, resThing42, opAccess),
//				},
//				[]uint64{role1},
//				resThing42,
//				opAccess,
//				Allow,
//			},
//			{ // deny wc and explictly allow 42
//				RuleSet{
//					DenyRule(role1, resThingWc, opAccess),
//					AllowRule(role1, resThing42, opAccess),
//				},
//				[]uint64{role1},
//				resThing13,
//				opAccess,
//				Deny,
//			},
//			{ // deny wc and and check if wc is denied
//				RuleSet{
//					DenyRule(role1, resThingWc, opAccess),
//					AllowRule(role1, resThing42, opAccess),
//				},
//				[]uint64{role1},
//				resThingWc,
//				opAccess,
//				Deny,
//			},
//			{ // allow wc and and check if wc is allowed
//				RuleSet{
//					AllowRule(role1, resThingWc, opAccess),
//					DenyRule(role1, resThing42, opAccess),
//				},
//				[]uint64{role1},
//				resThingWc,
//				opAccess,
//				Allow,
//			},
//		}
//	)
//
//	for c, sc := range sCases {
//		v := sc.rr.checkResource(sc.res, sc.op, sc.roles...)
//		r.Equalf(sc.expected, v, "Check test #%d failed, expected %s, got %s", c, sc.expected, v)
//	}
//}
//
//// Test role inheritance
//func TestRuleSet_Check(t *testing.T) {
//	var (
//		rr = RuleSet{
//			// 1st level
//			AllowRule(role1, resService1, opAccess),
//			DenyRule(role2, resService1, opAccess),
//			// 2nd level
//			DenyRule(EveryoneRoleID, resService2, opAccess),
//			AllowRule(EveryoneRoleID, resThing13, opAccess),
//			AllowRule(role1, resService2, opAccess),
//			// 3rd level
//			DenyRule(EveryoneRoleID, resThingWc, opAccess),
//			AllowRule(role1, resThing42, opAccess),
//		}
//
//		r = require.New(t)
//
//		sCases = []struct {
//			roles    []uint64
//			res      Resource
//			op       Operation
//			expected Access
//		}{
//			{[]uint64{role1}, resService1, opAccess, Allow},
//			{[]uint64{role2}, resService1, opAccess, Deny},
//			{[]uint64{role1}, resService2, opAccess, Allow},
//			{[]uint64{role2}, resService2, opAccess, Deny},
//			{[]uint64{role1}, resThing42, opAccess, Allow},
//			{[]uint64{role2}, resThing42, opAccess, Deny},
//			{[]uint64{}, resThing42, opAccess, Deny},
//			{[]uint64{}, resThing13, opAccess, Allow},
//		}
//	)
//
//	for c, sc := range sCases {
//		v := rr.Check(sc.res, sc.op, sc.roles...)
//		r.Equalf(sc.expected, v, "Check test #%d failed, expected %s, got %s", c, sc.expected, v)
//	}
//}
