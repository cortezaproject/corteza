package filter

import (
	"context"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	"github.com/stretchr/testify/require"
)

func Test_header(t *testing.T) {
	type (
		tf struct {
			name    string
			expr    string
			err     string
			headers http.Header
		}
	)

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
				err:     "could not validate headers",
			},
			{
				name:    "non matching key",
				expr:    `{"expr":"Foo1 == \"bar\""}`,
				headers: map[string][]string{"Foo": {"bar"}},
				err:     "could not validate headers: failed to select 'Foo1' on *expr.Vars: no such key 'Foo1'",
			},
			{
				name:    "regex matching key",
				expr:    `{"expr":"match(Foo, \"^b\\\\wr\\\\s.*$\")"}`,
				headers: map[string][]string{"Foo": {"bar "}},
			},
			// {
			// 	name:    "matching header with hyphen - TODO",
			// 	expr:    `{"expr":"Content-type == \"application/json\""}`,
			// 	headers: map[string][]string{"Content-type": {"application/json"}},
			// },
		}
	)

	for _, tc := range tcc {
		var (
			ctx = context.Background()
		)

		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			r, err := http.NewRequest(http.MethodGet, "/foo", http.NoBody)
			r.Header = tc.headers

			req.NoError(err)

			scope := &types.Scp{"request": r}

			h := NewHeader()
			h.Merge([]byte(tc.expr))

			err = h.Exec(ctx, scope)

			if tc.err != "" {
				req.EqualError(err, tc.err)
			} else {
				req.NoError(err)
			}
		})
	}
}

func Test_queryParam(t *testing.T) {
	type (
		tf struct {
			name string
			expr string
			err  string
			url  string
		}
	)

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
				err: "could not parse matching expression: parsing error: 	 - 1:1 unexpected EOF while scanning extensions",
			},
			{
				name: "matching simple query parameter - missing value",
				expr: `{"expr":"foo == \"bar\""}`,
				url:  "https://examp.le?foo=bar1",
				err:  "could not validate query params",
			},
			{
				name: "matching simple query parameter - missing value",
				expr: `{"expr":"foo == \"bar-baz\""}`,
				url:  "https://examp.le?foo=bar-baz",
			},
		}
	)

	for _, tc := range tcc {
		var (
			ctx = context.Background()
		)

		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			r, err := http.NewRequest(http.MethodGet, tc.url, http.NoBody)

			req.NoError(err)

			scope := &types.Scp{"request": r}

			h := NewQueryParam()
			h.Merge([]byte(tc.expr))

			err = h.Exec(ctx, scope)

			if tc.err != "" {
				req.EqualError(err, tc.err)
			} else {
				req.NoError(err)
			}
		})
	}
}

func Test_origin(t *testing.T) {
	type (
		tf struct {
			name string
			expr string
			err  string
			o    string
		}
	)

	var (
		tcc = []tf{
			{
				name: "matching simple origin value",
				expr: `{"expr":"origin == \"https://www.google.com\""}`,
				o:    "https://www.google.com",
			},
			{
				name: "matching simple nonexistent origin value",
				expr: `{"expr":"origin == \"https://www.google.com\""}`,
				o:    "",
				err:  "could not validate origin",
			},
			{
				name: "matching simple origin value - invalid expression key",
				expr: `{"expr1":"origin == \"https://www.google.com\""}`,
				o:    "",
				err:  "could not parse matching expression: parsing error: \t - 1:1 unexpected EOF while scanning extensions",
			},
			{
				name: "matching simple origin value - invalid expression key",
				expr: `{"expr1":"origin == \"https"}`,
				o:    "",
				err:  "could not parse matching expression: parsing error: \t - 1:1 unexpected EOF while scanning extensions",
			},
		}
	)

	for _, tc := range tcc {
		var (
			ctx = context.Background()
		)

		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			r, err := http.NewRequest(http.MethodGet, "/foo", http.NoBody)
			r.Header.Set("Origin", tc.o)

			req.NoError(err)

			scope := &types.Scp{"request": r}

			h := NewOrigin()
			h.Merge([]byte(tc.expr))

			err = h.Exec(ctx, scope)

			if tc.err != "" {
				req.EqualError(err, tc.err)
			} else {
				req.NoError(err)
			}
		})
	}
}
