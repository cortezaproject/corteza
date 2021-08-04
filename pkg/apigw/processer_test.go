package apigw

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_processerProxy(t *testing.T) {
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
			fn     func(*require.Assertions) mockRoundTripper
		}
	)

	var (
		tcc = []tf{
			{
				name: "proxy processer with auth headers",
				fn: func(req *require.Assertions) mockRoundTripper {
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
				fn: func(req *require.Assertions) mockRoundTripper {
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
				fn: func(req *require.Assertions) mockRoundTripper {
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
				err:    `could not parse destination location for proxying: parse "invalid url": invalid URI for request`,
			},
			{
				name: "proxy processer params request error",
				fn: func(req *require.Assertions) mockRoundTripper {
					return func(r *http.Request) (rs *http.Response, err error) {
						err = fmt.Errorf("error on client.Do")
						return
					}
				},
				params: `{"location": "https://example.com", "auth": {"type": "header", "params": {}}}`,
				err:    `could not proxy request: Post "https://example.com": error on client.Do`,
			},
			{
				name: "proxy processer hop headers removed",
				fn: func(req *require.Assertions) mockRoundTripper {
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
				fn: func(req *require.Assertions) mockRoundTripper {
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
				ctx = context.Background()
				req = require.New(t)
				c   = http.DefaultClient
				rq  = tc.rq
			)

			if tc.fn != nil {
				c.Transport = mockRoundTripper(tc.fn(req))
			}

			if rq == nil {
				rq, _ = http.NewRequest("POST", "/foo", strings.NewReader(`custom request body`))
			}

			proxy := NewProcesserProxy(zap.NewNop(), c, secureStorageTodo{})
			proxy.Merge([]byte(tc.params))

			scope := &scp{
				"request": rq,
				"writer":  httptest.NewRecorder(),
				"opts":    options.Apigw(),
			}

			err := proxy.Exec(ctx, scope)

			if tc.err != "" {
				req.EqualError(err, tc.err)
			} else {
				req.NoError(err)
				req.Equal(tc.exp.Header, scope.Writer().(*httptest.ResponseRecorder).Header())
				req.Equal(tc.exp.Body, scope.Writer().(*httptest.ResponseRecorder).Body)
			}
		})
	}
}

func Test_processerPayload(t *testing.T) {
	type (
		tf struct {
			name   string
			err    string
			params string
			exp    string
			rq     *http.Request
		}
	)

	var (
		tcc = []tf{
			{
				name: "payload processer",
				rq: &http.Request{
					Method: "POST",
					Body:   ioutil.NopCloser(strings.NewReader(`[1,2,3]`)),
				},
				exp: "2\n",
				params: prepareFuncPayload(`
				var b = JSON.parse(readRequestBody(input.Get('request').Body));
				return b[1];
				`),
			},
			{
				name: "payload processer js map",
				rq: &http.Request{
					Method: "POST",
					Body:   ioutil.NopCloser(strings.NewReader(`[{"name":"johnny", "surname":"mnemonic"},{"name":"johnny", "surname":"knoxville"}]`)),
				},
				exp: "{\"count\":2,\"results\":[{\"fullname\":\"Johnny Mnemonic\"},{\"fullname\":\"Johnny Knoxville\"}]}\n",
				params: prepareFuncPayload(`
				var b = JSON.parse(readRequestBody(input.Get('request').Body));

				return {
					"results":
						b.map(function({ name, surname }) {
							return {
								"fullname": name[0].toUpperCase() + name.substring(1) + " " + surname[0].toUpperCase() + surname.substring(1)
							}
						}),
					"count": b.length
				};
				`),
			},
			{
				name: "payload processer empty function",
				rq: &http.Request{
					Method: "POST",
					Body:   ioutil.NopCloser(strings.NewReader(`[{"name":"johnny", "surname":"mnemonic"},{"name":"johnny", "surname":"knoxville"}]`)),
				},
				params: prepareFuncPayload(``),
				err:    `function body empty`,
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				ctx = context.Background()
				req = require.New(t)
			)

			pp := NewProcesserPayload(zap.NewNop())
			pp.Merge([]byte(tc.params))

			scope := &scp{
				"request": tc.rq,
				"writer":  httptest.NewRecorder(),
				"opts":    options.Apigw(),
			}

			err := pp.Exec(ctx, scope)

			if tc.err != "" {
				req.EqualError(err, tc.err)
			} else {
				req.NoError(err)
				req.Equal(tc.exp, scope.Writer().(*httptest.ResponseRecorder).Body.String())
			}
		})
	}
}

func prepareFuncPayload(s string) string {
	return fmt.Sprintf(`{"func": "%s"}`, base64.StdEncoding.EncodeToString([]byte(s)))
}
