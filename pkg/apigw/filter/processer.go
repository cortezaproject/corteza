package filter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	atypes "github.com/cortezaproject/corteza-server/automation/types"
	agctx "github.com/cortezaproject/corteza-server/pkg/apigw/ctx"
	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	pe "github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/jsenv"
	"go.uber.org/zap"
)

type (
	workflow struct {
		types.FilterMeta
		d WfExecer

		params struct {
			Workflow uint64 `json:"workflow,string"`
		}
	}

	WfExecer interface {
		Exec(ctx context.Context, workflowID uint64, p atypes.WorkflowExecParams) (*expr.Vars, atypes.Stacktrace, error)
	}

	processerPayload struct {
		types.FilterMeta

		vm  jsenv.Vm
		log *zap.Logger

		params struct {
			Func   string `json:"jsfunc"`
			Encode bool   `json:"encode"`
		}
	}
)

func NewWorkflow(wf WfExecer) (p *workflow) {
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

func (h workflow) New() types.Handler {
	return NewWorkflow(h.d)
}

func (h workflow) String() string {
	return fmt.Sprintf("apigw filter %s (%s)", h.Name, h.Label)
}

func (h workflow) Meta() types.FilterMeta {
	return h.FilterMeta
}

func (f *workflow) Merge(params []byte) (types.Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&f.params)

	return f, err
}

func (h workflow) Handler() types.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) error {
		var (
			err   error
			ctx   = r.Context()
			scope = agctx.ScopeFromContext(ctx)
		)

		payload, err := scope.Get("payload")

		if err != nil {
			return pe.Internal("could not get payload: (%v)", err)
		}

		// setup scope for workflow
		vv := map[string]interface{}{
			"payload": payload,
			"request": r,
		}

		// get the request data and put it into vars
		in, err := expr.NewVars(vv)

		if err != nil {
			return pe.Internal("could not validate request data: (%v)", err)
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
			return pe.Internal("could not exec workflow: (%v)", err)
		}

		// merge out with scope
		merged, err := in.Merge(out)

		if err != nil {
			return pe.Internal("could not receive workflow results: (%v)", err)
		}

		mm, err := expr.CastToVars(merged)

		if err != nil {
			return pe.Internal("could not receive workflow results: (%v)", err)
		}

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

		scope.Set("request", r)
		scope.Set("payload", payload)

		return nil
	}
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

func (h processerPayload) New() types.Handler {
	return NewPayload(h.log)
}

func (h processerPayload) String() string {
	return fmt.Sprintf("apigw filter %s (%s)", h.Name, h.Label)
}

func (h processerPayload) Meta() types.FilterMeta {
	return h.FilterMeta
}

func (h *processerPayload) Merge(params []byte) (types.Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&h.params)

	if err != nil {
		return nil, err
	}

	if h.params.Func == "" {
		return nil, errors.New("could not register function, body empty")
	}

	return h, err
}

func (h processerPayload) Handler() types.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) (err error) {
		var (
			ctx   = r.Context()
			scope = agctx.ScopeFromContext(ctx)
		)

		scope.Set("request", r)

		fn, err := h.vm.RegisterFunction(h.params.Func)

		if err != nil {
			return pe.InvalidData("could not register function: (%v)", err)
		}

		out, err := fn.Exec(h.vm.New(scope))

		if err != nil {
			return pe.Internal("could not exec payload function: (%v)", err)
		}

		// add to scope, so next steps can get the structure
		scope.Set("payload", out)

		// check if string
		switch out.(type) {
		case string:
			// handling the newline, to keep the consistency with the json encoder
			// which automatically appends the newline
			_, err = rw.Write([]byte(fmt.Sprintf("%s\n", out)))
		default:
			err = json.NewEncoder(rw).Encode(out)
		}

		if err != nil {
			return pe.Internal("could not write to response body: (%v)", err)
		}

		return
	}
}

func (h processerPayload) VM() jsenv.Vm {
	return h.vm
}
