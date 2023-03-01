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

	"github.com/cortezaproject/corteza/server/automation/types"
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
		// Handle resources that don't belong to this decoder
		if !strings.HasPrefix(rt, "corteza::automation") {
			continue
		}

		wrappedFilters = append(wrappedFilters, filterWrap{rt: rt, f: f})
	}

	// Get all requested scopes
	scopedNodes := make(envoyx.NodeSet, len(p.Filter))

	// @note skipping scope logic since it's currently only supported within
	//       Compose resources.

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

			// @todo consider changing this.
			//       Currently it's required because the .decode may return some
			//       nested nodes as well.
			//       Consider a flag or a new function.
			aux = envoyx.NodesForResourceType(ref.ResourceType, aux...)
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
		case types.WorkflowResourceType:
			aux, err = d.decodeWorkflow(ctx, s, dl, d.makeWorkflowFilter(scopedNodes[i], refNodes[i], wf.f))
			if err != nil {
				return
			}
			for _, a := range aux {
				a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
				a.References = envoyx.MergeRefs(a.References, refRefs[i])
			}
			out = append(out, aux...)

		case types.TriggerResourceType:
			aux, err = d.decodeTrigger(ctx, s, dl, d.makeTriggerFilter(scopedNodes[i], refNodes[i], wf.f))
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
} // // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource workflow
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decodeWorkflow(ctx context.Context, s store.Storer, dl dal.FullService, f types.WorkflowFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchAutomationWorkflows(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		// Identifiers
		ii := envoyx.MakeIdentifiers(
			r.Handle,
			r.ID,
		)

		// Handle references
		// Omit any non-defined values
		refs := map[string]envoyx.Ref{}
		if r.CreatedBy > 0 {
			refs["CreatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.CreatedBy),
			}
		}
		if r.DeletedBy > 0 {
			refs["DeletedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.DeletedBy),
			}
		}
		if r.OwnedBy > 0 {
			refs["OwnedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.OwnedBy),
			}
		}
		if r.RunAs > 0 {
			refs["RunAs"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.RunAs),
			}
		}
		if r.UpdatedBy > 0 {
			refs["UpdatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.UpdatedBy),
			}
		}

		var scope envoyx.Scope

		out = append(out, &envoyx.Node{
			Resource: r,

			ResourceType: types.WorkflowResourceType,
			Identifiers:  ii,
			References:   refs,
			Scope:        scope,
		})
	}

	aux, err := d.extendedWorkflowDecoder(ctx, s, dl, f, out)
	if err != nil {
		return
	}
	out = append(out, aux...)

	return
}

// Resource should define a custom filter builder

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource trigger
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decodeTrigger(ctx context.Context, s store.Storer, dl dal.FullService, f types.TriggerFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchAutomationTriggers(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		// Identifiers
		ii := envoyx.MakeIdentifiers(
			r.ID,
		)

		// Handle references
		// Omit any non-defined values
		refs := map[string]envoyx.Ref{}
		if r.CreatedBy > 0 {
			refs["CreatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.CreatedBy),
			}
		}
		if r.DeletedBy > 0 {
			refs["DeletedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.DeletedBy),
			}
		}
		if r.OwnedBy > 0 {
			refs["OwnedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.OwnedBy),
			}
		}
		if r.UpdatedBy > 0 {
			refs["UpdatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.UpdatedBy),
			}
		}
		if r.WorkflowID > 0 {
			refs["WorkflowID"] = envoyx.Ref{
				ResourceType: "corteza::automation:workflow",
				Identifiers:  envoyx.MakeIdentifiers(r.WorkflowID),
			}
		}

		var scope envoyx.Scope

		out = append(out, &envoyx.Node{
			Resource: r,

			ResourceType: types.TriggerResourceType,
			Identifiers:  ii,
			References:   refs,
			Scope:        scope,
		})
	}

	return
}

// Resource should define a custom filter builder
