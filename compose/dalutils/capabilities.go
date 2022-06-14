package dalutils

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
)

func recCreateCapabilities(m *types.Module) (out capabilities.Set) {
	return capabilities.CreateCapabilities(m.ModelConfig.Capabilities...)
}

func recUpdateCapabilities(m *types.Module) (out capabilities.Set) {
	return capabilities.UpdateCapabilities(m.ModelConfig.Capabilities...)
}

func recDeleteCapabilities(m *types.Module) (out capabilities.Set) {
	return capabilities.DeleteCapabilities(m.ModelConfig.Capabilities...)
}

func recFilterCapabilities(f types.RecordFilter) (out capabilities.Set) {
	if f.PageCursor != nil {
		out = append(out, capabilities.Paging)
	}

	if f.IncPageNavigation {
		out = append(out, capabilities.Paging)
	}

	if f.IncTotal {
		out = append(out, capabilities.Stats)
	}

	if f.Sort != nil {
		out = append(out, capabilities.Sorting)
	}

	return
}

func recSearchCapabilities(m *types.Module, f types.RecordFilter) (out capabilities.Set) {
	return capabilities.SearchCapabilities(m.ModelConfig.Capabilities...).
		Union(recFilterCapabilities(f))
}

func recLookupCapabilities(m *types.Module) (out capabilities.Set) {
	return capabilities.LookupCapabilities(m.ModelConfig.Capabilities...)
}
