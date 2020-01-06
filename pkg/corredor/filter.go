package corredor

import (
	"strings"
)

type (
	Filter struct {
		ResourceTypes        []string
		EventTypes           []string
		ExcludeServerScripts bool
		ExcludeClientScripts bool

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
