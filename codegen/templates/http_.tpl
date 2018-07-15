package {package}

import (
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type {name} struct {}

func ({name}) New() *{name} {
	return &{name}{ldelim}{rdelim}
}

{foreach $calls as $call}
func (*{name}) {call.name|capitalize}(r *{name|lcfirst}{call.name|capitalize}Request) (interface{}, error) {
	return nil, errors.New("Not implemented: {name}.{call.name}")
}

{/foreach}
