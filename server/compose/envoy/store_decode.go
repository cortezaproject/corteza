package envoy

import (
	"context"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/store"
)

func (d StoreDecoder) extendNamespaceFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter, base types.NamespaceFilter) (out types.NamespaceFilter) {
	out = base

	if scope == nil {
		return
	}

	if scope.ResourceType == "" {
		return
	}

	// Overwrite it
	out.NamespaceID = []uint64{scope.Resource.GetID()}

	return
}

func (d StoreDecoder) extendModuleFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter, base types.ModuleFilter) (out types.ModuleFilter) {
	out = base

	if scope == nil {
		return
	}

	if scope.ResourceType == "" {
		return
	}

	// Overwrite it
	out.NamespaceID = scope.Resource.GetID()

	return
}

func (d StoreDecoder) extendPageFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter, base types.PageFilter) (out types.PageFilter) {
	out = base

	if scope == nil {
		return
	}

	if scope.ResourceType == "" {
		return
	}

	// Overwrite it
	out.NamespaceID = scope.Resource.GetID()

	return
}

func (d StoreDecoder) makeModuleFieldFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.ModuleFieldFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	// Refs
	var (
		ar *envoyx.Node
		ok bool
	)
	_ = ar
	_ = ok

	// ar, ok = refs["ModuleID"]
	// if ok {
	// 	out.ModuleID = ar.Resource.GetID()
	// }

	return
}

func (d StoreDecoder) extendedModuleDecoder(ctx context.Context, s store.Storer, dl dal.FullService, f types.ModuleFilter, base envoyx.NodeSet) (out envoyx.NodeSet, err error) {
	var ff types.ModuleFieldSet

	for _, b := range base {
		ff, _, err = store.SearchComposeModuleFields(ctx, s, types.ModuleFieldFilter{ModuleID: []uint64{b.Resource.GetID()}})
		if err != nil {
			return
		}

		// No need to assign them under the module since we're working with nodes now
		for _, f := range ff {
			out = append(out, &envoyx.Node{
				Resource: f,

				ResourceType: types.ModuleFieldResourceType,
				Identifiers:  envoyx.MakeIdentifiers(f.ID, f.Name),
				References: envoyx.MergeRefs(b.References, map[string]envoyx.Ref{
					"ModuleID": b.ToRef(),
				}),
				Scope: b.Scope,
			})
		}
	}

	return
}

func (d StoreDecoder) decodeChartRefs(c *types.Chart) (refs map[string]envoyx.Ref) {

	// @todo
	return
}
