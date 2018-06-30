package {package}

import (
	"github.com/pkg/errors"
)

{foreach $calls as $call}
func ({self} *{name}) {call.name|ucfirst}(r *{name|lcfirst}{call.name|ucfirst}Request) (interface{}, error) {
	return nil, errors.New("Not implemented: {name}.{call.name}")
}
{/foreach}
