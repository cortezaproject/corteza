package types

// Block option structs — the authoritative definition of what each block kind accepts.
// The agentic page block schema tool serializes these so agents know the correct structure.

type (
	BlockField struct {
		Name string `json:"name"`
	}

	RecordListBlockOptions struct {
		ModuleID                 string            `json:"moduleID"`
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
		ChartID string `json:"chartID"`
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
		Label        string `json:"label"`
		ModuleID     string `json:"moduleID"`
		Filter       string `json:"filter"`
		MetricField  string `json:"metricField"`
		Operation    string `json:"operation"`
		Prefix       string `json:"prefix"`
		Suffix       string `json:"suffix"`
		NumberFormat string `json:"numberFormat"`
	}

	ProgressValueOptions struct {
		Default   float64 `json:"default"`
		ModuleID  string  `json:"moduleID"`
		Filter    string  `json:"filter"`
		Field     string  `json:"field"`
		Operation string  `json:"operation"`
	}

	ProgressThreshold struct {
		Value   float64 `json:"value"`
		Variant string  `json:"variant"`
	}

	ProgressDisplayOptions struct {
		ShowValue    bool                `json:"showValue"`
		ShowRelative bool                `json:"showRelative"`
		ShowProgress bool                `json:"showProgress"`
		Variant      string              `json:"variant"`
		Thresholds   []ProgressThreshold `json:"thresholds"`
	}

	ProgressBlockOptions struct {
		Value    ProgressValueOptions   `json:"value"`
		MinValue ProgressValueOptions   `json:"minValue"`
		MaxValue ProgressValueOptions   `json:"maxValue"`
		Display  ProgressDisplayOptions `json:"display"`
	}

	CommentBlockOptions struct {
		ModuleID        string `json:"moduleID"`
		Filter          string `json:"filter"`
		TitleField      string `json:"titleField"`
		ContentField    string `json:"contentField"`
		ReplyField      string `json:"replyField"`
		ReferenceField  string `json:"referenceField"`
		AttachmentField string `json:"attachmentField"`
		ReactionsField  string `json:"reactionsField"`
		SortDirection   string `json:"sortDirection"`
	}

	CalendarBlockOptions struct {
		DefaultView string         `json:"defaultView"`
		Feeds       []CalendarFeed `json:"feeds"`
	}

	CalendarFeedOptions struct {
		ModuleID  string `json:"moduleID"`
		Color     string `json:"color"`
		Prefilter string `json:"prefilter"`
	}

	CalendarFeed struct {
		Resource   string              `json:"resource"` // "compose:record"
		StartField string              `json:"startField"`
		EndField   string              `json:"endField"`
		TitleField string              `json:"titleField"`
		AllDay     bool                `json:"allDay"`
		Options    CalendarFeedOptions `json:"options"`
	}

	RecordOrganizerBlockOptions struct {
		ModuleID   string `json:"moduleID"`
		GroupField string `json:"groupField"`
		// One block is one column of a board: the value of GroupField that this
		// block holds. Without it the block filters on the empty value and shows
		// nothing, so a board is one block per value.
		Group            string `json:"group"`
		LabelField       string `json:"labelField"`
		DescriptionField string `json:"descriptionField"`
		PositionField    string `json:"positionField"`
		Filter           string `json:"filter"`
	}

	ChatbotInboxBlockOptions struct {
		ChatbotIDs    []string `json:"chatbotIDs"`
		StatusFilter  []string `json:"statusFilter"`
		RefreshRate   int      `json:"refreshRate"`
		AutoOpenFirst bool     `json:"autoOpenFirst"`
		ShowFilter    bool     `json:"showFilter"`
	}
)

// PageBlockOptionSchemas maps each block kind to a zero-value of its options struct.
// The agentic blockSchema handler marshals these to show agents the correct field names and types.
// Slice fields hold one zero-value element so the serialized schema exposes
// the nested field names (a nil slice would marshal as null and hide them).
var PageBlockOptionSchemas = map[string]any{
	"RecordList": RecordListBlockOptions{
		Fields:           []BlockField{{}},
		SelectionButtons: []SelectionButton{{}},
	},
	"Record":     RecordBlockOptions{Fields: []BlockField{{}}},
	"Chart":      ChartBlockOptions{},
	"Automation": AutomationBlockOptions{Buttons: []AutomationButton{{}}},
	"Content":    ContentBlockOptions{},
	"SocialFeed": SocialFeedBlockOptions{},
	"Metric":     MetricBlockOptions{Metrics: []MetricItem{{}}},
	"Progress": ProgressBlockOptions{
		Display: ProgressDisplayOptions{Thresholds: []ProgressThreshold{{}}},
	},
	"Comment":         CommentBlockOptions{},
	"Calendar":        CalendarBlockOptions{Feeds: []CalendarFeed{{}}},
	"RecordOrganizer": RecordOrganizerBlockOptions{},
	"ChatbotInbox":    ChatbotInboxBlockOptions{ChatbotIDs: []string{}, StatusFilter: []string{}},
}
