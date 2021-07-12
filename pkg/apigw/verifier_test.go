package apigw

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_verifierQueryParam(t *testing.T) {
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

			scope := &scp{"request": r}

			h := NewVerifierQueryParam()
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

func Test_verifierOrigin(t *testing.T) {
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

			scope := &scp{"request": r}

			h := NewVerifierOrigin()
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
