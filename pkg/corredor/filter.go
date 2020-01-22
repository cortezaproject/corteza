package corredor

import (
	"strings"
)

type (
	Filter struct {
		ResourceTypes        []string `json:"resourceTypes"`
		EventTypes           []string `json:"eventTypes"`
		ExcludeServerScripts bool     `json:"excludeServerScripts"`
		ExcludeClientScripts bool     `json:"excludeClientScripts"`

		Page    uint `json:"page"`
		PerPage uint `json:"perPage"`
		Count   uint `json:"count"`
	}
)

// PrefixResources adds service string (if not already there) to all resources
func (f *Filter) PrefixResources(service string) {
	for i := range f.ResourceTypes {
		if !strings.HasPrefix(f.ResourceTypes[i], service) {
			f.ResourceTypes[i] = service + ":" + f.ResourceTypes[i]
		}
	}
}
