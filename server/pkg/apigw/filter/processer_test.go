package filter

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	atypes "github.com/cortezaproject/corteza/server/automation/types"
	agctx "github.com/cortezaproject/corteza/server/pkg/apigw/ctx"
	"github.com/cortezaproject/corteza/server/pkg/apigw/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	h "github.com/cortezaproject/corteza/server/pkg/http"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type (
	wfServicer struct {
		load   func(ctx context.Context) error
		exec   func(ctx context.Context, workflowID uint64, p atypes.WorkflowExecParams) (*expr.Vars, uint64, atypes.Stacktrace, error)
		search func(ctx context.Context, filter atypes.WorkflowFilter) (atypes.WorkflowSet, atypes.WorkflowFilter, error)
	}
)

func Test_processerWorkflow(t *testing.T) {
	type (
		tf struct {
			name   string
			err    string
			params string
			wfs    wfServicer
			exp    []string
		}
	)

	var (
		tcc = []tf{
			{
				name:   "workflow processer",
				exp:    []string{"opts", "request"},
				params: `{"workflow":"1"}`,
				err:    `could not exec workflow: mocked error`,
				wfs: wfServicer{
					load: func(ctx context.Context) error {
						return nil
					},
					exec: func(ctx context.Context, workflowID uint64, p atypes.WorkflowExecParams) (*expr.Vars, uint64, atypes.Stacktrace, error) {
						return must(expr.NewVars(map[string]interface{}{"foo": "bar"})), 0, make([]*wfexec.Frame, 0), fmt.Errorf("mocked error")
					},
				},
			},
			{
				name:   "workflow processer",
				exp:    []string{"foo", "request", "opts"},
				params: `{"workflow":"1"}`,
				wfs: wfServicer{
					load: func(ctx context.Context) error {
						return nil
					},
					exec: func(ctx context.Context, workflowID uint64, p atypes.WorkflowExecParams) (*expr.Vars, uint64, atypes.Stacktrace, error) {
						return must(expr.NewVars(map[string]interface{}{"foo": "bar"})), 0, make([]*wfexec.Frame, 0), nil
					},
				},
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req     = require.New(t)
				rc      = httptest.NewRecorder()
				rq, _   = http.NewRequest("POST", "/foo", http.NoBody)
				ar, err = h.NewRequest(rq)
				pp      = NewWorkflow(types.Config{}, tc.wfs)
			)

			_, err = pp.Merge([]byte(tc.params), types.Config{})
			req.NoError(err)

			scope := &types.Scp{
				"opts":    options.Apigw(),
				"request": ar,
			}

			rq = rq.WithContext(agctx.ScopeToContext(context.Background(), scope))

			hn := pp.Handler()
			err = hn(rc, rq)

			if tc.err != "" {
				req.EqualError(err, tc.err)
			} else {
				req.NoError(err)
			}

			req.ElementsMatch(tc.exp, agctx.ScopeFromContext(rq.Context()).Keys())
		})
	}
}

func Test_processerPayload(t *testing.T) {
	type (
		tf struct {
			name   string
			err    string
			errv   string
			params string
			exp    string
			rq     *http.Request
		}
	)

	var (
		tcc = []tf{
			{
				name: "payload processer parse request body",
				rq: &http.Request{
					Method: "POST",
					Body:   ioutil.NopCloser(strings.NewReader(`[1,2,3]`)),
				},
				exp: "2\n",
				params: prepareFuncPayload(t, `
				const b = JSON.parse(readRequestBody(input.Get('request')));
				return b[1];
				`),
			},
			{
				name: "payload processer js map request body",
				rq: &http.Request{
					Method: "POST",
					Body:   ioutil.NopCloser(strings.NewReader(`[{"name":"johnny", "surname":"mnemonic"},{"name":"johnny", "surname":"knoxville"}]`)),
				},
				exp: "{\"count\":2,\"results\":[{\"fullname\":\"Johnny Mnemonic\"},{\"fullname\":\"Johnny Knoxville\"}]}\n",
				params: prepareFuncPayload(t, `
				const readOnce = JSON.parse(readRequestBody(input.Get('request')));
				const readTwice = JSON.parse(readRequestBody(input.Get('request')));

				return {
					"results":
						readTwice.map(function({ name, surname }) {
							return {
								"fullname": name[0].toUpperCase() + name.substring(1) + " " + surname[0].toUpperCase() + surname.substring(1)
							}
						}),
					"count": readTwice.length
				};
				`),
			},
			{
				name: "payload processer empty function",
				rq: &http.Request{
					Method: "POST",
					Body:   ioutil.NopCloser(strings.NewReader(`[{"name":"johnny", "surname":"mnemonic"},{"name":"johnny", "surname":"knoxville"}]`)),
				},
				params: prepareFuncPayload(t, ``),
				errv:   `could not register function, body empty`,
			},
			{
				name: "payload processer invalid body",
				rq: &http.Request{
					Method: "POST",
					Body:   ioutil.NopCloser(strings.NewReader(`[{"name":"johnny", "surname":"mnemonic"},{"name":"johnny", "surname":"knoxville"}]`)),
				},
				params: prepareFuncPayload(t, `.foo`),
				errv:   `could not register function, invalid body: could not transform payload: Unexpected "."`,
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req     = require.New(t)
				rc      = httptest.NewRecorder()
				ar, err = h.NewRequest(tc.rq)

				cfg = types.Config{}
			)

			pp := NewPayload(cfg, zap.NewNop())
			_, err = pp.Merge([]byte(tc.params), cfg)

			if tc.errv != "" {
				req.EqualError(err, tc.errv)
				return
			} else {
				req.NoError(err)
			}

			scope := &types.Scp{
				"opts":    options.Apigw(),
				"request": ar,
			}

			tc.rq = tc.rq.WithContext(agctx.ScopeToContext(context.Background(), scope))

			hn := pp.Handler()
			err = hn(rc, tc.rq)

			if tc.err != "" {
				req.EqualError(err, tc.err)
			} else {
				req.NoError(err)
				req.Equal(tc.exp, rc.Body.String())
			}
		})
	}
}

func prepareFuncPayload(t *testing.T, s string) string {
	aux, err := json.Marshal(map[string]string{"jsfunc": s})
	if err != nil {
		t.Error(err)
	}
	return string(aux)
}

func (f wfServicer) Load(ctx context.Context) error {
	return f.load(ctx)
}

func (f wfServicer) Exec(ctx context.Context, workflowID uint64, p atypes.WorkflowExecParams) (*expr.Vars, uint64, atypes.Stacktrace, error) {
	return f.exec(ctx, workflowID, p)
}

func (f wfServicer) Search(ctx context.Context, filter atypes.WorkflowFilter) (atypes.WorkflowSet, atypes.WorkflowFilter, error) {
	return f.search(ctx, filter)
}

func must(v *expr.Vars, err error) *expr.Vars {
	if err != nil {
		return nil
	}

	return v
}
