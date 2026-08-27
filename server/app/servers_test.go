package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	city311Contract "github.com/cortezaproject/corteza/server/compose/types/city311"
	"github.com/stretchr/testify/require"
)

func TestHealthzReportsContractErrorUntilDatabaseIsAvailable(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	(&CortezaApp{}).healthz(response, request)
	require.Equal(t, http.StatusServiceUnavailable, response.Code)
	require.Equal(t, "application/json", response.Header().Get("Content-Type"))
	require.Contains(t, response.Body.String(), string(city311Contract.ErrorTemporarilyUnavailable))
	require.Contains(t, response.Body.String(), `"retryable":true`)
}
