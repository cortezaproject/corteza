package {package}

import (
	"context"
	"github.com/crusttech/crust/{project}/rest/request"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type {name} struct {}

func ({name}) New() *{name} {
	return &{name}{ldelim}{rdelim}
}

{foreach $calls as $call}
func (ctrl *{name}) {call.name|capitalize}(ctx context.Context, r *request.{name|ucfirst}{call.name|capitalize}) (interface{}, error) {
	return nil, errors.New("Not implemented: {name}.{call.name}")
}

{/foreach}
