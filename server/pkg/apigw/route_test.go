package apigw

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/apigw/pipeline"
	"github.com/cortezaproject/corteza/server/pkg/apigw/pipeline/chain"
	"github.com/cortezaproject/corteza/server/pkg/apigw/types"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_pl(t *testing.T) {
	type (
		tf struct {
			name       string
			method     string
			endpoint   string
			expError   string
			expStatus  int
			handler    *types.MockHandler
			errHandler *types.MockErrorHandler
		}
	)

	var (
		tcc = []tf{
			{
				name: "successful handler",
				handler: &types.MockHandler{
					Handler_: func(rw http.ResponseWriter, r *http.Request) error {
						rw.WriteHeader(http.StatusTemporaryRedirect)
						return nil
					},
				},
				method:    "POST",
				expStatus: http.StatusTemporaryRedirect,
				expError:  "",
			},
			{
				name: "unsuccessful handle custom error response",
				handler: &types.MockHandler{
					Handler_: func(rw http.ResponseWriter, r *http.Request) error {
						rw.WriteHeader(http.StatusTemporaryRedirect)
						return errors.New("test error")
					},
				},
				errHandler: &types.MockErrorHandler{
					Handler_: func(rw http.ResponseWriter, r *http.Request, err error) {
						rw.Write([]byte("custom error response: " + err.Error()))
					},
				},
				method:    "POST",
				expStatus: http.StatusTemporaryRedirect,
				expError:  "custom error response: test error",
			},
			{
				name: "unsuccessful handle default error response",
				handler: &types.MockHandler{
					Handler_: func(rw http.ResponseWriter, r *http.Request) error {
						rw.WriteHeader(http.StatusTemporaryRedirect)
						return errors.New("test error")
					},
				},
				method:    "POST",
				expStatus: http.StatusTemporaryRedirect,
				expError:  "{\"error\":{\"message\":\"test error\"}}\n",
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req  = require.New(t)
				rr   = httptest.NewRecorder()
				pipe = pipeline.NewPipeline(zap.NewNop(), chain.NewDefault())
			)

			r := httptest.NewRequest("POST", "/foo", http.NoBody)

			pipe.Add(&pipeline.Worker{
				Handler: tc.handler.Handler(),
			})

			if tc.errHandler != nil {
				pipe.ErrorHandler(tc.errHandler.Handler())
			}

			route := &route{
				method:     tc.method,
				endpoint:   tc.endpoint,
				log:        zap.NewNop(),
				opts:       *options.Apigw(),
				handler:    pipe.Handler(),
				errHandler: pipe.Error(),
			}

			route.ServeHTTP(rr, r)

			req.Equal(tc.expStatus, rr.Result().StatusCode)
			req.Equal(tc.expError, rr.Body.String())
		})
	}
}
