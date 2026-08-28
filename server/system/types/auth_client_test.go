package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Verify gates every token the oauth2 handlers issue, so the direction of its
// comparisons decides whether a validity window admits a client or locks it
// out. Both were the wrong way round: a client with an expiry in the future
// read as "expired" and one valid since yesterday read as "not yet valid",
// which made a window the one configuration that could never work.
func TestAuthClientVerifyWindow(t *testing.T) {
	var (
		past   = time.Now().Add(-time.Hour)
		future = time.Now().Add(time.Hour)
	)

	for _, tc := range []struct {
		name      string
		client    *AuthClient
		wantError string
	}{
		{"no window", &AuthClient{Enabled: true}, ""},
		{"inside the window", &AuthClient{Enabled: true, ValidFrom: &past, ExpiresAt: &future}, ""},
		{"valid since the past", &AuthClient{Enabled: true, ValidFrom: &past}, ""},
		{"expires in the future", &AuthClient{Enabled: true, ExpiresAt: &future}, ""},
		{"already expired", &AuthClient{Enabled: true, ExpiresAt: &past}, "expired"},
		{"not started yet", &AuthClient{Enabled: true, ValidFrom: &future}, "not yet valid"},
		{"disabled", &AuthClient{}, "disabled"},
		{"nil", nil, "disabled"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.client.Verify()
			if tc.wantError == "" {
				require.NoError(t, err)
				return
			}
			require.EqualError(t, err, tc.wantError)
		})
	}
}
