package store

import (
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/envoy"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
)

func newComposeModule(mod *types.Module) *composeModule {
	return &composeModule{
		mod: mod,
	}
}

func (mod *composeModule) MarshalEnvoy() ([]resource.Interface, error) {
	refNs := resource.MakeNamespaceRef(mod.mod.NamespaceID, "", "")
	refMod := resource.MakeModuleRef(mod.mod.ID, mod.mod.Handle, mod.mod.Name)

	rMod := resource.NewComposeModule(mod.mod, refNs)
	for _, f := range mod.mod.Fields {
		r := resource.NewComposeModuleField(f, refNs, refMod)
		rMod.AddField(r)
	}

	return envoy.CollectNodes(
		rMod,
	)
}
