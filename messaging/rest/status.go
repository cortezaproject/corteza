package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/messaging/rest/request"
)

var _ = errors.Wrap

type Status struct {
	// xxx service.XXXService
}

func (Status) New() *Status {
	return &Status{}
}

func (ctrl *Status) List(ctx context.Context, r *request.StatusList) (interface{}, error) {
	return nil, errors.New("Not implemented: Status.list")
}

func (ctrl *Status) Set(ctx context.Context, r *request.StatusSet) (interface{}, error) {
	return nil, errors.New("Not implemented: Status.set")
}

func (ctrl *Status) Delete(ctx context.Context, r *request.StatusDelete) (interface{}, error) {
	return nil, errors.New("Not implemented: Status.delete")
}
