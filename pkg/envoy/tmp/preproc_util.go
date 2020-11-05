package tmp

import (
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/store"
)

func CollectPreproc(is *encoderState, s store.Storer) []envoy.Processor {
	return []envoy.Processor{
		NewComposeModule(is, s),
		NewComposeNamespacePreproc(is, s),
		NewComposeRecordPreproc(is, s),
		// ...
	}
}
