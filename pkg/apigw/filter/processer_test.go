package filter

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
				params: prepareFuncPayload(`
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
				params: prepareFuncPayload(`
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
				params: prepareFuncPayload(``),
				err:    `function body empty`,
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				ctx = context.Background()
				req = require.New(t)
			)

			pp := NewPayload(zap.NewNop())
			pp.Merge([]byte(tc.params))

			scope := &types.Scp{
				"request": tc.rq,
				"writer":  httptest.NewRecorder(),
				"opts":    options.Apigw(),
			}

			err := pp.Exec(ctx, scope)

			if tc.err != "" {
				req.EqualError(err, tc.err)
			} else {
				req.NoError(err)
				req.Equal(tc.exp, scope.Writer().(*httptest.ResponseRecorder).Body.String())
			}
		})
	}
}

func prepareFuncPayload(s string) string {
	return fmt.Sprintf(`{"func": "%s"}`, base64.StdEncoding.EncodeToString([]byte(s)))
}
