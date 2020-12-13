package service

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/stretchr/testify/require"
)

type (
	RoundTripFunc func(req *http.Request) *http.Response

	testProcesser struct{}
)

func (p *testProcesser) Process(ctx context.Context, payload []byte) (ProcesserResponse, error) {
	return 0, nil
}

func TestSyncer_process(t *testing.T) {
	var (
		req    = require.New(t)
		syncer = &Syncer{}

		ctx = context.Background()
	)

	c := make(chan Url, 2)
	u := types.SyncerURI{
		Limit:    10,
		NextPage: "123",
	}
	tp := &testProcesser{}

	go syncer.Process(ctx, []byte(`{"response":{"filter":{"limit":1, "nextPage":"456"}}}`), c, u, tp)

	select {
	case url := <-c:
		req.Equal(url.Url.NextPage, "456")
		break
	}
}

func TestSyncer_parseHeader(t *testing.T) {
	var (
		req    = require.New(t)
		syncer = NewSyncer()

		ctx = context.Background()
	)

	c := make(chan Url, 2)
	u := types.SyncerURI{
		Limit:    10,
		NextPage: "123",
	}
	tp := &testProcesser{}

	n, err := syncer.Process(ctx, []byte(`{"response":{"filt`), c, u, tp)

	req.EqualError(err, "unexpected end of JSON input")
	req.Equal(n, 0)
}

func TestSyncer_response(t *testing.T) {
	var (
		req = require.New(t)
	)

	tests := []struct {
		name   string
		url    string
		expect string
		syncer *Syncer
		ctx    context.Context
	}{
		{
			name: "non 200 request",
			url:  "http://example.ltd",
			syncer: &Syncer{
				client: *NewHttpClient(func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: 500,
						Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":{"message":"common error response"}}`)),
						Header:     make(http.Header),
					}
				}),
			},
			ctx:    context.Background(),
			expect: "invalid return status: 500",
		},
		{
			name: "invalid hostname",
			url:  "http://examp le.ltd",
			syncer: &Syncer{
				client: *NewHttpClient(func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":{"message":"common error response"}}`)),
						Header:     make(http.Header),
					}
				}),
			},
			ctx:    context.Background(),
			expect: "parse \"http://examp le.ltd\": invalid character \" \" in host name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.syncer.Fetch(tt.ctx, tt.url)
			req.EqualError(err, tt.expect)
		})
	}

}

func TestSyncer_authTokenCorrectlySet(t *testing.T) {
	var (
		req               = require.New(t)
		expectedAuthToken = "TEST_JWT_TOKEN"
		actualAuthToken   = ""
	)

	ctx := context.WithValue(
		context.Background(),
		FederationUserToken,
		expectedAuthToken)

	syncer := &Syncer{
		client: *NewHttpClient(func(r *http.Request) *http.Response {

			actualAuthToken = r.Header.Get("Authorization")

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString("OK")),
				Header:     make(http.Header),
			}
		}),
	}

	_, err := syncer.Fetch(ctx, "http://example.ltd")

	req.Equal("Bearer "+expectedAuthToken, actualAuthToken)
	req.NoError(err)
}

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewHttpClient(rt RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(rt),
	}
}
