package store

import (
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/envoy"
	"github.com/crusttech/human/server/pkg/envoy/resource"
)

func newComposeChart(chr *types.Chart) *composeChart {
	return &composeChart{
		chr: chr,
	}
}

func (chr *composeChart) MarshalEnvoy() ([]resource.Interface, error) {
	refNs := resource.MakeNamespaceRef(chr.chr.NamespaceID, "", "")
	refMods := make(resource.RefSet, 0, 2)
	for _, r := range chr.chr.Config.Reports {
		refMods = append(refMods, resource.MakeModuleRef(r.ModuleID, "", ""))
	}

	return envoy.CollectNodes(
		resource.NewComposeChart(chr.chr, refNs, refMods),
	)
}
