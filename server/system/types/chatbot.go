package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/pkg/sql"
)

type (
	ChatbotFilter struct {
		ChatbotID []string `json:"chatbotID"`
		ProjectID uint64   `json:"projectID,string,omitempty"`
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
		Enabled    bool                     `json:"enabled"`
		Automation ChatbotHandoffAutomation `json:"automation"`
	}

	ChatbotHandoffAutomation struct {
		OnRequested ChatbotAutomationHook `json:"onRequested"`
		OnAccepted  ChatbotAutomationHook `json:"onAccepted"`
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
		ID         string                    `json:"id"`
		Name       string                    `json:"name,omitempty"`
		Type       string                    `json:"type"`
		AgentID    uint64                    `json:"agentID,string,omitempty"`
		Config     json.RawMessage           `json:"config,omitempty"`
		Automation ChatbotScenarioAutomation `json:"automation"`
	}

	ChatbotScenarioAutomation struct {
		Before ChatbotAutomationHook `json:"before"`
		After  ChatbotAutomationHook `json:"after"`
	}

	ChatbotAutomationHook struct {
		Automation string               `json:"automation"`
		Async      bool                 `json:"async"`
		Mappings   []ChatbotHookMapping `json:"mappings,omitempty"`
	}

	// ChatbotHookMapping maps a chatbot session state expression to a TAQ trigger input param.
	// StateExpression mirrors automation/types.Expr (circular import prevents direct use).
	ChatbotHookMapping struct {
		StateExpression ChatbotStateExpression `json:"stateExpression"`
		TriggerParam    string                 `json:"triggerParam"`
	}

	// ChatbotStateExpression holds an expression evaluated against chatbot session/step state.
	// Fields mirror automation/types.Expr Source/Expr/Value/Type for eval compatibility.
	ChatbotStateExpression struct {
		Source string      `json:"source,omitempty"`
		Expr   string      `json:"expr,omitempty"`
		Value  interface{} `json:"value,omitempty"`
		Type   string      `json:"type,omitempty"`
	}

	ChatbotScenarios []ChatbotScenario

	ChatbotAllowedOrigins []string
)

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

func ParseChatbotStyling(ss []string) (p ChatbotStyling, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
