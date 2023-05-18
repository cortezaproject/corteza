package envoy

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/compose/dalutils"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/id"
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
	refs = make(map[string]envoyx.Ref, len(c.Config.Reports))

	for i, r := range c.Config.Reports {
		if r.ModuleID == 0 {
			continue
		}

		refs[fmt.Sprintf("Config.Reports.%d.ModuleID", i)] = envoyx.Ref{
			ResourceType: types.ModuleResourceType,
			Identifiers:  envoyx.MakeIdentifiers(r.ModuleID),
		}
	}

	return
}

func decodeModuleFieldRefs(c *types.ModuleField) (refs map[string]envoyx.Ref) {
	refs = make(map[string]envoyx.Ref, 1)

	refs["NamespaceID"] = envoyx.Ref{
		ResourceType: types.NamespaceResourceType,
		Identifiers:  envoyx.MakeIdentifiers(c.NamespaceID),
	}

	id := c.Options.UInt64("moduleID")
	if id == 0 {
		return
	}

	refs["Options.ModuleID"] = envoyx.Ref{
		ResourceType: types.ModuleResourceType,
		Identifiers:  envoyx.MakeIdentifiers(id),
	}

	return
}

func decodePageRefs(p *types.Page) (refs map[string]envoyx.Ref) {
	refs = make(map[string]envoyx.Ref, len(p.Blocks)/2)

	for index, b := range p.Blocks {
		switch b.Kind {
		case "RecordList":
			refs = envoyx.MergeRefs(refs, getPageBlockRecordListRefs(b, index))

		case "Automation":
			refs = envoyx.MergeRefs(refs, getPageBlockAutomationRefs(b, index))

		case "RecordOrganizer":
			refs = envoyx.MergeRefs(refs, getPageBlockRecordOrganizerRefs(b, index))

		case "Chart":
			refs = envoyx.MergeRefs(refs, getPageBlockChartRefs(b, index))

		case "Calendar":
			refs = envoyx.MergeRefs(refs, getPageBlockCalendarRefs(b, index))

		case "Metric":
			refs = envoyx.MergeRefs(refs, getPageBlockMetricRefs(b, index))

		case "Comment":
			refs = envoyx.MergeRefs(refs, getPageBlockCommentRefs(b, index))

		case "Progress":
			refs = envoyx.MergeRefs(refs, getPageBlockProgressRefs(b, index))
		}
	}

	return
}

func (d StoreDecoder) extendDecoder(ctx context.Context, s store.Storer, dl dal.FullService, rt string, refs map[string]*envoyx.Node, rf envoyx.ResourceFilter) (out envoyx.NodeSet, err error) {
	switch rt {
	// @todo consider hooking into the regular record resource type as well
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
		// @todo consider providing defaults from the outside
		mapping: datasourceMapping{
			KeyField:    []string{"id"},
			Defaultable: true,
		},
	}

	rr := map[string]envoyx.Ref{
		"NamespaceID": namespaceNode.ToRef(),
		"ModuleID":    moduleNode.ToRef(),
	}

	out = append(out, &envoyx.Node{
		Datasource:   ou,
		ResourceType: ComposeRecordDatasourceAuxType,

		Identifiers: envoyx.MakeIdentifiers(module.ID, module.Handle),
		References:  rr,
		Scope:       moduleNode.Scope,
	})

	return
}
