package filter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	"github.com/stretchr/testify/require"
)

func Tesst_redirection(t *testing.T) {
	type (
		tf struct {
			name string
			expr string
			err  string
		}
	)

	var (
		tcc = []tf{
			{
				name: "simple redirection",
				expr: `{"status":"302", "location": "http://redire.ct/to"}`,
			},
			{
				name: "permanent redirection",
				expr: `{"status":"301", "location": "http://redire.ct/to"}`,
			},
			{
				name: "url validation",
				expr: `{"status":"301", "location": "invalid url"}`,
				err:  `could not redirect: parse "invalid url": invalid URI for request`,
			},
			{
				name: "invalid redirection status",
				expr: `{"status":"400", "location": "http://redire.ct/to"}`,
				err:  "could not redirect: wrong status 400",
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

			req.NoError(err)

			rc := httptest.NewRecorder()
			scope := &types.Scp{"request": r, "writer": rc}

			h := NewRedirection()
			h.Merge([]byte(tc.expr))

			err = h.Exec(ctx, scope)

			if tc.err != "" {
				req.EqualError(err, tc.err)
				return
			}

			req.NoError(err)
			req.Equal(h.params.Location, rc.Header().Get("Location"))
			req.Equal(h.params.HTTPStatus, rc.Code)
		})
	}
}
