package filter

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crusttech/human/server/pkg/apigw/types"
	"github.com/stretchr/testify/require"
)

func mustMerge(t *testing.T, h types.Handler, params any) types.Handler {
	t.Helper()
	b, err := json.Marshal(params)
	require.NoError(t, err)
	h, err = h.Merge(b, types.Config{})
	require.NoError(t, err)
	return h
}

func callHandler(h types.Handler, r *http.Request) error {
	return h.Handler()(httptest.NewRecorder(), r)
}

// ---- Merge validation -------------------------------------------------------

func TestWebhookAuth_MergeValidation(t *testing.T) {
	cases := []struct {
		name    string
		params  any
		wantErr string
	}{
		{
			name:    "missing type",
			params:  map[string]any{},
			wantErr: "webhookAuth: type is required",
		},
		{
			name:    "unknown type",
			params:  map[string]any{"type": "magic"},
			wantErr: `webhookAuth: unknown type "magic"`,
		},
		{
			name:    "hmac: missing secret",
			params:  map[string]any{"type": "hmac", "hmac": map[string]any{"header": "X-Sig"}},
			wantErr: "webhookAuth: hmac.secret is required",
		},
		{
			name:    "hmac: missing header",
			params:  map[string]any{"type": "hmac", "hmac": map[string]any{"secret": "s"}},
			wantErr: "webhookAuth: hmac.header is required",
		},
		{
			name:    "hmac: bad algorithm",
			params:  map[string]any{"type": "hmac", "hmac": map[string]any{"secret": "s", "header": "X-Sig", "algorithm": "md5"}},
			wantErr: "webhookAuth: hmac.algorithm must be sha256 or sha1",
		},
		{
			name:    "bearer: missing secret",
			params:  map[string]any{"type": "bearer", "bearer": map[string]any{"header": "Authorization"}},
			wantErr: "webhookAuth: bearer.secret is required",
		},
		{
			name:    "bearer: missing header",
			params:  map[string]any{"type": "bearer", "bearer": map[string]any{"secret": "tok"}},
			wantErr: "webhookAuth: bearer.header is required",
		},
		{
			name:    "query: missing secret",
			params:  map[string]any{"type": "query", "query": map[string]any{"param": "token"}},
			wantErr: "webhookAuth: query.secret is required",
		},
		{
			name:    "query: missing param",
			params:  map[string]any{"type": "query", "query": map[string]any{"secret": "s"}},
			wantErr: "webhookAuth: query.param is required",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b, _ := json.Marshal(tc.params)
			_, err := NewWebhookAuth(types.Config{}).Merge(b, types.Config{})
			require.EqualError(t, err, tc.wantErr)
		})
	}
}

// ---- noop -------------------------------------------------------------------

func TestWebhookAuth_Noop(t *testing.T) {
	h := mustMerge(t, NewWebhookAuth(types.Config{}), map[string]any{"type": "noop"})
	r := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	require.NoError(t, callHandler(h, r))
}

// ---- bearer -----------------------------------------------------------------

func TestWebhookAuth_Bearer(t *testing.T) {
	params := map[string]any{
		"type": "bearer",
		"bearer": map[string]any{
			"secret": "supersecret",
			"header": "Authorization",
			"prefix": "Bearer ",
		},
	}

	t.Run("valid token", func(t *testing.T) {
		h := mustMerge(t, NewWebhookAuth(types.Config{}), params)
		r := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
		r.Header.Set("Authorization", "Bearer supersecret")
		require.NoError(t, callHandler(h, r))
	})

	t.Run("wrong token", func(t *testing.T) {
		h := mustMerge(t, NewWebhookAuth(types.Config{}), params)
		r := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
		r.Header.Set("Authorization", "Bearer wrongsecret")
		require.Error(t, callHandler(h, r))
	})

	t.Run("missing header", func(t *testing.T) {
		h := mustMerge(t, NewWebhookAuth(types.Config{}), params)
		r := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
		require.Error(t, callHandler(h, r))
	})
}

// ---- query ------------------------------------------------------------------

func TestWebhookAuth_Query(t *testing.T) {
	params := map[string]any{
		"type": "query",
		"query": map[string]any{
			"secret": "mytoken",
			"param":  "token",
		},
	}

	t.Run("valid param", func(t *testing.T) {
		h := mustMerge(t, NewWebhookAuth(types.Config{}), params)
		r := httptest.NewRequest(http.MethodPost, "/?token=mytoken", http.NoBody)
		require.NoError(t, callHandler(h, r))
	})

	t.Run("wrong param value", func(t *testing.T) {
		h := mustMerge(t, NewWebhookAuth(types.Config{}), params)
		r := httptest.NewRequest(http.MethodPost, "/?token=wrong", http.NoBody)
		require.Error(t, callHandler(h, r))
	})

	t.Run("missing param", func(t *testing.T) {
		h := mustMerge(t, NewWebhookAuth(types.Config{}), params)
		r := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
		require.Error(t, callHandler(h, r))
	})
}

// ---- hmac -------------------------------------------------------------------

func signBody(secret, body []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func TestWebhookAuth_HMAC(t *testing.T) {
	secret := []byte("webhooksecret")
	body := []byte(`{"event":"push"}`)
	sig := signBody(secret, body)

	params := map[string]any{
		"type": "hmac",
		"hmac": map[string]any{
			"secret":    "webhooksecret",
			"header":    "X-Hub-Signature-256",
			"algorithm": "sha256",
			"prefix":    "sha256=",
		},
	}

	t.Run("valid signature", func(t *testing.T) {
		h := mustMerge(t, NewWebhookAuth(types.Config{}), params)
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		r.Header.Set("X-Hub-Signature-256", "sha256="+sig)
		require.NoError(t, callHandler(h, r))
	})

	t.Run("wrong signature", func(t *testing.T) {
		h := mustMerge(t, NewWebhookAuth(types.Config{}), params)
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		r.Header.Set("X-Hub-Signature-256", "sha256=deaddead")
		require.Error(t, callHandler(h, r))
	})

	t.Run("missing header", func(t *testing.T) {
		h := mustMerge(t, NewWebhookAuth(types.Config{}), params)
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		require.Error(t, callHandler(h, r))
	})

	t.Run("body restored after read", func(t *testing.T) {
		h := mustMerge(t, NewWebhookAuth(types.Config{}), params)
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		r.Header.Set("X-Hub-Signature-256", "sha256="+sig)
		require.NoError(t, callHandler(h, r))

		// body must still be readable downstream
		got, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, body, got)
	})
}
