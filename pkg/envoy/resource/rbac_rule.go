package resource

import (
	"fmt"

	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	RbacRule struct {
		*base
		Res *rbac.Rule

		RefRole     *Ref
		RefResource *Ref
		RefPath     []*Ref
	}
)

func NewRbacRule(res *rbac.Rule, refRole string, resRef *Ref, refPath ...*Ref) *RbacRule {
	r := &RbacRule{base: &base{}}
	r.SetResourceType(RbacResourceType)
	r.Res = res

	r.RefRole = r.AddRef(types.RoleResourceType, refRole)

	if resRef != nil {
		r.RefResource = r.AddRef(resRef.ResourceType, resRef.Identifiers.StringSlice()...)
	}

	// any additional constraints
	for _, rp := range refPath {
		r.RefPath = append(r.RefPath, r.AddRef(rp.ResourceType, rp.Identifiers.StringSlice()...))
	}

	// ComposeRecords are internally grouped and identified with the module identifier.
	if res.Resource == composeTypes.RecordResourceType && resRef != nil && len(refPath) == 2 {
		r.AddRef(res.Resource, refPath[1].Identifiers.StringSlice()...)
	}

	return r
}

func RbacResourceErrNotFound(ii Identifiers) error {
	return fmt.Errorf("rbac resource not found %v", ii.StringSlice())
}
