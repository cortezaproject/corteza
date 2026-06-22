package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/sql"
	"time"
)

type Chatbot struct {
	ID             uint64                           `json:"chatbotID,string"`
	TenantID       uint64                           `json:"tenantID,string,omitempty"`
	ProjectID      uint64                           `json:"projectID,string,omitempty"`
	Handle         string                           `json:"handle"`
	Name           string                           `json:"name"`
	Enabled        bool                             `json:"enabled"`
	WidgetKey      string                           `json:"widgetKey,omitempty"`
	AllowedOrigins ChatbotAllowedOrigins            `json:"allowedOrigins,omitempty"`
	SessionTTL     string                           `json:"sessionTTL,omitempty"`
	Handoff        ChatbotHandoff                   `json:"handoff"`
	Styling        ChatbotStyling                   `json:"styling"`
	Scenarios      ChatbotScenarios                 `json:"scenarios"`
	CreatedAt      time.Time                        `json:"createdAt,omitempty"`
	UpdatedAt      *time.Time                       `json:"updatedAt,omitempty"`
	DeletedAt      *time.Time                       `json:"deletedAt,omitempty"`
	CreatedBy      uint64                           `json:"createdBy,string"`
	UpdatedBy      uint64                           `json:"updatedBy,string,omitempty"`
	DeletedBy      uint64                           `json:"deletedBy,string,omitempty"`
	Labels         map[string]labelTypes.LabelValue `json:"labels,omitempty"`
}

func (r Chatbot) Clone() *Chatbot {
	dup := r
	if b, err := json.Marshal(r); err == nil {
		_ = json.Unmarshal(b, &dup)
	}
	return &dup
}

func (m *ChatbotAllowedOrigins) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ChatbotAllowedOrigins) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseChatbotAllowedOrigins(ss []string) (p ChatbotAllowedOrigins, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *ChatbotHandoff) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ChatbotHandoff) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseChatbotHandoff(ss []string) (p ChatbotHandoff, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *ChatbotScenarios) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ChatbotScenarios) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseChatbotScenarios(ss []string) (p ChatbotScenarios, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
