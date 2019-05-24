package permissions

import (
	"reflect"
	"testing"

	"github.com/cortezaproject/corteza-server/internal/test"
)

// Test role inheritance
func TestRuleSet_merge(t *testing.T) {
	var (
		assert = test.Assert

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

	for c, sc := range sCases {
		// Apply changed and get update candidates
		mrg := sc.old.merge(sc.new...)
		del, upd := mrg.dirty()

		// Clear dirty flag so that we do not confuse DeepEqual
		del.clear()
		upd.clear()

		assert(t, len(del) == len(sc.del), "Check test #%d failed, expected delete list length %d, got %d", c, len(sc.del), len(del))
		assert(t, len(upd) == len(sc.upd), "Check test #%d failed, expected update list length %d, got %d", c, len(sc.upd), len(upd))
		assert(t, reflect.DeepEqual(del, sc.del), "Check test #%d failed for delete list, reflect.DeepEqual == false", c)
		assert(t, reflect.DeepEqual(upd, sc.upd), "Check test #%d failed for update list, reflect.DeepEqual == false", c)
	}
}
