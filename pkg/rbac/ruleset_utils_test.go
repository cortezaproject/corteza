package rbac

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test role inheritance
func TestRuleSet_merge(t *testing.T) {
	var (
		req = require.New(t)

		role1 uint64 = 1
		role2 uint64 = 2
		role3 uint64 = 3

		resService1 = "res1"
		resService2 = "res2"
		opAccess    = "access"
		resThing42  = "42"
		resThingWc  = "*"

		sCases = []struct {
			old RuleSet
			new RuleSet
			del RuleSet
			upd RuleSet
		}{
			{
				RuleSet{AllowRule(role1, resService1, opAccess)},
				RuleSet{AllowRule(role1, resService1, opAccess)},
				RuleSet{},
				RuleSet{},
			},
			{
				RuleSet{AllowRule(role1, resService1, opAccess)},
				RuleSet{DenyRule(role1, resService1, opAccess)},
				RuleSet{},
				RuleSet{DenyRule(role1, resService1, opAccess)},
			},
			{
				RuleSet{AllowRule(role1, resService1, opAccess)},
				RuleSet{InheritRule(role1, resService1, opAccess)},
				RuleSet{InheritRule(role1, resService1, opAccess)},
				RuleSet{},
			},
			{
				RuleSet{AllowRule(role1, resService1, opAccess)},
				RuleSet{AllowRule(role1, resService1, opAccess)},
				RuleSet{},
				RuleSet{},
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
					// AllowRule(role1, resService1, opAccess),
					// DenyRule(role2, resService1, opAccess),
					// DenyRule(EveryoneRoleID, resService2, opAccess),
					// AllowRule(role1, resService2, opAccess),
					DenyRule(role3, resThingWc, opAccess),
					AllowRule(role1, resThing42, opAccess),
				},
			},
		}
	)

	for _, sc := range sCases {
		// Apply changed and get update candidates
		mrg := merge(sc.old, sc.new...)
		del, upd := flushable(mrg)

		// Clear dirty flag so that we do not confuse DeepEqual
		clear(del)
		clear(upd)

		req.Equal(len(sc.del), len(del))
		req.Equal(len(sc.upd), len(upd))
		req.Equal(sc.del, del)
		req.Equal(sc.upd, upd)
	}
}
