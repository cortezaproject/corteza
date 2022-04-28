package resource

import (
	"fmt"

	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
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

func NewRbacRule(res *rbac.Rule, refRole, refRes *Ref, refResource string, refPath ...*Ref) *RbacRule {
	r := &RbacRule{base: &base{}}
	r.SetResourceType(RbacResourceType)
	r.Res = res

	r.RefRole = r.addRef(refRole)

	r.RefResource = refResource
	if refRes != nil {
		r.RefRes = r.addRef(refRes)
	}

	// any additional constraints
	for i, rp := range refPath {
		ref := MakeRef(rp.ResourceType, rp.Identifiers)

		if r.RefRes != nil {
			r.RefRes.Constraint(ref)
		}

		if i > 0 {
			for j := i - 1; j < i; j++ {
				aux := refPath[j]
				ref = ref.Constraint(MakeRef(aux.ResourceType, aux.Identifiers))
			}
		}

		r.RefPath = append(r.RefPath, r.addRef(ref))
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

func (r *RbacRule) Resource() interface{} {
	return r.Res
}

func (r *RbacRule) ReRef(old RefSet, new RefSet) {
	r.base.ReRef(old, new)

	// Care for wildcards
	if r.RefRes != nil {
		for i, o := range old {
			if o.equals(r.RefRes) {
				r.RefRes = new[i]
				break
			}
		}
	}

	for i, o := range old {
		if RefSet(r.RefPath).findRef(o) > -1 {
			r.RefPath = RefSet(r.RefPath).replaceRef(o, new[i])
		}
	}
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
