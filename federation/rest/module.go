package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/federation/rest/request"
	"github.com/cortezaproject/corteza-server/federation/service"
)

type (
	payload struct{}
	Module  struct{}
)

func (Module) New() *Module {
	return &Module{}
}

func (ctrl Module) Read(ctx context.Context, r *request.ModuleRead) (interface{}, error) {
	// use filtering and call structure sync service
	s := service.ExposedModule()

	// find the correct node (from request) and use it here
	mod, err := s.FindByID(context.Background(), 0, r.GetModuleID())

	return mod, err
}
