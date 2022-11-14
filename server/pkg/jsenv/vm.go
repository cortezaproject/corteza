package jsenv

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/dop251/goja"
)

const (
	// placeholder for functions in js env
	exportFuncDescriptor = "_expFunc_"
)

type (
	Vm struct {
		g *goja.Runtime
		t Transformer
	}
)

func New(t Transformer) Vm {
	return Vm{
		g: goja.New(),
		t: t,
	}
}

// New creates a goja internal type
func (ss Vm) New(i interface{}) goja.Value {
	return ss.g.ToValue(i)
}

// Register the value in vm
func (ss Vm) Register(key string, i interface{}) error {
	return ss.g.Set(key, i)
}

// Fetch value from vm
func (ss Vm) Fetch(key string) goja.Value {
	return ss.g.Get(key)
}

// RegisterFunction registers the function to the vm and returns the
// function that can be used in go
func (ss Vm) RegisterFunction(s string, wrapperFn ...func() string) (f *Fn, err error) {
	if len(wrapperFn) > 0 {
		for _, wfn := range wrapperFn {
			s = fmt.Sprintf(wfn(), s)
		}
	} else {
		s = fmt.Sprintf("function (input) { %s }", s)
	}

	desc := ss.funcDescriptor(s)
	run := fmt.Sprintf("var %s=%s;", desc, s)

	err = ss.Eval(run)

	if err != nil {
		return
	}

	internalF := ss.Fetch(desc)

	if internalF == nil {
		err = errors.New("could not fetch registered value")
		return
	}

	fnn, ok := goja.AssertFunction(internalF)

	if !ok {
		err = errors.New("could not assert function")
		return
	}

	return &Fn{
		f: fnn,
	}, nil
}

// Eval transforms the input js to the specified
// version and evals in vm
func (ss Vm) Eval(p string) (err error) {
	tr, err := ss.t.Transform(p)

	if err != nil {
		return
	}

	_, err = ss.g.RunString(string(tr))
	return
}

// Compile is used only when parsing the
// input evaluation without actually running it
func (ss Vm) Compile(p string) (err error) {
	_, err = goja.Parse("", p)
	return
}

func (ss Vm) genID(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

func (ss Vm) funcDescriptor(s string) string {
	return fmt.Sprintf("%s%s", exportFuncDescriptor, ss.genID(s))
}
