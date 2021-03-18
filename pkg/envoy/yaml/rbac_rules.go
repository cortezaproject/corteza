package yaml

import (
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	rbacRule struct {
		res *rbac.Rule

		// To help us construct the resource
		resource rbac.Resource

		refResource string
		refRes      *resource.Ref
		relResource resource.Interface

		refRole string
		relRole *types.Role
	}
	rbacRuleSet []*rbacRule
)

func rbacRuleFromResource(r *resource.RbacRule, cfg *EncoderConfig) *rbacRule {
	rr := string(r.Res.Resource)
	rr = strings.TrimRight(rr, ":")

	return &rbacRule{
		res:         r.Res,
		refRes:      r.RefResource,
		refResource: rr,
		refRole:     r.RefRole.Identifiers.First(),
	}
}

func (rr rbacRuleSet) groupByAccess() []rbacRuleSet {
	rtr := make([]rbacRuleSet, 2)

	for _, r := range rr {
		if r.res.Access == rbac.Allow {
			rtr[0] = append(rtr[0], r)
		} else if r.res.Access == rbac.Deny {
			rtr[1] = append(rtr[1], r)
		}
	}

	return rtr
}

func (rr rbacRuleSet) groupByRole() []rbacRuleSet {
	rolx := make(map[string]rbacRuleSet)

	for _, r := range rr {
		if _, has := rolx[r.refRole]; !has {
			rolx[r.refRole] = make(rbacRuleSet, 0, 100)
		}

		rolx[r.refRole] = append(rolx[r.refRole], r)
	}

	rtr := make([]rbacRuleSet, 0, len(rolx))
	for _, rr := range rolx {
		rtr = append(rtr, rr)
	}

	return rtr
}

func (rr rbacRuleSet) groupByResource() []rbacRuleSet {
	rolx := make(map[string]rbacRuleSet)

	for _, r := range rr {
		k := r.res.Resource.String()
		if r.relResource != nil {
			if ri, is := r.relResource.(resource.RefableInterface); is {
				k += ri.Ref()
			}
		}

		if _, has := rolx[k]; !has {
			rolx[k] = make(rbacRuleSet, 0, 100)
		}

		rolx[k] = append(rolx[k], r)
	}

	rtr := make([]rbacRuleSet, 0, len(rolx))
	for _, rr := range rolx {
		rtr = append(rtr, rr)
	}

	return rtr
}
