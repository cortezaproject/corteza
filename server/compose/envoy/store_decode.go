package envoy

import (
	"context"
	"fmt"
	"sort"

	"github.com/crusttech/human/server/compose/dalutils"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/envoyx"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/resourceref"
	"github.com/crusttech/human/server/store"
	"github.com/spf13/cast"
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
	out.NamespaceID = id.Strings(scope.Resource.GetID())

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

// extendChartFilter narrows a chart decode to the namespace it was scoped to.
// Without it the generated makeChartFilter reads only refs["NamespaceID"] and
// silently ignores the scope, so a scoped decode — a namespace clone, say —
// pulls in every chart in the store. Those foreign charts then encode with
// their ORIGINAL namespace_id and handle, tripping the unique_handle index
// (compose/chart.cue) as "not unique". Same shape as the module and page
// extenders below/above; see CloneFromStore for the sibling ModuleField case.
func (d StoreDecoder) extendChartFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter, base types.ChartFilter) (out types.ChartFilter) {
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

// extendPageLayoutFilter is extendChartFilter's twin: layouts are cloned with
// the same namespace scope and were ignoring it in the same way. It bit less
// visibly only because a layout's unique index also spans page_id and layouts
// routinely carry an empty handle, which the index's predicate excludes — so
// instead of failing, a clone quietly dragged in every other namespace's
// layouts.
func (d StoreDecoder) extendPageLayoutFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter, base types.PageLayoutFilter) (out types.PageLayoutFilter) {
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
	var ff envoyx.NodeSet

	for _, b := range base {
		mod := b.Resource.(*types.Module)

		// Get all of the related module fields, append them to the output and
		// the original module (so other code can have access to the related fields)
		ff, err = d.decodeModuleField(ctx, s, dl, types.ModuleFieldFilter{ModuleID: []uint64{mod.ID}})
		if err != nil {
			return
		}

		// decodeModuleField (SearchComposeModuleFields) has no ORDER BY, so ff
		// arrives in whatever order the store happened to return -- not
		// necessarily Place order, and for a module whose fields were all
		// supplied in one createModule call (which, unlike updateModuleFields,
		// never assigns Place) not necessarily distinguishable by Place either,
		// since every field is Place=0. A raw read of the source namespace only
		// "looks" correctly ordered because it's never been rewritten since
		// those inserts; re-encoding these fields into a brand new namespace
		// (CloneFromStore's use of this decoder) is a fresh set of writes that
		// carries no such accidental guarantee, so the destination can come
		// back in a different order even though every field's Place value
		// round-trips unchanged.
		//
		// Stable-sort by the existing Place first (so a module whose fields do
		// carry a real, distinct Place is reproduced exactly), falling back to
		// this fetch's own order for ties -- i.e. reproducing whatever order a
		// read of the source shows right now, degenerate Place or not -- and
		// then re-sequence Place itself to a clean 0..n-1 walk of that order.
		// This is what makes the order stick on the far side of the clone: any
		// later read sorts by Place again, and Place is no longer degenerate.
		sort.SliceStable(ff, func(i, j int) bool {
			return ff[i].Resource.(*types.ModuleField).Place < ff[j].Resource.(*types.ModuleField).Place
		})
		for i, n := range ff {
			n.Resource.(*types.ModuleField).Place = i
		}

		for _, f := range ff {
			f.Scope = b.Scope
			f.References = envoyx.MergeRefs(f.References, b.References, map[string]envoyx.Ref{
				"ModuleID": b.ToRef(),
			})
			for k, ref := range f.References {
				ref.Scope = b.Scope
				f.References[k] = ref
			}

			mod.Fields = append(mod.Fields, f.Resource.(*types.ModuleField))
		}

		out = append(out, ff...)
	}

	return
}

func decodeChartRefs(c *types.Chart) (refs map[string]envoyx.Ref) {
	return toEnvoyRefs(c.ResourceRefs())
}

