package guard

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crusttech/human/server/system/types"
)

func TestParseLlamaGuardOutput(t *testing.T) {
	tests := []struct {
		name       string
		output     string
		blocked    bool
		categories []string
	}{
		{name: "safe", output: "safe", blocked: false},
		{name: "safe uppercase", output: "SAFE", blocked: false},
		{name: "safe with whitespace", output: "  safe  ", blocked: false},
		{name: "unsafe single", output: "unsafe\nS1", blocked: true, categories: []string{"violent_crimes"}},
		{name: "unsafe multiple", output: "unsafe\nS1,S10", blocked: true, categories: []string{"violent_crimes", "hate"}},
		{name: "unsafe with whitespace", output: "unsafe\n S1 , S10 ", blocked: true, categories: []string{"violent_crimes", "hate"}},
		{name: "unsafe no category", output: "unsafe", blocked: true},
		{name: "unsafe unknown category", output: "unsafe\nS99", blocked: true, categories: []string{"S99"}},
		{name: "unexpected format", output: "something unexpected", blocked: false},
	}

	// parseLlamaGuardOutput reads lg.thresholds, so it is a method rather than
	// the free function this test used to call — which is why the package would
	// not build. With no thresholds configured every category blocks
	// unconditionally, which is what the expectations below assume.
	lg := NewLlamaGuard(&types.LlmProvider{}, "")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lg.parseLlamaGuardOutput(tt.output)
			if result.Blocked != tt.blocked {
				t.Errorf("blocked = %v, want %v", result.Blocked, tt.blocked)
			}
			for _, cat := range tt.categories {
				if _, ok := result.Categories[cat]; !ok {
					t.Errorf("expected category %q, got: %v", cat, result.Categories)
				}
			}
		})
	}
}

func TestLlamaGuard_CheckInput(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		blocked    bool
		wantErr    bool
	}{
		{
			name:       "safe response",
			response:   `{"choices":[{"message":{"content":"safe"}}]}`,
			statusCode: 200,
			blocked:    false,
		},
		{
			name:       "unsafe response",
			response:   `{"choices":[{"message":{"content":"unsafe\nS1"}}]}`,
			statusCode: 200,
			blocked:    true,
		},
		{
			name:       "server error",
			response:   `{"error":"internal"}`,
			statusCode: 500,
			wantErr:    true,
		},
		{
			name:       "empty choices",
			response:   `{"choices":[]}`,
			statusCode: 200,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request format
				var req llamaGuardRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Errorf("failed to decode request: %v", err)
				}
				if req.Model != "llama-guard3:8b" {
					t.Errorf("model = %q, want llama-guard3:8b", req.Model)
				}

				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			provider := &types.LlmProvider{
				Config: types.LLMProviderConfig{
					PromptURL: server.URL,
					Guard: &types.LLMProviderGuardConfig{
						Enabled:  true,
						Provider: "llama-guard",
						Model:    "llama-guard3:8b",
					},
				},
			}

			lg := NewLlamaGuard(provider, "")
			result, err := lg.CheckInput(context.Background(), "test input", nil)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Blocked != tt.blocked {
				t.Errorf("blocked = %v, want %v", result.Blocked, tt.blocked)
			}
		})
	}
}
