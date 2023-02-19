package envoy

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/store"
)

type (
	// StoreDecoder is responsible for fetching already stored Corteza resources
	// which are then managed by envoy and imported via an encoder.
	StoreDecoder struct{}
)

// Decode returns a set of envoy nodes based on the provided params
//
// StoreDecoder expects the DecodeParam of `storer` and `dal` which conform
// to the store.Storer and dal.FullService interfaces.
func (d StoreDecoder) Decode(ctx context.Context, p envoyx.DecodeParams) (out envoyx.NodeSet, err error) {
	var (
		s  store.Storer
		dl dal.FullService
	)

	// @todo we can optionally not require them based on what we're doing
	if auxS, ok := p.Params["storer"]; ok {
		s = auxS.(store.Storer)
	}
	if auxDl, ok := p.Params["dal"]; ok {
		dl = auxDl.(dal.FullService)
	}

	return d.decode(ctx, s, dl, p)
}

func (d StoreDecoder) decode(ctx context.Context, s store.Storer, dl dal.FullService, p envoyx.DecodeParams) (out envoyx.NodeSet, err error) {
	// Transform passed filters into an ordered structure
	type (
		filterWrap struct {
			rt string
			f  envoyx.ResourceFilter
		}
	)
	wrappedFilters := make([]filterWrap, 0, len(p.Filter))
	for rt, f := range p.Filter {
		wrappedFilters = append(wrappedFilters, filterWrap{rt: rt, f: f})
	}

	// Get all requested scopes
	scopedNodes := make(envoyx.NodeSet, len(p.Filter))

	for i, a := range wrappedFilters {
		if a.f.Scope.ResourceType == "" {
			continue
		}

		// For now the scope can only point to namespace so this will do
		var nn envoyx.NodeSet
		nn, err = d.decodeNamespace(ctx, s, dl, d.makeNamespaceFilter(nil, nil, envoyx.ResourceFilter{Identifiers: a.f.Scope.Identifiers}))
		if err != nil {
			return
		}
		if len(nn) > 1 {
			err = fmt.Errorf("ambiguous scope %v", a.f.Scope)
			return
		}
		if len(nn) == 0 {
			err = fmt.Errorf("invalid scope: resource not found %v", a.f)
			return
		}

		scopedNodes[i] = nn[0]
	}

	// Get all requested references
	//
	// Keep an index for the Node and one for the reference to make our
	// lives easier.
	refNodes := make([]map[string]*envoyx.Node, len(p.Filter))
	refRefs := make([]map[string]envoyx.Ref, len(p.Filter))
	for i, a := range wrappedFilters {
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
				return nil, err
			}
			if len(aux) == 0 {
				return nil, fmt.Errorf("invalid reference %v", ref)
			}
			if len(aux) > 1 {
				return nil, fmt.Errorf("ambiguous reference: too many resources returned %v", a.f)
			}

			auxr[field] = aux[0]
			auxa[field] = aux[0].ToRef()
		}

		refNodes[i] = auxr
		refRefs[i] = auxa
	}

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

		}
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
		// Identifiers
		ii := envoyx.MakeIdentifiers(
			r.Handle,
			r.ID,
		)

		refs := map[string]envoyx.Ref{
			// Handle references
			"NamespaceID": envoyx.Ref{
				ResourceType: "corteza::compose:namespace",
				Identifiers:  envoyx.MakeIdentifiers(r.NamespaceID),
			},
		}

		refs = envoyx.MergeRefs(refs, d.decodeChartRefs(r))

		var scope envoyx.Scope

		scope = envoyx.Scope{
			ResourceType: refs["NamespaceID"].ResourceType,
			Identifiers:  refs["NamespaceID"].Identifiers,
		}
		for k, ref := range refs {
			ref.Scope = scope
			refs[k] = ref
		}

		out = append(out, &envoyx.Node{
			Resource: r,

			ResourceType: types.ChartResourceType,
			Identifiers:  ii,
			References:   refs,
			Scope:        scope,
		})
	}

	return
}

