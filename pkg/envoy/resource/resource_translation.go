package resource

import (
	"fmt"

	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/system/types"
	"golang.org/x/text/language"
)

type (
	ResourceTranslation struct {
		*base
		Res types.ResourceTranslationSet

		RefResource string
		RefRes      *Ref
		RefPath     []*Ref

		Priority int
	}
)

var (
	defaultLanguage = language.English
)

func NewResourceTranslation(res types.ResourceTranslationSet, refResource string, refRes *Ref, refPath ...*Ref) *ResourceTranslation {
	r := &ResourceTranslation{base: &base{}}
	r.SetResourceType(ResourceTranslationType)
	r.Res = res

	r.RefResource = refResource
	r.RefRes = r.AddRef(refRes.ResourceType, refRes.Identifiers.StringSlice()...)

	// any additional constraints
	for i, rp := range refPath {
		ref := MakeRef(rp.ResourceType, rp.Identifiers)

		// @todo generalize when needed; for now only module field resource translations require this
		if i == 1 && refRes.ResourceType == composeTypes.ModuleFieldResourceType {
			ref = ref.Constraint(r.RefPath[0])
		}

		r.RefPath = append(r.RefPath, r.addRef(ref))
	}

	return r
}

func (r *ResourceTranslation) Resource() interface{} {
	return r.Res
}

func (r *ResourceTranslation) ReRef(old RefSet, new RefSet) {
	r.base.ReRef(old, new)

	for i, o := range old {
		if o.equals(r.RefRes) {
			r.RefRes = new[i]
			break
		}
	}

	for i, o := range old {
		if RefSet(r.RefPath).findRef(o) > -1 {
			r.RefPath = RefSet(r.RefPath).replaceRef(o, new[i])
		}
	}
}

func (l *ResourceTranslation) MarkDefault() {
	l.Priority = 1
}

func (l *ResourceTranslation) IsDefault() bool {
	return l.Priority == 1
}

func (l *ResourceTranslation) MarkGeneric() {
	l.Priority = 0
}

func (l *ResourceTranslation) IsGeneric() bool {
	return l.Priority == 0
}

func ResourceTranslationErrNotFound(ii Identifiers) error {
	return fmt.Errorf("resource translation not found %v", ii.StringSlice())
}
