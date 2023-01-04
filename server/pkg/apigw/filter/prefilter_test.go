package filter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	agctx "github.com/cortezaproject/corteza/server/pkg/apigw/ctx"
	prf "github.com/cortezaproject/corteza/server/pkg/apigw/profiler"
	"github.com/cortezaproject/corteza/server/pkg/apigw/types"
	h "github.com/cortezaproject/corteza/server/pkg/http"
	"github.com/stretchr/testify/require"
)

type (
	tf struct {
		name    string
		expr    string
		err     string
		url     string
		o       string
		headers http.Header
	}
)

func Test_headerMerge(t *testing.T) {
	var (
		tcc = []tf{
			{
				name:    "non matching key",
				expr:    `{"expr":"Foo1 == bar\""}`,
				headers: map[string][]string{"Foo": {"bar"}},
				err:     "could not validate origin parameters: parsing error: Foo1 == bar\"\t:1:12 - 1:13 unexpected String while scanning operator",
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, testMerge(NewHeader(types.Config{}), tc))
	}
}

func Test_headerHandle(t *testing.T) {
	var (
		tcc = []tf{
			{
				name:    "matching simple",
				expr:    `{"expr":"foo == \"bar\""}`,
				headers: map[string][]string{"foo": {"bar"}},
			},
			{
				name:    "matching case",
				expr:    `{"expr":"Foo == \"bar\""}`,
				headers: map[string][]string{"Foo": {"bar"}},
			},
			{
				name:    "non matching value",
				expr:    `{"expr":"Foo == \"bar1\""}`,
				headers: map[string][]string{"Foo": {"bar"}},
				err:     `could not validate headers: validation failed`,
			},
			{
				name:    "non matching key",
				expr:    `{"expr":"Foo1 == \"bar\""}`,
				headers: map[string][]string{"Foo": {"bar"}},
				err:     `could not validate headers: failed to select 'Foo1' on *expr.Vars: no such key 'Foo1'`,
			},
			{
				name:    "regex matching key",
				expr:    `{"expr":"match(Foo, \"^b\\\\wr\\\\s.*$\")"}`,
				headers: map[string][]string{"Foo": {"bar "}},
			},
			{
				name:    "matching header with hyphen",
				expr:    `{"expr":"headers[\"Content-type\"] == \"application/json\""}`,
				headers: map[string][]string{"Content-type": {"application/json"}},
			},
		}
	)

	for _, tc := range tcc {
		r := httptest.NewRequest(http.MethodGet, "/foo", http.NoBody)
		r.Header = tc.headers

		t.Run(tc.name, testHandle(NewHeader(types.Config{}), r, tc))
	}
}

func Test_queryParamMerge(t *testing.T) {
	var (
		tcc = []tf{
			{
				name: "matching simple query parameter",
				expr: `{"expr":"foo == \"bar\""}`,
				url:  "https://examp.le?foo=bar",
			},
			{
				name: "matching simple query parameter - invalid expression key",
				expr: `{"expr1":"foo == \"bar\""}`,
				url:  "https://examp.le?foo=bar",
				err:  "could not validate query parameters: parsing error: 	 - 1:1 unexpected EOF while scanning extensions",
			},
			{
				name: "matching simple query parameter - missing value",
				expr: `{"expr":"foo == \"bar\""}`,
				url:  "https://examp.le?foo=bar1",
			},
			{
				name: "matching simple query parameter - missing value",
				expr: `{"expr":"foo == \"bar-baz\""}`,
				url:  "https://examp.le?foo=bar-baz",
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, testMerge(NewQueryParam(types.Config{}), tc))
	}
}

func Test_queryParamHandle(t *testing.T) {
	var (
		tcc = []tf{
			{
				name: "matching simple query parameter",
				expr: `{"expr":"foo == \"bar\""}`,
				url:  "https://examp.le?foo=bar",
			},
			{
				name: "matching simple query parameter - missing value",
				expr: `{"expr":"foo == \"bar\""}`,
				url:  "https://examp.le?foo=bar1",
				err:  `could not validate query parameters: validation failed`,
			},
			{
				name: "matching query parameter",
				expr: `{"expr":"foo == \"bar-baz\""}`,
				url:  "https://examp.le?foo=bar-baz",
			},
			{
				name: "matching query parameter",
				expr: `{"expr":"params[\"foo-bar\"] == \"bar-baz\""}`,
				url:  "https://examp.le?foo-bar=bar-baz",
			},
		}
	)

	for _, tc := range tcc {
		r := httptest.NewRequest(http.MethodGet, tc.url, http.NoBody)
		t.Run(tc.name, testHandle(NewQueryParam(types.Config{}), r, tc))
	}
}

func Test_profilerHandle_profilerGlobal(t *testing.T) {
	type (
		tfp struct {
			name string
			cfg  types.Config
			r    *http.Request
			exp  *h.Request
		}
	)

	var (
		rr  = httptest.NewRequest(http.MethodGet, "/foo", http.NoBody)
		tcc = []tfp{
			{
				name: "skip profiler hit on profiler global = true",
				cfg:  types.Config{ProfilerGlobal: true},
				r:    rr,
				exp:  nil,
			},
			{
				name: "add profiler hit on profiler global = false",
				cfg:  types.Config{ProfilerGlobal: false},
				r:    rr,
				exp:  createRequest(rr),
			},
		}
	)

	for _, tc := range tcc {
		var (
			req = require.New(t)

			ph  = NewProfiler(tc.cfg)
			hfn = ph.Handler()

			hr  = createRequest(tc.r)
			hit = &prf.Hit{}
		)

		req.Nil(hit.R)

		tc.r = tc.r.WithContext(agctx.ScopeToContext(tc.r.Context(), &types.Scp{"request": hr}))
		tc.r = tc.r.WithContext(agctx.ProfilerToContext(tc.r.Context(), hit))

		err := hfn(httptest.NewRecorder(), tc.r)
		req.NoError(err)

		scoped := agctx.ProfilerFromContext(tc.r.Context())

		req.Equal(tc.exp, scoped.R)
	}
}

func createRequest(r *http.Request) (hr *h.Request) {
	hr, _ = h.NewRequest(r)
	return
}

func testMerge(h types.Handler, tc tf) func(t *testing.T) {
	return func(t *testing.T) {
		var (
			req = require.New(t)
		)

		_, err := h.Merge([]byte(tc.expr), types.Config{})

		if tc.err != "" {
			req.EqualError(err, tc.err)
		} else {
			req.NoError(err)
		}
	}
}

func testHandle(h types.Handler, r *http.Request, tc tf) func(t *testing.T) {
	return func(t *testing.T) {
		var (
			req = require.New(t)
		)

		h, err := h.Merge([]byte(tc.expr), types.Config{})

		req.NoError(err)

		hfn := h.Handler()

		err = hfn(httptest.NewRecorder(), r)

		if tc.err != "" {
			req.EqualError(err, tc.err)
		} else {
			req.NoError(err)
		}
	}
}
