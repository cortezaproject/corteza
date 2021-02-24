package resource

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

type (
	RbacRule struct {
		*base
		Res *rbac.Rule

		RefRole     *Ref
		RefResource *Ref
	}
)

func NewRbacRule(res *rbac.Rule, refRole string, resRef *Ref) *RbacRule {
	r := &RbacRule{base: &base{}}
	r.SetResourceType(RBAC_RESOURCE_TYPE)
	r.Res = res

	r.RefRole = r.AddRef(ROLE_RESOURCE_TYPE, refRole)

	if resRef != nil {
		r.RefResource = r.AddRef(resRef.ResourceType, resRef.Identifiers.StringSlice()...)
	}

	return r
}

func RbacResourceErrNotFound(ii Identifiers) error {
	return fmt.Errorf("rbac resource not found %v", ii.StringSlice())
}
