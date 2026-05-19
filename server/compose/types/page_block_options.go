package types

// Block option structs — the authoritative definition of what each block kind accepts.
// The agentic page block schema tool serializes these so agents know the correct structure.

type (
	BlockField struct {
		Name string `json:"name"`
	}

	RecordListBlockOptions struct {
		Module                   string            `json:"module"`
		Fields                   []BlockField      `json:"fields"`
		Prefilter                string            `json:"prefilter"`
		Presort                  string            `json:"presort"`
		PerPage                  int               `json:"perPage"`
		Selectable               bool              `json:"selectable"`
		SelectMode               string            `json:"selectMode"`
		Editable                 bool              `json:"editable"`
		Draggable                bool              `json:"draggable"`
		AllowExport              bool              `json:"allowExport"`
		HideRecordViewButton     bool              `json:"hideRecordViewButton"`
		HideRecordEditButton     bool              `json:"hideRecordEditButton"`
		HideRecordCloneButton    bool              `json:"hideRecordCloneButton"`
		HideRecordReminderButton bool              `json:"hideRecordReminderButton"`
		SelectionButtons         []SelectionButton `json:"selectionButtons"`
	}

	SelectionButton struct {
		Label    string `json:"label"`
		ButtonID string `json:"buttonID"`
		Enabled  bool   `json:"enabled"`
	}

	RecordBlockOptions struct {
		Fields []BlockField `json:"fields"`
	}

	ChartBlockOptions struct {
		Chart string `json:"chart"`
	}

	AutomationBlockOptions struct {
		Buttons []AutomationButton `json:"buttons"`
		Sealed  bool               `json:"sealed"`
	}

	AutomationButton struct {
		Label        string `json:"label"`
		Enabled      bool   `json:"enabled"`
		ResourceType string `json:"resourceType"`
		Script       string `json:"script"`
		WorkflowID   string `json:"workflowID"`
	}

	ContentBlockOptions struct {
		Body string `json:"body"`
	}

	SocialFeedBlockOptions struct {
		ProfileSourceField string `json:"profileSourceField"`
		ProfileUrl         string `json:"profileUrl"`
	}

	MetricBlockOptions struct {
		Metrics []MetricItem `json:"metrics"`
	}

	MetricItem struct {
		Label    string `json:"label"`
		ModuleID string `json:"moduleID"`
		Filter   string `json:"filter"`
		Field    string `json:"field"`
		Reduce   string `json:"reduce"`
		Prefix   string `json:"prefix"`
		Suffix   string `json:"suffix"`
		Format   string `json:"format"`
	}

	ProgressBlockOptions struct {
		ModuleID string  `json:"moduleID"`
		Filter   string  `json:"filter"`
		Field    string  `json:"field"`
		MinValue float64 `json:"minValue"`
		MaxValue float64 `json:"maxValue"`
		Value    float64 `json:"value"`
		Label    string  `json:"label"`
	}

	CommentBlockOptions struct {
		ModuleID     string `json:"moduleID"`
		CommentField string `json:"commentField"`
		TitleField   string `json:"titleField"`
	}

	CalendarBlockOptions struct {
		DefaultView string         `json:"defaultView"`
		Feeds       []CalendarFeed `json:"feeds"`
	}

	CalendarFeed struct {
		ModuleID   string `json:"moduleID"`
		StartField string `json:"startField"`
		EndField   string `json:"endField"`
		TitleField string `json:"titleField"`
		Color      string `json:"color"`
	}

	RecordOrganizerBlockOptions struct {
		ModuleID         string `json:"moduleID"`
		GroupField       string `json:"groupField"`
		LabelField       string `json:"labelField"`
		DescriptionField string `json:"descriptionField"`
		Filter           string `json:"filter"`
		Sort             string `json:"sort"`
	}

	ChatbotInboxBlockOptions struct {
		ChatbotIDs    []string `json:"chatbotIDs"`
		StatusFilter  []string `json:"statusFilter"`
		RefreshRate   int      `json:"refreshRate"`
		AutoOpenFirst bool     `json:"autoOpenFirst"`
	}
)

// PageBlockOptionSchemas maps each block kind to a zero-value of its options struct.
// The agentic blockSchema handler marshals these to show agents the correct field names and types.
var PageBlockOptionSchemas = map[string]any{
	"RecordList":      RecordListBlockOptions{},
	"Record":          RecordBlockOptions{},
	"Chart":           ChartBlockOptions{},
	"Automation":      AutomationBlockOptions{},
	"Content":         ContentBlockOptions{},
	"SocialFeed":      SocialFeedBlockOptions{},
	"Metric":          MetricBlockOptions{},
	"Progress":        ProgressBlockOptions{},
	"Comment":         CommentBlockOptions{},
	"Calendar":        CalendarBlockOptions{},
	"RecordOrganizer": RecordOrganizerBlockOptions{},
	"ChatbotInbox":    ChatbotInboxBlockOptions{},
}
