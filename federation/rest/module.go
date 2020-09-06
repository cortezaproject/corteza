package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/federation/rest/request"
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
	return &struct{}{}, nil
}
