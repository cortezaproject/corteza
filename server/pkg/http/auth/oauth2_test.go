package auth

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_oauth2(t *testing.T) {
	type (
		tf struct {
			name   string
			err    string
			exp    string
			params Oauth2Params
		}
	)

	var (
		tcc = []tf{
			{
				name:   "match oauth2 fail client validation",
				err:    "invalid param client",
				params: Oauth2Params{},
			},
			{
				name: "match oauth2 fail secret key validation",
				err:  "invalid param secret",
				params: Oauth2Params{
					Client: "client_ID",
				},
			},
			{
				name: "match oauth2 fail url validation",
				err:  "invalid param token url",
				params: Oauth2Params{
					Client: "client_ID",
					Secret: "secret_KEY",
				},
			},
			{
				name: "match oauth2 fail url validation",
				err:  "invalid param token url",
				params: Oauth2Params{
					Client:   "client_ID",
					Secret:   "secret_KEY",
					TokenUrl: &url.URL{},
				},
			},
			{
				name: "match oauth2 fail validation",
				params: Oauth2Params{
					Client:   "client_ID",
					Secret:   "secret_KEY",
					TokenUrl: generateURL("http://example.com"),
				},
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req    = require.New(t)
				c      = http.DefaultClient
				_, err = NewOauth2(tc.params, c, struct{}{})
			)

			if tc.err != "" {
				req.EqualError(err, tc.err)
				return
			}
		})
	}
}

func generateURL(s string) (u *url.URL) {
	u, _ = url.Parse(s)
	return
}
