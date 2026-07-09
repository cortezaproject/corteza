package envoy

import (
	"context"
	"fmt"

	"github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/envoyx"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/resourceref"
	"github.com/crusttech/human/server/store"
	"github.com/spf13/cast"
)

func (e StoreEncoder) prepare(ctx context.Context, p envoyx.EncodeParams, s store.Storer, rt string, nn envoyx.NodeSet) (err error) {
	return
}

func (d StoreDecoder) extendDecoder(ctx context.Context, s store.Storer, dl dal.FullService, _ envoyx.DecodeParams, rt string, nodes map[string]*envoyx.Node, f envoyx.ResourceFilter) (out envoyx.NodeSet, err error) {
	return
}

func (d StoreDecoder) makeWorkflowFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.WorkflowFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.IdentsAsStrings()
	_ = ids
	_ = hh

	out.WorkflowID = ids

	if len(hh) > 0 {
		out.Handle = hh[0]
	}

	return
}

func (d StoreDecoder) makeTriggerFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.TriggerFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.TriggerID = id.Strings(ids...)

	return
}

// decodeWorkflowRefs returns path-keyed envoy refs for workflow step arguments
// that reference external resources (modules, namespaces, roles, etc.).
// Called by the generated store decoder when extendedRefDecoder is true.
func decodeWorkflowRefs(wf *types.Workflow) map[string]envoyx.Ref {
	refs := make(map[string]envoyx.Ref)
	for i, s := range wf.Steps {
		if s == nil {
			continue
		}
		paramKinds := types.WorkflowStepParamKinds(s)
		if len(paramKinds) == 0 {
			continue
		}
		for _, a := range s.Arguments {
			if a == nil || a.Value == nil {
				continue
			}
			kind, ok := paramKinds[a.Target]
			if !ok {
				continue
			}
			val := cast.ToString(a.Value)
			ref := resourceref.MakeIdent(kind, val, resourceref.ReasonStepArgument)
			if ref.IsEmpty() {
				continue
			}
			key := fmt.Sprintf("Steps.%d.Arguments.%s", i, a.Target)
			var ident any
			if id := ref.ID(); id > 0 {
				ident = id
			} else {
				ident = ref.Label
			}
			refs[key] = envoyx.Ref{
				ResourceType: ref.Kind(),
				Identifiers:  envoyx.MakeIdentifiers(ident),
			}
		}
	}
	return refs
}

func (d StoreDecoder) extendedWorkflowDecoder(ctx context.Context, s store.Storer, dl dal.FullService, f types.WorkflowFilter, base envoyx.NodeSet) (out envoyx.NodeSet, err error) {
	for _, b := range base {
		wf := b.Resource.(*types.Workflow)

		filters, err := d.decodeTrigger(ctx, s, dl, types.TriggerFilter{
			WorkflowID: id.Strings(wf.ID),
		})
		if err != nil {
			return nil, err
		}

		out = append(out, filters...)
	}

	return
}
