package http

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPClient(t *testing.T) {
	handler := &Fortune{}
	server := httptest.NewServer(handler)
	defer server.Close()

	client, err := New(&Config{
		Timeout: 10,
	})
	require.True(t, err == nil, "%+v", err)
	client.Debug(FULL)

	req, err := client.Get(server.URL)
	require.True(t, err == nil, "%+v", err)

	resp, err := client.Do(req)
	require.True(t, err == nil, "%+v", err)

	err = func() error {
		defer resp.Body.Close()
		switch resp.StatusCode {
		case 200:
			return nil
		default:
			return ToError(resp)
		}
	}()

	require.True(t, err == nil, "Invalid response: %+v", err)
}
