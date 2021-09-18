package resource

import (
	"fmt"

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
	for _, rp := range refPath {
		r.RefPath = append(r.RefPath, r.AddRef(rp.ResourceType, rp.Identifiers.StringSlice()...))
	}

	return r
}

func (l *ResourceTranslation) MarkDefault() {
	l.Priority = 1
}

func (l *ResourceTranslation) MarkGeneric() {
	l.Priority = 0
}

func ResourceTranslationErrNotFound(ii Identifiers) error {
	return fmt.Errorf("resource translation not found %v", ii.StringSlice())
}
