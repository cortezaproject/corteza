package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	types "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
)

type workflow struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        workflowAccessController
	services  *workflowServices
}

func (svc *workflow) Search(ctx context.Context, filter types.WorkflowFilter) (set types.WorkflowSet, f types.WorkflowFilter, err error) {
	var (
		aProps = &workflowActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Workflow) (bool, error) {
		if !svc.ac.CanReadWorkflow(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if !svc.ac.CanSearchWorkflows(ctx) {
			return WorkflowErrNotAllowedToSearch()
		}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.Workflow{}.LabelResourceKind(),
				filter.Labels,
			)
			if err != nil {
				return err
			}

			// labels specified but no labeled resources found
			if len(filter.LabeledIDs) == 0 {
				return nil
			}
		}

		if set, f, err = store.SearchAutomationWorkflows(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledWorkflows(set)...); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, WorkflowActionSearch, err)
}

// toLabeledWorkflows converts to []label.LabeledResource
func toLabeledWorkflows(set []*types.Workflow) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}

func (svc *workflow) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}

func (svc *workflow) Exec(ctx context.Context, workflowID uint64, p types.WorkflowExecParams) (results *expr.Vars, sessionID uint64, stacktrace types.Stacktrace, err error) {
	var (
		aProps = &workflowActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		results, sessionID, stacktrace, err = svc.onExec(ctx, aProps, workflowID, p)
		return err
	}()

	return results, sessionID, stacktrace, svc.recordAction(ctx, aProps, WorkflowActionExecute, err)
}
