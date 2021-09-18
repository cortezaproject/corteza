package yaml

import (
	"context"
	"fmt"
	"strings"

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
	n.role = resource.FindRole(state.ParentResources, rl.RefRole.Identifiers)
	if n.role == nil {
		return resource.RoleErrUnresolved(rl.RefRole.Identifiers)
	}

	refRes := n.refRbacRes
	if refRes == nil {
		return
	}

	// For now we will only allow resource specific RBAC rules if that resource is
	// also present.
	// Here we check if we can find it in case we're handling a resource specific rule.
	found := false
	switch refRes.ResourceType {
	case composeTypes.RecordResourceType:
		if found = n.handleComposeRecord(state.ParentResources); !found {
			return resource.RbacResourceErrNotFound(refRes.Identifiers)
		}

	case composeTypes.ModuleFieldResourceType:
		if found = n.handleComposeModuleField(state.ParentResources); !found {
			return resource.RbacResourceErrNotFound(refRes.Identifiers)
		}

	default:
		for _, r := range state.ParentResources {
			if refRes.ResourceType == r.ResourceType() && r.Identifiers().HasAny(refRes.Identifiers) {
				return
			}
		}
		// We couldn't find it...
		return resource.RbacResourceErrNotFound(refRes.Identifiers)
	}

	return nil
}

func (r *rbacRule) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
	// @todo Improve RBAC rule placement
	//
	// In cases where a specific rule is created for a specific resource, nest the rule
	// under the related namespace.
	// For now all rules will be nested under a root node for simplicity sake.

	refResource, err := r.makeRBACResource(state)
	if err != nil {
		return err
	}

	r.refRbacResource = refResource
	r.res.Resource = refResource

	doc.addRbacRule(r)

	return nil
}

