package yaml

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	resourceTranslation struct {
		// set of locales
		locales types.ResourceTranslationSet

		// point to the locale resource
		refResourceTranslation string
		refLocaleRes           *resource.Ref

		// PathRes and PathResource slices hold parent resources we should nest the rule by
		refPathRes []*resource.Ref
	}

	resourceTranslationSet []*resourceTranslation
)

func (ll resourceTranslationSet) groupByResourceTranslation() (out resourceTranslationSet) {
	resMap := make(map[string]*resourceTranslation)

	for _, _l := range ll {
		l := _l

		if _, ok := resMap[l.refResourceTranslation]; !ok {
			resMap[l.refResourceTranslation] = l
			continue
		}

		resMap[l.refResourceTranslation].locales = append(resMap[l.refResourceTranslation].locales, l.locales...)
	}

	for _, l := range resMap {
		out = append(out, l)
	}

	return
}

func resourceTranslationFromResource(r *resource.ResourceTranslation, cfg *EncoderConfig) *resourceTranslation {
	if len(r.Res) == 0 {
		return nil
	}

	return &resourceTranslation{
		locales: r.Res,

		refResourceTranslation: r.Res[0].Resource,
		refLocaleRes:           r.RefRes,

		refPathRes: r.RefPath,
	}
}
