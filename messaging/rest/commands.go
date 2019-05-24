package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/messaging/rest/request"
	"github.com/cortezaproject/corteza-server/messaging/types"
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
			Name:        "me",
			Description: "Illeism"},
		&types.Command{
			Name:        "shrug",
			Description: "It does exactly what it says on the tin"},
		&types.Command{
			Name:        "tableflip",
			Description: "Flatten a table in anger"},
		&types.Command{
			Name:        "unflip",
			Description: "Put the table back from a flip"},
	}, nil
}
