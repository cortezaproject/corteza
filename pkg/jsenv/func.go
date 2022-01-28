package jsenv

import (
	"fmt"

	"github.com/dop251/goja"
)

type (
	Fn struct {
		f goja.Callable
	}
)

func (f Fn) Exec(i ...goja.Value) (interface{}, error) {
	ret, err := f.f(goja.Undefined(), i...)

	if err != nil {
		return nil, fmt.Errorf("could not run function: %s", err)
	}

	return ret.Export(), nil
}
