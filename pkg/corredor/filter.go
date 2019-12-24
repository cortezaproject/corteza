package corredor

import (
	"strings"
)

type (
	ManualScriptFilter struct {
		ResourceTypes        []string
		EventTypes           []string
		ExcludeServerScripts bool
		ExcludeClientScripts bool
	}
)

// PrefixResource adds service string (if not already there) to all resources
func (f *ManualScriptFilter) PrefixResource(service string) {
	for i := range f.ResourceTypes {
		if !strings.HasPrefix(f.ResourceTypes[i], service) {
			f.ResourceTypes[i] = service + ":" + f.ResourceTypes[i]
		}
	}
}
