package yaml

import (
	"context"
	"fmt"

	automationTypes "github.com/cortezaproject/corteza-server/automation/types"
	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	federationTypes "github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
)

func (n *rbacRule) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	rl, ok := state.Res.(*resource.RbacRule)
	if !ok {
		return encoderErrInvalidResource(resource.RbacResourceType, state.Res.ResourceType())
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

	res := state.Res.(*resource.RbacRule)

	p0ID := "*"
	p1ID := "*"

	switch r.res.Resource {
	case composeTypes.ComponentResourceType,
		systemTypes.ComponentResourceType,
		automationTypes.ComponentResourceType,
		federationTypes.ComponentResourceType:
		r.res.Resource += "/"
		goto next
	}

	// @todo the following stuff I'm not too sure about
	switch r.res.Resource {
	case composeTypes.NamespaceResourceType:
		if res.RefResource != nil {
			p1 := resource.FindComposeNamespace(state.ParentResources, res.RefResource.Identifiers)
			if p1 == nil {
				return resource.ComposeNamespaceErrUnresolved(res.RefResource.Identifiers)
			}
			p1ID = p1.Slug
		}

		r.res.Resource = fmt.Sprintf(composeTypes.NamespaceRbacResourceTpl(), composeTypes.NamespaceResourceType, p1ID)
	case composeTypes.ModuleResourceType:
		if len(res.RefPath) > 0 {
			p0 := resource.FindComposeNamespace(state.ParentResources, res.RefPath[0].Identifiers)
			if p0 == nil {
				return resource.ComposeNamespaceErrUnresolved(res.RefPath[0].Identifiers)
			}
			p0ID = p0.Slug
		}

		if res.RefResource != nil {
			p1 := resource.FindComposeModule(state.ParentResources, res.RefResource.Identifiers)
			if p1 == nil {
				return resource.ComposeModuleErrUnresolved(res.RefResource.Identifiers)
			}
			p1ID = p1.Handle
		}

		r.res.Resource = fmt.Sprintf(composeTypes.ModuleRbacResourceTpl(), composeTypes.ModuleResourceType, p0ID, p1ID)
	case composeTypes.ChartResourceType:
		if len(res.RefPath) > 0 {
			p0 := resource.FindComposeNamespace(state.ParentResources, res.RefPath[0].Identifiers)
			if p0 == nil {
				return resource.ComposeNamespaceErrUnresolved(res.RefPath[0].Identifiers)
			}
			p0ID = p0.Slug
		}

		if res.RefResource != nil {
			p1 := resource.FindComposeChart(state.ParentResources, res.RefResource.Identifiers)
			if p1 == nil {
				return resource.ComposeChartErrUnresolved(res.RefResource.Identifiers)
			}
			p1ID = p1.Handle
		}

		r.res.Resource = fmt.Sprintf(composeTypes.ChartRbacResourceTpl(), composeTypes.ChartResourceType, p0ID, p1ID)
	case composeTypes.PageResourceType:
		if len(res.RefPath) > 0 {
			p0 := resource.FindComposeNamespace(state.ParentResources, res.RefPath[0].Identifiers)
			if p0 == nil {
				return resource.ComposeNamespaceErrUnresolved(res.RefPath[0].Identifiers)
			}
			p0ID = p0.Slug
		}

		if res.RefResource != nil {
			p1 := resource.FindComposePage(state.ParentResources, res.RefResource.Identifiers)
			if p1 == nil {
				return resource.ComposePageErrUnresolved(res.RefResource.Identifiers)
			}
			p1ID = p1.Handle
		}

		r.res.Resource = fmt.Sprintf(composeTypes.PageRbacResourceTpl(), composeTypes.PageResourceType, p0ID, p1ID)
	case composeTypes.RecordResourceType:
		return fmt.Errorf("importing rbac rules on record level is not supported")
	case composeTypes.ModuleFieldResourceType:
		return fmt.Errorf("importing rbac rules on record level is not supported")
	case systemTypes.UserResourceType:
		if res.RefResource != nil {
			p1 := resource.FindUser(state.ParentResources, res.RefResource.Identifiers)
			if p1 == nil {
				return resource.UserErrUnresolved(res.RefResource.Identifiers)
			}
			p1ID = p1.Handle
		}

		r.res.Resource = fmt.Sprintf(systemTypes.UserRbacResourceTpl(), systemTypes.UserResourceType, p1ID)
	case systemTypes.RoleResourceType:
		if res.RefResource != nil {
			p1 := resource.FindRole(state.ParentResources, res.RefResource.Identifiers)
			if p1 == nil {
				return resource.RoleErrUnresolved(res.RefResource.Identifiers)
			}
			p1ID = p1.Handle
		}

		r.res.Resource = fmt.Sprintf(systemTypes.RoleRbacResourceTpl(), systemTypes.RoleResourceType, p1ID)
	case systemTypes.ApplicationResourceType:
		if res.RefResource != nil {
			p1 := resource.FindApplication(state.ParentResources, res.RefResource.Identifiers)
			if p1 == nil {
				return resource.ApplicationErrUnresolved(res.RefResource.Identifiers)
			}
			p1ID = p1.Name
		}

		r.res.Resource = fmt.Sprintf(systemTypes.ApplicationRbacResourceTpl(), systemTypes.ApplicationResourceType, p1ID)
	default:
		return fmt.Errorf("unsupported resource type '%s' for RBAC YAML encode", r.res.Resource)
	}

next:
	doc.addRbacRule(r)

	return nil
}

func (rr rbacRuleSet) MarshalYAML() (interface{}, error) {
	if rr == nil || len(rr) == 0 {
		return nil, nil
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
					opNode, err = addSeq(opNode, rule.res.Operation)

					if err != nil {
						return nil, err
					}
				}

				resNode, err = addMap(resNode,
					resRules[0].res.Resource, opNode,
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
	return r.res.Operation, nil
}
