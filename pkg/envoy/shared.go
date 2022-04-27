package envoy

import "github.com/cortezaproject/corteza-server/pkg/envoy/resource"

// procResSet is a little utility to run some op over given resources
//
// Helps cover special cases such as modules & module fields
func procResSet(resources resource.InterfaceSet, fn func(r resource.Interface)) {
	for _, res := range resources {
		fn(res)

		// Special case for modules since it has
		if modR, ok := res.(*resource.ComposeModule); ok {
			for _, f := range modR.ResFields {
				fn(f)
			}
		}
	}
}