func decodeModuleFieldRefs(c *types.ModuleField) (refs map[string]envoyx.Ref) {
	refs = toEnvoyRefs(c.ResourceRefs())

	refs["NamespaceID"] = envoyx.Ref{
		ResourceType: types.NamespaceResourceType,
		Identifiers:  envoyx.MakeIdentifiers(c.NamespaceID),
	}

	return
}

func decodePageRefs(p *types.Page) (refs map[string]envoyx.Ref) {
	// only block refs; the page's own ModuleID ref is added by the generated decoder
	return toEnvoyRefs(p.Blocks.ResourceRefs())
}

// toEnvoyRefs maps shared resourceref extractor output to envoy references,
// keyed by the location of the reference within the resource
func toEnvoyRefs(rr []resourceref.Ref) (refs map[string]envoyx.Ref) {
	refs = make(map[string]envoyx.Ref, len(rr))

	for _, r := range rr {
		var ident any = r.Label
		if id := r.ID(); id > 0 {
			ident = id
		}

		refs[r.Resource] = envoyx.Ref{
			ResourceType: r.Kind(),
			Identifiers:  envoyx.MakeIdentifiers(ident),
		}
	}

	return
}

func (d StoreDecoder) extendDecoder(ctx context.Context, s store.Storer, dl dal.FullService, p envoyx.DecodeParams, rt string, refs map[string]*envoyx.Node, rf envoyx.ResourceFilter) (out envoyx.NodeSet, err error) {
	switch rt {
	// @todo consider hooking into the regular record resource type as well
	case ComposeRecordDatasourceAuxType:
		return d.decodeRecordDatasource(ctx, s, dl, p, refs, rf)
	}

	return
}

func (d StoreDecoder) decodeRecordDatasource(ctx context.Context, s store.Storer, dl dal.FullService, p envoyx.DecodeParams, refs map[string]*envoyx.Node, rf envoyx.ResourceFilter) (out envoyx.NodeSet, err error) {
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
		err = fmt.Errorf("missing NamespaceID reference")
		return
	}
	namespace = namespaceNode.Resource.(*types.Namespace)

	moduleNode, ok = refs["ModuleID"]
	if !ok {
		err = fmt.Errorf("missing ModuleID reference")
		return
	}
	module = moduleNode.Resource.(*types.Module)

	// Get the iterator
	iter, _, err := dalutils.ComposeRecordsIterator(ctx, dl, module, types.RecordFilter{
		Query:       rf.Query,
		ModuleID:    module.ID,
		NamespaceID: namespace.ID,
		Paging: filter.Paging{
			Limit: rf.Limit,
		},
	})

	if err != nil {
		return
	}

	mv := make(map[string]bool)
	for _, f := range module.Fields {
		if !f.Multi {
			continue
		}

		mv[f.Name] = true
	}

	rds := &RecordDatasource{
		Provider:    &iteratorProvider{iter: iter},
		multivalues: mv,

		refToID:     make(map[string]uint64),
		existingIDs: make(map[uint64]bool),
		// @todo consider providing defaults from the outside
		Mapping: envoyx.DatasourceMapping{
			KeyField:    []string{"id"},
			Defaultable: true,
		},
	}

	rds.Provider, err = mkIteratorProvider(ctx, s, dl, iter, module, cast.ToBool(p.Params["resolveRefs"]))
	if err != nil {
		return
	}

	rr := map[string]envoyx.Ref{
		"NamespaceID": namespaceNode.ToRef(),
		"ModuleID":    moduleNode.ToRef(),
	}

	out = append(out, &envoyx.Node{
		Datasource:   rds,
		ResourceType: ComposeRecordDatasourceAuxType,

		Identifiers: envoyx.MakeIdentifiers(module.ID, module.Handle),
		References:  rr,
		Scope:       moduleNode.Scope,
	})

	return
}
