package store

import (
	"context"
	"fmt"

	automationTypes "github.com/cortezaproject/corteza-server/automation/types"
	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	federationTypes "github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
)

var (
	gRbacRules map[string]bool

	rbacRuleIndex = func(r *rbac.Rule) string {
		return fmt.Sprintf("%d_%s_%s", r.RoleID, r.Resource, r.Operation)
	}
)

func newRbacRuleFromResource(res *resource.RbacRule, cfg *EncoderConfig) resourceState {
	return &rbacRule{
		cfg: mergeConfig(cfg, res.Config()),

		res:  res,
		rule: res.Res,
	}
}

func (n *rbacRule) Prepare(ctx context.Context, pl *payload) (err error) {
	// Init global state
	if gRbacRules == nil {
		gRbacRules = make(map[string]bool)
		rr, _, err := store.SearchRbacRules(ctx, pl.s, rbac.RuleFilter{})
		if err != store.ErrNotFound && err != nil {
			return err
		}
		for _, r := range rr {
			gRbacRules[rbacRuleIndex(r)] = true
		}
	}

	// Related role
	n.relRole, err = findRole(ctx, pl.s, pl.state.ParentResources, n.res.RefRole.Identifiers)
	if err != nil {
		return err
	}
	if n.relRole == nil {
		return resource.RoleErrUnresolved(n.res.RefRole.Identifiers)
	}

	// For now we will only allow resource specific RBAC rules if that resource is
	// also present.
	// Here we check if we can find it in case we're handling a resource specific rule.
	refRes := n.res.RefResource
	if refRes != nil && len(refRes.Identifiers) > 0 {
		for _, r := range pl.state.ParentResources {
			if refRes.ResourceType == r.ResourceType() && r.Identifiers().HasAny(refRes.Identifiers) {
				return
			}
		}
		// We couldn't find it...
		return resource.RbacResourceErrNotFound(refRes.Identifiers)
	}

	return
}

