package jsenv

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_registerFunction(t *testing.T) {
	type (
		tf struct {
			name      string
			err       string
			fn        string
			wrapperFn func() string
		}
	)

	var (
		tcc = []tf{
			{
				name: "register js func",
				fn:   "return 1;",
			},
			{
				name:      "register js custom wrapper fn",
				fn:        "return second;",
				wrapperFn: func() string { return "function (first, second, third) { %s }" },
			},
			{
				name: "register js func error",
				fn:   "function () {return 1",
				err:  "SyntaxError: SyntaxError: (anonymous): Line 1:69 Unexpected end of input (and 1 more errors)",
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				f   *fn
				err error

				req = require.New(t)
				tr  = NewTransformer(LoaderJS, TargetNoop)
				vm  = New(tr)
			)

			if tc.wrapperFn != nil {
				f, err = vm.RegisterFunction(tc.fn, tc.wrapperFn)
			} else {
				f, err = vm.RegisterFunction(tc.fn)
			}

			if tc.err == "" {
				req.NoError(err)
				req.Equal(reflect.Func, reflect.TypeOf(f.f).Kind())
				req.Equal("Callable", reflect.TypeOf(f.f).Name())
			} else {
				req.Errorf(err, tc.err)
				req.Nil(f)
			}

		})
	}
}

func Test_compile(t *testing.T) {
	var (
		req = require.New(t)
		tr  = NewTransformer(LoaderJS, TargetNoop)
		vm  = New(tr)
	)

	err := vm.Compile(`function () {return 1`)

	req.Errorf(err, "SyntaxError: SyntaxError: (anonymous): Line 1:69 Unexpected end of input (and 1 more errors)")
}

func Test_funcDescriptor(t *testing.T) {
	var (
		req = require.New(t)
		tr  = NewTransformer(LoaderJS, TargetNoop)
		vm  = New(tr)
	)

	req.Equal("_expFunc_1a6126e35863d2e16ba8e40f40668fdb", vm.funcDescriptor(`function () {return 1}`))
}

func Test_fetchEval(t *testing.T) {
	var (
		req = require.New(t)
		tr  = NewTransformer(LoaderJS, TargetNoop)
		vm  = New(tr)
	)

	req.NoError(vm.Eval(`foo = 'bar'`))
	req.Equal("bar", vm.Fetch("foo").Export())
	req.Nil(vm.Fetch("bar"))
}
