package rest

import (
	"context"
	"github.com/cortezaproject/corteza-server/automation/rest/request"
	"github.com/cortezaproject/corteza-server/automation/service"
	"github.com/cortezaproject/corteza-server/automation/types/fn"
)

type (
	Function struct {
		svc interface {
			RegisteredFn() []*fn.Function
		}
	}

	functionSetPayload struct {
		Set []*fn.Function `json:"set"`
	}
)

func (Function) New() *Function {
	ctrl := &Function{}
	ctrl.svc = service.DefaultWorkflow
	return ctrl
}

func (ctrl Function) List(ctx context.Context, r *request.FunctionList) (interface{}, error) {
	return functionSetPayload{Set: ctrl.svc.RegisteredFn()}, nil
}
