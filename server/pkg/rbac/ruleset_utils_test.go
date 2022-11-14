package rbac

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test role inheritance
func TestRuleSet_merge(t *testing.T) {
	var (
		role1 uint64 = 1
		role2 uint64 = 2
		role3 uint64 = 3

		resService1 = "res1"
		resService2 = "res2"
		opAccess    = "access"
		resThing42  = "42"
		resThingWc  = "*"

		sCases = []struct {
			old       RuleSet
			new       RuleSet
			deletable RuleSet
			updatable RuleSet
			final     RuleSet
		}{
			{
				RuleSet{AllowRule(role1, resService1, opAccess)},
				RuleSet{AllowRule(role1, resService1, opAccess)},
				RuleSet{},
				RuleSet{},
				RuleSet{AllowRule(role1, resService1, opAccess)},
			},
			{
				RuleSet{AllowRule(role1, resService1, opAccess)},
				RuleSet{DenyRule(role1, resService1, opAccess)},
				RuleSet{},
				RuleSet{DenyRule(role1, resService1, opAccess)},
				RuleSet{DenyRule(role1, resService1, opAccess)},
			},
			{
				RuleSet{AllowRule(role1, resService1, opAccess)},
				RuleSet{InheritRule(role1, resService1, opAccess)},
				RuleSet{InheritRule(role1, resService1, opAccess)},
				RuleSet{},
				RuleSet{},
			},
			{
				RuleSet{AllowRule(role1, resService1, opAccess)},
				RuleSet{AllowRule(role1, resService1, opAccess)},
				RuleSet{},
				RuleSet{},
				RuleSet{AllowRule(role1, resService1, opAccess)},
			},
			{
				RuleSet{
					AllowRule(role1, resService1, opAccess),
					DenyRule(role2, resService1, opAccess),
					DenyRule(role3, resService2, opAccess),
					AllowRule(role1, resService2, opAccess),
					AllowRule(role2, resThing42, opAccess),
				},
				RuleSet{
					DenyRule(role3, resThingWc, opAccess),
					AllowRule(role1, resService2, opAccess),
					AllowRule(role1, resThing42, opAccess),
					InheritRule(role2, resThing42, opAccess),
				},
				RuleSet{
					InheritRule(role2, resThing42, opAccess),
				},
				RuleSet{
					DenyRule(role3, resThingWc, opAccess),
					AllowRule(role1, resThing42, opAccess),
				},
				RuleSet{
					AllowRule(role1, resService1, opAccess),
					DenyRule(role2, resService1, opAccess),
					DenyRule(role3, resService2, opAccess),
					AllowRule(role1, resService2, opAccess),
					DenyRule(role3, resThingWc, opAccess),
					AllowRule(role1, resThing42, opAccess),
				},
			},
		}
	)

	for _, sc := range sCases {
		t.Run("", func(t *testing.T) {
			var (
				req = require.New(t)

				// Apply changed and get update candidates
				mrg                         = merge(sc.old, sc.new...)
				deletable, updatable, final = flushable(mrg)
			)

			req.Equal(len(sc.deletable), len(deletable), "rule count for deletable do not match")
			req.Equal(len(sc.updatable), len(updatable), "rule count for updatable do not match")
			req.Equal(len(sc.final), len(final), "rule count for final do not match")

			// Clear dirty flag so that we do not confuse comparison test
			clear(deletable)
			clear(updatable)
			clear(final)

			req.Equal(sc.deletable, deletable, "deletable rules do not match")
			req.Equal(sc.updatable, updatable, "updatable rules do not match")
			req.Equal(sc.final, final, "final rules do not match")
		})
	}
}

// Test role inheritance
func TestRuleSet_sigRoles(t *testing.T) {
	const (
		roleA uint64 = 1
		roleB uint64 = 2
		roleC uint64 = 3
		roleD uint64 = 4
		roleE uint64 = 5

		opRead  = "read"
		opWrite = "write"
		resD    = "res:foo:1"
		resI    = "res:foo:*"
	)

	var (
		rr     = func(rr ...uint64) []uint64 { return rr }
		sCases = []struct {
			set RuleSet
			arr []uint64
			drr []uint64
		}{
			{
				RuleSet{
					AllowRule(roleA, resD, opRead),
					DenyRule(roleA, resD, opWrite),
					AllowRule(roleB, resI, opRead),
					DenyRule(roleA, resI, opRead),
					DenyRule(roleC, resI, opRead),
					DenyRule(roleD, resI, opRead),
					DenyRule(roleE, resI, opRead),
					DenyRule(roleE, resI, opWrite),
				},
				rr(roleB),
				rr(roleA, roleC, roleD, roleE),
			},
			{
				RuleSet{
					AllowRule(roleA, resD, opRead),
					DenyRule(roleA, resD, opRead),
				},
				rr(),
				rr(roleA),
			},
			{
				RuleSet{
					AllowRule(roleA, resD, opRead),
					DenyRule(roleA, resD, opRead),
					AllowRule(roleA, resI, opRead),
					DenyRule(roleA, resI, opRead),
				},
				rr(),
				rr(roleA),
			},
		}
	)

	for _, sc := range sCases {
		t.Run("a", func(t *testing.T) {
			req := require.New(t)
			arr, drr := sc.set.sigRoles(resD, opRead)

			sort.Slice(sc.arr, func(i, j int) bool { return sc.arr[i] < sc.arr[j] })
			sort.Slice(sc.drr, func(i, j int) bool { return sc.drr[i] < sc.drr[j] })
			sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
			sort.Slice(drr, func(i, j int) bool { return drr[i] < drr[j] })

			req.Equal(sc.arr, arr)
			req.Equal(sc.drr, drr)
		})
	}
}
