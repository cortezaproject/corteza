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

type ChatbotSession struct {
	ID          uint64              `json:"id,string"`
	TenantID    uint64              `json:"tenantID,string,omitempty"`
	ProjectID   uint64              `json:"projectID,string,omitempty"`
	ChatbotID   uint64              `json:"chatbotID,string"`
	Status      string              `json:"status"`
	CurrentStep int                 `json:"currentStep"`
	State       ChatbotSessionState `json:"state,omitempty"`
	CreatedAt   time.Time           `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time          `json:"updatedAt,omitempty"`
	DeletedAt   *time.Time          `json:"deletedAt,omitempty"`
	CreatedBy   uint64              `json:"createdBy,string"`
	UpdatedBy   uint64              `json:"updatedBy,string,omitempty"`
	DeletedBy   uint64              `json:"deletedBy,string,omitempty"`
}

func (m *ChatbotSessionState) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ChatbotSessionState) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseChatbotSessionState(ss []string) (p ChatbotSessionState, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
