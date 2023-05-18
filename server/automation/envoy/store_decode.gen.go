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
	}()
	if err != nil {
		err = errors.Wrap(err, "failed to decode filters")
		return
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
		var n *envoyx.Node
		n, err = WorkflowToEnvoyNode(r)
		if err != nil {
			return
		}
		out = append(out, n)
	}

	aux, err := d.extendedWorkflowDecoder(ctx, s, dl, f, out)
	if err != nil {
		return
	}
	out = append(out, aux...)

	return
}

func WorkflowToEnvoyNode(r *types.Workflow) (node *envoyx.Node, err error) {
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

	node = &envoyx.Node{
		Resource: r,

		ResourceType: types.WorkflowResourceType,
		Identifiers:  ii,
		References:   refs,
		Scope:        scope,
	}
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
		var n *envoyx.Node
		n, err = TriggerToEnvoyNode(r)
		if err != nil {
			return
		}
		out = append(out, n)
	}

	return
}

func TriggerToEnvoyNode(r *types.Trigger) (node *envoyx.Node, err error) {
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

	node = &envoyx.Node{
		Resource: r,

		ResourceType: types.TriggerResourceType,
		Identifiers:  ii,
		References:   refs,
		Scope:        scope,
	}
	return
}

// Resource should define a custom filter builder

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
		if !strings.HasPrefix(rt, "corteza::automation") {
			continue
		}

		out = append(out, filterWrap{rt: rt, f: f})
	}

	return
}

func (d StoreDecoder) getScopeNodes(ctx context.Context, s store.Storer, dl dal.FullService, ff []filterWrap) (scopes envoyx.NodeSet, err error) {
	// Get all requested scopes
	scopes = make(envoyx.NodeSet, len(ff))

	// @note skipping scope logic since it's currently only supported within
	//       Compose resources.

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
