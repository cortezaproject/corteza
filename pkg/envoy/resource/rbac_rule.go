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

		RefResource string
		RefRes      *Ref
		RefRole     *Ref
		RefPath     []*Ref
	}
)

func NewRbacRule(res *rbac.Rule, refRole string, refRes *Ref, refResource string, refPath ...*Ref) *RbacRule {
	r := &RbacRule{base: &base{}}
	r.SetResourceType(RbacResourceType)
	r.Res = res

	r.RefRole = r.AddRef(types.RoleResourceType, refRole)

	r.RefResource = refResource
	if refRes != nil {
		r.RefRes = r.AddRef(refRes.ResourceType, refRes.Identifiers.StringSlice()...)
	}

	// any additional constraints
	for _, rp := range refPath {
		r.RefPath = append(r.RefPath, r.AddRef(rp.ResourceType, rp.Identifiers.StringSlice()...))
	}

	// Handle cases for nested resources
	if refRes != nil {
		switch refRes.ResourceType {
		case composeTypes.RecordResourceType:
			r.handleComposeRecord(res, refPath)
		case composeTypes.ModuleFieldResourceType:
			r.handleComposeModuleField(res, refPath)
		}
	}

	return r
}

func (r *RbacRule) handleComposeRecord(res *rbac.Rule, refPath []*Ref) {

	// records are grouped under module

	if len(refPath) < 2 {
		return
	}

	r.AddRef(composeTypes.RecordResourceType, refPath[1].Identifiers.StringSlice()...)
}

func (r *RbacRule) handleComposeModuleField(res *rbac.Rule, refPath []*Ref) {

	// module fields are grouped under module

	if len(refPath) < 2 {
		return
	}

	r.AddRef(composeTypes.ModuleFieldResourceTranslationType, refPath[1].Identifiers.StringSlice()...)
}

func RbacResourceErrNotFound(ii Identifiers) error {
	return fmt.Errorf("rbac resource not found %v", ii.StringSlice())
}
