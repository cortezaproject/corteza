package filter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	agctx "github.com/cortezaproject/corteza-server/pkg/apigw/ctx"
	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	pe "github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/options"
)

type (
	header struct {
		types.FilterMeta
		eval   expr.Evaluable
		params struct {
			Expr string `json:"expr"`
		}
	}

	queryParam struct {
		types.FilterMeta
		eval   expr.Evaluable
		params struct {
			Expr string `json:"expr"`
		}
	}

	origin struct {
		types.FilterMeta
		eval   expr.Evaluable
		params struct {
			Expr string `json:"expr"`
		}
	}

	profiler struct {
		opts options.ApigwOpt
		types.FilterMeta
	}
)

func NewHeader(opts options.ApigwOpt) (v *header) {
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

	return
}

func (h header) New(opts options.ApigwOpt) types.Handler {
	return NewHeader(opts)
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

func (h *header) Merge(params []byte) (types.Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&h.params)

	if err != nil {
		return nil, err
	}

	parser := expr.NewParser()
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
			return pe.InvalidData("could not validate headers: %v", errors.New("validation failed"))
		}

		return nil
	}
}

func NewQueryParam(opts options.ApigwOpt) (v *queryParam) {
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

	return
}

func (qp queryParam) New(opts options.ApigwOpt) types.Handler {
	return NewQueryParam(opts)
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

func (qp *queryParam) Merge(params []byte) (types.Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&qp.params)

	if err != nil {
		return nil, err
	}

	parser := expr.NewParser()
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
			return pe.InvalidData("could not validate query parameters: %v", errors.New("validation failed"))
		}

		return nil
	}
}

func NewProfiler(opts options.ApigwOpt) (pp *profiler) {
	pp = &profiler{}

	pp.opts = opts
	pp.Name = "profiler"
	pp.Label = "Profiler"
	pp.Kind = types.PreFilter

	return
}

func (pr profiler) New(opts options.ApigwOpt) types.Handler {
	return NewProfiler(opts)
}

func (pr profiler) Enabled() bool {
	return pr.opts.ProfilerEnabled
}

func (pr profiler) String() string {
	return fmt.Sprintf("apigw filter %s (%s)", pr.Name, pr.Label)
}

func (pr profiler) Meta() types.FilterMeta {
	return pr.FilterMeta
}

func (pr *profiler) Merge(params []byte) (types.Handler, error) {
	return pr, nil
}

func (pr *profiler) Handler() types.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) (err error) {
		var (
			ctx   = r.Context()
			scope = agctx.ScopeFromContext(ctx)
		)

		if pr.opts.ProfilerGlobal {
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
