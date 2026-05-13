package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/sql"
)

type (
	Chatbot struct {
		ID      uint64 `json:"chatbotID,string"`
		Handle  string `json:"handle"`
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`

		WidgetKey      string                `json:"widgetKey,omitempty"`
		AllowedOrigins ChatbotAllowedOrigins `json:"allowedOrigins,omitempty"`
		SessionTTL     string                `json:"sessionTTL,omitempty"`

		Handoff   ChatbotHandoff   `json:"handoff"`
		Styling   ChatbotStyling   `json:"styling"`
		Scenarios ChatbotScenarios `json:"scenarios"`

		Labels map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
	}

	ChatbotFilter struct {
		ChatbotID []string `json:"chatbotID"`
		Handle    string   `json:"handle"`
		WidgetKey string   `json:"widgetKey"`
		Query     string   `json:"query"`

		Deleted filter.State `json:"deleted"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		Check func(*Chatbot) (bool, error) `json:"-"`

		filter.Sorting
		filter.Paging
	}

	ChatbotHandoff struct {
		Enabled     bool     `json:"enabled"`
		TargetRoles []uint64 `json:"targetRoles,omitempty"`
	}

	ChatbotStyling struct {
		LogoAttachmentID uint64           `json:"logoAttachmentID,string,omitempty"`
		LogoURL          string           `json:"logoURL,omitempty"` // computed at read time; not persisted
		FontFamily       string           `json:"fontFamily,omitempty"`
		FontSizes        ChatbotFontSizes `json:"fontSizes,omitempty"`
		Colors           ChatbotColors    `json:"colors,omitempty"`
		Launcher         ChatbotLauncher  `json:"launcher,omitempty"`
	}

	ChatbotFontSizes struct {
		Base    string `json:"base,omitempty"`
		Small   string `json:"small,omitempty"`
		Heading string `json:"heading,omitempty"`
	}

	ChatbotColors struct {
		Primary     string `json:"primary,omitempty"`
		PrimaryText string `json:"primaryText,omitempty"`
		Header      string `json:"header,omitempty"`
		HeaderText  string `json:"headerText,omitempty"`
		Background  string `json:"background,omitempty"`
		Text        string `json:"text,omitempty"`
		UserBubble  string `json:"userBubble,omitempty"`
		AgentBubble string `json:"agentBubble,omitempty"`
	}

	ChatbotLauncher struct {
		IconURL          string `json:"iconURL,omitempty"`
		IconAttachmentID uint64 `json:"iconAttachmentID,string,omitempty"`
		IconVisible      bool   `json:"iconVisible"`
		Label            string `json:"label,omitempty"`
		ButtonLabel      string `json:"buttonLabel,omitempty"`
		Size             string `json:"size,omitempty"`
		Shape            string `json:"shape,omitempty"`
		Position         string `json:"position,omitempty"`
		StartOpen        bool   `json:"startOpen,omitempty"`
	}

	ChatbotScenario struct {
		ID                   string          `json:"id"`
		Name                 string          `json:"name,omitempty"`
		Type                 string          `json:"type"`
		AgentID              uint64          `json:"agentID,string,omitempty"`
		Config               json.RawMessage `json:"config,omitempty"`
		BeforeAutomationID   *uint64         `json:"beforeAutomationID,string,omitempty"`
		AfterAutomationID    *uint64         `json:"afterAutomationID,string,omitempty"`
	}

	ChatbotScenarios []ChatbotScenario

	ChatbotAllowedOrigins []string
)

func (m *ChatbotHandoff) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ChatbotHandoff) Value() (driver.Value, error) { return json.Marshal(m) }

func (m *ChatbotStyling) Scan(src any) error { return sql.ParseJSON(src, m) }

// Value strips computed URL fields before persisting. LogoURL/IconURL are
// derived from the attachmentID at read time; keeping stale values in the DB
// causes drift when attachments are swapped or regenerated.
func (m ChatbotStyling) Value() (driver.Value, error) {
	out := m
	out.LogoURL = ""
	out.Launcher.IconURL = ""
	return json.Marshal(out)
}

func (m *ChatbotScenarios) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ChatbotScenarios) Value() (driver.Value, error) { return json.Marshal(m) }

func (m *ChatbotAllowedOrigins) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ChatbotAllowedOrigins) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseChatbotAllowedOrigins(ss []string) (p ChatbotAllowedOrigins, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func ParseChatbotHandoff(ss []string) (p ChatbotHandoff, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func ParseChatbotStyling(ss []string) (p ChatbotStyling, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func ParseChatbotScenarios(ss []string) (p ChatbotScenarios, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
