package store

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
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

// Prepare prepares the rbacRule to be encoded
//
// Any validation, additional constraining should be performed here.
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
	n.relRole, err = findRoleRS(ctx, pl.s, pl.state.ParentResources, n.res.RefRole.Identifiers)
	if err != nil {
		return err
	}
	if n.relRole == nil {
		return resource.RoleErrUnresolved(n.res.RefRole.Identifiers)
	}

	// For now we will only allow resource specific RBAC rules if that resource is
	// also present.
	// Here we check if we can find it in case we're handling a resource specific rule.

	// For now we will only allow resource specific RBAC rules if that resource is
	// also present.
	// Here we check if we can find it in case we're handling a resource specific rule.
	refRes := n.res.RefResource
	if refRes != nil && len(refRes.Identifiers) > 0 {
		for _, r := range pl.state.ParentResources {
			if refRes.ResourceType == r.ResourceType() && r.Identifiers().HasAny(refRes.Identifiers) {
				goto pass
			}
		}

		// We couldn't find it...
		return resource.RbacResourceErrNotFound(refRes.Identifiers)
	}
pass:

	return nil
}

// Encode encodes the rbacRule to the store
//
// Encode is allowed to do some data manipulation, but no resource constraints
// should be changed.
func (n *rbacRule) Encode(ctx context.Context, pl *payload) (err error) {
	res := n.res.Res

	// Related role
	res.RoleID = n.relRole.ID
	if res.RoleID <= 0 {
		rl := resource.FindRole(pl.state.ParentResources, n.res.RefRole.Identifiers)
		res.RoleID = rl.ID
	}
	if res.RoleID == 0 {
		return resource.RoleErrUnresolved(n.res.RefRole.Identifiers)
	}

	// Related resource
	refRes := n.res.RefResource
	if refRes != nil && len(refRes.Identifiers) > 0 {
		var relRes resource.Interface
		for _, r := range pl.state.ParentResources {
			if n.res.RefResource.ResourceType == r.ResourceType() && r.Identifiers().HasAny(n.res.RefResource.Identifiers) {
				relRes = r
				break
			}
		}
		relResI, ok := relRes.(resource.IdentifiableInterface)
		if !ok {
			return rbacResourceErrUnidentifiable(relRes.Identifiers())
		}
		res.Resource = res.Resource.AppendID(relResI.SysID())
	} else if res.Resource.IsAppendable() {
		res.Resource = res.Resource.AppendWildcard()
	}

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
