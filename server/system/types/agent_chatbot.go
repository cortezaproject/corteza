package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/crusttech/human/server/pkg/sql"
)

type (
	AgentChatbot struct {
		Enabled        bool                   `json:"enabled"`
		WidgetKey      string                 `json:"widgetKey,omitempty"`
		AllowedOrigins []string               `json:"allowedOrigins,omitempty"`
		SessionTTL     string                 `json:"sessionTTL,omitempty"`
		Handoff        AgentChatbotHandoff    `json:"handoff,omitempty"`
		Styling        AgentChatbotStyling    `json:"styling,omitempty"`
		Scenarios      []AgentChatbotScenario `json:"scenarios,omitempty"`
	}

	AgentChatbotHandoff struct {
		Enabled     bool     `json:"enabled"`
		TargetRoles []uint64 `json:"targetRoles,omitempty"`
	}

	AgentChatbotStyling struct {
		LogoURL    string                `json:"logoUrl,omitempty"`
		FontFamily string                `json:"fontFamily,omitempty"`
		FontSizes  AgentChatbotFontSizes `json:"fontSizes,omitempty"`
		Colors     AgentChatbotColors    `json:"colors,omitempty"`
		Launcher   AgentChatbotLauncher  `json:"launcher,omitempty"`
	}

	AgentChatbotFontSizes struct {
		Base    string `json:"base,omitempty"`
		Small   string `json:"small,omitempty"`
		Heading string `json:"heading,omitempty"`
	}

	AgentChatbotColors struct {
		Primary     string `json:"primary,omitempty"`
		PrimaryText string `json:"primaryText,omitempty"`
		Background  string `json:"background,omitempty"`
		Text        string `json:"text,omitempty"`
		UserBubble  string `json:"userBubble,omitempty"`
		AgentBubble string `json:"agentBubble,omitempty"`
	}

	AgentChatbotLauncher struct {
		IconURL     string `json:"iconUrl,omitempty"`
		IconVisible bool   `json:"iconVisible"`
		Label       string `json:"label,omitempty"`
		ButtonLabel string `json:"buttonLabel,omitempty"`
		Size        string `json:"size,omitempty"`
		Shape       string `json:"shape,omitempty"`
	}

	AgentChatbotScenario struct {
		ID     string          `json:"id"`
		Name   string          `json:"name,omitempty"`
		Type   string          `json:"type"`
		Config json.RawMessage `json:"config,omitempty"`
	}
)

func (m *AgentChatbot) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m AgentChatbot) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseAgentChatbot(ss []string) (p AgentChatbot, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
