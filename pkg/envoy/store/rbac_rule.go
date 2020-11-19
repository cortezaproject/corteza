package store

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	rbacRuleState struct {
		cfg *EncoderConfig

		res  *resource.RbacRule
		rule *rbac.Rule

		relRole     *types.Role
		relResource resource.Interface
	}
)

var (
	// Use this as a cache, so we don't have to constantly fetch them
	rbacRules rbac.RuleSet = nil
)

func NewRbacRuleState(res *resource.RbacRule, cfg *EncoderConfig) resourceState {
	return &rbacRuleState{
		cfg: cfg,

		res: res,
	}
}

func (n *rbacRuleState) Prepare(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
	// Preload rbac rules if needed
	if rbacRules == nil {
		rbacRules, _, err = store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
		if err != store.ErrNotFound && err != nil {
			return err
		}
	}

	// Check for role
	n.relRole, err = findRoleRS(ctx, s, state.ParentResources, n.res.RefRole.Identifiers)
	if err != nil {
		return err
	}
	if n.relRole == nil {
		return roleErrUnresolved(n.res.RefRole.Identifiers)
	}

	// Check for resource if there is any
	if n.res.RefResource != nil {
		n.relResource = n.findResourceR(ctx, state.ParentResources, n.res.RefResource.ResourceType, n.res.RefResource.Identifiers)
		if n.relResource == nil {
			// Try to find it in the store
			// @todo...
		}
	}

	return nil
}

func (n *rbacRuleState) Encode(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
	rule := n.res.Res

	// Role
	rule.RoleID = n.relRole.ID
	if rule.RoleID <= 0 {
		rl := findRoleR(state.ParentResources, n.res.RefRole.Identifiers)
		rule.RoleID = rl.ID
	}

	// Resource?
	if n.res.RefResource != nil {
		if ir, ok := n.relResource.(resource.IdentifiableInterface); !ok {
			return rbacResourceErrUnidentifiable(n.relResource.Identifiers())
		} else {
			rule.Resource = rule.Resource.AppendID(ir.SysID())
		}
	} else if rule.Resource.IsAppendable() {
		rule.Resource = rule.Resource.AppendWildcard()
	}

	ee, _ := rbacRules.Filter(func(r *rbac.Rule) (bool, error) {
		return (r.RoleID == rule.RoleID && r.Resource == rule.Resource && r.Operation == rule.Operation), nil
	})

	if len(ee) <= 0 {
		rbacRules = rbacRules.Merge(rule)
	} else {
		// There isn't anything to merge really, so Skip & MergeLeft skip it;
		// Replace & MergeRight replace it.
		switch n.cfg.OnExisting {
		case Skip,
			MergeLeft:
			return nil
		}

		rbacRules = rbacRules.Merge(rule)
	}

	d, u := rbacRules.Dirty()
	err = store.DeleteRbacRule(ctx, s, d...)
	if err != nil {
		return
	}

	err = store.UpsertRbacRule(ctx, s, u...)
	if err != nil {
		return
	}

	rbacRules.Clear()
	return nil
}

func (n *rbacRuleState) findResourceR(ctx context.Context, rr resource.InterfaceSet, rt string, ii resource.Identifiers) (rtr resource.Interface) {
	for _, r := range rr {
		if r.ResourceType() != rt {
			continue
		}

		if !r.Identifiers().HasAny(ii) {
			continue
		}

		return r
	}
	return nil
}

func rbacResourceErrUnidentifiable(ii resource.Identifiers) error {
	return fmt.Errorf("rbac resource unidentifiable %v", ii.StringSlice())
}