func (n *rbacRule) Encode(ctx context.Context, pl *payload) (err error) {
	res := n.res.Res

	// Related role
	res.RoleID = n.relRole.ID
	if res.RoleID == 0 {
		rl := resource.FindRole(pl.state.ParentResources, n.res.RefRole.Identifiers)
		res.RoleID = rl.ID
	}
	if res.RoleID == 0 {
		return resource.RoleErrUnresolved(n.res.RefRole.Identifiers)
	}

	p0ID := uint64(0)
	p1ID := uint64(0)
	p2ID := uint64(0)

	switch n.rule.Resource {
	case composeTypes.ComponentResourceType:
		res.Resource = composeTypes.ComponentRbacResource()
		goto store
	case systemTypes.ComponentResourceType:
		res.Resource = systemTypes.ComponentRbacResource()
		goto store
	case automationTypes.ComponentResourceType:
		res.Resource = automationTypes.ComponentRbacResource()
		goto store
	case federationTypes.ComponentResourceType:
		res.Resource = federationTypes.ComponentRbacResource()
		goto store
	}

	switch n.rule.Resource {
	case automationTypes.WorkflowResourceType:
		if n.res.RefResource != nil {
			p1 := resource.FindAutomationWorkflow(pl.state.ParentResources, n.res.RefResource.Identifiers)
			if p1 == nil {
				return resource.AutomationWorkflowErrUnresolved(n.res.RefResource.Identifiers)
			}
			p1ID = p1.ID
		}

		res.Resource = automationTypes.WorkflowRbacResource(p1ID)

	case composeTypes.NamespaceResourceType:
		if n.res.RefResource != nil {
			p1 := resource.FindComposeNamespace(pl.state.ParentResources, n.res.RefResource.Identifiers)
			if p1 == nil {
				return resource.ComposeNamespaceErrUnresolved(n.res.RefResource.Identifiers)
			}
			p1ID = p1.ID
		}

		res.Resource = composeTypes.NamespaceRbacResource(p1ID)
	case composeTypes.ModuleResourceType:
		if len(n.res.RefPath) > 0 {
			p0 := resource.FindComposeNamespace(pl.state.ParentResources, n.res.RefPath[0].Identifiers)
			if p0 == nil {
				return resource.ComposeNamespaceErrUnresolved(n.res.RefPath[0].Identifiers)
			}
			p0ID = p0.ID
		}

		if n.res.RefResource != nil {
			p1 := resource.FindComposeModule(pl.state.ParentResources, n.res.RefResource.Identifiers)
			if p1 == nil {
				return resource.ComposeModuleErrUnresolved(n.res.RefResource.Identifiers)
			}
			p1ID = p1.ID
		}

		res.Resource = composeTypes.ModuleRbacResource(p0ID, p1ID)
	case composeTypes.ChartResourceType:
		if len(n.res.RefPath) > 0 {
			p0 := resource.FindComposeNamespace(pl.state.ParentResources, n.res.RefPath[0].Identifiers)
			if p0 == nil {
				return resource.ComposeNamespaceErrUnresolved(n.res.RefPath[0].Identifiers)
			}
			p0ID = p0.ID
		}

		if n.res.RefResource != nil {
			p1 := resource.FindComposeChart(pl.state.ParentResources, n.res.RefResource.Identifiers)
			if p1 == nil {
				return resource.ComposeChartErrUnresolved(n.res.RefResource.Identifiers)
			}
			p1ID = p1.ID
		}

		res.Resource = composeTypes.ChartRbacResource(p0ID, p1ID)
	case composeTypes.PageResourceType:
		if len(n.res.RefPath) > 0 {
			p0 := resource.FindComposeNamespace(pl.state.ParentResources, n.res.RefPath[0].Identifiers)
			if p0 == nil {
				return resource.ComposeNamespaceErrUnresolved(n.res.RefPath[0].Identifiers)
			}
			p0ID = p0.ID
		}

		if n.res.RefResource != nil {
			p1 := resource.FindComposePage(pl.state.ParentResources, n.res.RefResource.Identifiers)
			if p1 == nil {
				return resource.ComposePageErrUnresolved(n.res.RefResource.Identifiers)
			}
			p1ID = p1.ID
		}

		res.Resource = composeTypes.PageRbacResource(p0ID, p1ID)
	case composeTypes.RecordResourceType:
		if len(n.res.RefPath) > 0 {
			p0 := resource.FindComposeNamespace(pl.state.ParentResources, n.res.RefPath[0].Identifiers)
			if p0 == nil {
				return resource.ComposeNamespaceErrUnresolved(n.res.RefPath[0].Identifiers)
			}
			p0ID = p0.ID
		}
		if len(n.res.RefPath) > 1 {
			p1 := resource.FindComposeModule(pl.state.ParentResources, n.res.RefPath[1].Identifiers)
			if p1 == nil {
				return resource.ComposeModuleErrUnresolved(n.res.RefPath[1].Identifiers)
			}
			p1ID = p1.ID
		}

		// @todo specific record RBAC

		res.Resource = composeTypes.RecordRbacResource(p0ID, p1ID, p2ID)

	case composeTypes.ModuleFieldResourceType:
		if len(n.res.RefPath) > 0 {
			p0 := resource.FindComposeNamespace(pl.state.ParentResources, n.res.RefPath[0].Identifiers)
			if p0 == nil {
				return resource.ComposeNamespaceErrUnresolved(n.res.RefPath[0].Identifiers)
			}
			p0ID = p0.ID
		}
		if len(n.res.RefPath) > 1 {
			p1 := resource.FindComposeModule(pl.state.ParentResources, n.res.RefPath[1].Identifiers)
			if p1 == nil {
				return resource.ComposeModuleErrUnresolved(n.res.RefPath[1].Identifiers)
			}
			p1ID = p1.ID
		}

		// @todo specific ModuleField RBAC

		res.Resource = composeTypes.ModuleFieldRbacResource(p0ID, p1ID, p2ID)

	case systemTypes.UserResourceType:
		if n.res.RefResource != nil {
			p1 := resource.FindUser(pl.state.ParentResources, n.res.RefResource.Identifiers)
			if p1 == nil {
				return resource.UserErrUnresolved(n.res.RefResource.Identifiers)
			}
			p1ID = p1.ID
		}

		res.Resource = systemTypes.UserRbacResource(p1ID)
	case systemTypes.RoleResourceType:
		if n.res.RefResource != nil {
			p1 := resource.FindRole(pl.state.ParentResources, n.res.RefResource.Identifiers)
			if p1 == nil {
				return resource.RoleErrUnresolved(n.res.RefResource.Identifiers)
			}
			p1ID = p1.ID
		}

		res.Resource = systemTypes.RoleRbacResource(p1ID)
	case systemTypes.ApplicationResourceType:
		if n.res.RefResource != nil {
			p1 := resource.FindApplication(pl.state.ParentResources, n.res.RefResource.Identifiers)
			if p1 == nil {
				return resource.ApplicationErrUnresolved(n.res.RefResource.Identifiers)
			}
			p1ID = p1.ID
		}

		res.Resource = systemTypes.ApplicationRbacResource(p1ID)
	case systemTypes.ApigwRouteResourceType:
		res.Resource = systemTypes.ApigwRouteRbacResource(p1ID)
	case systemTypes.ApigwFilterResourceType:
		res.Resource = systemTypes.ApigwFilterRbacResource(p1ID)
	case systemTypes.AuthClientResourceType:
		// @todo add support for importing rbac rules for specific client
		res.Resource = systemTypes.AuthClientRbacResource(p1ID)
	case systemTypes.TemplateResourceType:
		// @todo add support for importing rbac rules for specific template
		res.Resource = systemTypes.TemplateRbacResource(p1ID)
	case systemTypes.ReportResourceType:
		res.Resource = systemTypes.ReportRbacResource(p1ID)
	case messagebus.QueueResourceType:
		// @todo add support for importing rbac rules for specific queue
		res.Resource = messagebus.QueueRbacResource(p1ID)

	case federationTypes.NodeResourceType:
		// @todo add support for importing rbac rules for specific queue
		res.Resource = federationTypes.NodeRbacResource(p1ID)

	case federationTypes.SharedModuleResourceType:
		// @todo add support for importing rbac rules for specific queue
		res.Resource = federationTypes.SharedModuleRbacResource(p0ID, p1ID)

	case federationTypes.ExposedModuleResourceType:
		// @todo add support for importing rbac rules for specific queue
		res.Resource = federationTypes.ExposedModuleRbacResource(p0ID, p1ID)

	default:
		// @todo if we wish to support rbac for external stuff, this needs to pass through.
		//       this also requires some tweaks in the path ID thing.
		return fmt.Errorf("unsupported resource type '%s' for RBAC store encode", n.rule.Resource)
	}

store:

	if _, exists := gRbacRules[rbacRuleIndex(res)]; !exists {
		return store.CreateRbacRule(ctx, pl.s, res)
	}

	// On existing rbac rule, replace/merge right basically overwrites the existing rule;
	// otherwise, the new rule is ignored.
	switch n.cfg.OnExisting {
	case resource.Replace,
		resource.MergeRight:
		return store.UpdateRbacRule(ctx, pl.s, res)
	}

	return nil
}
