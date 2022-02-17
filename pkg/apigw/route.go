package apigw

import (
	"fmt"
	"net/http"
	"time"

	actx "github.com/cortezaproject/corteza-server/pkg/apigw/ctx"
	"github.com/cortezaproject/corteza-server/pkg/apigw/profiler"
	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	h "github.com/cortezaproject/corteza-server/pkg/http"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"go.uber.org/zap"
)

type (
	route struct {
		ID       uint64
		endpoint string
		method   string
		meta     routeMeta

		opts *options.ApigwOpt
		log  *zap.Logger
		pr   *profiler.Profiler

		handler    http.Handler
		errHandler types.ErrorHandlerFunc
	}

	routeMeta struct {
		debug bool
		async bool
	}
)

func (r route) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var (
		ctx   = auth.SetIdentityToContext(req.Context(), auth.ServiceUser())
		start = time.Now()
		scope = types.Scp{}
		hit   = &profiler.Hit{}
	)

	r.log.Debug("started serving route")

	ar, err := h.NewRequest(req)

	if err != nil {
		r.log.Error("could not get initial request", zap.Error(err))
	}

	scope.Set("opts", r.opts)
	scope.Set("request", ar)

	// use profiler, override any profiling prefilter
	if r.opts.ProfilerEnabled {
		// add request to profiler
		hit = r.pr.Hit(ar)
		hit.Route = r.ID
	}

	req = req.WithContext(actx.ScopeToContext(ctx, &scope))
	req = req.WithContext(actx.ProfilerToContext(req.Context(), hit))

	r.handler.ServeHTTP(w, req)

	if r.opts.ProfilerEnabled {
		r.pr.Push(hit)
	} else {
		if hit = actx.ProfilerFromContext(req.Context()).(*profiler.Hit); hit != nil && hit.R != nil {
			// updated hit from a possible prefilter
			// we need to push route ID even if the profiler is disabled
			hit.Route = r.ID
			r.pr.Push(hit)
		}
	}

	r.log.Debug("finished serving route",
		zap.Duration("duration", time.Since(start)),
	)
}

func (r route) String() string {
	return fmt.Sprintf("%s %s", r.method, r.endpoint)
}
