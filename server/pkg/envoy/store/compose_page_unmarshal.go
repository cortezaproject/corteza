package store

import (
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/envoy"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
)

func newComposePage(pg *types.Page) *composePage {
	return &composePage{
		pg: pg,
	}
}

func (pg *composePage) MarshalEnvoy() ([]resource.Interface, error) {
	var (
		refNs     = resource.MakeNamespaceRef(pg.pg.NamespaceID, "", "")
		refMod    *resource.Ref
		refParent *resource.Ref
	)
	if pg.pg.ModuleID > 0 {
		refMod = resource.MakeModuleRef(pg.pg.ModuleID, "", "")
	}
	if pg.pg.SelfID > 0 {
		refParent = resource.MakePageRef(pg.pg.SelfID, "", "")
	}

	return envoy.CollectNodes(
		resource.NewComposePage(pg.pg, refNs, refMod, refParent),
	)
}
