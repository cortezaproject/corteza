package permissions

import (
	"testing"

	"github.com/cortezaproject/corteza-server/internal/test"
)

const (
	role1 uint64 = 10001
	role2 uint64 = 10002

	resService1 = Resource("service1")
	resService2 = Resource("service2")

	resThingWc = Resource("some:answer:*")
	resThing13 = Resource("some:answer:13")
	resThing42 = Resource("some:answer:42")

	opAccess = "access"
	opRead   = "read"
	opWrite  = "write"
)

func TestRuleSet_check(t *testing.T) {
	var (
		assert = test.Assert

		rr = RuleSet{
			AllowRule(role1, resThing42, opRead),
			DenyRule(role1, resThing13, opWrite),
			AllowRule(role2, resThing13, opWrite),
		}

		sCases = []struct {
			roles    []uint64
			res      Resource
			op       Operation
			expected Access
		}{
			{[]uint64{role1}, resThing42, opRead, Allow},
			{[]uint64{role1}, resThing42, opWrite, Inherit},
			{[]uint64{role1}, resThing13, opWrite, Deny},
			{[]uint64{role2}, resThing13, opWrite, Allow},
			{[]uint64{role1, role2}, resThing13, opWrite, Deny},
			{[]uint64{role1, role2}, resThing42, opRead, Allow},
		}
	)

	for c, sc := range sCases {
		v := rr.check(sc.res, sc.op, sc.roles...)
		assert(t, v == sc.expected, "Check test #%d failed, expected %s, got %s", c, sc.expected, v)
	}
}

// Test resource inheritance
func TestRuleSet_checkResource(t *testing.T) {
	const (
		role1 uint64 = 10001

		resService1 = Resource("service1")
		resService2 = Resource("service2")

		resThingWc = Resource("some:answer:*")
		resThing13 = Resource("some:answer:13")
		resThing42 = Resource("some:answer:42")

		opAccess = "access"
	)

	var (
		assert = test.Assert

		sCases = []struct {
			rr       RuleSet
			roles    []uint64
			res      Resource
			op       Operation
			expected Access
		}{
			{
				RuleSet{
					AllowRule(role1, resService1, opAccess),
				},
				[]uint64{role1},
				resService1,
				opAccess,
				Allow,
			},
			{
				RuleSet{
					AllowRule(role1, resThingWc, opAccess),
				},
				[]uint64{role1},
				resThing42,
				opAccess,
				Allow,
			},
			{ // deny wc and explictly allow 42
				RuleSet{
					DenyRule(role1, resThingWc, opAccess),
					AllowRule(role1, resThing42, opAccess),
				},
				[]uint64{role1},
				resThing42,
				opAccess,
				Allow,
			},
			{ // deny wc and explictly allow 42
				RuleSet{
					DenyRule(role1, resThingWc, opAccess),
					AllowRule(role1, resThing42, opAccess),
				},
				[]uint64{role1},
				resThing13,
				opAccess,
				Deny,
			},
		}
	)

	for c, sc := range sCases {
		v := sc.rr.checkResource(sc.res, sc.op, sc.roles...)
		assert(t, v == sc.expected, "Check test #%d failed, expected %s, got %s", c, sc.expected, v)
	}
}

// Test role inheritance
func TestRuleSet_Check(t *testing.T) {
	var (
		rr = RuleSet{
			// 1st level
			AllowRule(role1, resService1, opAccess),
			DenyRule(role2, resService1, opAccess),
			// 2nd level
			DenyRule(EveryoneRoleID, resService2, opAccess),
			AllowRule(EveryoneRoleID, resThing13, opAccess),
			AllowRule(role1, resService2, opAccess),
			// 3rd level
			DenyRule(EveryoneRoleID, resThingWc, opAccess),
			AllowRule(role1, resThing42, opAccess),
		}

		assert = test.Assert

		sCases = []struct {
			roles    []uint64
			res      Resource
			op       Operation
			expected Access
		}{
			{[]uint64{role1}, resService1, opAccess, Allow},
			{[]uint64{role2}, resService1, opAccess, Deny},
			{[]uint64{role1}, resService2, opAccess, Allow},
			{[]uint64{role2}, resService2, opAccess, Deny},
			{[]uint64{role1}, resThing42, opAccess, Allow},
			{[]uint64{role2}, resThing42, opAccess, Deny},
			{[]uint64{}, resThing42, opAccess, Deny},
			{[]uint64{}, resThing13, opAccess, Allow},
		}
	)

	for c, sc := range sCases {
		v := rr.Check(sc.res, sc.op, sc.roles...)
		assert(t, v == sc.expected, "Check test #%d failed, expected %s, got %s", c, sc.expected, v)
	}
}
