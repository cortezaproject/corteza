package resource

import (
	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

const (
	RBAC_RESOURCE_TYPE = "rbacRule"
)

type (
	RbacRule struct {
		*base
		Res *rbac.Rule

		// Perhaps?
		RefRole     string
		RefResource string
	}
)

func NewRbacRule(res *rbac.Rule, refRole string, resRef *Ref) *RbacRule {
	r := &RbacRule{base: &base{}}
	r.SetResourceType(RBAC_RESOURCE_TYPE)
	r.Res = res
	r.RefRole = refRole

	r.AddRef(ROLE_RESOURCE_TYPE, refRole)
	r.AddRef(resRef.ResourceType, resRef.Identifiers.StringSlice()...)

	// @todo identifiers?
	// Combination of resID, operation, rule?

	return r
}

func (r *RbacRule) SearchQuery() rbac.RuleFilter {
	f := rbac.RuleFilter{}

	// @todo?

	return f
}
