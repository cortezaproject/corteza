package envoy

import "github.com/cortezaproject/corteza-server/pkg/envoy/resource"

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
