package store

import (
	"strconv"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

func newComposeChart(chr *types.Chart) *composeChart {
	return &composeChart{
		chr: chr,
	}
}

// MarshalEnvoy converts the chart struct to a resource
func (chr *composeChart) MarshalEnvoy() ([]resource.Interface, error) {
	refNs := strconv.FormatUint(chr.chr.NamespaceID, 10)
	refMods := make([]string, 0, 2)
	for _, r := range chr.chr.Config.Reports {
		refMods = append(refMods, strconv.FormatUint(r.ModuleID, 10))
	}

	return envoy.CollectNodes(
		resource.NewComposeChart(chr.chr, refNs, refMods),
	)
}
