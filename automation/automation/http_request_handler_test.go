package automation

import (
	"context"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestHttpRequestMaker(t *testing.T) {
	validateBody := func(r *require.Assertions, req *http.Request, expected string) {
		reader, err := req.GetBody()
		r.NoError(err)
		body, err := ioutil.ReadAll(reader)
		r.NoError(err)

		r.Equal(expected, string(body))

	}

	t.Run("basic get", func(t *testing.T) {
		var (
			r = require.New(t)

			req, err = httpRequestHandler{}.makeRequest(context.Background(), &httpRequestSendArgs{
				Url:    "http://localhost/test",
				Method: "GET",
			})
		)

		r.NoError(err)
		r.Equal("GET", req.Method)
		r.Equal("http://localhost/test", req.URL.String())
	})

	t.Run("post form", func(t *testing.T) {
		var (
			r  = require.New(t)
			in = &httpRequestSendArgs{
				Form: url.Values(map[string][]string{
					"a": {"a"},
					"b": {"b", "b"},
					"i": {"42"},
				}),
			}
			req, err = httpRequestHandler{}.makeRequest(context.Background(), in)
		)

		r.NoError(err)
		r.Equal("POST", req.Method)
		validateBody(r, req, "a=a&b=b&b=b&i=42")
	})
}
