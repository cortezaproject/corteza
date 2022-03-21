package filter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cortezaproject/corteza-server/automation/automation"
	atypes "github.com/cortezaproject/corteza-server/automation/types"
	agctx "github.com/cortezaproject/corteza-server/pkg/apigw/ctx"
	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	pe "github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/jsenv"
	"github.com/cortezaproject/corteza-server/pkg/options"
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
		Load(ctx context.Context) error
		Exec(ctx context.Context, workflowID uint64, p atypes.WorkflowExecParams) (*expr.Vars, atypes.Stacktrace, error)
	}

	processerPayload struct {
		types.FilterMeta

		vm  jsenv.Vm
		fn  *jsenv.Fn
		log *zap.Logger

		params struct {
			Func   string `json:"jsfunc"`
			Encode bool   `json:"encode"`
		}
	}
)

func NewWorkflow(opts options.ApigwOpt, wf WfExecer) (p *workflow) {
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

func (h workflow) New(opts options.ApigwOpt) types.Handler {
	return NewWorkflow(opts, h.d)
}

func (h workflow) Enabled() bool {
	return true
}

func (h workflow) String() string {
	return fmt.Sprintf("apigw filter %s (%s)", h.Name, h.Label)
}

func (h workflow) Meta() types.FilterMeta {
	return h.FilterMeta
}

func (h *workflow) Merge(params []byte) (types.Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&h.params)

	if err != nil {
		return h, err
	}

	// preload workflow cache
	return h, h.d.Load(context.Background())
}

func (h workflow) Handler() types.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) error {
		var (
			err   error
			ctx   = r.Context()
			scope = agctx.ScopeFromContext(ctx)
		)
		// cleanup scope for wf
		scp := filterScope(scope, "opts")

		in, err := expr.NewVars(scp.Dict())

		if err != nil {
			return pe.Internal("could not validate request data: %v", err)
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
			return pe.Internal("could not exec workflow: %v", err)
		}

		// merge out with scope
		merged, err := in.Merge(out)

		if err != nil {
			return pe.Internal("could not receive workflow results: %v", err)
		}

		mm, err := expr.CastToVars(merged)

		if err != nil {
			return pe.Internal("could not receive workflow results: %v", err)
		}

		for k, v := range mm {
			scope.Set(k, v)
		}

		scope = filterScope(scope, "eventType", "resourceType", "invoker")

		// update scope for next items in pipeline
		r.WithContext(agctx.ScopeToContext(ctx, scope))

		return nil
	}
}

func NewPayload(opts options.ApigwOpt, l *zap.Logger) (p *processerPayload) {
	p = &processerPayload{}

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
	p.vm.Register("readRequestBody", automation.ReadRequestBody)

	return
}

func (h processerPayload) New(opts options.ApigwOpt) types.Handler {
	return NewPayload(opts, h.log)
}

func (h processerPayload) Enabled() bool {
	return true
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

	h.fn, _ = h.vm.RegisterFunction(h.params.Func)

	return h, err
}

func (h processerPayload) Handler() types.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) (err error) {
		var (
			ctx   = r.Context()
			scope = agctx.ScopeFromContext(ctx)
		)

		// cleanup scope for js
		scp := filterScope(scope, "opts")

		// check fn type
		if err != nil {
			return pe.InvalidData("could not register function: %v", err)
		}

		// need to find a consistent approach to the workflow jsenv function
		// wf: input expr Var and the resulting variable in the jsenv is `input`
		// apigw: input types.Scope and the resulting variable in the jsenv is `input['some_var']`
		out, err := h.fn.Exec(h.vm.New(scp))

		if err != nil {
			return pe.Internal("could not exec payload function: %v", err)
		}

		// add to scope, so next steps can get the structure
		scope.Set("payload", out)

		// check if string
		switch out.(type) {
		case string, []byte:
			// handling the newline, to keep the consistency with the json encoder
			// which automatically appends the newline
			_, err = rw.Write([]byte(fmt.Sprintf("%s\n", out)))
		default:
			err = json.NewEncoder(rw).Encode(out)
		}

		if err != nil {
			return pe.Internal("could not write to response body: %v", err)
		}

		return
	}
}

func (h processerPayload) VM() jsenv.Vm {
	return h.vm
}

func filterScope(scope *types.Scp, kk ...string) (s *types.Scp) {
	s = scope.Filter(func(k string, v interface{}) bool {
		for _, v := range kk {
			if k == v {
				return false
			}
		}

		return true
	})

	return
}
