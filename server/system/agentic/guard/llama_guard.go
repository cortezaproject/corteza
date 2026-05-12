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
	defaultLlamaGuardTimeout = 5 * time.Second
)

// LlamaGuard is a guard adapter that uses Meta's Llama Guard 3 model
// via an Ollama-compatible OpenAI API endpoint.
type LlamaGuard struct {
	endpoint   string             // e.g. "http://ollama-guard:11434/v1"
	model      string             // e.g. "llama-guard3:8b"
	apiKey     string             // credential from LlmProvider
	timeout    time.Duration
	thresholds map[string]float64 // per-category block thresholds; absent = always block
}

// NewLlamaGuard creates a Llama Guard adapter from an LlmProvider's config.
func NewLlamaGuard(provider *types.LlmProvider, apiKey string) *LlamaGuard {
	model := "llama-guard3:8b"
	var thresholds map[string]float64
	if provider.Config.Guard != nil {
		if provider.Config.Guard.Model != "" {
			model = provider.Config.Guard.Model
		}
		thresholds = provider.Config.Guard.Thresholds
	}

	timeout := defaultLlamaGuardTimeout
	if provider.Config.Timeout != "" {
		if d, err := time.ParseDuration(provider.Config.Timeout); err == nil && d > 0 {
			timeout = d
		}
	}

	return &LlamaGuard{
		endpoint:   provider.Config.PromptURL,
		model:      model,
		apiKey:     apiKey,
		timeout:    timeout,
		thresholds: thresholds,
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
	ctx, cancel := context.WithTimeout(ctx, lg.timeout)
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

	return lg.parseLlamaGuardOutput(lgResp.Choices[0].Message.Content), nil
}

// parseLlamaGuardOutput parses Llama Guard's text response.
// Safe response: "safe"
// Unsafe response: "unsafe\nS1" or "unsafe\nS1,S2"
// If thresholds are configured, a category only triggers a block when its
// score (always 1.0 from Llama Guard) meets or exceeds the threshold.
// Categories with no threshold entry are blocked unconditionally.
func (lg *LlamaGuard) parseLlamaGuardOutput(output string) *GuardResult {
	output = strings.TrimSpace(output)

	if strings.EqualFold(output, "safe") {
		return safe()
	}

	lines := strings.SplitN(output, "\n", 2)
	if !strings.EqualFold(strings.TrimSpace(lines[0]), "unsafe") {
		// Unexpected format — treat as safe
		return safe()
	}

	categories := make(map[string]float64)
	var blockedNames []string

	if len(lines) > 1 {
		cats := strings.Split(strings.TrimSpace(lines[1]), ",")
		for _, c := range cats {
			c = strings.TrimSpace(c)
			name := c
			if mapped, ok := llamaGuardCategories[c]; ok {
				name = mapped
			}
			if name == "" {
				continue
			}
			categories[name] = 1.0
			// Block if no threshold configured, or score meets threshold.
			if threshold, ok := lg.thresholds[name]; !ok || 1.0 >= threshold {
				blockedNames = append(blockedNames, name)
			}
		}
	}

	if len(blockedNames) == 0 {
		return safe()
	}

	return &GuardResult{
		Safe:       false,
		Blocked:    true,
		Reason:     "content classified as unsafe: " + strings.Join(blockedNames, ", "),
		Categories: categories,
	}
}
