package apigw

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"go.uber.org/zap"
)

type (
	processerProxy struct {
		functionMeta
		a   ProxyAuthServicer
		s   SecureStorager
		c   *http.Client
		log *zap.Logger

		params struct {
			Location string          `json:"location"`
			Auth     proxyAuthParams `json:"auth"`
		}
	}
)

func NewProcesserProxy(l *zap.Logger, c *http.Client, s SecureStorager) (p *processerProxy) {
	p = &processerProxy{}

	p.c = c
	p.s = s
	p.log = l

	p.Step = 2
	p.Name = "processerProxy"
	p.Label = "Proxy processer"
	p.Kind = FunctionKindProcesser

	p.Args = []*functionMetaArg{
		{
			Type:    "text",
			Label:   "location",
			Options: map[string]interface{}{},
		},
	}

	return
}

func (h processerProxy) String() string {
	return fmt.Sprintf("apigw function %s (%s)", h.Name, h.Label)
}

func (h processerProxy) Meta() functionMeta {
	return h.functionMeta
}

func (f *processerProxy) Merge(params []byte) (Handler, error) {
	err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&f.params)

	if err != nil {
		return nil, err
	}

	// get the auth mechanism
	f.a, err = NewProxyAuthServicer(f.c, f.params.Auth, f.s)

	if err != nil {
		return nil, fmt.Errorf("could not load auth servicer for proxying: %s", err)
	}

	return f, err
}

func (h processerProxy) Exec(ctx context.Context, scope *scp) (err error) {
	ctx, cancel := context.WithTimeout(ctx, scope.Opts().ProxyOutboundTimeout)
	defer cancel()

	req := scope.Request()
	log := h.log.With(zap.String("ref", h.Name))

	outreq := req.Clone(ctx)

	l, err := url.ParseRequestURI(h.params.Location)

	if err != nil {
		return fmt.Errorf("could not parse destination location for proxying: %s", err)
	}

	outreq.URL = l
	outreq.RequestURI = ""
	outreq.Method = req.Method
	outreq.Host = l.Hostname()

	// use authservicer, set any additional headers
	err = h.a.Do(outreq)

	if err != nil {
		return fmt.Errorf("errors setting auth for proxying: %s", err)
	}

	// merge the old query params to the new request
	// do not overwrite old ones
	// do it after the authServicer, since we also may add them there
	mergeQueryParams(req, outreq)

	if scope.Opts().ProxyEnableDebugLog {
		o, _ := httputil.DumpRequestOut(outreq, false)
		log.Debug("proxy outbound request", zap.Any("request", string(o)))
	}

	// temporary metrics before the proper functionality
	startTime := time.Now()

	// todo - disable / enable follow redirects, already
	// added to options
	resp, err := h.c.Do(outreq)

	if err != nil {
		return fmt.Errorf("could not proxy request: %s", err)
	}

	if scope.Opts().ProxyEnableDebugLog {
		o, _ := httputil.DumpResponse(resp, false)
		log.Debug("proxy outbound response", zap.Any("request", string(o)), zap.Duration("duration", time.Since(startTime)))
	}

	b, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("could not read get body on proxy request: %s", err)
	}

	mergeHeaders(resp.Header, scope.Writer().Header())

	// add to writer
	scope.Writer().Write(b)

	return nil
}
