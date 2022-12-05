package filter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	agctx "github.com/cortezaproject/corteza/server/pkg/apigw/ctx"
	"github.com/cortezaproject/corteza/server/pkg/apigw/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/stretchr/testify/require"
)

type (
	mockHandlerRegistry struct{}
)

func Test_redirectionMerge(t *testing.T) {
	var (
		tcc = []tf{
			{
				name: "url validation",
				expr: `{"status":"301", "location": "invalid url"}`,
				err:  `could not validate parameters, invalid URL: parse "invalid url": invalid URI for request`,
			},
			{
				name: "invalid redirection status",
				expr: `{"status":"400", "location": "http://redire.ct/to"}`,
				err:  "could not validate parameters, wrong status 400",
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, testMerge(NewRedirection(options.ApigwOpt{}), tc))
	}
}

func Test_redirection(t *testing.T) {
	type (
		tf struct {
			name string
			expr string
			err  string
			loc  string
			code int
		}
	)

	var (
		tcc = []tf{
			{
				name: "simple redirection",
				expr: `{"status":"302", "location": "http://redire.ct/to"}`,
				loc:  "http://redire.ct/to",
				code: 302,
			},
			{
				name: "permanent redirection",
				expr: `{"status":"301", "location": "http://redire.ct/to"}`,
				loc:  "http://redire.ct/to",
				code: 301,
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req = require.New(t)
				r   = httptest.NewRequest(http.MethodGet, "/foo", http.NoBody)
				rc  = httptest.NewRecorder()
			)

			h := getHandler(NewRedirection(options.ApigwOpt{}))
			h, err := h.Merge([]byte(tc.expr))

			req.NoError(err)

			hn := h.Handler()
			err = hn(rc, r)

			if tc.err != "" {
				req.EqualError(err, tc.err)
				return
			}

			req.NoError(err)
			req.Equal(tc.loc, rc.Header().Get("Location"))
			req.Equal(tc.code, rc.Code)
		})
	}
}

func Test_jsonResponse(t *testing.T) {
	type (
		tf struct {
			name  string
			expr  string
			err   string
			exp   string
			scope interface{}
		}

		aux struct {
			Name    string `json:"name"`
			Surname string `json:"surname"`
		}
	)

	var (
		tcc = []tf{
			{
				name:  "String response as JSON",
				expr:  `{"header":{"content-type":["application/json"]},"input":{"expr": "records", "type": "String"}}`,
				scope: expr.Must(expr.Any{}.Cast("foobar")),
				exp:   `foobar`,
			},
			{
				name:  "Array response as JSON",
				expr:  `{"header":{"content-type":["application/json"]},"input":{"expr": "records", "type": "Array"}}`,
				scope: expr.Must(expr.Any{}.Cast([]float64{3.14, 42.690})),
				exp:   `[3.14,42.69]`,
			},
			{
				name:  "Array response as text",
				expr:  `{"input":{"expr": "records", "type": "Array"}}`,
				scope: expr.Must(expr.Any{}.Cast([]float64{3.14, 42.690})),
				exp:   `[3.14 42.69]`,
			},
			{
				name:  "Any response as JSON",
				expr:  `{"header":{"content-type":["application/json"]},"input": {"expr": "records", "type": "Any"}}`,
				scope: expr.Must(expr.Any{}.Cast(map[string]string{"foo": "bar", "baz": "bzz"})),
				exp:   `{"baz":"bzz","foo":"bar"}`,
			},
			{
				name:  "Any response as text",
				expr:  `{"input": {"expr": "records", "type": "Any"}}`,
				scope: expr.Must(expr.Any{}.Cast(map[string]string{"foo": "bar", "baz": "bzz"})),
				exp:   `map[baz:bzz foo:bar]`,
			},
			{
				name:  "struct array response as JSON",
				expr:  `{"input":{"expr": "toJSON(records)", "type": "String"}}`,
				scope: []aux{{"First", "Last"}, {"Foo", "bar"}},
				exp:   `[{"name":"First","surname":"Last"},{"name":"Foo","surname":"bar"}]`,
			},
			{
				name:  "struct array response as text",
				expr:  `{"input":{"expr": "records", "type": "String"}}`,
				scope: []aux{{"First", "Last"}, {"Foo", "bar"}},
				exp:   `[{First Last} {Foo bar}]`,
			},
			{
				name:  "string csv response as text",
				expr:  `{"header":{"content-type":["application/octet-stream"],"Content-Disposition":["attachment; filename=foo.txt"],"Content-Transfer-Encoding":["binary"]},"input":{"expr": "records", "type": "String"}}`,
				scope: "\"header 1\",\"header 2\"\nvalue 1,value 2\nvalue 3, value 4",
				exp:   "\"header 1\",\"header 2\"\nvalue 1,value 2\nvalue 3, value 4",
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req   = require.New(t)
				r     = httptest.NewRequest(http.MethodGet, "/foo", http.NoBody)
				rc    = httptest.NewRecorder()
				scope = &types.Scp{"records": tc.scope}
			)

			r = r.WithContext(agctx.ScopeToContext(context.Background(), scope))

			h := getHandler(NewResponse(options.ApigwOpt{}, &mockHandlerRegistry{}))
			h, err := h.Merge([]byte(tc.expr))

			req.NoError(err)

			hn := h.Handler()
			err = hn(rc, r)

			if tc.err != "" {
				req.EqualError(err, tc.err)
				return
			}

			req.NoError(err)
			req.Equal(tc.exp, strings.TrimSuffix(rc.Body.String(), "\n"))
		})
	}
}

// hackity hack
func getHandler(h types.Handler) types.Handler {
	return h
}

func (r *mockHandlerRegistry) Type(ref string) expr.Type {
	return expr.Any{}
}
