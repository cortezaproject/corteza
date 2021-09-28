package filter

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	agctx "github.com/cortezaproject/corteza-server/pkg/apigw/ctx"
	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

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
				name: "payload processer",
				rq: &http.Request{
					Method: "POST",
					Body:   ioutil.NopCloser(strings.NewReader(`[1,2,3]`)),
				},
				exp: "2\n",
				params: prepareFuncPayload(t, `
				var b = JSON.parse(readRequestBody(input.Get('request').Body));
				return b[1];
				`),
			},
			{
				name: "payload processer js map",
				rq: &http.Request{
					Method: "POST",
					Body:   ioutil.NopCloser(strings.NewReader(`[{"name":"johnny", "surname":"mnemonic"},{"name":"johnny", "surname":"knoxville"}]`)),
				},
				exp: "{\"count\":2,\"results\":[{\"fullname\":\"Johnny Mnemonic\"},{\"fullname\":\"Johnny Knoxville\"}]}\n",
				params: prepareFuncPayload(t, `
				var b = JSON.parse(readRequestBody(input.Get('request').Body));

				return {
					"results":
						b.map(function({ name, surname }) {
							return {
								"fullname": name[0].toUpperCase() + name.substring(1) + " " + surname[0].toUpperCase() + surname.substring(1)
							}
						}),
					"count": b.length
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
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req = require.New(t)
				rc  = httptest.NewRecorder()
			)

			pp := NewPayload(zap.NewNop())
			_, err := pp.Merge([]byte(tc.params))

			if tc.errv != "" {
				req.EqualError(err, tc.errv)
				return
			} else {
				req.NoError(err)
			}

			scope := &types.Scp{
				"opts": options.Apigw(),
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
