package apigw

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_validatorHeader(t *testing.T) {
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
				name:    "matching header with hyphen - TODO",
				expr:    `{"expr":"Content-type == \"application/json\""}`,
				headers: map[string][]string{"Content-type": {"application/json"}},
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
			r.Header = tc.headers

			req.NoError(err)

			scope := &scp{"request": r}

			h := NewValidatorHeader()
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
