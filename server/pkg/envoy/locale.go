package envoy

import (
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	resutil "github.com/cortezaproject/corteza-server/pkg/resource"
)

// NormalizeResourceTranslations takes the provided resource.ResourceTranslation
// and merges duplicates based on the Priority parameter
func NormalizeResourceTranslations(rr ...resource.Interface) []resource.Interface {
	out := make([]resource.Interface, 0, len(rr))
	locales := make([]*resource.ResourceTranslation, 0, len(rr))

	// - collect all locale resources
	for _, r := range rr {
		if l, ok := r.(*resource.ResourceTranslation); ok {
			locales = append(locales, l)
		} else {
			out = append(out, r)
		}
	}

	// make an index
	var byResource map[string][2]*resource.ResourceTranslation
	byResource = make(map[string][2]*resource.ResourceTranslation)

	for _, locale := range locales {
		pp := byResource[locale.RefResource]
		pp[locale.Priority] = locale
		byResource[locale.RefResource] = pp
	}

	// squash index by priority ascending
	for _, pp := range byResource {
		var aux *resource.ResourceTranslation
		seen := make(map[string]bool)

		for _, p := range pp {
			if p == nil {
				continue
			}

			if aux == nil {
				aux = p
			} else {
				for _, r := range p.Res {
					if seen[r.Lang.String()+r.K] {
						continue
					}

					aux.Res = append(aux.Res, r)
				}
			}

			for _, r := range p.Res {
				if r.K != "" {
					seen[r.Lang.String()+r.K] = true
				}
			}
		}

		if aux != nil {
			out = append(out, aux)
		}
	}

	return out
}

func appendRefSet(a resource.RefSet, b *resource.Ref) resource.RefSet {
	return append(a, b)
}

// FilterRequiredResourceTranslations returns only resource translations relevant for the given resources
func FilterRequiredResourceTranslations(request resource.InterfaceSet, translations []*resource.ResourceTranslation) (out []*resource.ResourceTranslation) {
	out = make([]*resource.ResourceTranslation, 0, 10)

	rtrIndex := resutil.NewIndex()
	for _, r := range translations {
		rtrIndex.Add(r, r.IndexPath()...)
	}

	// Filter
	procResSet(request, func(r resource.Interface) {
		if r.Placeholder() {
			return
		}

		rtrRes, ok := r.(resource.LocaleInterface)
		if !ok {
			return
		}

		_, ref, pp := rtrRes.ResourceTranslationParts()
		// @todo this string replace should eventually be improved; we didn't prefix
		// resource translation's resources with corteza::
		resPath := refsToIndexPath(strings.Replace(r.ResourceType(), "corteza::", "", -1), appendRefSet(pp, ref)...)
		for _, _restr := range rtrIndex.Collect(resPath...) {
			restr := _restr.(*resource.ResourceTranslation)

			out = append(out, restr)
		}
	})

	return
}
