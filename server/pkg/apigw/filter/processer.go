package filter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cortezaproject/corteza/server/automation/automation"
	atypes "github.com/cortezaproject/corteza/server/automation/types"
	agctx "github.com/cortezaproject/corteza/server/pkg/apigw/ctx"
	"github.com/cortezaproject/corteza/server/pkg/apigw/types"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	pe "github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/jsenv"
	"go.uber.org/zap"
)

type (
	workflow struct {
		types.FilterMeta
		d   WfExecer
		cfg types.Config

		params struct {
			Workflow uint64 `json:"workflow,string"`
		}
	}

	WfExecer interface {
		Load(ctx context.Context) error
		Exec(ctx context.Context, workflowID uint64, p atypes.WorkflowExecParams) (*expr.Vars, uint64, atypes.Stacktrace, error)
		Search(ctx context.Context, filter atypes.WorkflowFilter) (atypes.WorkflowSet, atypes.WorkflowFilter, error)
	}

	processerPayload struct {
		types.FilterMeta

		vm  jsenv.Vm
		fn  *jsenv.Fn
		log *zap.Logger
		cfg types.Config

		params struct {
			Func   string `json:"jsfunc"`
			Encode bool   `json:"encode"`
		}
	}
)

func NewWorkflow(cfg types.Config, wf WfExecer) (p *workflow) {
	p = &workflow{}

	p.d = wf
	p.cfg = cfg

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

func (h workflow) New(cfg types.Config) types.Handler {
	return NewWorkflow(cfg, h.d)
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

func (h *workflow) Merge(params []byte, cfg types.Config) (types.Handler, error) {
	var (
		t = struct {
			Workflow string `json:"workflow"`
		}{}

		ctx = auth.SetIdentityToContext(context.Background(), auth.ServiceUser())
	)

	h.cfg = cfg

	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&t)
	if err != nil {
		return h, err
	}

	i, err := strconv.Atoi(t.Workflow)

	if err == nil {
		h.params.Workflow = uint64(i)
		return h, h.d.Load(ctx)
	}

	if wf, _, err := h.d.Search(ctx, atypes.WorkflowFilter{Query: t.Workflow}); err != nil {
		return h, err
	} else {
		h.params.Workflow = wf[0].ID
	}

	// preload workflow cache
	return h, h.d.Load(ctx)
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

		out, _, _, err := h.d.Exec(ctx, h.params.Workflow, wp)

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
		_ = r.WithContext(agctx.ScopeToContext(ctx, scope))

		return nil
	}
}

func NewPayload(cfg types.Config, l *zap.Logger) (p *processerPayload) {
	p = &processerPayload{}

	p.vm = jsenv.New(jsenv.NewTransformer(jsenv.LoaderJS, jsenv.TargetES2016))
	p.log = l
	p.cfg = cfg

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

func (h processerPayload) New(cfg types.Config) types.Handler {
	return NewPayload(cfg, h.log)
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

func (h *processerPayload) Merge(params []byte, cfg types.Config) (types.Handler, error) {
	h.cfg = cfg

	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&h.params)

	if err != nil {
		return nil, err
	}

	if h.params.Func == "" {
		return nil, fmt.Errorf("could not register function, body empty")
	}

	if h.fn, err = h.vm.RegisterFunction(h.params.Func); err != nil {
		return nil, fmt.Errorf("could not register function, invalid body: %s", err)
	}

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
