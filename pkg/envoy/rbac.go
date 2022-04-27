package envoy

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

// FilterRequestedRBACRules returns only RBAC rules relevant for the given resources
func FilterRequestedRBACRules(request resource.InterfaceSet, rules []*resource.RbacRule) (out []*resource.RbacRule) {
	out = make([]*resource.RbacRule, 0, 10)

	// Filter
	dupRuleIndex := make(map[string]bool)
	procResSet(request, func(r resource.Interface) {
		rbacRes, ok := r.(resource.RBACInterface)
		if !ok {
			return
		}

		_, ref, pp := rbacRes.RBACParts()
		resourceRefSet := appendRefSet(pp, ref)

		for _, rule := range rules {
			k := fmt.Sprintf("%s, %s, %d; %d", rule.Res.Resource, rule.Res.Operation, rule.Res.Access, rule.Res.RoleID)
			if dupRuleIndex[k] {
				continue
			}
			dupRuleIndex[k] = true
			ruleRefSet := appendRefSet(rule.RefPath, rule.RefRes)
			// Checking if rule is <= resource since wildflags can be used
			if !ruleRefSet.IsSubset(resourceRefSet) {
				continue
			}

			out = append(out, rule)
		}
	})

	return
}
