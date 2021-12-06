package filter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cortezaproject/corteza-server/automation/service"
	agctx "github.com/cortezaproject/corteza-server/pkg/apigw/ctx"
	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/stretchr/testify/require"
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
		t.Run(tc.name, testMerge(NewRedirection(), tc))
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

			h := getHandler(NewRedirection())
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
				name:  "Any response as JSON",
				expr:  `{"expr": "records", "type": "KV"}`,
				scope: expr.Must(expr.Any{}.Cast([]float64{3.14, 42.690})),
				exp:   `[3.14,42.69]`,
			},
			{
				name:  "KV response as JSON",
				expr:  `{"expr": "records", "type": "KV"}`,
				scope: map[string]string{"foo": "bar", "baz": "bzz"},
				exp:   `{"baz":"bzz","foo":"bar"}`,
			},
			{
				name:  "struct array response as JSON",
				expr:  `{"expr": "toJSON(records)", "type": "String"}`,
				scope: []aux{{"First", "Last"}, {"Foo", "bar"}},
				exp:   `[{"name":"First","surname":"Last"},{"name":"Foo","surname":"bar"}]`,
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

			h := getHandler(NewJsonResponse(service.Registry()))
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
