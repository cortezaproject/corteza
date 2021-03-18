package store

import (
	"strconv"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

func newComposePage(pg *types.Page) *composePage {
	return &composePage{
		pg: pg,
	}
}

func (pg *composePage) MarshalEnvoy() ([]resource.Interface, error) {
	refNs := strconv.FormatUint(pg.pg.NamespaceID, 10)
	refMod := ""
	refParent := ""
	if pg.pg.ModuleID > 0 {
		refMod = strconv.FormatUint(pg.pg.ModuleID, 10)
	}
	if pg.pg.SelfID > 0 {
		refParent = strconv.FormatUint(pg.pg.SelfID, 10)
	}

	return envoy.CollectNodes(
		resource.NewComposePage(pg.pg, refNs, refMod, refParent),
	)
}
