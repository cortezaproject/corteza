package envoy

import (
	"context"

	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
)

func (e StoreEncoder) prepare(ctx context.Context, p envoyx.EncodeParams, s store.Storer, rt string, nn envoyx.NodeSet) (err error) {
	return
}

func (d StoreDecoder) extendDecoder(ctx context.Context, s store.Storer, dl dal.FullService, rt string, nodes map[string]*envoyx.Node, f envoyx.ResourceFilter) (out envoyx.NodeSet, err error) {
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
