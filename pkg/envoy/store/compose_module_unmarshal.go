package store

import (
	"strconv"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

func newComposeModule(mod *types.Module) *composeModule {
	return &composeModule{
		mod: mod,
	}
}

func (mod *composeModule) MarshalEnvoy() ([]resource.Interface, error) {
	refNs := strconv.FormatUint(mod.mod.NamespaceID, 10)
	refMod := strconv.FormatUint(mod.mod.ID, 10)

	rMod := resource.NewComposeModule(mod.mod, refNs)
	for _, f := range mod.mod.Fields {
		r := resource.NewComposeModuleField(f, refNs, refMod)
		rMod.AddField(r)
	}

	return envoy.CollectNodes(
		rMod,
	)
}
