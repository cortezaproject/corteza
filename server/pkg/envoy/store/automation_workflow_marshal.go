package store

import (
	"context"

	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
)

func newAutomationWorkflowFromResource(res *resource.AutomationWorkflow, cfg *EncoderConfig) resourceState {
	return &automationWorkflow{
		cfg: mergeConfig(cfg, res.Config()),

		res: res,
	}
}

func (n *automationWorkflow) Prepare(ctx context.Context, pl *payload) (err error) {
	err = n.prepareWorkflows(ctx, pl)
	if err != nil {
		return err
	}

	return n.prepareTriggers(ctx, pl)
}

func (n *automationWorkflow) prepareWorkflows(ctx context.Context, pl *payload) (err error) {
	// Reset old identifiers
	n.res.Res.ID = 0

	// Try to get the original workflow
	n.wf, err = findAutomationWorkflowStore(ctx, pl.s, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	if n.wf != nil {
		n.res.Res.ID = n.wf.ID
	}
	return nil
}

func (n *automationWorkflow) prepareTriggers(ctx context.Context, pl *payload) (err error) {
	if n.wf == nil || n.wf.ID == 0 {
		return nil
	}

	// Reset old identifiers
	for _, t := range n.res.Triggers {
		t.Res.ID = 0
		t.Res.WorkflowID = 0
	}

	// Try to find any related triggers for this workflow
	tt, _, err := store.SearchAutomationTriggers(ctx, pl.s, types.TriggerFilter{
		WorkflowID: []uint64{n.wf.ID},
		Disabled:   filter.StateInclusive,
	})
	if err != nil {
		return err
	}

	n.tt = tt
	return nil
}

func (n *automationWorkflow) Encode(ctx context.Context, pl *payload) (err error) {
	err = n.encodeWorkflow(ctx, pl)
	if err != nil {
		return err
	}

	return n.encodeTriggers(ctx, pl)
}

func (n *automationWorkflow) encodeWorkflow(ctx context.Context, pl *payload) (err error) {
	res := n.res.Res
	exists := n.wf != nil && n.wf.ID > 0

	// Determine the ID
	if res.ID <= 0 && exists {
		res.ID = n.wf.ID
	}
	if res.ID <= 0 {
		res.ID = NextID()
	}

	// Sys users
	us, err := resolveUserstamps(ctx, pl.s, pl.state.ParentResources, n.res.Userstamps())
	if err != nil {
		return err
	}

	ts := n.res.Timestamps()
	if ts != nil {
		if ts.CreatedAt != nil {
			res.CreatedAt = *ts.CreatedAt.T
		} else {
			res.CreatedAt = *now()
		}
		if ts.UpdatedAt != nil {
			res.UpdatedAt = ts.UpdatedAt.T
		}
		if ts.DeletedAt != nil {
			res.DeletedAt = ts.DeletedAt.T
		}
	}

	res.CreatedBy = pl.invokerID
	if us != nil {
		if us.OwnedBy != nil {
			res.OwnedBy = us.OwnedBy.UserID
		}
		if us.CreatedBy != nil {
			res.CreatedBy = us.CreatedBy.UserID
		}
		if us.UpdatedBy != nil {
			res.UpdatedBy = us.UpdatedBy.UserID
		}
		if us.DeletedBy != nil {
			res.DeletedBy = us.DeletedBy.UserID
		}
		if us.RunAs != nil {
			res.RunAs = us.RunAs.UserID
		}
	}

	res.Steps = make(types.WorkflowStepSet, 0, 100)
	for _, sres := range n.res.Steps {
		res.Steps = append(res.Steps, sres.Res)
	}

	res.Paths = make(types.WorkflowPathSet, 0, 100)
	for _, pres := range n.res.Paths {
		p := pres.Res
		res.Paths = append(res.Paths, p)
	}

	// Evaluate the resource skip expression
	// @todo expand available parameters; similar implementation to automation/types/record@Dict
	if skip, err := basicSkipEval(ctx, n.cfg, !exists); err != nil {
		return err
	} else if skip {
		return nil
	}

	// Create a fresh workflow
	if !exists {
		return store.CreateAutomationWorkflow(ctx, pl.s, res)
	}

	// Update existing workflow
	switch n.cfg.OnExisting {
	case resource.Skip:
		return nil

	case resource.MergeLeft:
		res = mergeAutomationWorkflows(n.wf, res)

	case resource.MergeRight:
		res = mergeAutomationWorkflows(res, n.wf)
	}

	err = store.UpdateAutomationWorkflow(ctx, pl.s, res)
	if err != nil {
		return err
	}

	n.res.Res = res
	return nil
}

func (n *automationWorkflow) encodeTriggers(ctx context.Context, pl *payload) (err error) {
	exists := len(n.tt) > 0
	rr := make([]*types.Trigger, 0, len(n.res.Triggers))

	for _, tr := range n.res.Triggers {
		res := tr.Res
		res.WorkflowID = n.res.Res.ID
		res.ID = NextID()

		// Sys users
		us, err := resolveUserstamps(ctx, pl.s, pl.state.ParentResources, tr.Userstamps())
		if err != nil {
			return err
		}

		ts := tr.Timestamps()
		if ts != nil {
			if ts.CreatedAt != nil {
				res.CreatedAt = *ts.CreatedAt.T
			} else {
				res.CreatedAt = *now()
			}
			if ts.UpdatedAt != nil {
				res.UpdatedAt = ts.UpdatedAt.T
			}
			if ts.DeletedAt != nil {
				res.DeletedAt = ts.DeletedAt.T
			}
		}
		res.CreatedBy = pl.invokerID
		if us != nil {
			if us.OwnedBy != nil {
				res.OwnedBy = us.OwnedBy.UserID
			}
			if us.CreatedBy != nil {
				res.CreatedBy = us.CreatedBy.UserID
			}
			if us.UpdatedBy != nil {
				res.UpdatedBy = us.UpdatedBy.UserID
			}
			if us.DeletedBy != nil {
				res.DeletedBy = us.DeletedBy.UserID
			}
		}

		rr = append(rr, res)
	}

	// Create a fresh workflow
	if !exists {
		return store.CreateAutomationTrigger(ctx, pl.s, rr...)
	}

	// If these triggers already exist and we wish to modify them,
	// remove the old ones and create new ones
	switch n.cfg.OnExisting {
	case resource.Skip,
		resource.MergeLeft:
		return nil
	}

	err = store.DeleteAutomationTrigger(ctx, pl.s, n.tt...)
	return store.CreateAutomationTrigger(ctx, pl.s, rr...)
}
