package envoy

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/pkg/errors"
)

type (
	// StoreDecoder is responsible for generating Envoy nodes from already stored
	// resources which can then be managed by Envoy and imported via an encoder.
	StoreDecoder struct{}

	filterWrap struct {
		rt string
		f  envoyx.ResourceFilter
	}
)

const (
	paramsKeyStorer = "storer"
	paramsKeyDAL    = "dal"
)

var (
	// @todo temporary fix to make unused pkg/id not throw errors
	_ = id.Next
)

// Decode returns a set of envoy nodes based on the provided params
//
// StoreDecoder expects the DecodeParam of `storer` and `dal` which conform
// to the store.Storer and dal.FullService interfaces.
func (d StoreDecoder) Decode(ctx context.Context, p envoyx.DecodeParams) (out envoyx.NodeSet, err error) {
	// @todo we can optionally not require them based on what we're doing
	s, err := d.getStorer(p)
	if err != nil {
		return
	}
	dl, err := d.getDal(p)
	if err != nil {
		return
	}

	return d.decode(ctx, s, dl, p)
}

func (d StoreDecoder) decode(ctx context.Context, s store.Storer, dl dal.FullService, p envoyx.DecodeParams) (out envoyx.NodeSet, err error) {
	// Preprocessing and basic filtering (to omit what this decoder can't handle)
	wrappedFilters := d.prepFilters(p.Filter)

	// Get all scoped nodes
	scopedNodes, err := d.getScopeNodes(ctx, s, dl, wrappedFilters)
	if err != nil {
		return
	}

	// Get all reference nodes
	refNodes, refRefs, err := d.getReferenceNodes(ctx, s, dl, wrappedFilters)
	if err != nil {
		return
	}

	// Process filters to get the envoy nodes
	err = func() (err error) {
		var aux envoyx.NodeSet
		for i, wf := range wrappedFilters {
			switch wf.rt {
			case types.ChartResourceType:
				aux, err = d.decodeChart(ctx, s, dl, d.makeChartFilter(scopedNodes[i], refNodes[i], wf.f))
				if err != nil {
					return
				}
				for _, a := range aux {
					a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
					a.References = envoyx.MergeRefs(a.References, refRefs[i])
				}
				out = append(out, aux...)

			case types.ModuleResourceType:
				aux, err = d.decodeModule(ctx, s, dl, d.makeModuleFilter(scopedNodes[i], refNodes[i], wf.f))
				if err != nil {
					return
				}
				for _, a := range aux {
					a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
					a.References = envoyx.MergeRefs(a.References, refRefs[i])
				}
				out = append(out, aux...)

			case types.ModuleFieldResourceType:
				aux, err = d.decodeModuleField(ctx, s, dl, d.makeModuleFieldFilter(scopedNodes[i], refNodes[i], wf.f))
				if err != nil {
					return
				}
				for _, a := range aux {
					a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
					a.References = envoyx.MergeRefs(a.References, refRefs[i])
				}
				out = append(out, aux...)

			case types.NamespaceResourceType:
				aux, err = d.decodeNamespace(ctx, s, dl, d.makeNamespaceFilter(scopedNodes[i], refNodes[i], wf.f))
				if err != nil {
					return
				}
				for _, a := range aux {
					a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
					a.References = envoyx.MergeRefs(a.References, refRefs[i])
				}
				out = append(out, aux...)

			case types.PageResourceType:
				aux, err = d.decodePage(ctx, s, dl, d.makePageFilter(scopedNodes[i], refNodes[i], wf.f))
				if err != nil {
					return
				}
				for _, a := range aux {
					a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
					a.References = envoyx.MergeRefs(a.References, refRefs[i])
				}
				out = append(out, aux...)

			case types.PageLayoutResourceType:
				aux, err = d.decodePageLayout(ctx, s, dl, d.makePageLayoutFilter(scopedNodes[i], refNodes[i], wf.f))
				if err != nil {
					return
				}
				for _, a := range aux {
					a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
					a.References = envoyx.MergeRefs(a.References, refRefs[i])
				}
				out = append(out, aux...)

			default:
				aux, err = d.extendDecoder(ctx, s, dl, wf.rt, refNodes[i], wf.f)
				if err != nil {
					return
				}
				for _, a := range aux {
					a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
					a.References = envoyx.MergeRefs(a.References, refRefs[i])
				}
				out = append(out, aux...)
			}
		}
		return
	}()
	if err != nil {
		err = errors.Wrap(err, "failed to decode filters")
		return
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource chart
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decodeChart(ctx context.Context, s store.Storer, dl dal.FullService, f types.ChartFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchComposeCharts(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		var n *envoyx.Node
		n, err = ChartToEnvoyNode(r)
		if err != nil {
			return
		}
		out = append(out, n)
	}

	return
}

func ChartToEnvoyNode(r *types.Chart) (node *envoyx.Node, err error) {
	// Identifiers
	ii := envoyx.MakeIdentifiers(
		r.Handle,
		r.ID,
	)

	// Handle references
	// Omit any non-defined values
	refs := map[string]envoyx.Ref{}
	if r.NamespaceID > 0 {
		refs["NamespaceID"] = envoyx.Ref{
			ResourceType: "corteza::compose:namespace",
			Identifiers:  envoyx.MakeIdentifiers(r.NamespaceID),
		}
	}

	refs = envoyx.MergeRefs(refs, decodeChartRefs(r))

	var scope envoyx.Scope

	scope = envoyx.Scope{
		ResourceType: refs["NamespaceID"].ResourceType,
		Identifiers:  refs["NamespaceID"].Identifiers,
	}
	for k, ref := range refs {
		// @todo temporary solution to not needlessly scope resources.
		//       Optimally, this would be selectively handled by codegen.
		if !strings.HasPrefix(ref.ResourceType, "corteza::compose") {
			continue
		}

		ref.Scope = scope
		refs[k] = ref
	}

	node = &envoyx.Node{
		Resource: r,

		ResourceType: types.ChartResourceType,
		Identifiers:  ii,
		References:   refs,
		Scope:        scope,
	}
	return
}

func (d StoreDecoder) makeChartFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.ChartFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.ChartID = id.Strings(ids...)

	if len(hh) > 0 {
		out.Handle = hh[0]
	}

	// Refs
	var (
		ar *envoyx.Node
		ok bool
	)
	_ = ar
	_ = ok

	ar, ok = refs["NamespaceID"]
	if ok {
		out.NamespaceID = ar.Resource.GetID()
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource module
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decodeModule(ctx context.Context, s store.Storer, dl dal.FullService, f types.ModuleFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchComposeModules(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		var n *envoyx.Node
		n, err = ModuleToEnvoyNode(r)
		if err != nil {
			return
		}
		out = append(out, n)
	}

	aux, err := d.extendedModuleDecoder(ctx, s, dl, f, out)
	if err != nil {
		return
	}
	out = append(out, aux...)

	return
}

func ModuleToEnvoyNode(r *types.Module) (node *envoyx.Node, err error) {
	// Identifiers
	ii := envoyx.MakeIdentifiers(
		r.Handle,
		r.ID,
	)

	// Handle references
	// Omit any non-defined values
	refs := map[string]envoyx.Ref{}
	if r.NamespaceID > 0 {
		refs["NamespaceID"] = envoyx.Ref{
			ResourceType: "corteza::compose:namespace",
			Identifiers:  envoyx.MakeIdentifiers(r.NamespaceID),
		}
	}

	var scope envoyx.Scope

	scope = envoyx.Scope{
		ResourceType: refs["NamespaceID"].ResourceType,
		Identifiers:  refs["NamespaceID"].Identifiers,
	}
	for k, ref := range refs {
		// @todo temporary solution to not needlessly scope resources.
		//       Optimally, this would be selectively handled by codegen.
		if !strings.HasPrefix(ref.ResourceType, "corteza::compose") {
			continue
		}

		ref.Scope = scope
		refs[k] = ref
	}

	node = &envoyx.Node{
		Resource: r,

		ResourceType: types.ModuleResourceType,
		Identifiers:  ii,
		References:   refs,
		Scope:        scope,
	}
	return
}

func (d StoreDecoder) makeModuleFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.ModuleFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.ModuleID = id.Strings(ids...)

	if len(hh) > 0 {
		out.Handle = hh[0]
	}

	// Refs
	var (
		ar *envoyx.Node
		ok bool
	)
	_ = ar
	_ = ok

	ar, ok = refs["NamespaceID"]
	if ok {
		out.NamespaceID = ar.Resource.GetID()
	}

	out = d.extendModuleFilter(scope, refs, auxf, out)
	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource moduleField
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decodeModuleField(ctx context.Context, s store.Storer, dl dal.FullService, f types.ModuleFieldFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchComposeModuleFields(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		var n *envoyx.Node
		n, err = ModuleFieldToEnvoyNode(r)
		if err != nil {
			return
		}
		out = append(out, n)
	}

	return
}

func ModuleFieldToEnvoyNode(r *types.ModuleField) (node *envoyx.Node, err error) {
	// Identifiers
	ii := envoyx.MakeIdentifiers(
		r.ID,
		r.Name,
	)

	// Handle references
	// Omit any non-defined values
	refs := map[string]envoyx.Ref{}
	if r.ModuleID > 0 {
		refs["ModuleID"] = envoyx.Ref{
			ResourceType: "corteza::compose:module",
			Identifiers:  envoyx.MakeIdentifiers(r.ModuleID),
		}
	}

	refs = envoyx.MergeRefs(refs, decodeModuleFieldRefs(r))

	var scope envoyx.Scope

	scope = envoyx.Scope{
		ResourceType: refs["NamespaceID"].ResourceType,
		Identifiers:  refs["NamespaceID"].Identifiers,
	}
	for k, ref := range refs {
		// @todo temporary solution to not needlessly scope resources.
		//       Optimally, this would be selectively handled by codegen.
		if !strings.HasPrefix(ref.ResourceType, "corteza::compose") {
			continue
		}

		ref.Scope = scope
		refs[k] = ref
	}

	node = &envoyx.Node{
		Resource: r,

		ResourceType: types.ModuleFieldResourceType,
		Identifiers:  ii,
		References:   refs,
		Scope:        scope,
	}
	return
}

// Resource should define a custom filter builder

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource namespace
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decodeNamespace(ctx context.Context, s store.Storer, dl dal.FullService, f types.NamespaceFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchComposeNamespaces(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		var n *envoyx.Node
		n, err = NamespaceToEnvoyNode(r)
		if err != nil {
			return
		}
		out = append(out, n)
	}

	return
}

func NamespaceToEnvoyNode(r *types.Namespace) (node *envoyx.Node, err error) {
	// Identifiers
	ii := envoyx.MakeIdentifiers(
		r.ID,
		r.Slug,
	)

	// Handle references
	// Omit any non-defined values
	refs := map[string]envoyx.Ref{}

	var scope envoyx.Scope

	scope = envoyx.Scope{
		ResourceType: types.NamespaceResourceType,
		Identifiers:  ii,
	}

	node = &envoyx.Node{
		Resource: r,

		ResourceType: types.NamespaceResourceType,
		Identifiers:  ii,
		References:   refs,
		Scope:        scope,
	}
	return
}

func (d StoreDecoder) makeNamespaceFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.NamespaceFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.NamespaceID = id.Strings(ids...)

	if len(hh) > 0 {
		out.Slug = hh[0]
	}

	// Refs
	var (
		ar *envoyx.Node
		ok bool
	)
	_ = ar
	_ = ok

	out = d.extendNamespaceFilter(scope, refs, auxf, out)
	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource page
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decodePage(ctx context.Context, s store.Storer, dl dal.FullService, f types.PageFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchComposePages(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		var n *envoyx.Node
		n, err = PageToEnvoyNode(r)
		if err != nil {
			return
		}
		out = append(out, n)
	}

	return
}

func PageToEnvoyNode(r *types.Page) (node *envoyx.Node, err error) {
	// Identifiers
	ii := envoyx.MakeIdentifiers(
		r.Handle,
		r.ID,
	)

	// Handle references
	// Omit any non-defined values
	refs := map[string]envoyx.Ref{}
	if r.ModuleID > 0 {
		refs["ModuleID"] = envoyx.Ref{
			ResourceType: "corteza::compose:module",
			Identifiers:  envoyx.MakeIdentifiers(r.ModuleID),
		}
	}
	if r.NamespaceID > 0 {
		refs["NamespaceID"] = envoyx.Ref{
			ResourceType: "corteza::compose:namespace",
			Identifiers:  envoyx.MakeIdentifiers(r.NamespaceID),
		}
	}
	if r.SelfID > 0 {
		refs["SelfID"] = envoyx.Ref{
			ResourceType: "corteza::compose:page",
			Identifiers:  envoyx.MakeIdentifiers(r.SelfID),
		}
	}

	refs = envoyx.MergeRefs(refs, decodePageRefs(r))

	var scope envoyx.Scope

	scope = envoyx.Scope{
		ResourceType: refs["NamespaceID"].ResourceType,
		Identifiers:  refs["NamespaceID"].Identifiers,
	}
	for k, ref := range refs {
		// @todo temporary solution to not needlessly scope resources.
		//       Optimally, this would be selectively handled by codegen.
		if !strings.HasPrefix(ref.ResourceType, "corteza::compose") {
			continue
		}

		ref.Scope = scope
		refs[k] = ref
	}

	node = &envoyx.Node{
		Resource: r,

		ResourceType: types.PageResourceType,
		Identifiers:  ii,
		References:   refs,
		Scope:        scope,
	}
	return
}

func (d StoreDecoder) makePageFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.PageFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.PageID = id.Strings(ids...)

	if len(hh) > 0 {
		out.Handle = hh[0]
	}

	// Refs
	var (
		ar *envoyx.Node
		ok bool
	)
	_ = ar
	_ = ok

	ar, ok = refs["ModuleID"]
	if ok {
		out.ModuleID = ar.Resource.GetID()
	}

	ar, ok = refs["NamespaceID"]
	if ok {
		out.NamespaceID = ar.Resource.GetID()
	}

	ar, ok = refs["SelfID"]
	if ok {
		out.ParentID = ar.Resource.GetID()
	}

	out = d.extendPageFilter(scope, refs, auxf, out)
	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource pageLayout
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decodePageLayout(ctx context.Context, s store.Storer, dl dal.FullService, f types.PageLayoutFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchComposePageLayouts(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		var n *envoyx.Node
		n, err = PageLayoutToEnvoyNode(r)
		if err != nil {
			return
		}
		out = append(out, n)
	}

	return
}

func PageLayoutToEnvoyNode(r *types.PageLayout) (node *envoyx.Node, err error) {
	// Identifiers
	ii := envoyx.MakeIdentifiers(
		r.Handle,
		r.ID,
	)

	// Handle references
	// Omit any non-defined values
	refs := map[string]envoyx.Ref{}
	if r.NamespaceID > 0 {
		refs["NamespaceID"] = envoyx.Ref{
			ResourceType: "corteza::compose:namespace",
			Identifiers:  envoyx.MakeIdentifiers(r.NamespaceID),
		}
	}
	if r.OwnedBy > 0 {
		refs["OwnedBy"] = envoyx.Ref{
			ResourceType: "corteza::system:user",
			Identifiers:  envoyx.MakeIdentifiers(r.OwnedBy),
		}
	}
	if r.PageID > 0 {
		refs["PageID"] = envoyx.Ref{
			ResourceType: "corteza::compose:page",
			Identifiers:  envoyx.MakeIdentifiers(r.PageID),
		}
	}
	if r.ParentID > 0 {
		refs["ParentID"] = envoyx.Ref{
			ResourceType: "corteza::compose:page-layout",
			Identifiers:  envoyx.MakeIdentifiers(r.ParentID),
		}
	}

	var scope envoyx.Scope

	scope = envoyx.Scope{
		ResourceType: refs["NamespaceID"].ResourceType,
		Identifiers:  refs["NamespaceID"].Identifiers,
	}
	for k, ref := range refs {
		// @todo temporary solution to not needlessly scope resources.
		//       Optimally, this would be selectively handled by codegen.
		if !strings.HasPrefix(ref.ResourceType, "corteza::compose") {
			continue
		}

		ref.Scope = scope
		refs[k] = ref
	}

	node = &envoyx.Node{
		Resource: r,

		ResourceType: types.PageLayoutResourceType,
		Identifiers:  ii,
		References:   refs,
		Scope:        scope,
	}
	return
}

func (d StoreDecoder) makePageLayoutFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.PageLayoutFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.PageLayoutID = id.Strings(ids...)

	if len(hh) > 0 {
		out.Handle = hh[0]
	}

	// Refs
	var (
		ar *envoyx.Node
		ok bool
	)
	_ = ar
	_ = ok

	ar, ok = refs["NamespaceID"]
	if ok {
		out.NamespaceID = ar.Resource.GetID()
	}

	ar, ok = refs["PageID"]
	if ok {
		out.PageID = ar.Resource.GetID()
	}

	ar, ok = refs["ParentID"]
	if ok {
		out.ParentID = ar.Resource.GetID()
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utilities
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) getStorer(p envoyx.DecodeParams) (s store.Storer, err error) {
	aux, ok := p.Params[paramsKeyStorer]
	if ok {
		s, ok = aux.(store.Storer)
		if ok {
			return
		}
	}

	err = errors.Errorf("store decoder expects a storer conforming to store.Storer interface")
	return
}

func (d StoreDecoder) getDal(p envoyx.DecodeParams) (dl dal.FullService, err error) {
	aux, ok := p.Params[paramsKeyDAL]
	if ok {
		dl, ok = aux.(dal.FullService)
		if ok {
			return
		}
	}

	err = errors.Errorf("store decoder expects a DAL conforming to dal.FullService interface")
	return
}

func (d StoreDecoder) prepFilters(ff map[string]envoyx.ResourceFilter) (out []filterWrap) {
	out = make([]filterWrap, 0, len(ff))
	for rt, f := range ff {
		// Handle resources that don't belong to this decoder
		if !strings.HasPrefix(rt, "corteza::compose") {
			continue
		}

		out = append(out, filterWrap{rt: rt, f: f})
	}

	return
}

func (d StoreDecoder) getScopeNodes(ctx context.Context, s store.Storer, dl dal.FullService, ff []filterWrap) (scopes envoyx.NodeSet, err error) {
	// Get all requested scopes
	scopes = make(envoyx.NodeSet, len(ff))

	err = func() (err error) {
		for i, fw := range ff {
			if fw.f.Scope.ResourceType == "" {
				continue
			}

			// For now the scope can only point to namespace so this will do
			var nn envoyx.NodeSet
			nn, err = d.decodeNamespace(ctx, s, dl, d.makeNamespaceFilter(nil, nil, envoyx.ResourceFilter{Identifiers: fw.f.Scope.Identifiers}))
			if err != nil {
				return
			}
			if len(nn) > 1 {
				err = fmt.Errorf("ambiguous scope %v: matches multiple resources", fw.f.Scope)
				return
			}
			if len(nn) == 0 {
				err = fmt.Errorf("invalid scope %v: resource not found", fw.f)
				return
			}

			scopes[i] = nn[0]
		}
		return
	}()
	if err != nil {
		err = errors.Wrap(err, "failed to decode node scopes")
		return
	}

	return
}

// getReferenceNodes returns all of the nodes referenced by the nodes defined by the filters
//
// The nodes are provided as a slice (the same order as the filters) and as a map for easier lookups.
func (d StoreDecoder) getReferenceNodes(ctx context.Context, s store.Storer, dl dal.FullService, ff []filterWrap) (nodes []map[string]*envoyx.Node, refs []map[string]envoyx.Ref, err error) {
	nodes = make([]map[string]*envoyx.Node, len(ff))
	refs = make([]map[string]envoyx.Ref, len(ff))
	err = func() (err error) {
		for i, a := range ff {
			if len(a.f.Refs) == 0 {
				continue
			}

			auxr := make(map[string]*envoyx.Node, len(a.f.Refs))
			auxa := make(map[string]envoyx.Ref)
			for field, ref := range a.f.Refs {
				f := ref.ResourceFilter()
				aux, err := d.decode(ctx, s, dl, envoyx.DecodeParams{
					Type:   envoyx.DecodeTypeStore,
					Filter: f,
				})
				if err != nil {
					return err
				}

				// @todo consider changing this.
				//       Currently it's required because the .decode may return some
				//       nested nodes as well.
				//       Consider a flag or a new function.
				aux = envoyx.NodesForResourceType(ref.ResourceType, aux...)
				if len(aux) == 0 {
					return fmt.Errorf("invalid reference %v", ref)
				}
				if len(aux) > 1 {
					return fmt.Errorf("ambiguous reference: too many resources returned %v", a.f)
				}

				auxr[field] = aux[0]
				auxa[field] = aux[0].ToRef()
			}

			nodes[i] = auxr
			refs[i] = auxa
		}
		return
	}()
	if err != nil {
		err = errors.Wrap(err, "failed to decode node references")
		return
	}

	return
}
