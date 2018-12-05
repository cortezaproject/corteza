package rest

import (
	"context"
	"github.com/crusttech/crust/crm/rest/request"
	"github.com/crusttech/crust/crm/service"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Workflow struct {
		workflow service.WorkflowService
	}
)

func (Workflow) New(workflowSvc service.WorkflowService) *Workflow {
	return &Workflow{workflowSvc}
}

func (ctrl *Workflow) List(ctx context.Context, r *request.WorkflowList) (interface{}, error) {
	return nil, errors.New("Not implemented: Workflow.list")
}

func (ctrl *Workflow) Create(ctx context.Context, r *request.WorkflowCreate) (interface{}, error) {
	return nil, errors.New("Not implemented: Workflow.create")
}

func (ctrl *Workflow) Get(ctx context.Context, r *request.WorkflowGet) (interface{}, error) {
	return nil, errors.New("Not implemented: Workflow.get")
}

func (ctrl *Workflow) Update(ctx context.Context, r *request.WorkflowUpdate) (interface{}, error) {
	return nil, errors.New("Not implemented: Workflow.update")
}

func (ctrl *Workflow) Delete(ctx context.Context, r *request.WorkflowDelete) (interface{}, error) {
	return nil, errors.New("Not implemented: Workflow.delete")
}
