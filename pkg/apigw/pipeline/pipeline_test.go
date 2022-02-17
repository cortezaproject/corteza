package pipeline

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/apigw/pipeline/chain"
	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	mockEmptyHandler = func(rw http.ResponseWriter, r *http.Request) (err error) { return }
)

func NewPl() *Pl {
	return NewPipeline(zap.NewNop(), chain.NewDefault())
}

func Test_pipelineAdd(t *testing.T) {
	var (
		req = require.New(t)
	)

	p := NewPl()
	p.Add(&Worker{
		Handler: mockEmptyHandler,
		Weight:  0,
		Name:    "mockWorker",
	})

	req.Len(p.workers, 1)
}

func Test_pipelineHandleMultiple(t *testing.T) {
	var (
		req = require.New(t)
		rr  = httptest.NewRecorder()
		p   = NewPl()

		first = types.MockHandler{
			Handler_: func(rw http.ResponseWriter, r *http.Request) error {
				rw.Write([]byte(`first`))
				return nil
			},
		}

		second = types.MockHandler{
			Handler_: func(rw http.ResponseWriter, r *http.Request) error {
				rw.Write([]byte(`second`))
				return nil
			},
		}
	)

	p.Add(&Worker{
		Handler: first.Handler(),
		Weight:  5,
		Name:    "mockHandler",
	})

	p.Add(&Worker{
		Handler: second.Handler(),
		Weight:  0,
		Name:    "mockHandler",
	})

	p.Handler().ServeHTTP(rr, &http.Request{})

	req.Equal(`secondfirst`, rr.Body.String())
}

func Test_pipelineExecErr(t *testing.T) {
	type (
		tf struct {
			name string
			mh   types.MockHandler
			w    *Worker
			wgt  int
			exp  string
		}
	)

	var (
		tcc = []tf{
			{
				name: "matching simple",
				mh: types.MockHandler{
					Handler_: func(rw http.ResponseWriter, r *http.Request) error {
						return errors.New("triggered")
					}},
				w:   &Worker{},
				exp: `{"error":{"message":"triggered"}}` + "\n",
			},
			{
				name: "matching simple",
				mh: types.MockHandler{
					Handler_: func(rw http.ResponseWriter, r *http.Request) error {
						rw.Write([]byte(`foobar`))
						return nil
					}},
				w:   &Worker{},
				exp: `foobar`,
			},
		}
	)

	for _, tc := range tcc {
		var (
			p   = NewPl()
			rr  = httptest.NewRecorder()
			req = require.New(t)
		)

		tc.w.Handler = tc.mh.Handler()

		p.Add(tc.w)
		p.Handler().ServeHTTP(rr, &http.Request{})

		req.Equal(tc.exp, rr.Body.String())
	}

}
