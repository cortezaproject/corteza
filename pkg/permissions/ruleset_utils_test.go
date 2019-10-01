package permissions

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test role inheritance
func TestRuleSet_merge(t *testing.T) {
	var (
		req = require.New(t)

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
					DenyRule(EveryoneRoleID, resService2, opAccess),
					AllowRule(role1, resService2, opAccess),
					AllowRule(role2, resThing42, opAccess),
				},
				RuleSet{
					DenyRule(EveryoneRoleID, resThingWc, opAccess),
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
					DenyRule(EveryoneRoleID, resThingWc, opAccess),
					AllowRule(role1, resThing42, opAccess),
				},
			},
		}
	)

	for _, sc := range sCases {
		// Apply changed and get update candidates
		mrg := sc.old.merge(sc.new...)
		del, upd := mrg.dirty()

		// Clear dirty flag so that we do not confuse DeepEqual
		del.clear()
		upd.clear()

		req.Equal(len(sc.del), len(del))
		req.Equal(len(sc.upd), len(upd))
		req.Equal(sc.del, del)
		req.Equal(sc.upd, upd)
	}
}