func (r *rbacRule) makeRBACResource(state *envoy.ResourceState) (string, error) {
	res := state.Res.(*resource.RbacRule)

	p0ID := "*"
	p1ID := "*"
	p2ID := "*"

	rt := strings.Split(r.refRbacResource, "/")[0]

	switch rt {
	case composeTypes.ComponentResourceType,
		systemTypes.ComponentResourceType,
		automationTypes.ComponentResourceType,
		federationTypes.ComponentResourceType:

		return rt, nil
	}

	switch rt {
	case composeTypes.NamespaceResourceType:
		if res.RefRes != nil {
			p1 := resource.FindComposeNamespace(state.ParentResources, res.RefRes.Identifiers)
			if p1 == nil {
				return "", resource.ComposeNamespaceErrUnresolved(res.RefRes.Identifiers)
			}
			p1ID = p1.Slug
		}

		return fmt.Sprintf(composeTypes.NamespaceRbacResourceTpl(), composeTypes.NamespaceResourceType, p1ID), nil
	case composeTypes.ModuleResourceType:
		if len(res.RefPath) > 0 {
			p0 := resource.FindComposeNamespace(state.ParentResources, res.RefPath[0].Identifiers)
			if p0 == nil {
				return "", resource.ComposeNamespaceErrUnresolved(res.RefPath[0].Identifiers)
			}
			p0ID = p0.Slug
		}

		if res.RefRes != nil {
			p1 := resource.FindComposeModule(state.ParentResources, res.RefRes.Identifiers)
			if p1 == nil {
				return "", resource.ComposeModuleErrUnresolved(res.RefRes.Identifiers)
			}
			p1ID = p1.Handle
		}

		return fmt.Sprintf(composeTypes.ModuleRbacResourceTpl(), composeTypes.ModuleResourceType, p0ID, p1ID), nil
	case composeTypes.ChartResourceType:
		if len(res.RefPath) > 0 {
			p0 := resource.FindComposeNamespace(state.ParentResources, res.RefPath[0].Identifiers)
			if p0 == nil {
				return "", resource.ComposeNamespaceErrUnresolved(res.RefPath[0].Identifiers)
			}
			p0ID = p0.Slug
		}

		if res.RefRes != nil {
			p1 := resource.FindComposeChart(state.ParentResources, res.RefRes.Identifiers)
			if p1 == nil {
				return "", resource.ComposeChartErrUnresolved(res.RefRes.Identifiers)
			}
			p1ID = p1.Handle
		}

		return fmt.Sprintf(composeTypes.ChartRbacResourceTpl(), composeTypes.ChartResourceType, p0ID, p1ID), nil
	case composeTypes.PageResourceType:
		if len(res.RefPath) > 0 {
			p0 := resource.FindComposeNamespace(state.ParentResources, res.RefPath[0].Identifiers)
			if p0 == nil {
				return "", resource.ComposeNamespaceErrUnresolved(res.RefPath[0].Identifiers)
			}
			p0ID = p0.Slug
		}

		if res.RefRes != nil {
			p1 := resource.FindComposePage(state.ParentResources, res.RefRes.Identifiers)
			if p1 == nil {
				return "", resource.ComposePageErrUnresolved(res.RefRes.Identifiers)
			}
			p1ID = p1.Handle
		}

		return fmt.Sprintf(composeTypes.PageRbacResourceTpl(), composeTypes.PageResourceType, p0ID, p1ID), nil
	case composeTypes.RecordResourceType:
		if len(res.RefPath) > 0 {
			p0 := resource.FindComposeNamespace(state.ParentResources, res.RefPath[0].Identifiers)
			if p0 == nil {
				return "", resource.ComposeNamespaceErrUnresolved(res.RefPath[0].Identifiers)
			}
			p0ID = p0.Slug
		}
		if len(res.RefPath) > 1 {
			p1 := resource.FindComposeModule(state.ParentResources, res.RefPath[1].Identifiers)
			if p1 == nil {
				return "", resource.ComposeModuleErrUnresolved(res.RefPath[1].Identifiers)
			}
			p1ID = p1.Handle
		}

		if res.RefRes != nil {
			ref := res.RefPath[1]

			p2 := resource.FindComposeRecordResource(state.ParentResources, ref.Identifiers)
			if p2 == nil {
				return "", resource.ComposeRecordErrUnresolved(res.RefRes.Identifiers)
			}

			for i := range res.RefRes.Identifiers {
				if _, ok := p2.IDMap[i]; ok {
					p2ID = res.RefRes.Identifiers.First()
					break
				}
			}

			if p2ID == "" {
				return "", resource.ComposeRecordErrUnresolved(res.RefRes.Identifiers)
			}
		}

		return fmt.Sprintf(composeTypes.RecordRbacResourceTpl(), composeTypes.RecordResourceType, p0ID, p1ID, p2ID), nil

	case composeTypes.ModuleFieldResourceType:
		if len(res.RefPath) > 0 {
			p0 := resource.FindComposeNamespace(state.ParentResources, res.RefPath[0].Identifiers)
			if p0 == nil {
				return "", resource.ComposeNamespaceErrUnresolved(res.RefPath[0].Identifiers)
			}
			p0ID = p0.Slug
		}
		if len(res.RefPath) > 1 {
			p1 := resource.FindComposeModule(state.ParentResources, res.RefPath[1].Identifiers)
			if p1 == nil {
				return "", resource.ComposeModuleErrUnresolved(res.RefPath[1].Identifiers)
			}
			p1ID = p1.Handle
		}

		if r.refRbacRes != nil {
			modRef := r.refPathRes[1]
			p2 := resource.FindComposeModuleField(state.ParentResources, modRef.Identifiers, r.refRbacRes.Identifiers)
			if p2 == nil {
				return "", resource.ComposeModuleFieldErrUnresolved(r.refRbacRes.Identifiers)
			}

			p2ID = p2.Name
		}

		// @todo specific module field RBAC
		return fmt.Sprintf(composeTypes.ModuleFieldRbacResourceTpl(), composeTypes.ModuleFieldResourceType, p0ID, p1ID, p2ID), nil
	case systemTypes.UserResourceType:
		if res.RefRes != nil {
			p1 := resource.FindUser(state.ParentResources, res.RefRes.Identifiers)
			if p1 == nil {
				return "", resource.UserErrUnresolved(res.RefRes.Identifiers)
			}
			p1ID = p1.Handle
		}

		return fmt.Sprintf(systemTypes.UserRbacResourceTpl(), systemTypes.UserResourceType, p1ID), nil
	case systemTypes.RoleResourceType:
		if res.RefRes != nil {
			p1 := resource.FindRole(state.ParentResources, res.RefRes.Identifiers)
			if p1 == nil {
				return "", resource.RoleErrUnresolved(res.RefRes.Identifiers)
			}
			p1ID = p1.Handle
		}

		return fmt.Sprintf(systemTypes.RoleRbacResourceTpl(), systemTypes.RoleResourceType, p1ID), nil
	case systemTypes.ApplicationResourceType:
		if res.RefRes != nil {
			p1 := resource.FindApplication(state.ParentResources, res.RefRes.Identifiers)
			if p1 == nil {
				return "", resource.ApplicationErrUnresolved(res.RefRes.Identifiers)
			}
			p1ID = p1.Name
		}

		return fmt.Sprintf(systemTypes.ApplicationRbacResourceTpl(), systemTypes.ApplicationResourceType, p1ID), nil

		// // @todo
		// case systemTypes.ApigwRouteResourceType:
		// case systemTypes.ApigwFilterResourceType:
	}

	return "", fmt.Errorf("unsupported resource type '%s' for RBAC YAML encode", r.res.Resource)
}

func (n *rbacRule) handleComposeRecord(pp []resource.Interface) bool {
	for _, p := range pp {
		if p.ResourceType() == composeTypes.RecordResourceType && p.Identifiers().HasAny(n.refPathRes[1].Identifiers) {
			return true
		}
	}

	return false
}

func (n *rbacRule) handleComposeModuleField(pp []resource.Interface) bool {
	for _, p := range pp {
		if p.ResourceType() == composeTypes.ModuleResourceType && p.Identifiers().HasAny(n.refPathRes[1].Identifiers) {
			return true
		}
	}

	return false
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

			roleNode, err = addMap(roleNode,
				roleRules[0].getRoleKey(), resNode,
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
