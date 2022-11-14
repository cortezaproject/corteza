package store

import (
	"context"

	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	automationWorkflow struct {
		cfg *EncoderConfig

		res *resource.AutomationWorkflow
		wf  *types.Workflow
		tt  types.TriggerSet

		ux *userIndex
	}
	automationWorkflowSet []*automationWorkflow

	automationTrigger struct {
		cfg *EncoderConfig

		res *resource.AutomationTrigger
		tr  *types.Trigger
	}
	automationTriggerSet []*automationTrigger
)

// mergeAutomationWorkflows merges b into a, prioritising a
func mergeAutomationWorkflows(a, b *types.Workflow) *types.Workflow {
	c := a

	if c.Handle == "" {
		c.Handle = b.Handle
	}
	if c.Meta == nil {
		c.Meta = b.Meta
	}

	if c.Scope == nil {
		c.Scope = b.Scope
	}
	if c.Steps == nil {
		c.Steps = b.Steps
	}
	if c.Paths == nil {
		c.Paths = b.Paths
	}

	if c.RunAs == 0 {
		c.RunAs = b.RunAs
	}
	if c.OwnedBy == 0 {
		c.OwnedBy = b.OwnedBy
	}
	if c.CreatedBy == 0 {
		c.CreatedBy = b.CreatedBy
	}
	if c.UpdatedBy == 0 {
		c.UpdatedBy = b.UpdatedBy
	}
	if c.DeletedBy == 0 {
		c.DeletedBy = b.DeletedBy
	}

	if c.CreatedAt.IsZero() {
		c.CreatedAt = b.CreatedAt
	}

	if c.UpdatedAt == nil {
		c.UpdatedAt = b.UpdatedAt
	}

	if c.DeletedAt == nil {
		c.DeletedAt = b.DeletedAt
	}

	return c
}

// findAutomationWorkflow looks for the workflow in the resources & the store
//
// Provided resources are prioritized.
func findAutomationWorkflow(ctx context.Context, s store.Storer, rr resource.InterfaceSet, ii resource.Identifiers) (wf *types.Workflow, err error) {
	wf = resource.FindAutomationWorkflow(rr, ii)
	if wf != nil {
		return wf, nil
	}

	return findAutomationWorkflowStore(ctx, s, makeGenericFilter(ii))
}

// findAutomationWorkflowStore looks for the workflow in the store
func findAutomationWorkflowStore(ctx context.Context, s store.Storer, gf genericFilter) (wf *types.Workflow, err error) {
	if gf.id > 0 {
		wf, err = store.LookupAutomationWorkflowByID(ctx, s, gf.id)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if wf != nil {
			return
		}
	}

	for _, i := range gf.identifiers {
		wf, err = store.LookupAutomationWorkflowByHandle(ctx, s, i)
		if err == store.ErrNotFound {
			var nn types.WorkflowSet
			nn, _, err = store.SearchAutomationWorkflows(ctx, s, types.WorkflowFilter{
				Query: i,
				Paging: filter.Paging{
					Limit: 2,
				},
			})
			if len(nn) > 1 {
				return nil, resourceErrIdentifierNotUnique(i)
			}
			if len(nn) == 1 {
				wf = nn[0]
			}
		}

		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if wf != nil {
			return
		}
	}

	return nil, nil
}
