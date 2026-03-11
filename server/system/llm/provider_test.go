package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sysTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func anthropicProvider(url string, providerMaxTokens int) (*sysTypes.LlmProvider, *sysTypes.Credential) {
	return &sysTypes.LlmProvider{
		Provider: "anthropic",
		Status:   "active",
		Config: sysTypes.LLMProviderConfig{
			PromptURL: url,
			Model:     "claude-3-5-sonnet-20241022",
			MaxTokens: providerMaxTokens,
		},
	}, &sysTypes.Credential{Credentials: "test-key"}
}

func openaiProvider(url string, providerMaxTokens int) (*sysTypes.LlmProvider, *sysTypes.Credential) {
	return &sysTypes.LlmProvider{
		Provider: "openai",
		Status:   "active",
		Config: sysTypes.LLMProviderConfig{
			PromptURL: url,
			Model:     "gpt-4o",
			MaxTokens: providerMaxTokens,
		},
	}, &sysTypes.Credential{Credentials: "test-key"}
}

func mockAnthropicServer(t *testing.T, capturedMaxTokens *int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			MaxTokens int `json:"max_tokens"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		*capturedMaxTokens = body.MaxTokens

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"content":     []map[string]any{{"type": "text", "text": "ok"}},
			"stop_reason": "end_turn",
			"usage":       map[string]any{"input_tokens": 10, "output_tokens": 5},
		})
	}))
}

func mockOpenAIServer(t *testing.T, capturedMaxTokens *int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			MaxTokens *int `json:"max_tokens"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		if body.MaxTokens != nil {
			*capturedMaxTokens = *body.MaxTokens
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

func TestAnthropicMaxTokens_FromCall(t *testing.T) {
	var captured int
	srv := mockAnthropicServer(t, &captured)
	defer srv.Close()

	p, cred := anthropicProvider(srv.URL, 0)
	_, err := promptAnthropic(context.Background(), p, cred, "", 8000, nil, nil, "2023-06-01")
	require.NoError(t, err)
	assert.Equal(t, 8000, captured)
}

func TestAnthropicMaxTokens_FromProvider(t *testing.T) {
	var captured int
	srv := mockAnthropicServer(t, &captured)
	defer srv.Close()

	p, cred := anthropicProvider(srv.URL, 2000)
	_, err := promptAnthropic(context.Background(), p, cred, "", 0, nil, nil, "2023-06-01")
	require.NoError(t, err)
	assert.Equal(t, 2000, captured)
}

func TestAnthropicMaxTokens_Default(t *testing.T) {
	var captured int
	srv := mockAnthropicServer(t, &captured)
	defer srv.Close()

	p, cred := anthropicProvider(srv.URL, 0)
	_, err := promptAnthropic(context.Background(), p, cred, "", 0, nil, nil, "2023-06-01")
	require.NoError(t, err)
	assert.Equal(t, 4096, captured)
}

func TestAnthropicMaxTokens_CallOverridesProvider(t *testing.T) {
	var captured int
	srv := mockAnthropicServer(t, &captured)
	defer srv.Close()

	p, cred := anthropicProvider(srv.URL, 2000)
	_, err := promptAnthropic(context.Background(), p, cred, "", 8000, nil, nil, "2023-06-01")
	require.NoError(t, err)
	assert.Equal(t, 8000, captured)
}

func TestOpenAIMaxTokens_FromCall(t *testing.T) {
	var captured int
	srv := mockOpenAIServer(t, &captured)
	defer srv.Close()

	p, cred := openaiProvider(srv.URL, 0)
	_, err := promptOpenAI(context.Background(), p, cred, "", 8000, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, 8000, captured)
}

func TestOpenAIMaxTokens_FromProvider(t *testing.T) {
	var captured int
	srv := mockOpenAIServer(t, &captured)
	defer srv.Close()

	p, cred := openaiProvider(srv.URL, 2000)
	_, err := promptOpenAI(context.Background(), p, cred, "", 0, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, 2000, captured)
}

func TestOpenAIMaxTokens_CallOverridesProvider(t *testing.T) {
	var captured int
	srv := mockOpenAIServer(t, &captured)
	defer srv.Close()

	p, cred := openaiProvider(srv.URL, 2000)
	_, err := promptOpenAI(context.Background(), p, cred, "", 8000, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, 8000, captured)
}