func (d StoreDecoder) makeChartFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.ChartFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.ChartID = ids

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
		// Identifiers
		ii := envoyx.MakeIdentifiers(
			r.Handle,
			r.ID,
		)

		refs := map[string]envoyx.Ref{
			// Handle references
			"NamespaceID": envoyx.Ref{
				ResourceType: "corteza::compose:namespace",
				Identifiers:  envoyx.MakeIdentifiers(r.NamespaceID),
			},
		}

		var scope envoyx.Scope

		scope = envoyx.Scope{
			ResourceType: refs["NamespaceID"].ResourceType,
			Identifiers:  refs["NamespaceID"].Identifiers,
		}
		for k, ref := range refs {
			ref.Scope = scope
			refs[k] = ref
		}

		out = append(out, &envoyx.Node{
			Resource: r,

			ResourceType: types.ModuleResourceType,
			Identifiers:  ii,
			References:   refs,
			Scope:        scope,
		})
	}

	aux, err := d.extendedModuleDecoder(ctx, s, dl, f, out)
	if err != nil {
		return
	}
	out = append(out, aux...)

	return
}

func (d StoreDecoder) makeModuleFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.ModuleFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.ModuleID = ids

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
		// Identifiers
		ii := envoyx.MakeIdentifiers(
			r.ID,
			r.Name,
		)

		refs := map[string]envoyx.Ref{
			// Handle references
			"ModuleID": envoyx.Ref{
				ResourceType: "corteza::compose:module",
				Identifiers:  envoyx.MakeIdentifiers(r.ModuleID),
			},
		}

		var scope envoyx.Scope

		scope = envoyx.Scope{
			ResourceType: refs["NamespaceID"].ResourceType,
			Identifiers:  refs["NamespaceID"].Identifiers,
		}
		for k, ref := range refs {
			ref.Scope = scope
			refs[k] = ref
		}

		out = append(out, &envoyx.Node{
			Resource: r,

			ResourceType: types.ModuleFieldResourceType,
			Identifiers:  ii,
			References:   refs,
			Scope:        scope,
		})
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
		// Identifiers
		ii := envoyx.MakeIdentifiers(
			r.ID,
			r.Slug,
		)

		refs := map[string]envoyx.Ref{}

		var scope envoyx.Scope

		scope = envoyx.Scope{
			ResourceType: types.NamespaceResourceType,
			Identifiers:  ii,
		}

		out = append(out, &envoyx.Node{
			Resource: r,

			ResourceType: types.NamespaceResourceType,
			Identifiers:  ii,
			References:   refs,
			Scope:        scope,
		})
	}

	return
}

func (d StoreDecoder) makeNamespaceFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.NamespaceFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.NamespaceID = ids

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
		// Identifiers
		ii := envoyx.MakeIdentifiers(
			r.Handle,
			r.ID,
		)

		refs := map[string]envoyx.Ref{
			// Handle references
			"ModuleID": envoyx.Ref{
				ResourceType: "corteza::compose:module",
				Identifiers:  envoyx.MakeIdentifiers(r.ModuleID),
			},
			// Handle references
			"NamespaceID": envoyx.Ref{
				ResourceType: "corteza::compose:namespace",
				Identifiers:  envoyx.MakeIdentifiers(r.NamespaceID),
			},
			// Handle references
			"SelfID": envoyx.Ref{
				ResourceType: "corteza::compose:page",
				Identifiers:  envoyx.MakeIdentifiers(r.SelfID),
			},
		}

		var scope envoyx.Scope

		scope = envoyx.Scope{
			ResourceType: refs["NamespaceID"].ResourceType,
			Identifiers:  refs["NamespaceID"].Identifiers,
		}
		for k, ref := range refs {
			ref.Scope = scope
			refs[k] = ref
		}

		out = append(out, &envoyx.Node{
			Resource: r,

			ResourceType: types.PageResourceType,
			Identifiers:  ii,
			References:   refs,
			Scope:        scope,
		})
	}

	return
}

func (d StoreDecoder) makePageFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.PageFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.PageID = ids

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
