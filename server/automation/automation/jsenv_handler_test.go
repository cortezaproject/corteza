package automation

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/expr"
	h "github.com/cortezaproject/corteza/server/pkg/http"
	"github.com/stretchr/testify/require"
)

func Test_jsenvHandler(t *testing.T) {
	type (
		exp struct {
			s string
			i int64
			a interface{}
		}

		tf struct {
			name   string
			exp    *exp
			err    error
			params *jsenvExecuteArgs
		}
	)

	var (
		handler = &jsenvHandler{}
		tcc     = []tf{
			{
				name: "jsenv handler check payload",
				err:  errors.New(`could not process payload, scope missing`),
				params: &jsenvExecuteArgs{
					hasScope:  false,
					hasSource: true,
				},
			},
			{
				name: "jsenv handler check payload",
				err:  errors.New(`could not process payload, function missing`),
				params: &jsenvExecuteArgs{
					hasScope:  true,
					hasSource: false,
				},
			},
			{
				name: "jsenv handler invalid function",
				err:  errors.New(`could not register jsenv function: SyntaxError: SyntaxError: (anonymous): Line 1:74 Unexpected token function (and 1 more errors)`),
				params: &jsenvExecuteArgs{
					hasScope:  true,
					hasSource: true,
					Source:    `invalid function here...`,
				},
			},
			{
				name: "jsenv handler check payload",
				err:  errors.New(``),
				params: &jsenvExecuteArgs{
					hasScope:  true,
					hasSource: true,
					Source:    `return nonexistent`,
				},
			},
			{
				name: "jsenv handler parse request body json to string",
				exp: &exp{
					s: `bar`,
					i: 0,
				},
				params: &jsenvExecuteArgs{
					hasScope:  true,
					hasSource: true,
					Scope:     mustAny(expr.NewAny(makeRequest(t, `{"foo":"bar"}`))),
					Source:    `const b = JSON.parse(readRequestBody(input)); return b.foo;`,
				},
			},
			{
				name: "jsenv handler parse request body json to int",
				exp: &exp{
					s: ``,
					i: 42,
				},
				params: &jsenvExecuteArgs{
					hasScope:  true,
					hasSource: true,
					Scope:     mustAny(expr.NewAny(makeRequest(t, `{"foo":42}`))),
					Source:    `const b = JSON.parse(readRequestBody(input)); return b.foo;`,
				},
			},
			{
				name: "jsenv handler parse request body json to int",
				exp: &exp{
					s: ``,
					i: 42,
				},
				params: &jsenvExecuteArgs{
					hasScope:  true,
					hasSource: true,
					Scope:     mustAny(expr.NewAny(makeRequest(t, `42`))),
					Source:    `const b = readRequestBody(input); return parseInt(b);`,
				},
			},
			{
				name: "jsenv handler parse request body try catch",
				exp: &exp{
					s: `caught`,
					i: 0,
				},
				params: &jsenvExecuteArgs{
					hasScope:  true,
					hasSource: true,
					Scope:     mustAny(expr.NewAny(makeRequest(t, `42`))),
					Source:    `try { const b = readRequestBody(input_NONEXISTENT); } catch (e) { return 'caught'; }`,
				},
			},
			{
				name: "jsenv handler parse request body json to float",
				exp: &exp{
					s: ``,
					i: 0,
					a: 42.690,
				},
				params: &jsenvExecuteArgs{
					hasScope:  true,
					hasSource: true,
					Scope:     mustAny(expr.NewAny(makeRequest(t, `42.690`))),
					Source:    `const b = readRequestBody(input); return parseFloat(b);`,
				},
			},
			{
				name: "jsenv handler input scope",
				exp: &exp{
					s: ``,
					i: 41,
				},
				params: &jsenvExecuteArgs{
					hasScope:  true,
					hasSource: true,
					Scope:     mustAny(expr.NewAny(makeRequest(t, `42`))),
					Source:    `const b = readRequestBody(input); return b - 1;`,
				},
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req = require.New(t)
				ctx = context.Background()
			)

			out, err := handler.execute(ctx, tc.params)

			if tc.err == nil {
				req.NoError(err)
			} else {
				req.Error(err)
			}

			if tc.exp != nil {
				req.Equal(tc.exp.s, out.ResultString)
				req.Equal(tc.exp.i, out.ResultInt)

				if tc.exp.a != nil {
					req.Equal(tc.exp.a, out.ResultAny)
				}
			}
		})
	}
}

func makeRequest(t *testing.T, b string) *h.Request {
	r, err := http.NewRequest("POST", "/foo", ioutil.NopCloser(strings.NewReader(b)))

	if err != nil {
		t.Error(err)
	}

	ar, err := h.NewRequest(r)

	if err != nil {
		t.Error(err)
	}

	return ar
}

func mustAny(v *expr.Any, err error) *expr.Any {
	if err != nil {
		panic(err)
	}
	return v
}
