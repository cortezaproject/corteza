package filter

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/jsenv"
	"go.uber.org/zap"
)

type (
	workflow struct {
		types.FilterMeta
		d types.WfExecer

		params struct {
			Workflow uint64 `json:"workflow,string"`
		}
	}

	processerPayload struct {
		types.FilterMeta
		vm  jsenv.Vm
		log *zap.Logger

		params struct {
			Func   string `json:"func"`
			Encode bool   `json:"encode"`
		}
	}
)

func NewWorkflow(wf types.WfExecer) (p *workflow) {
	p = &workflow{}

	p.d = wf

	p.Name = "workflow"
	p.Label = "Workflow processer"
	p.Kind = types.Processer

	p.Args = []*types.FilterMetaArg{
		{
			Type:    "workflow",
			Label:   "workflow",
			Options: map[string]interface{}{},
		},
	}

	return
}

func (h workflow) String() string {
	return fmt.Sprintf("apigw function %s (%s)", h.Name, h.Label)
}

func (h workflow) Type() types.FilterKind {
	return h.Kind
}

func (h workflow) Meta() types.FilterMeta {
	return h.FilterMeta
}

func (h workflow) Weight() int {
	return h.Wgt
}

func (f *workflow) Merge(params []byte) (types.Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&f.params)

	return f, err
}

func (h workflow) Exec(ctx context.Context, scope *types.Scp) error {
	var (
		err error
	)

	payload, err := scope.Get("payload")

	if err != nil {
		return err
	}

	rr, err := scope.Get("request")

	if err != nil {
		return err
	}

	// setup scope for workflow
	vv := map[string]interface{}{
		"payload": payload,
		"request": rr,
	}

	// get the request data and put it into vars
	in, err := expr.NewVars(vv)

	if err != nil {
		return err
	}

	wp := atypes.WorkflowExecParams{
		Trace: true,
		// todo depending on settings per-route
		Async: false,
		// todo depending on settings per-route
		Wait:  true,
		Input: in,
	}

	out, _, err := h.d.Exec(ctx, h.params.Workflow, wp)

	if err != nil {
		return err
	}

	// merge out with scope
	merged, err := in.Merge(out)

	if err != nil {
		return err
	}

	mm, err := expr.CastToVars(merged)

	for k, v := range mm {
		scope.Set(k, v)
	}

	ss := scope.Filter(func(k string, v interface{}) bool {
		if k == "eventType" || k == "resourceType" {
			return false
		}

		return true
	})

	scope = ss

	// spew.Dump(scope.Get("payload"))
	// trgt, _ := ss.Get("trgt")

	// scope.Writer().Write([]byte(trgt.(*expr.String).GetValue()))

	return err
}

func NewPayload(l *zap.Logger) (p *processerPayload) {
	p = &processerPayload{}

	// todo - check the consequences of doing this here
	p.vm = jsenv.New(jsenv.NewTransformer(jsenv.LoaderJS, jsenv.TargetES2016))
	p.log = l

	p.Name = "payload"
	p.Label = "Payload processer"
	p.Kind = types.Processer

	p.Args = []*types.FilterMetaArg{
		{
			Type:    "jsfunc",
			Label:   "jsfunc",
			Options: map[string]interface{}{},
		},
	}

	// register a request body reader
	// since it's a readcloser, it can be read only once
	p.vm.Register("readRequestBody", func(rc io.ReadCloser) string {
		b, _ := io.ReadAll(rc)
		return string(b)
	})

	return
}

func (h processerPayload) String() string {
	return fmt.Sprintf("apigw function %s (%s)", h.Name, h.Label)
}

func (h processerPayload) Type() types.FilterKind {
	return h.Kind
}

func (h processerPayload) Meta() types.FilterMeta {
	return h.FilterMeta
}

func (h processerPayload) Weight() int {
	return h.Wgt
}

func (f *processerPayload) Merge(params []byte) (types.Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&f.params)

	if err != nil {
		return nil, err
	}

	fn, err := base64.StdEncoding.DecodeString(f.params.Func)

	if err != nil {
		return nil, fmt.Errorf("could not decode js func: %s", err)
	}

	f.params.Func = string(fn)

	return f, err
}

func (h processerPayload) Exec(ctx context.Context, scope *types.Scp) (err error) {
	log := h.log.With(zap.String("function", h.String()))

	if h.params.Func == "" {
		err = errors.New("function body empty")
		log.Debug("could not register function", zap.Error(err))
		return
	}

	fn, err := h.vm.RegisterFunction(h.params.Func)

	if err != nil {
		log.Debug("could not register function", zap.Error(err))
		return
	}

	out, err := fn.Exec(h.vm.New(scope))

	if err != nil {
		log.Debug("could not exec payload function", zap.Error(err))
		return
	}

	// add to scope, so next steps can get the structure
	scope.Set("payload", out)

	// check if string
	switch out.(type) {
	case string:
		// handling the newline, to keep the consistency with the json encoder
		// which automatically appends the newline
		_, err = scope.Writer().Write([]byte(fmt.Sprintf("%s\n", out)))
	default:
		err = json.NewEncoder(scope.Writer()).Encode(out)
	}

	if err != nil {
		log.Debug("could not write to body", zap.Error(err))
		return
	}

	return
}

func (h processerPayload) VM() jsenv.Vm {
	return h.vm
}
