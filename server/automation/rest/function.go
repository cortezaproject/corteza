package rest

import (
	"context"
	"github.com/crusttech/human/server/automation/rest/request"
	"github.com/crusttech/human/server/automation/service"
	"github.com/crusttech/human/server/automation/types"
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
