package filter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	agctx "github.com/cortezaproject/corteza/server/pkg/apigw/ctx"
	"github.com/cortezaproject/corteza/server/pkg/apigw/types"
	pe "github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/expr"
)

type (
	header struct {
		types.FilterMeta
		eval expr.Evaluable
		cfg  types.Config

		params struct {
			Expr string `json:"expr"`
		}
	}

	queryParam struct {
		types.FilterMeta
		eval expr.Evaluable
		cfg  types.Config

		params struct {
			Expr string `json:"expr"`
		}
	}

	origin struct {
		types.FilterMeta
		eval expr.Evaluable
		cfg  types.Config

		params struct {
			Expr string `json:"expr"`
		}
	}

	profiler struct {
		cfg types.Config
		types.FilterMeta
	}
)

func NewHeader(cfg types.Config) (v *header) {
	v = &header{}

	v.Name = "header"
	v.Label = "Header"
	v.Kind = types.PreFilter

	v.Args = []*types.FilterMetaArg{
		{
			Type:    "expr",
			Label:   "expr",
			Options: map[string]interface{}{},
		},
	}

	v.cfg = cfg

	return
}

func (h header) New(cfg types.Config) types.Handler {
	return NewHeader(cfg)
}

func (h header) Enabled() bool {
	return true
}

func (h header) String() string {
	return fmt.Sprintf("apigw filter %s (%s)", h.Name, h.Label)
}

func (h header) Meta() types.FilterMeta {
	return h.FilterMeta
}

func (h *header) Merge(params []byte, cfg types.Config) (types.Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&h.params)

	if err != nil {
		return nil, err
	}

	parser := expr.NewParser()

	h.cfg = cfg
	h.eval, err = parser.Parse(h.params.Expr)

	if err != nil {
		return nil, fmt.Errorf("could not validate origin parameters: %s", err)
	}

	return h, err
}

func (h header) Handler() types.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) error {
		var (
			ctx = r.Context()
		)

		vv := map[string]interface{}{}
		headers := r.Header

		for k, v := range headers {
			vv[k] = v[0]
		}

		// get the request data and put it into vars
		out, err := expr.NewVars(vv)
		if err != nil {
			return pe.Internal("could not validate headers: %v", err)
		}

		err = out.Set("headers", vv)
		if err != nil {
			return pe.Internal("could not set headers: %v", err)
		}

		b, err := h.eval.Test(ctx, out)
		if err != nil {
			return pe.InvalidData("could not validate headers: %v", err)
		}

		if !b {
			return pe.InvalidData("could not validate headers: %v", fmt.Errorf("validation failed"))
		}

		return nil
	}
}

func NewQueryParam(cfg types.Config) (v *queryParam) {
	v = &queryParam{}

	v.Name = "queryParam"
	v.Label = "Query parameters"
	v.Kind = types.PreFilter

	v.Args = []*types.FilterMetaArg{
		{
			Type:    "expr",
			Label:   "expr",
			Options: map[string]interface{}{},
		},
	}

	v.cfg = cfg

	return
}

func (qp queryParam) New(cfg types.Config) types.Handler {
	return NewQueryParam(cfg)
}

func (qp queryParam) Enabled() bool {
	return true
}

func (qp queryParam) String() string {
	return fmt.Sprintf("apigw filter %s (%s)", qp.Name, qp.Label)
}

func (qp queryParam) Meta() types.FilterMeta {
	return qp.FilterMeta
}

func (qp *queryParam) Merge(params []byte, cfg types.Config) (types.Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&qp.params)

	if err != nil {
		return nil, err
	}

	parser := expr.NewParser()

	qp.cfg = cfg
	qp.eval, err = parser.Parse(qp.params.Expr)

	if err != nil {
		return nil, fmt.Errorf("could not validate query parameters: %s", err)
	}

	return qp, err
}

func (qp *queryParam) Handler() types.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) error {
		var (
			ctx = r.Context()
		)

		vv := map[string]interface{}{}
		vals := r.URL.Query()

		for k, v := range vals {
			vv[k] = v[0]
		}

		// get the request data and put it into vars
		out, err := expr.NewVars(vv)
		if err != nil {
			return pe.Internal("could not validate query parameters: %v", err)
		}

		err = out.Set("params", vv)
		if err != nil {
			return pe.Internal("could not set params: %v", err)
		}

		b, err := qp.eval.Test(ctx, out)
		if err != nil {
			return pe.InvalidData("could not validate query parameters: %v", err)
		}

		if !b {
			return pe.InvalidData("could not validate query parameters: %v", fmt.Errorf("validation failed"))
		}

		return nil
	}
}

func NewProfiler(cfg types.Config) (pp *profiler) {
	pp = &profiler{}

	pp.cfg = cfg
	pp.Name = "profiler"
	pp.Label = "Profiler"
	pp.Kind = types.PreFilter

	return
}

func (pr *profiler) New(cfg types.Config) types.Handler {
	return NewProfiler(cfg)
}

func (pr *profiler) Enabled() bool {
	return pr.cfg.Profiler.Enabled
}

func (pr *profiler) String() string {
	return fmt.Sprintf("apigw filter %s (%s)", pr.Name, pr.Label)
}

func (pr *profiler) Meta() types.FilterMeta {
	return pr.FilterMeta
}

func (pr *profiler) Merge(params []byte, cfg types.Config) (types.Handler, error) {
	pr.cfg = cfg
	return pr, nil
}

func (pr *profiler) Handler() types.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) (err error) {
		var (
			ctx   = r.Context()
			scope = agctx.ScopeFromContext(ctx)
		)

		if pr.cfg.Profiler.Global {
			// profiler global overrides any profiling prefilter
			// the hit is registered on lower level
			return
		}

		var (
			n   = time.Now()
			hit = agctx.ProfilerFromContext(r.Context())
		)

		if hit == nil {
			return
		}

		hit.Ts = &n
		hit.R = scope.Request()

		r = r.WithContext(agctx.ProfilerToContext(r.Context(), hit))

		return
	}
}
