package rest

import (
	"context"

	"github.com/cortezaproject/corteza/server/automation/rest/request"
	"github.com/cortezaproject/corteza/server/automation/service"
	"github.com/cortezaproject/corteza/server/automation/types"
)

type (
	ConstructLibrary struct {
		reg interface {
			Functions() []types.ConstructFunction
			Triggers() []types.ConstructTrigger
		}
	}

	cFunctionSetPayload struct {
		Set []types.ConstructFunction `json:"set"`
	}
	cTriggerSetPayload struct {
		Set []types.ConstructTrigger `json:"set"`
	}
)

func (ConstructLibrary) New() *ConstructLibrary {
	ctrl := &ConstructLibrary{reg: service.ConstructLibrary()}
	return ctrl
}

func (ctrl *ConstructLibrary) Functions(_ context.Context, _ *request.ConstructLibraryFunctions) (interface{}, error) {
	return cFunctionSetPayload{Set: ctrl.reg.Functions()}, nil
}

func (ctrl *ConstructLibrary) Triggers(_ context.Context, _ *request.ConstructLibraryTriggers) (interface{}, error) {
	return cTriggerSetPayload{Set: ctrl.reg.Triggers()}, nil
}
