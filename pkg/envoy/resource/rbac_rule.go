package resource

import (
	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

const (
	RBAC_RESOURCE_TYPE = "rbacRule"
)

type (
	rbacRule struct {
		*base
		Res *rbac.Rule

		// Perhaps?
		RefRole string
	}
)

func RbacRule(res *rbac.Rule) *rbacRule {
	r := &rbacRule{base: &base{}}
	r.SetResourceType(RBAC_RESOURCE_TYPE)
	r.Res = res

	// @todo identifiers?
	// Combination of resID, operation, rule?

	return r
}

func (r *rbacRule) SearchQuery() rbac.RuleFilter {
	f := rbac.RuleFilter{}

	// @todo?

	return f
}
