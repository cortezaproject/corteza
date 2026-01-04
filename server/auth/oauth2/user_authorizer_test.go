package oauth2

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_redirectWithOIDCError(t *testing.T) {
	testCases := []struct {
		name         string
		redirectURI  string
		state        string
		errorCode    string
		errorDesc    string
		expectStatus int
		expectParams map[string]string
	}{
		{
			name:         "login_required error with state",
			redirectURI:  "https://client.example.com/callback",
			state:        "abc123",
			errorCode:    "login_required",
			errorDesc:    "User is not authenticated",
			expectStatus: http.StatusFound,
			expectParams: map[string]string{
				"error":             "login_required",
				"error_description": "User is not authenticated",
				"state":             "abc123",
			},
		},
		{
			name:         "interaction_required error without state",
			redirectURI:  "https://client.example.com/callback",
			state:        "",
			errorCode:    "interaction_required",
			errorDesc:    "User consent is required",
			expectStatus: http.StatusFound,
			expectParams: map[string]string{
				"error":             "interaction_required",
				"error_description": "User consent is required",
			},
		},
		{
			name:         "preserves existing query params",
			redirectURI:  "https://client.example.com/callback?existing=param",
			state:        "xyz",
			errorCode:    "access_denied",
			errorDesc:    "Access denied",
			expectStatus: http.StatusFound,
			expectParams: map[string]string{
				"error":             "access_denied",
				"error_description": "Access denied",
				"state":             "xyz",
				"existing":          "param",
			},
		},
		{
			name:         "invalid redirect_uri returns bad request",
			redirectURI:  "://invalid-uri",
			state:        "",
			errorCode:    "server_error",
			errorDesc:    "Something went wrong",
			expectStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			req := httptest.NewRequest(http.MethodGet, "/authorize", nil)
			rec := httptest.NewRecorder()

			redirectWithOIDCError(rec, req, tc.redirectURI, tc.state, tc.errorCode, tc.errorDesc)

			rq.Equal(tc.expectStatus, rec.Code)

			if tc.expectStatus == http.StatusFound {
				location := rec.Header().Get("Location")
				rq.NotEmpty(location)

				parsedURL, err := url.Parse(location)
				rq.NoError(err)

				query := parsedURL.Query()
				for key, expectedValue := range tc.expectParams {
					rq.Equal(expectedValue, query.Get(key), "query param %s mismatch", key)
				}

				// State should not be present if it was empty
				if tc.state == "" {
					rq.Empty(query.Get("state"))
				}
			}
		})
	}
}

func Test_UserIDSerializer(t *testing.T) {
	testCases := []struct {
		name     string
		userID   uint64
		roles    []uint64
		expected string
	}{
		{
			name:     "user ID only",
			userID:   12345,
			roles:    nil,
			expected: "12345",
		},
		{
			name:     "user ID with one role",
			userID:   12345,
			roles:    []uint64{100},
			expected: "12345 100",
		},
		{
			name:     "user ID with multiple roles",
			userID:   12345,
			roles:    []uint64{100, 200, 300},
			expected: "12345 100 200 300",
		},
		{
			name:     "empty roles slice",
			userID:   99999,
			roles:    []uint64{},
			expected: "99999",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)
			result := UserIDSerializer(tc.userID, tc.roles...)
			rq.Equal(tc.expected, result)
		})
	}
}
