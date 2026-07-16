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

	req, err := client.Get(server.URL)
	require.True(t, err == nil, "%+v", err)

	resp, err := client.Do(req)
	require.True(t, err == nil, "%+v", err)
	defer resp.Body.Close()

	require.Equal(t, 200, resp.StatusCode)
}
