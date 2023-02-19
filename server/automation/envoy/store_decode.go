package envoy

import (
	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
)

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

	out.TriggerID = ids

	return
}
