package {package}

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/{project}/rest/request"
)

var _ = errors.Wrap

type {name|expose} struct {
	// xxx service.XXXService
}

func ({name|expose}) New() *{name|expose} {
	return &{name|expose}{ldelim}{rdelim}
}

{foreach $calls as $call}
func (ctrl *{name|expose}) {call.name|capitalize}(ctx context.Context, r *request.{name|expose}{call.name|capitalize}) (interface{}, error) {
	return nil, errors.New("Not implemented: {name|expose}.{call.name}")
}

{/foreach}
