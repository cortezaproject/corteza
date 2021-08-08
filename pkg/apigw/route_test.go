package apigw

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/apigw/pipeline"
	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_pl(t *testing.T) {
	type (
		tf struct {
			name       string
			handler    pipeline.Worker
			method     string
			errHandler types.ErrorHandler
			expStatus  int
			expError   error
		}
	)

	var (
		tcc = []tf{
			{
				name: "successful exec",
				handler: types.MockExecer{
					Exec_: func(c context.Context, s *types.Scp) (err error) {
						s.Writer().WriteHeader(http.StatusTemporaryRedirect)
						return
					},
				},
				errHandler: types.MockErrorExecer{
					Exec_: func(c context.Context, s *types.Scp, e error) {
						s.Writer().Write([]byte(e.Error()))
					},
				},
				method:    "POST",
				expStatus: http.StatusTemporaryRedirect,
				expError:  nil,
			},
			{
				name: "unsuccessful exec",
				handler: types.MockExecer{
					Exec_: func(c context.Context, s *types.Scp) (err error) {
						s.Writer().WriteHeader(http.StatusTemporaryRedirect)
						return errors.New("test error")
					},
				},
				errHandler: types.MockErrorExecer{
					Exec_: func(c context.Context, s *types.Scp, e error) {
						s.Writer().WriteHeader(http.StatusInternalServerError)
						s.Writer().Write([]byte(e.Error()))
					},
				},
				method:    "POST",
				expStatus: http.StatusTemporaryRedirect,
				expError:  errors.New("test error"),
			},
			{
				name: "request method validation fail",
				handler: types.MockExecer{
					Exec_: func(c context.Context, s *types.Scp) (err error) {
						s.Writer().WriteHeader(http.StatusTemporaryRedirect)
						return errors.New("test error")
					},
				},
				errHandler: types.MockErrorExecer{
					Exec_: func(c context.Context, s *types.Scp, e error) {
						s.Writer().WriteHeader(http.StatusInternalServerError)
						s.Writer().Write([]byte(e.Error()))
					},
				},
				method:    "GET",
				expStatus: http.StatusInternalServerError,
				expError:  errors.New("could not validate request: invalid method POST"),
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req  = require.New(t)
				rr   = httptest.NewRecorder()
				pipe = pipeline.NewPipeline(zap.NewNop())
			)

			r, err := http.NewRequest("POST", "/foo", http.NoBody)
			req.NoError(err)

			pipe.Add(tc.handler)
			pipe.ErrorHandler(tc.errHandler)

			route := &route{
				method: tc.method,
				pipe:   pipe,
				log:    zap.NewNop(),
				opts:   options.Apigw(),
			}

			route.ServeHTTP(rr, r)

			expError := ""
			if tc.expError != nil {
				expError = tc.expError.Error()
			}

			req.Equal(tc.expStatus, rr.Result().StatusCode)
			req.Equal(expError, rr.Body.String())
		})
	}
}
