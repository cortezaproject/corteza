package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/store"
)

func (svc resourceTranslation) loadWorkflow(ctx context.Context, s store.Storer, ID uint64) (*types.Workflow, error) {
	return loadWorkflow(ctx, s, ID)
}
