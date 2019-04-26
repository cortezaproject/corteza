package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/messaging/rest/request"
	"github.com/crusttech/crust/messaging/types"
)

var _ = errors.Wrap

type Commands struct {
	// xxx service.XXXService
}

func (Commands) New() *Commands {
	return &Commands{}
}

func (ctrl *Commands) List(ctx context.Context, r *request.CommandsList) (interface{}, error) {
	return types.CommandSet{
		&types.Command{
			Name:        "echo",
			Description: "It does exactly what it says on the tin"},
		&types.Command{
			Name:        "shrug",
			Description: "It does exactly what it says on the tin"},
	}, nil
}
