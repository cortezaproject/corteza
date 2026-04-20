package envoy

import (
	"context"
	"time"

	"github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/envoyx"
	"github.com/crusttech/human/server/store"
)

func (e StoreEncoder) encode(ctx context.Context, p envoyx.EncodeParams, s store.Storer, rt string, nn envoyx.NodeSet, tree envoyx.Traverser) (err error) {
	return
}

func (e StoreEncoder) setWorkflowDefaults(res *types.Workflow) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}
	return
}

func (e StoreEncoder) validateWorkflow(res *types.Workflow) (err error) {
	return
}

func (e StoreEncoder) setTriggerDefaults(res *types.Trigger) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}
	return
}

func (e StoreEncoder) validateTrigger(res *types.Trigger) (err error) {
	return
}
