package apigw

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/cortezaproject/corteza-server/pkg/options"
	"go.uber.org/zap"
)

type (
	route struct {
		ID       uint64
		endpoint string
		method   string

		opts *options.ApigwOpt
		log  *zap.Logger
		pipe *pl
	}
)

func (r route) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var (
		ctx   = req.Context()
		scope = scp{}
	)

	scope.Set("request", req)
	scope.Set("writer", w)
	scope.Set("opts", r.opts)

	if err := r.validate(req); err != nil {
		r.log.Debug("error validating request on route", zap.Error(err))
		r.pipe.err.Exec(ctx, &scope, fmt.Errorf("could not validate request: %s", err))
		return
	}

	if r.opts.LogEnabled {
		o, _ := httputil.DumpRequest(req, false)
		r.log.Debug("incoming request", zap.Any("request", string(o)))
	}

	err := r.pipe.Exec(ctx, &scope)

	if err != nil {
		// call the error handler
		r.log.Debug("calling default error handler on error")
		r.pipe.err.Exec(ctx, &scope, err)
	}
}

func (r route) validate(req *http.Request) (err error) {
	if req.Method != r.method {
		err = fmt.Errorf("invalid method %s", req.Method)
	}

	return
}

func (r route) String() string {
	return fmt.Sprintf("%s %s", r.method, r.endpoint)
}
