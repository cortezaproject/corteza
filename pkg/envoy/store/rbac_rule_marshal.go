package store

import (
	"context"
	"fmt"
	"strings"

	automationTypes "github.com/cortezaproject/corteza-server/automation/types"
	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	federationTypes "github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
)

var (
	gRbacRules map[string]bool
)

func newRbacRuleFromResource(res *resource.RbacRule, cfg *EncoderConfig) resourceState {
	return &rbacRule{
		cfg: mergeConfig(cfg, res.Config()),

		rule: res.Res,

		refRbacResource: res.RefResource,
		refRbacRes:      res.RefRes,

		refPathRes: res.RefPath,

		refRole: res.RefRole,
	}
}

func (n *rbacRule) Prepare(ctx context.Context, pl *payload) (err error) {
	// Init global state
	if gRbacRules == nil {
		err = n.initGlobalIndex(ctx, pl.s)
		if err != nil {
			return
		}
	}

	// Related role
	n.role, err = findRole(ctx, pl.s, pl.state.ParentResources, n.refRole.Identifiers)
	if err != nil {
		return err
	}
	if n.role == nil {
		return resource.RoleErrUnresolved(n.refRole.Identifiers)
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
		if found = n.handleComposeRecord(pl.state.ParentResources); !found {
			return resource.RbacResourceErrNotFound(refRes.Identifiers)
		}

	case composeTypes.ModuleFieldResourceType:
		if found = n.handleComposeModuleField(pl.state.ParentResources); !found {
			return resource.RbacResourceErrNotFound(refRes.Identifiers)
		}

	default:
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
	// Assure correct resource
	n.rule.Resource, err = n.makeRBACResource(pl)
	if err != nil {
		return err
	}

	// Assure correct role
	n.rule.RoleID = n.role.ID
	if n.rule.RoleID == 0 {
		rl := resource.FindRole(pl.state.ParentResources, n.refRole.Identifiers)
		n.rule.RoleID = rl.ID
	}
	if n.rule.RoleID == 0 {
		return resource.RoleErrUnresolved(n.refRole.Identifiers)
	}

	// Upsert
	if _, exists := gRbacRules[n.rbacRuleIndex((n.rule))]; !exists {
		return store.CreateRbacRule(ctx, pl.s, n.rule)
	}

	// On existing rbac rule, replace/merge right basically overwrites the existing rule;
	// otherwise, the new rule is ignored.
	switch n.cfg.OnExisting {
	case resource.Replace,
		resource.MergeRight:
		return store.UpdateRbacRule(ctx, pl.s, n.rule)
	}

	return nil
}

func (n *rbacRule) rbacRuleIndex(r *rbac.Rule) string {
	return fmt.Sprintf("%d_%s_%s", r.RoleID, r.Resource, r.Operation)
}

func (n *rbacRule) initGlobalIndex(ctx context.Context, s store.Storer) error {
	gRbacRules = make(map[string]bool)

	rr, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
	if err != store.ErrNotFound && err != nil {
		return err
	}

	for _, r := range rr {
		gRbacRules[n.rbacRuleIndex(r)] = true
	}
	return nil
}

func (n *rbacRule) makeRBACResource(pl *payload) (string, error) {
	p0ID := uint64(0)
	p1ID := uint64(0)
	p2ID := uint64(0)

	rt := strings.Split(n.refRbacResource, "/")[0]

	// Component level stuff
	switch rt {
	case composeTypes.ComponentResourceType:
		return composeTypes.ComponentRbacResource(), nil
	case systemTypes.ComponentResourceType:
		return systemTypes.ComponentRbacResource(), nil
	case automationTypes.ComponentResourceType:
		return automationTypes.ComponentRbacResource(), nil
	case federationTypes.ComponentResourceType:
		return federationTypes.ComponentRbacResource(), nil
	}

	// Component resource stuff
	switch rt {
	case automationTypes.WorkflowResourceType:
		if n.refRbacRes != nil {
			p1 := resource.FindAutomationWorkflow(pl.state.ParentResources, n.refRbacRes.Identifiers)
			if p1 == nil {
				return "", resource.AutomationWorkflowErrUnresolved(n.refRbacRes.Identifiers)
			}
			p1ID = p1.ID
		}
		return automationTypes.WorkflowRbacResource(p1ID), nil

	case composeTypes.NamespaceResourceType:
		if n.refRbacRes != nil {
			p1 := resource.FindComposeNamespace(pl.state.ParentResources, n.refRbacRes.Identifiers)
			if p1 == nil {
				return "", resource.ComposeNamespaceErrUnresolved(n.refRbacRes.Identifiers)
			}
			p1ID = p1.ID
		}
		return composeTypes.NamespaceRbacResource(p1ID), nil

	case composeTypes.ModuleResourceType:
		if len(n.refPathRes) > 0 {
			p0 := resource.FindComposeNamespace(pl.state.ParentResources, n.refPathRes[0].Identifiers)
			if p0 == nil {
				return "", resource.ComposeNamespaceErrUnresolved(n.refPathRes[0].Identifiers)
			}
			p0ID = p0.ID
		}

		if n.refRbacRes != nil {
			p1 := resource.FindComposeModule(pl.state.ParentResources, n.refRbacRes.Identifiers)
			if p1 == nil {
				return "", resource.ComposeModuleErrUnresolved(n.refRbacRes.Identifiers)
			}
			p1ID = p1.ID
		}
		return composeTypes.ModuleRbacResource(p0ID, p1ID), nil

	case composeTypes.ChartResourceType:
		if len(n.refPathRes) > 0 {
			p0 := resource.FindComposeNamespace(pl.state.ParentResources, n.refPathRes[0].Identifiers)
			if p0 == nil {
				return "", resource.ComposeNamespaceErrUnresolved(n.refPathRes[0].Identifiers)
			}
			p0ID = p0.ID
		}

		if n.refRbacRes != nil {
			p1 := resource.FindComposeChart(pl.state.ParentResources, n.refRbacRes.Identifiers)
			if p1 == nil {
				return "", resource.ComposeChartErrUnresolved(n.refRbacRes.Identifiers)
			}
			p1ID = p1.ID
		}
		return composeTypes.ChartRbacResource(p0ID, p1ID), nil

	case composeTypes.PageResourceType:
		if len(n.refPathRes) > 0 {
			p0 := resource.FindComposeNamespace(pl.state.ParentResources, n.refPathRes[0].Identifiers)
			if p0 == nil {
				return "", resource.ComposeNamespaceErrUnresolved(n.refPathRes[0].Identifiers)
			}
			p0ID = p0.ID
		}

		if n.refRbacRes != nil {
			p1 := resource.FindComposePage(pl.state.ParentResources, n.refRbacRes.Identifiers)
			if p1 == nil {
				return "", resource.ComposePageErrUnresolved(n.refRbacRes.Identifiers)
			}
			p1ID = p1.ID
		}
		return composeTypes.PageRbacResource(p0ID, p1ID), nil

	case composeTypes.RecordResourceType:
		if len(n.refPathRes) > 0 {
			p0 := resource.FindComposeNamespace(pl.state.ParentResources, n.refPathRes[0].Identifiers)
			if p0 == nil {
				return "", resource.ComposeNamespaceErrUnresolved(n.refPathRes[0].Identifiers)
			}
			p0ID = p0.ID
		}
		if len(n.refPathRes) > 1 {
			p1 := resource.FindComposeModule(pl.state.ParentResources, n.refPathRes[1].Identifiers)
			if p1 == nil {
				return "", resource.ComposeModuleErrUnresolved(n.refPathRes[1].Identifiers)
			}
			p1ID = p1.ID
		}
		if n.refRbacRes != nil {
			ref := n.refPathRes[1]

			p2 := resource.FindComposeRecordResource(pl.state.ParentResources, ref.Identifiers)
			if p2 == nil {
				return "", resource.ComposeRecordErrUnresolved(n.refRbacRes.Identifiers)
			}

			for _, i := range n.refRbacRes.Identifiers {
				if p2ID = p2.IDMap[i]; p2ID > 0 {
					break
				}
			}
			if p2ID == 0 {
				return "", resource.ComposeRecordErrUnresolved(n.refRbacRes.Identifiers)
			}
		}
		return composeTypes.RecordRbacResource(p0ID, p1ID, p2ID), nil

	case composeTypes.ModuleFieldResourceType:
		if len(n.refPathRes) > 0 {
			p0 := resource.FindComposeNamespace(pl.state.ParentResources, n.refPathRes[0].Identifiers)
			if p0 == nil {
				return "", resource.ComposeNamespaceErrUnresolved(n.refPathRes[0].Identifiers)
			}
			p0ID = p0.ID
		}
		if len(n.refPathRes) > 1 {
			p1 := resource.FindComposeModule(pl.state.ParentResources, n.refPathRes[1].Identifiers)
			if p1 == nil {
				return "", resource.ComposeModuleErrUnresolved(n.refPathRes[1].Identifiers)
			}
			p1ID = p1.ID
		}

		if n.refRbacRes != nil {
			modRef := n.refPathRes[1]
			p2 := resource.FindComposeModuleField(pl.state.ParentResources, modRef.Identifiers, n.refRbacRes.Identifiers)
			if p2 == nil {
				return "", resource.ComposeModuleFieldErrUnresolved(n.refRbacRes.Identifiers)
			}

			p2ID = p2.ID
		}

		return composeTypes.ModuleFieldRbacResource(p0ID, p1ID, p2ID), nil

	case systemTypes.UserResourceType:
		if n.refRbacRes != nil {
			p1 := resource.FindUser(pl.state.ParentResources, n.refRbacRes.Identifiers)
			if p1 == nil {
				return "", resource.UserErrUnresolved(n.refRbacRes.Identifiers)
			}
			p1ID = p1.ID
		}
		return systemTypes.UserRbacResource(p1ID), nil

	case systemTypes.RoleResourceType:
		if n.refRbacRes != nil {
			p1 := resource.FindRole(pl.state.ParentResources, n.refRbacRes.Identifiers)
			if p1 == nil {
				return "", resource.RoleErrUnresolved(n.refRbacRes.Identifiers)
			}
			p1ID = p1.ID
		}
		return systemTypes.RoleRbacResource(p1ID), nil

	case systemTypes.ApplicationResourceType:
		if n.refRbacRes != nil {
			p1 := resource.FindApplication(pl.state.ParentResources, n.refRbacRes.Identifiers)
			if p1 == nil {
				return "", resource.ApplicationErrUnresolved(n.refRbacRes.Identifiers)
			}
			p1ID = p1.ID
		}
		return systemTypes.ApplicationRbacResource(p1ID), nil

	case systemTypes.ApigwRouteResourceType:
		if n.refRbacRes != nil {
			p1 := resource.FindAPIGateway(pl.state.ParentResources, n.refRbacRes.Identifiers)
			if p1 == nil {
				return "", resource.APIGatewayErrUnresolved(n.refRbacRes.Identifiers)
			}
			p1ID = p1.ID
		}
		return systemTypes.ApigwRouteRbacResource(p1ID), nil

	case systemTypes.AuthClientResourceType:
		// @todo add support for importing rbac rules for specific client
		return systemTypes.AuthClientRbacResource(p1ID), nil

	case systemTypes.TemplateResourceType:
		if n.refRbacRes != nil {
			p1 := resource.FindTemplate(pl.state.ParentResources, n.refRbacRes.Identifiers)
			if p1 == nil {
				return "", resource.TemplateErrUnresolved(n.refRbacRes.Identifiers)
			}
			p1ID = p1.ID
		}
		return systemTypes.TemplateRbacResource(p1ID), nil

	case systemTypes.ReportResourceType:
		if n.refRbacRes != nil {
			p1 := resource.FindReport(pl.state.ParentResources, n.refRbacRes.Identifiers)
			if p1 == nil {
				return "", resource.ReportErrUnresolved(n.refRbacRes.Identifiers)
			}
			p1ID = p1.ID
		}
		return systemTypes.ReportRbacResource(p1ID), nil

	case systemTypes.QueueResourceType:
		// @todo add support for importing rbac rules for specific queue
		return systemTypes.QueueRbacResource(p1ID), nil

	case systemTypes.DataPrivacyRequestResourceType:
		return systemTypes.DataPrivacyRequestRbacResource(p1ID), nil

	case federationTypes.NodeResourceType:
		return federationTypes.NodeRbacResource(p1ID), nil
	case federationTypes.SharedModuleResourceType:
		return federationTypes.SharedModuleRbacResource(p0ID, p1ID), nil
	case federationTypes.ExposedModuleResourceType:
		return federationTypes.ExposedModuleRbacResource(p0ID, p1ID), nil
	}

	// @todo if we wish to support rbac for external stuff, this needs to pass through.
	//       this also requires some tweaks in the path ID thing.
	return "", fmt.Errorf("unsupported resource type '%s' for RBAC store encode", n.rule.Resource)
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
