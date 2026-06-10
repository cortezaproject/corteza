package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/sql"
)

type (
	LlmProvider struct {
		ID           uint64 `json:"llmProviderID,string"`
		TenantID     uint64 `json:"tenantID,string,omitempty"`
		ProjectID    uint64 `json:"projectID,string,omitempty"`
		Handle       string `json:"handle"`
		Status       string `json:"status"`
		Provider     string `json:"provider"`
		CredentialID uint64 `json:"credentialID,string"`

		Meta   LLMProviderMeta   `json:"meta"`
		Config LLMProviderConfig `json:"config"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
	}

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

func ParseLLMProviderConfig(ss []string) (m LLMProviderConfig, err error) {
	if len(ss) == 0 {
		return
	}
	err = json.Unmarshal([]byte(ss[0]), &m)
	return
}

func ParseLLMProviderMeta(ss []string) (m LLMProviderMeta, err error) {
	if len(ss) == 0 {
		return
	}
	err = json.Unmarshal([]byte(ss[0]), &m)
	return
}

func (nm *LLMProviderConfig) Scan(src any) error          { return sql.ParseJSON(src, nm) }
func (nm LLMProviderConfig) Value() (driver.Value, error) { return json.Marshal(nm) }

func (nm *LLMProviderMeta) Scan(src any) error          { return sql.ParseJSON(src, nm) }
func (nm LLMProviderMeta) Value() (driver.Value, error) { return json.Marshal(nm) }
