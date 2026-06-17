package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/crusttech/human/server/pkg/sql"
	"time"
)

type LlmProvider struct {
	ID           uint64            `json:"llmProviderID,string"`
	TenantID     uint64            `json:"tenantID,string,omitempty"`
	ProjectID    uint64            `json:"projectID,string,omitempty"`
	Handle       string            `json:"handle"`
	Status       string            `json:"status"`
	Provider     string            `json:"provider"`
	CredentialID uint64            `json:"credentialID,string"`
	Meta         LLMProviderMeta   `json:"meta"`
	Config       LLMProviderConfig `json:"config"`
	CreatedAt    time.Time         `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time        `json:"updatedAt,omitempty"`
	DeletedAt    *time.Time        `json:"deletedAt,omitempty"`
	CreatedBy    uint64            `json:"createdBy,string"`
	UpdatedBy    uint64            `json:"updatedBy,string,omitempty"`
	DeletedBy    uint64            `json:"deletedBy,string,omitempty"`
}

func (m *LLMProviderMeta) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m LLMProviderMeta) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseLLMProviderMeta(ss []string) (p LLMProviderMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *LLMProviderConfig) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m LLMProviderConfig) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseLLMProviderConfig(ss []string) (p LLMProviderConfig, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
