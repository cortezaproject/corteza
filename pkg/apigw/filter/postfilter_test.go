package filter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
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

// hackity hack
func getHandler(h types.Handler) types.Handler {
	return h
}
