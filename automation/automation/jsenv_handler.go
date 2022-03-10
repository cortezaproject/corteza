package automation

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/jsenv"
)

type (
	jsenvHandler struct {
		reg jsenvHandlerRegistry
	}
)

func JsenvHandler(reg queueHandlerRegistry) *jsenvHandler {
	h := &jsenvHandler{
		reg: reg,
	}

	h.register()

	return h
}

func initVm() jsenv.Vm {
	tr := jsenv.NewTransformer(jsenv.LoaderJS, jsenv.TargetNoop)
	vm := jsenv.New(tr)

	// register a request body reader
	vm.Register("readRequestBody", ReadRequestBody)

	return vm
}

func (h jsenvHandler) execute(ctx context.Context, args *jsenvExecuteArgs) (res *jsenvExecuteResults, err error) {
	vm := initVm()

	res = &jsenvExecuteResults{}

	if !args.hasSource {
		err = fmt.Errorf("could not process payload, function missing")
		return
	}

	if !args.hasScope {
		err = fmt.Errorf("could not process payload, scope missing")
		return
	}

	fn, err := vm.RegisterFunction(args.Source)

	if err != nil {
		err = fmt.Errorf("could not register jsenv function: %s", err)
		return
	}

	out, err := fn.Exec(vm.New(expr.UntypedValue(args.Scope)))

	if err != nil {
		err = fmt.Errorf("could not exec jsenv function: %s", err)
		return
	}

	switch vv := out.(type) {
	case int64:
		res.ResultInt = int64(vv)
	case string:
		res.ResultString = string(vv)
	case bool:
		res.ResultBool = vv
	default:
		res.ResultAny = vv
	}

	return
}
