package rest

import (
	"context"

	"github.com/cortezaproject/corteza/server/automation/rest/request"
	"github.com/cortezaproject/corteza/server/automation/service"
	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	triggerDefinition struct {
		svc service.TriggerDefinitionService
	}
)

func TriggerDefinition() *triggerDefinition {
	return &triggerDefinition{
		svc: service.DefaultTriggerDefinition,
	}
}

func (ctrl *triggerDefinition) List(ctx context.Context, r *request.TriggerDefinitionList) (interface{}, error) {
	var (
		err error
		f   = types.TriggerDefinitionFilter{
			TriggerDefinitionID: r.TriggerDefinitionID,
			Query:               r.Query,
			Deleted:             filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, f, err := ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, f, err)
}

func (ctrl *triggerDefinition) Create(ctx context.Context, r *request.TriggerDefinitionCreate) (interface{}, error) {
	var (
		def = &types.TriggerDefinition{
			Handle:       r.Handle,
			SkipEventBus: r.SkipEventBus,
			InputSchema:  r.InputSchema,
			OutputSchema: r.OutputSchema,
			Meta:         r.Meta,
		}
	)

	return ctrl.svc.Create(ctx, def)
}

func (ctrl *triggerDefinition) Update(ctx context.Context, r *request.TriggerDefinitionUpdate) (interface{}, error) {
	var (
		def = &types.TriggerDefinition{
			ID:           r.TriggerDefinitionID,
			Handle:       r.Handle,
			SkipEventBus: r.SkipEventBus,
			InputSchema:  r.InputSchema,
			OutputSchema: r.OutputSchema,
			Meta:         r.Meta,
		}
	)

	return ctrl.svc.Update(ctx, def)
}

func (ctrl *triggerDefinition) Read(ctx context.Context, r *request.TriggerDefinitionRead) (interface{}, error) {
	return ctrl.svc.LookupByID(ctx, r.TriggerDefinitionID)
}

func (ctrl *triggerDefinition) Delete(ctx context.Context, r *request.TriggerDefinitionDelete) (interface{}, error) {
	return nil, ctrl.svc.DeleteByID(ctx, r.TriggerDefinitionID)
}

func (ctrl *triggerDefinition) Undelete(ctx context.Context, r *request.TriggerDefinitionUndelete) (interface{}, error) {
	return nil, ctrl.svc.UndeleteByID(ctx, r.TriggerDefinitionID)
}

func (ctrl *triggerDefinition) makeFilterPayload(_ context.Context, set types.TriggerDefinitionSet, f types.TriggerDefinitionFilter, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	return set, nil
}
