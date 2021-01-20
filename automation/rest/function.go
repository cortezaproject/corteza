package rest

import (
	"context"
	"github.com/cortezaproject/corteza-server/automation/rest/request"
	"github.com/cortezaproject/corteza-server/automation/service"
	"github.com/cortezaproject/corteza-server/automation/types"
)

type (
	Function struct {
		reg interface {
			Functions() []*types.Function
		}
	}

	functionSetPayload struct {
		Set []*types.Function `json:"set"`
	}
)

func (Function) New() *Function {
	ctrl := &Function{reg: service.Registry()}
	return ctrl
}

func (ctrl Function) List(_ context.Context, _ *request.FunctionList) (interface{}, error) {
	return functionSetPayload{Set: ctrl.reg.Functions()}, nil
}
