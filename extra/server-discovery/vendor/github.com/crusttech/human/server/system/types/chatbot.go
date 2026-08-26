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
