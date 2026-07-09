package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func anthropicProvider(url string) (*sysTypes.LlmProvider, *sysTypes.Credential) {
	return &sysTypes.LlmProvider{
		Provider: "anthropic",
		Status:   "active",
		Config: sysTypes.LLMProviderConfig{
			PromptURL: url,
			Model:     "claude-3-5-sonnet-20241022",
		},
	}, &sysTypes.Credential{Credentials: "test-key"}
}

func openaiProvider(url string) (*sysTypes.LlmProvider, *sysTypes.Credential) {
	return &sysTypes.LlmProvider{
		Provider: "openai",
		Status:   "active",
		Config: sysTypes.LLMProviderConfig{
			PromptURL: url,
			Model:     "gpt-4o",
		},
	}, &sysTypes.Credential{Credentials: "test-key"}
}

func mockAnthropicServer(t *testing.T, capturedOutputTokens *int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			MaxTokens int `json:"max_tokens"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		*capturedOutputTokens = body.MaxTokens

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"content":     []map[string]any{{"type": "text", "text": "ok"}},
			"stop_reason": "end_turn",
			"usage":       map[string]any{"input_tokens": 10, "output_tokens": 5},
		})
	}))
}

func mockOpenAIServer(t *testing.T, capturedOutputTokens *int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			MaxTokens *int `json:"max_tokens"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		if body.MaxTokens != nil {
			*capturedOutputTokens = *body.MaxTokens
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]any{"role": "assistant", "content": "ok"}, "finish_reason": "stop"},
			},
			"usage": map[string]any{"prompt_tokens": 10, "completion_tokens": 5, "total_tokens": 15},
		})
	}))
}

func TestAnthropicOutputTokens_Set(t *testing.T) {
	var captured int
	srv := mockAnthropicServer(t, &captured)
	defer srv.Close()

	p, cred := anthropicProvider(srv.URL)
	_, err := promptAnthropic(context.Background(), p, cred, "", nil, 8000, nil, nil, "2023-06-01")
	require.NoError(t, err)
	assert.Equal(t, 8000, captured)
}

func TestAnthropicOutputTokens_Zero(t *testing.T) {
	var captured int
	srv := mockAnthropicServer(t, &captured)
	defer srv.Close()

	p, cred := anthropicProvider(srv.URL)
	_, err := promptAnthropic(context.Background(), p, cred, "", nil, 0, nil, nil, "2023-06-01")
	require.NoError(t, err)
	assert.Equal(t, 0, captured)
}

func TestOpenAIOutputTokens_Set(t *testing.T) {
	var captured int
	srv := mockOpenAIServer(t, &captured)
	defer srv.Close()

	p, cred := openaiProvider(srv.URL)
	_, err := promptOpenAI(context.Background(), p, cred, "", nil, 8000, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, 8000, captured)
}

func TestOpenAIOutputTokens_NotSentWhenZero(t *testing.T) {
	var captured int
	srv := mockOpenAIServer(t, &captured)
	defer srv.Close()

	p, cred := openaiProvider(srv.URL)
	_, err := promptOpenAI(context.Background(), p, cred, "", nil, 0, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, 0, captured)
}
