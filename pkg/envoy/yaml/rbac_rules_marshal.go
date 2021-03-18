package yaml

import (
	"context"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

func (n *rbacRule) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	rl, ok := state.Res.(*resource.RbacRule)
	if !ok {
		return encoderErrInvalidResource(resource.RBAC_RESOURCE_TYPE, state.Res.ResourceType())
	}

	// Get the related role
	n.relRole = resource.FindRole(state.ParentResources, rl.RefRole.Identifiers)
	if n.relRole == nil {
		return resource.RoleErrUnresolved(rl.RefRole.Identifiers)
	}
	if n.relRole == nil {
		return resource.RoleErrUnresolved(rl.RefRole.Identifiers)
	}

	// For now we will only allow resource specific RBAC rules if that resource is
	// also present.
	// Here we check if we can find it in case we're handling a resource specific rule.
	refRes := n.refRes
	if refRes != nil && len(refRes.Identifiers) > 0 {
		for _, r := range state.ParentResources {
			if refRes.ResourceType == r.ResourceType() && r.Identifiers().HasAny(refRes.Identifiers) {
				n.relResource = r
				break
			}
		}

		if n.relResource == nil {
			// We couldn't find it...
			return resource.RbacResourceErrNotFound(refRes.Identifiers)
		}
	}

	return nil
}

func (r *rbacRule) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
	// @todo Improve RBAC rule placement
	//
	// In cases where a specific rule is created for a specific resource, nest the rule
	// under the related namespace.
	// For now all rules will be nested under a root node for simplicity sake.

	doc.addRbacRule(r)

	return nil
}

func (rr rbacRuleSet) MarshalYAML() (interface{}, error) {
	if rr == nil || len(rr) == 0 {
		return nil, nil
	}

	addRef := func(r *rbacRule, base rbac.Resource) string {
		rtr := base.TrimID().String()

		if r.relResource == nil {
			return rtr
		}

		refr, ok := r.relResource.(resource.RefableInterface)
		if !ok {
			return rtr
		}

		return rtr + refr.Ref()
	}

	var err error
	accNode, _ := makeMap()

	for i, accRules := range rr.groupByAccess() {
		roleNode, _ := makeMap()

		for _, roleRules := range accRules.groupByRole() {
			resNode, _ := makeMap()
			for _, resRules := range roleRules.groupByResource() {
				opNode, _ := makeSeq()

				for _, rule := range resRules {
					opNode, err = addSeq(opNode, rule.res.Operation.String())
					if err != nil {
						return nil, err
					}
				}

				resNode, err = addMap(resNode,
					strings.TrimRight(addRef(resRules[0], resRules[0].res.Resource), ":"), opNode,
				)
				if err != nil {
					return nil, err
				}
			}

			rr := roleRules[0].relRole
			rk := rr.Handle
			if rk == "" {
				rk = rr.Name
			}
			roleNode, err = addMap(roleNode,
				rk, resNode,
			)
			if err != nil {
				return nil, err
			}
		}

		if i == 0 {
			accNode, err = addMap(accNode,
				"allow", roleNode,
			)
		} else {
			accNode, err = addMap(accNode,
				"deny", roleNode,
			)
		}
		if err != nil {
			return nil, err
		}
	}

	return accNode, nil
}

func (r *rbacRule) MarshalYAML() (interface{}, error) {
	return r.res.Operation.String(), nil
}
