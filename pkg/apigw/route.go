package apigw

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	actx "github.com/cortezaproject/corteza-server/pkg/apigw/ctx"
	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
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

	b, _ := io.ReadAll(req.Body)
	body := string(b)

	// write again
	req.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	scope.Set("opts", r.opts)
	scope.Set("payload", body)

	if r.opts.LogEnabled {
		o, _ := httputil.DumpRequest(req, r.opts.LogRequestBody)
		r.log.Debug("incoming request", zap.Any("request", string(o)))
	}

	req = req.WithContext(actx.ScopeToContext(ctx, &scope))

	r.handler.ServeHTTP(w, req)

	r.log.Debug("finished serving route",
		zap.Duration("duration", time.Since(start)),
	)
}

func (r route) String() string {
	return fmt.Sprintf("%s %s", r.method, r.endpoint)
}
