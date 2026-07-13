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
	execTypes "github.com/crusttech/human/server/pkg/automation_exec/types"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
)

type ngAutomation struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        ngAutomationAccessController
	services  *ngAutomationServices
}

func (svc *ngAutomation) Search(ctx context.Context, filter types.NgAutomationFilter) (set types.NgAutomationSet, f types.NgAutomationFilter, err error) {
	var (
		aProps = &ngAutomationActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.NgAutomation) (bool, error) {
		if !svc.ac.CanReadNgAutomation(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if !svc.ac.CanSearchNgAutomations(ctx) {
			return NgAutomationErrNotAllowedToSearch()
		}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.NgAutomation{}.LabelResourceKind(),
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

		if set, f, err = store.SearchAutomationNgAutomations(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledNgAutomations(set)...); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, NgAutomationActionSearch, err)
}

// toLabeledNgAutomations converts to []label.LabeledResource
func toLabeledNgAutomations(set []*types.NgAutomation) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
func (svc *ngAutomation) guard(_ context.Context, _ *types.NgAutomation) error { return nil }

func (svc *ngAutomation) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}

func (svc *ngAutomation) Exec(ctx context.Context, automationID uint64, p types.NgAutomationExecParams) (executionID id.ID, err error) {
	var (
		aProps = &ngAutomationActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		executionID, err = svc.onExec(ctx, aProps, automationID, p)
		return err
	}()

	return executionID, svc.recordAction(ctx, aProps, NgAutomationActionExec, err)
}

func (svc *ngAutomation) ExecAndWait(ctx context.Context, automationID uint64, p types.NgAutomationExecParams) (out *execTypes.ExecutionResult, err error) {
	var (
		aProps = &ngAutomationActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		out, err = svc.onExecAndWait(ctx, aProps, automationID, p)
		return err
	}()

	return out, svc.recordAction(ctx, aProps, NgAutomationActionExecAndWait, err)
}

func (svc *ngAutomation) GetExecutions(ctx context.Context, automationID uint64) (out []*execTypes.ExecutionResult, err error) {
	var (
		aProps = &ngAutomationActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}

		out, err = svc.onGetExecutions(ctx, aProps, automationID)
		return err
	}()

	return out, svc.recordAction(ctx, aProps, NgAutomationActionGetExecutions, err)
}

func (svc *ngAutomation) GetExecutionTrace(ctx context.Context, exeID uint64, executionID uint64, rev int) (out []execTypes.StackFrame, err error) {
	var (
		aProps = &ngAutomationActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}

		out, err = svc.onGetExecutionTrace(ctx, aProps, exeID, executionID, rev)
		return err
	}()

	return out, svc.recordAction(ctx, aProps, NgAutomationActionGetExecutionTrace, err)
}

func (svc *ngAutomation) GetAllExecutions(ctx context.Context, f execTypes.ExecutionFilter) (out []*execTypes.ExecutionResult, err error) {
	var (
		aProps = &ngAutomationActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}

		out, err = svc.onGetAllExecutions(ctx, aProps, f)
		return err
	}()

	return out, svc.recordAction(ctx, aProps, NgAutomationActionGetAllExecutions, err)
}

func (svc *ngAutomation) Load(ctx context.Context) (err error) {
	var (
		aProps = &ngAutomationActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onLoad(ctx, aProps)
		return err
	}()

	return svc.recordAction(ctx, aProps, NgAutomationActionLoad, err)
}
