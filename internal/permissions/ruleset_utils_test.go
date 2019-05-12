package permissions

import (
	"reflect"
	"testing"

	"github.com/crusttech/crust/internal/test"
)

// Test role inheritance
func TestRuleSet_merge(t *testing.T) {
	var (
		assert = test.Assert

		sCases = []struct {
			old RuleSet
			in  RuleSet
			exp RuleSet
		}{
			{
				RuleSet{
					&Rule{role1, resService1, opAccess, Allow},
					&Rule{role2, resService1, opAccess, Deny},
					&Rule{EveryoneRoleID, resService2, opAccess, Deny},
					&Rule{role1, resService2, opAccess, Allow},
				},
				RuleSet{
					&Rule{EveryoneRoleID, resThingWc, opAccess, Deny},
					&Rule{role1, resThing42, opAccess, Allow},
					&Rule{role1, resThing42, opAccess, Inherit},
				},
				RuleSet{
					&Rule{role1, resService1, opAccess, Allow},
					&Rule{role2, resService1, opAccess, Deny},
					&Rule{EveryoneRoleID, resService2, opAccess, Deny},
					&Rule{role1, resService2, opAccess, Allow},
					&Rule{EveryoneRoleID, resThingWc, opAccess, Deny},
					&Rule{role1, resThing42, opAccess, Allow},
					&Rule{role1, resThing42, opAccess, Inherit},
				},
			},
		}
	)

	for c, sc := range sCases {
		out, _ := sc.old.merge(sc.in...)

		assert(t, len(out) == len(sc.exp), "Check test #%d failed, expected length %d, got %d", c, len(out), len(sc.exp))
		assert(t, reflect.DeepEqual(out, sc.exp), "Check test #%d failed, reflect.DeepEqual == false", c)

	}
}

// Test role inheritance
func TestRuleSet_split(t *testing.T) {
	var (
		assert = test.Assert

		sCases = []struct {
			set RuleSet
			i   RuleSet
			r   RuleSet
		}{
			{
				RuleSet{
					&Rule{role1, resService1, opAccess, Allow},
					&Rule{role2, resService1, opAccess, Deny},
					&Rule{EveryoneRoleID, resService2, opAccess, Inherit},
				},
				RuleSet{
					&Rule{EveryoneRoleID, resService2, opAccess, Inherit},
				},
				RuleSet{
					&Rule{role1, resService1, opAccess, Allow},
					&Rule{role2, resService1, opAccess, Deny},
				},
			},
		}
	)

	for c, sc := range sCases {
		i, r := sc.set.split()

		assert(t, len(i) == len(sc.i), "Check test #%d failed, expected length %d, got %d", c, len(i), len(sc.i))
		assert(t, len(r) == len(sc.r), "Check test #%d failed, expected length %d, got %d", c, len(r), len(sc.r))
		assert(t, reflect.DeepEqual(i, sc.i), "Check test #%d failed, reflect.DeepEqual == false", c)
		assert(t, reflect.DeepEqual(r, sc.r), "Check test #%d failed, reflect.DeepEqual == false", c)

	}
}
