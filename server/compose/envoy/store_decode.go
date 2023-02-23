package envoy

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/compose/dalutils"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/filter"
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
		mod := b.Resource.(*types.Module)
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

			mod.Fields = append(mod.Fields, f)
		}
	}

	return
}

func (d StoreDecoder) decodeChartRefs(c *types.Chart) (refs map[string]envoyx.Ref) {

	// @todo
	return
}

func (d StoreDecoder) extendDecoder(ctx context.Context, s store.Storer, dl dal.FullService, rt string, refs map[string]*envoyx.Node, rf envoyx.ResourceFilter) (out envoyx.NodeSet, err error) {
	switch rt {
	case ComposeRecordDatasourceAuxType:
		return d.decodeRecordDatasource(ctx, s, dl, refs, rf)
	}

	return
}

func (d StoreDecoder) decodeRecordDatasource(ctx context.Context, s store.Storer, dl dal.FullService, refs map[string]*envoyx.Node, rf envoyx.ResourceFilter) (out envoyx.NodeSet, err error) {
	var (
		ok bool

		module        *types.Module
		moduleNode    *envoyx.Node
		namespace     *types.Namespace
		namespaceNode *envoyx.Node
	)

	// Get the refs
	namespaceNode, ok = refs["NamespaceID"]
	if !ok {
		err = fmt.Errorf("namespace ref not found")
		return
	}
	namespace = namespaceNode.Resource.(*types.Namespace)

	moduleNode, ok = refs["ModuleID"]
	if !ok {
		err = fmt.Errorf("module ref not found")
		return
	}
	module = moduleNode.Resource.(*types.Module)

	// Get the iterator
	iter, _, err := dalutils.ComposeRecordsIterator(ctx, dl, module, types.RecordFilter{
		ModuleID:    module.ID,
		NamespaceID: namespace.ID,
		Paging: filter.Paging{
			Limit: rf.Limit,
		},
	})

	if err != nil {
		return
	}

	ou := &RecordDatasource{
		provider: &iteratorProvider{iter: iter},
		refToID:  make(map[string]uint64),
	}

	rr := map[string]envoyx.Ref{
		"NamespaceID": namespaceNode.ToRef(),
		"ModuleID":    moduleNode.ToRef(),
	}

	// @todo add refs based on module fields
	// for _, f := range module.Fields {
	// 	if f.Kind != "Record" {
	// 		continue
	// 	}

	// 	rr[fmt.Sprintf("%s.module", f.Name)] = r
	// 	rr[fmt.Sprintf("%s.datasource", f.Name)] = envoyx.Ref{
	// 		ResourceType: ComposeRecordDatasourceAuxType,
	// 		Identifiers:  r.Identifiers,
	// 		Scope:        r.Scope,
	// 		Optional:     true,
	// 	}
	// }

	out = append(out, &envoyx.Node{
		Datasource:   ou,
		ResourceType: ComposeRecordDatasourceAuxType,

		Identifiers: envoyx.MakeIdentifiers(module.ID, module.Handle),
		References:  rr,
		Scope:       moduleNode.Scope,
	})

	return
}
