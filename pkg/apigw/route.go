package apigw

import (
	"fmt"
	"net/http"
	"time"

	actx "github.com/cortezaproject/corteza-server/pkg/apigw/ctx"
	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/system/automation"
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
		scope = types.Scp{}
		start = time.Now()
	)

	r.log.Debug("started serving route")

	// create a new automation HttpRequest
	ar, err := automation.NewHttpRequest(req)

	if err != nil {
		r.log.Error("could not prepare a request holder", zap.Error(err))
		return
	}

	scope.Set("opts", r.opts)
	scope.Set("request", ar)

	req = req.WithContext(actx.ScopeToContext(ctx, &scope))

	r.handler.ServeHTTP(w, req)

	r.log.Debug("finished serving route",
		zap.Duration("duration", time.Since(start)),
	)
}

func (r route) String() string {
	return fmt.Sprintf("%s %s", r.method, r.endpoint)
}
