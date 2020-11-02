package tmp

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/store"
)

func CollectPreproc(is *importState, s store.Storer) []envoy.Processor {
	return []envoy.Processor{
		NewComposeModulePreproc(is, s),
		NewComposeNamespacePreproc(is, s),
		NewComposeRecordSetPreproc(is, s),
		// ...
	}
}
