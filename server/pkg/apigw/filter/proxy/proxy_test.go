package proxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	agctx "github.com/cortezaproject/corteza/server/pkg/apigw/ctx"
	"github.com/cortezaproject/corteza/server/pkg/apigw/types"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_proxy(t *testing.T) {
	type (
		exp struct {
			Status int
			Header http.Header
			Body   *bytes.Buffer
		}

		tf struct {
			name   string
			err    string
			params string
			exp    exp
			rq     *http.Request
			fn     func(*require.Assertions) types.MockRoundTripper
		}
	)

	var (
		tcc = []tf{
			{
				name: "proxy processer with auth headers",
				fn: func(req *require.Assertions) types.MockRoundTripper {
					return func(r *http.Request) (rs *http.Response, err error) {
						rs = &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader("default response")),
						}

						return
					}
				},
				params: `{"location": "/foo", "auth": {"type": "header", "params": {"access-token": "123", "client": "456"}}}`,
				exp: exp{
					Status: http.StatusOK,
					Header: http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}},
					Body:   bytes.NewBufferString("default response"),
				},
			},
			{
				name: "proxy processer with auth query params",
				fn: func(req *require.Assertions) types.MockRoundTripper {
					return func(r *http.Request) (rs *http.Response, err error) {
						rs = &http.Response{}
						req.Equal("access-param=123%2B456", r.URL.RawQuery)
						return
					}
				},
				params: `{"location": "/foo", "auth": {"type": "query", "params": {"access-param": "123+456"}}}`,
				exp: exp{
					Status: http.StatusOK,
					Header: http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}},
					Body:   bytes.NewBuffer(nil),
				},
			},
			{
				name: "proxy processer with auth headers unauthorized",
				fn: func(req *require.Assertions) types.MockRoundTripper {
					return func(r *http.Request) (rs *http.Response, err error) {
						rs = &http.Response{
							StatusCode: http.StatusUnauthorized,
							Body:       io.NopCloser(strings.NewReader("unauthorized response")),
						}

						return
					}
				},
				params: `{"location": "/foo", "auth": {"type": "header", "params": {"access-token": "123", "client": "456"}}}`,
				exp: exp{
					Status: http.StatusUnauthorized,
					Header: http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}},
					Body:   bytes.NewBufferString("unauthorized response"),
				},
			},
			{
				name:   "proxy processer params parse error",
				params: `{"location": "invalid url", "auth": {"type": "header", "params": {}}}`,
				err:    `could not parse destination location for proxying: (parse "invalid url": invalid URI for request)`,
			},
			{
				name: "proxy processer params request error",
				fn: func(req *require.Assertions) types.MockRoundTripper {
					return func(r *http.Request) (rs *http.Response, err error) {
						err = fmt.Errorf("error on client.Do")
						return
					}
				},
				params: `{"location": "https://example.com", "auth": {"type": "header", "params": {}}}`,
				err:    `could not proxy request: (Post "https://example.com": error on client.Do)`,
			},
			{
				name: "proxy processer hop headers removed",
				fn: func(req *require.Assertions) types.MockRoundTripper {
					return func(r *http.Request) (rs *http.Response, err error) {
						rs = &http.Response{
							Header: http.Header{
								"Proxy-Authenticate": []string{`Basic realm="Access to the internal site"`},
								"Content-Type":       []string{"application/json; charset=utf-8"},
							},
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader("default response")),
						}

						return
					}
				},
				params: `{"location": "https://example.com", "auth": {"type": "header", "params": {}}}`,
				exp: exp{
					Status: http.StatusUnauthorized,
					Header: http.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
					Body:   bytes.NewBufferString("default response"),
				},
			},
			{
				name: "proxy processer query parameters merged",
				fn: func(req *require.Assertions) types.MockRoundTripper {
					return func(r *http.Request) (rs *http.Response, err error) {
						rs = &http.Response{}
						req.Equal("access-param=123%2B456&addedCustomQueryParam=true", r.URL.RawQuery)
						return
					}
				},
				params: `{"location": "https://example.com", "auth": {"type": "query", "params": {"access-param": "123+456"}}}`,
				exp: exp{
					Status: http.StatusUnauthorized,
					Header: http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}},
					Body:   bytes.NewBuffer(nil),
				},
				rq: &http.Request{
					Header: http.Header{},
					URL:    &url.URL{Path: "/foo", RawQuery: "addedCustomQueryParam=true"},
					Body:   http.NoBody,
					Method: "POST",
				},
			},
		}
	)

	for _, tc := range tcc {

		t.Run(tc.name, func(t *testing.T) {
			var (
				rc  = httptest.NewRecorder()
				req = require.New(t)
				c   = http.DefaultClient
				rq  = tc.rq
				cfg = types.Config{}
			)

			if tc.fn != nil {
				c.Transport = types.MockRoundTripper(tc.fn(req))
			}

			if rq == nil {
				rq = httptest.NewRequest("POST", "/foo", strings.NewReader(`custom request body`))
			}

			proxy := New(cfg, zap.NewNop(), c, struct{}{})
			_, err := proxy.Merge([]byte(tc.params), cfg)
			req.NoError(err)

			scope := &types.Scp{
				"opts": options.Apigw(),
			}

			ctx := agctx.ScopeToContext(context.Background(), scope)
			rq = rq.WithContext(ctx)

			hn := proxy.Handler()
			err = hn(rc, rq)

			if tc.err != "" {
				req.EqualError(err, tc.err)
			} else {
				req.NoError(err)
				req.Equal(tc.exp.Header, rc.Header())
				req.Equal(tc.exp.Body, rc.Body)
			}
		})
	}
}
