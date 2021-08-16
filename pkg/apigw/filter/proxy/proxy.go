package proxy

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

	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	"go.uber.org/zap"
)

var (
	hopHeaders = []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	}
)

type (
	proxy struct {
		types.FilterMeta
		a   ProxyAuthServicer
		s   types.SecureStorager
		c   *http.Client
		log *zap.Logger

		params struct {
			Location string          `json:"location"`
			Auth     proxyAuthParams `json:"auth"`
		}
	}
)

func New(l *zap.Logger, c *http.Client, s types.SecureStorager) (p *proxy) {
	p = &proxy{}

	p.c = c
	p.s = s
	p.log = l

	p.Name = "proxy"
	p.Label = "Proxy processer"
	p.Kind = types.Processer

	p.Args = []*types.FilterMetaArg{
		{
			Type:    "text",
			Label:   "location",
			Options: map[string]interface{}{},
		},
	}

	return
}

func (h proxy) String() string {
	return fmt.Sprintf("apigw function %s (%s)", h.Name, h.Label)
}

func (h proxy) Type() types.FilterKind {
	return h.Kind
}

func (h proxy) Meta() types.FilterMeta {
	return h.FilterMeta
}

func (h proxy) Weight() int {
	return h.Wgt
}

func (f *proxy) Merge(params []byte) (types.Handler, error) {
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

func (h proxy) Exec(ctx context.Context, scope *types.Scp) (err error) {
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

func mergeHeaders(orig, dest http.Header) {
OUTER:
	for name, values := range orig {
		// skip headers that need to be omitted
		// when proxying
		for _, v := range hopHeaders {
			if v == name {
				continue OUTER
			}
		}
		dest[name] = values
	}
}

func mergeQueryParams(orig, dest *http.Request) {
	origValues := dest.URL.Query()

	for k, qp := range orig.URL.Query() {
		// skip existing
		if dest.URL.Query().Get(k) != "" {
			continue
		}

		for _, v := range qp {
			origValues.Add(k, v)
		}
	}

	dest.URL.RawQuery = origValues.Encode()
}
