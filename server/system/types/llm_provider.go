package types

import (
	"github.com/crusttech/human/server/pkg/filter"
)

type (
	LLMProviderMeta struct {
		Short       string `json:"short"`
		Description string `json:"description"`
	}

	LLMProviderConfig struct {
		PromptURL   string                  `json:"promptURL"`
		Model       string                  `json:"model"`
		Temperature *float64                `json:"temperature,omitempty"`
		Timeout     string                  `json:"timeout"`
		Guard       *LLMProviderGuardConfig `json:"guard,omitempty"`
	}

	LLMProviderGuardConfig struct {
		Enabled    bool               `json:"enabled"`
		Provider   string             `json:"provider"` // "llama-guard"
		Model      string             `json:"model"`    // "llama-guard3:8b"
		Thresholds map[string]float64 `json:"thresholds,omitempty"`
	}

	LlmProviderFilter struct {
		LlmProviderID []uint64     `json:"llmProviderID"`
		Handle        string       `json:"handle"`
		Status        string       `json:"status"`
		Provider      string       `json:"provider"`
		Deleted       filter.State `json:"deleted"`

		Check func(*LlmProvider) (bool, error) `json:"-"`

		filter.Paging
		filter.Sorting
	}
)
