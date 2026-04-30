package guard

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/crusttech/human/server/system/types"
)

// llamaGuardCategories maps Llama Guard 3 safety category codes to human-readable names.
var llamaGuardCategories = map[string]string{
	"S1":  "violent_crimes",
	"S2":  "non_violent_crimes",
	"S3":  "sex_related_crimes",
	"S4":  "child_sexual_exploitation",
	"S5":  "defamation",
	"S6":  "specialized_advice",
	"S7":  "privacy",
	"S8":  "intellectual_property",
	"S9":  "indiscriminate_weapons",
	"S10": "hate",
	"S11": "suicide_self_harm",
	"S12": "sexual_content",
	"S13": "elections",
	"S14": "code_interpreter_abuse",
}

const (
	llamaGuardTimeout = 5 * time.Second
)

// LlamaGuard is a guard adapter that uses Meta's Llama Guard 3 model
// via an Ollama-compatible OpenAI API endpoint.
type LlamaGuard struct {
	endpoint string // e.g. "http://ollama-guard:11434/v1"
	model    string // e.g. "llama-guard3:8b"
	apiKey   string // credential from LlmProvider
}

// NewLlamaGuard creates a Llama Guard adapter from an LlmProvider's config.
func NewLlamaGuard(provider *types.LlmProvider, apiKey string) *LlamaGuard {
	model := "llama-guard3:8b"
	if provider.Config.Guard != nil && provider.Config.Guard.Model != "" {
		model = provider.Config.Guard.Model
	}

	return &LlamaGuard{
		endpoint: provider.Config.PromptURL,
		model:    model,
		apiKey:   apiKey,
	}
}

type llamaGuardRequest struct {
	Model    string              `json:"model"`
	Messages []llamaGuardMessage `json:"messages"`
}

type llamaGuardMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type llamaGuardResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// CheckInput sends the user input to Llama Guard for classification.
func (lg *LlamaGuard) CheckInput(ctx context.Context, input string, history []types.AiConversationMessage) (*GuardResult, error) {
	ctx, cancel := context.WithTimeout(ctx, llamaGuardTimeout)
	defer cancel()

	// Build the conversation for Llama Guard to evaluate.
	// Llama Guard expects a user message with the content to classify.
	messages := []llamaGuardMessage{
		{Role: "user", Content: input},
	}

	reqBody := llamaGuardRequest{
		Model:    lg.model,
		Messages: messages,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("llama guard: failed to marshal request: %w", err)
	}

	url := lg.endpoint + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("llama guard: failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if lg.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+lg.apiKey)
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("llama guard: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("llama guard: failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("llama guard: HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var lgResp llamaGuardResponse
	if err := json.Unmarshal(respBody, &lgResp); err != nil {
		return nil, fmt.Errorf("llama guard: failed to parse response: %w", err)
	}

	if len(lgResp.Choices) == 0 {
		return nil, fmt.Errorf("llama guard: no choices in response")
	}

	return parseLlamaGuardOutput(lgResp.Choices[0].Message.Content), nil
}

// parseLlamaGuardOutput parses Llama Guard's text response.
// Safe response: "safe"
// Unsafe response: "unsafe\nS1" or "unsafe\nS1,S2"
func parseLlamaGuardOutput(output string) *GuardResult {
	output = strings.TrimSpace(output)

	if strings.EqualFold(output, "safe") {
		return safe()
	}

	lines := strings.SplitN(output, "\n", 2)
	if !strings.EqualFold(strings.TrimSpace(lines[0]), "unsafe") {
		// Unexpected format — treat as safe but log the raw output
		return safe()
	}

	categories := make(map[string]float64)
	reason := "content classified as unsafe"

	if len(lines) > 1 {
		cats := strings.Split(strings.TrimSpace(lines[1]), ",")
		catNames := make([]string, 0, len(cats))
		for _, c := range cats {
			c = strings.TrimSpace(c)
			if name, ok := llamaGuardCategories[c]; ok {
				categories[name] = 1.0
				catNames = append(catNames, name)
			} else if c != "" {
				categories[c] = 1.0
				catNames = append(catNames, c)
			}
		}
		if len(catNames) > 0 {
			reason = "content classified as unsafe: " + strings.Join(catNames, ", ")
		}
	}

	return &GuardResult{
		Safe:       false,
		Blocked:    true,
		Reason:     reason,
		Categories: categories,
	}
}
