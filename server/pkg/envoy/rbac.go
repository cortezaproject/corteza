package envoy

import (
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	resutil "github.com/cortezaproject/corteza/server/pkg/resource"
)

// FilterRequestedRBACRules returns only RBAC rules relevant for the given resources
func FilterRequestedRBACRules(request resource.InterfaceSet, rules []*resource.RbacRule) (out []*resource.RbacRule) {
	out = make([]*resource.RbacRule, 0, 10)

	ruleIndex := resutil.NewIndex()
	for _, r := range rules {
		ruleIndex.Add(r, r.IndexPath()...)
	}

	// Filter
	dupRuleIndex := make(map[string]bool)
	procResSet(request, func(r resource.Interface) {
		if r.Placeholder() {
			return
		}

		rbacRes, ok := r.(resource.RBACInterface)
		if !ok {
			return
		}

		_, ref, pp := rbacRes.RBACParts()
		resPath := refsToIndexPath(r.ResourceType(), appendRefSet(pp, ref)...)
		for _, _rule := range ruleIndex.Collect(resPath...) {
			// All will be *resource.RbacRule so no further checks are needed
			rule := _rule.(*resource.RbacRule)

			k := fmt.Sprintf("%s, %s, %d; %d", rule.Res.Resource, rule.Res.Operation, rule.Res.Access, rule.Res.RoleID)
			if dupRuleIndex[k] {
				continue
			}
			dupRuleIndex[k] = true

			out = append(out, rule)
		}
	})

	return
}

func refsToIndexPath(rt string, pp ...*resource.Ref) (out [][]string) {
	out = append(out, []string{rt})
	for _, p := range pp {
		out = append(out, p.Identifiers)
	}

	return
}
