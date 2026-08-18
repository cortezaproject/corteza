package types

// Block option structs — the authoritative definition of what each block kind accepts.
// The agentic page block schema tool serializes these so agents know the correct structure.

type (
	BlockField struct {
		Name string `json:"name"`
	}

	// RefreshOptions are the keys every refreshable block carries. They are
	// embedded rather than repeated so a block kind cannot quietly advertise a
	// different spelling of the same three options.
	RefreshOptions struct {
		RefreshRate   int    `json:"refreshRate"`
		ShowRefresh   bool   `json:"showRefresh"`
		MagnifyOption string `json:"magnifyOption"`
	}

	// DrillDown opens a record list from a chart segment or a metric tile.
	DrillDown struct {
		Enabled           bool                `json:"enabled"`
		BlockID           string              `json:"blockID"`
		RecordListOptions DrillDownRecordList `json:"recordListOptions"`
	}

	DrillDownRecordList struct {
		Fields []BlockField `json:"fields"`
	}

	RecordListBlockOptions struct {
		RefreshOptions
		ModuleID  string       `json:"moduleID"`
		Fields    []BlockField `json:"fields"`
		Prefilter string       `json:"prefilter"`
		Presort   string       `json:"presort"`
		PerPage   int          `json:"perPage"`

		Selectable       bool              `json:"selectable"`
		SelectMode       string            `json:"selectMode"`
		SelectionButtons []SelectionButton `json:"selectionButtons"`

		Editable                      bool         `json:"editable"`
		Draggable                     bool         `json:"draggable"`
		EditFields                    []BlockField `json:"editFields"`
		InlineEditFields              []BlockField `json:"inlineEditFields"`
		InlineRecordEditEnabled       bool         `json:"inlineRecordEditEnabled"`
		InlineRecordEditAllowAddField bool         `json:"inlineRecordEditAllowAddField"`
		InlineRecordCopyEnabled       bool         `json:"inlineRecordCopyEnabled"`
		InlineValueFiltering          bool         `json:"inlineValueFiltering"`
		BulkRecordEditEnabled         bool         `json:"bulkRecordEditEnabled"`
		OpenRecordInEditMode          bool         `json:"openRecordInEditMode"`

		HideHeader                  bool `json:"hideHeader"`
		HideAddButton               bool `json:"hideAddButton"`
		HideImportButton            bool `json:"hideImportButton"`
		HideConfigureFieldsButton   bool `json:"hideConfigureFieldsButton"`
		HideSearch                  bool `json:"hideSearch"`
		HidePaging                  bool `json:"hidePaging"`
		HideSorting                 bool `json:"hideSorting"`
		HideFiltering               bool `json:"hideFiltering"`
		HideRecordViewButton        bool `json:"hideRecordViewButton"`
		HideRecordEditButton        bool `json:"hideRecordEditButton"`
		HideRecordCloneButton       bool `json:"hideRecordCloneButton"`
		HideRecordReminderButton    bool `json:"hideRecordReminderButton"`
		HideRecordPermissionsButton bool `json:"hideRecordPermissionsButton"`
		HideRecordDeleteButton      bool `json:"hideRecordDeleteButton"`

		AllowExport                bool     `json:"allowExport"`
		EnableRecordPageNavigation bool     `json:"enableRecordPageNavigation"`
		FullPageNavigation         bool     `json:"fullPageNavigation"`
		ShowTotalCount             bool     `json:"showTotalCount"`
		ShowRecordPerPageOption    bool     `json:"showRecordPerPageOption"`
		ShowDeletedRecordsOption   bool     `json:"showDeletedRecordsOption"`
		SearchableFields           []string `json:"searchableFields"`
		SearchSubmitMode           string   `json:"searchSubmitMode"`
		LinkToParent               bool     `json:"linkToParent"`
		OpenInNewTab               bool     `json:"openInNewTab"`
		PositionField              string   `json:"positionField"`
		RefField                   string   `json:"refField"`

		RecordDisplayOption         string `json:"recordDisplayOption"`
		RecordSelectorDisplayOption string `json:"recordSelectorDisplayOption"`
		AddRecordDisplayOption      string `json:"addRecordDisplayOption"`

		CustomFilterPresets bool                     `json:"customFilterPresets"`
		FilterPresets       []RecordListFilterPreset `json:"filterPresets"`
		CustomSummaries     bool                     `json:"customSummaries"`
		Summaries           []RecordListSummary      `json:"summaries"`
		TextStyles          RecordListTextStyles     `json:"textStyles"`
	}

	SelectionButton struct {
		Label    string `json:"label"`
		ButtonID string `json:"buttonID"`
		Enabled  bool   `json:"enabled"`
	}

	RecordListFilterPreset struct {
		Name   string   `json:"name"`
		Filter []any    `json:"filter"`
		Roles  []string `json:"roles"`
	}

	// Metric is one of min, max, avg, sum, emptyCount, notEmptyCount,
	// uniqueCount, earliest, latest.
	RecordListSummary struct {
		Label  string   `json:"label"`
		Field  []string `json:"field"`
		Metric string   `json:"metric"`
		Roles  []string `json:"roles"`
	}

	RecordListTextStyles struct {
		WrappedFields []string `json:"wrappedFields"`
	}

	RecordBlockOptions struct {
		Fields                               []BlockField     `json:"fields"`
		FieldConditions                      []FieldCondition `json:"fieldConditions"`
		ClearConditionalFieldsOnHide         bool             `json:"clearConditionalFieldsOnHide"`
		RecordSelectorShowAddRecordButton    bool             `json:"recordSelectorShowAddRecordButton"`
		RecordSelectorDisplayOption          string           `json:"recordSelectorDisplayOption"`
		RecordSelectorAddRecordDisplayOption string           `json:"recordSelectorAddRecordDisplayOption"`
		ReferenceField                       string           `json:"referenceField"`
		ReferenceModuleID                    string           `json:"referenceModuleID"`
		InlineRecordEditEnabled              bool             `json:"inlineRecordEditEnabled"`
		InlineRecordCopyEnabled              bool             `json:"inlineRecordCopyEnabled"`
		HorizontalFieldLayoutEnabled         bool             `json:"horizontalFieldLayoutEnabled"`
		RecordFieldLayoutOption              string           `json:"recordFieldLayoutOption"`
		MagnifyOption                        string           `json:"magnifyOption"`
	}

	// Condition is an expression; when it evaluates false the field is hidden.
	FieldCondition struct {
		Field       string `json:"field"`
		Condition   string `json:"condition"`
		ClearOnHide bool   `json:"clearOnHide"`
	}

	ChartBlockOptions struct {
		RefreshOptions
		ChartID           string    `json:"chartID"`
		LiveFilterEnabled bool      `json:"liveFilterEnabled"`
		DrillDown         DrillDown `json:"drillDown"`
	}

	AutomationBlockOptions struct {
		Buttons       []AutomationButton `json:"buttons"`
		Sealed        bool               `json:"sealed"`
		MagnifyOption string             `json:"magnifyOption"`
	}

	// A button runs one of three things: AutomationID for a TAQ, WorkflowID for
	// a workflow with an onManual trigger, or Script for a Corredor script.
	AutomationButton struct {
		Label        string `json:"label"`
		Enabled      bool   `json:"enabled"`
		ResourceType string `json:"resourceType"`
		Script       string `json:"script"`
		WorkflowID   string `json:"workflowID"`
		AutomationID string `json:"automationID"`
		StepID       string `json:"stepID"`
		Variant      string `json:"variant"`
	}

	ContentBlockOptions struct {
		Body          string `json:"body"`
		MagnifyOption string `json:"magnifyOption"`
	}

	IFrameBlockOptions struct {
		RefreshOptions
		Src string `json:"src"`
		// A record field holding the URL, for an iframe on a record page.
		SrcField string `json:"srcField"`
	}

	FileBlockOptions struct {
		RefreshOptions
		// list or gallery
		Mode            string   `json:"mode"`
		Attachments     []string `json:"attachments"`
		HideFileName    bool     `json:"hideFileName"`
		Height          string   `json:"height"`
		Width           string   `json:"width"`
		MaxHeight       string   `json:"maxHeight"`
		MaxWidth        string   `json:"maxWidth"`
		BorderRadius    string   `json:"borderRadius"`
		Margin          string   `json:"margin"`
		BackgroundColor string   `json:"backgroundColor"`
		ClickToView     bool     `json:"clickToView"`
		EnableDownload  bool     `json:"enableDownload"`
	}

	NavigationBlockOptions struct {
		Display         NavigationDisplay `json:"display"`
		NavigationItems []NavigationItem  `json:"navigationItems"`
		MagnifyOption   string            `json:"magnifyOption"`
	}

	NavigationDisplay struct {
		Appearance string `json:"appearance"`
		Alignment  string `json:"alignment"`
		Justify    string `json:"justify"`
	}

	// Type decides what the item renders as and which of Options.Item it reads:
	// "compose" follows pageID, "url" follows url, "dropdown" opens the
	// dropdown items, "text-section" is a plain label. An item with no type
	// matches none of them and renders nothing.
	NavigationItem struct {
		Type    string                `json:"type"`
		Options NavigationItemOptions `json:"options"`
	}

	NavigationItemOptions struct {
		Enabled         bool                      `json:"enabled"`
		TextColor       string                    `json:"textColor"`
		BackgroundColor string                    `json:"backgroundColor"`
		Item            NavigationItemDestination `json:"item"`
	}

	// PageID links to a page in this namespace; URL leaves the app.
	NavigationItemDestination struct {
		Label           string             `json:"label"`
		URL             string             `json:"url"`
		Target          string             `json:"target"`
		Delimiter       bool               `json:"delimiter"`
		PageID          string             `json:"pageID"`
		PageLayoutID    string             `json:"pageLayoutID"`
		ModuleID        string             `json:"moduleID"`
		DisplaySubPages bool               `json:"displaySubPages"`
		Align           string             `json:"align"`
		Dropdown        NavigationDropdown `json:"dropdown"`
	}

	NavigationDropdown struct {
		Items []NavigationDropdownItem `json:"items"`
	}

	NavigationDropdownItem struct {
		Label     string `json:"label"`
		URL       string `json:"url"`
		Delimiter bool   `json:"delimiter"`
		Target    string `json:"target"`
	}

	// A tab draws a sibling block of the same page by its blockID; give those
	// blocks meta.hidden so the grid does not draw them a second time.
	TabsBlockOptions struct {
		Style         TabsStyle `json:"style"`
		Tabs          []Tab     `json:"tabs"`
		MagnifyOption string    `json:"magnifyOption"`
	}

	TabsStyle struct {
		Appearance  string `json:"appearance"`
		Alignment   string `json:"alignment"`
		Justify     string `json:"justify"`
		Orientation string `json:"orientation"`
		Position    string `json:"position"`
	}

	Tab struct {
		BlockID string `json:"blockID"`
		Title   string `json:"title"`
		Lazy    bool   `json:"lazy"`
	}

	RecordRevisionsBlockOptions struct {
		RefreshOptions
		Preload bool `json:"preload"`
		// Empty shows every field.
		DisplayedFields []string `json:"displayedFields"`
		SortDirection   string   `json:"sortDirection"`
	}

	GeometryBlockOptions struct {
		RefreshOptions
		DefaultView   string         `json:"defaultView"`
		Center        []float64      `json:"center"`
		Feeds         []GeometryFeed `json:"feeds"`
		ZoomStarting  int            `json:"zoomStarting"`
		ZoomMin       int            `json:"zoomMin"`
		ZoomMax       int            `json:"zoomMax"`
		Bounds        [][]float64    `json:"bounds"`
		LockBounds    bool           `json:"lockBounds"`
		DisplayOption string         `json:"displayOption"`
		HideGeoSearch bool           `json:"hideGeoSearch"`
	}

	GeometryFeed struct {
		Resource       string              `json:"resource"` // "compose:record"
		TitleField     string              `json:"titleField"`
		GeometryField  string              `json:"geometryField"`
		DisplayMarker  bool                `json:"displayMarker"`
		DisplayPolygon bool                `json:"displayPolygon"`
		Options        GeometryFeedOptions `json:"options"`
	}

	GeometryFeedOptions struct {
		ModuleID  string `json:"moduleID"`
		Color     string `json:"color"`
		Prefilter string `json:"prefilter"`
	}

	AgentChatBlockOptions struct {
		// Empty allows every agent the viewer may reach.
		AllowedAgentIDs []string `json:"allowedAgentIDs"`
		DefaultAgentID  string   `json:"defaultAgentID"`
		AutoResume      bool     `json:"autoResume"`
	}

	SocialFeedBlockOptions struct {
		ProfileSourceField string `json:"profileSourceField"`
		ProfileUrl         string `json:"profileUrl"`
	}

	MetricBlockOptions struct {
		RefreshOptions
		Metrics []MetricItem `json:"metrics"`
	}

	MetricItem struct {
		Label          string           `json:"label"`
		ModuleID       string           `json:"moduleID"`
		Filter         string           `json:"filter"`
		MetricField    string           `json:"metricField"`
		Operation      string           `json:"operation"`
		Prefix         string           `json:"prefix"`
		Suffix         string           `json:"suffix"`
		NumberFormat   string           `json:"numberFormat"`
		DimensionField string           `json:"dimensionField"`
		DateFormat     string           `json:"dateFormat"`
		BucketSize     string           `json:"bucketSize"`
		TransformFx    string           `json:"transformFx"`
		ValueStyle     MetricValueStyle `json:"valueStyle"`
		DrillDown      DrillDown        `json:"drillDown"`
		Comparison     MetricComparison `json:"comparison"`
	}

	MetricValueStyle struct {
		Color           string `json:"color"`
		BackgroundColor string `json:"backgroundColor"`
	}

	// Shows the change against the same metric one period ago.
	MetricComparison struct {
		Enabled      bool   `json:"enabled"`
		Type         string `json:"type"`
		Period       string `json:"period"`
		CustomFilter string `json:"customFilter"`
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
		RefreshOptions
		Value    ProgressValueOptions   `json:"value"`
		MinValue ProgressValueOptions   `json:"minValue"`
		MaxValue ProgressValueOptions   `json:"maxValue"`
		Display  ProgressDisplayOptions `json:"display"`
	}

	CommentBlockOptions struct {
		RefreshOptions
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
		RefreshOptions
		DefaultView        string         `json:"defaultView"`
		Feeds              []CalendarFeed `json:"feeds"`
		Header             CalendarHeader `json:"header"`
		Locale             string         `json:"locale"`
		EventDisplayOption string         `json:"eventDisplayOption"`
	}

	CalendarHeader struct {
		Hide         bool     `json:"hide"`
		Views        []string `json:"views"`
		HidePrevNext bool     `json:"hidePrevNext"`
		HideToday    bool     `json:"hideToday"`
		HideTitle    bool     `json:"hideTitle"`
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
		RefreshOptions
		ModuleID   string `json:"moduleID"`
		GroupField string `json:"groupField"`
		// One block is one column of a board: the value of GroupField that this
		// block holds, and what a card dropped here is set to. Without it the
		// block filters on the ungrouped value, so a board is one block per value.
		Group                  string `json:"group"`
		LabelField             string `json:"labelField"`
		DescriptionField       string `json:"descriptionField"`
		PositionField          string `json:"positionField"`
		Filter                 string `json:"filter"`
		DisplayOption          string `json:"displayOption"`
		AddRecordDisplayOption string `json:"addRecordDisplayOption"`
	}

	ChatbotInboxBlockOptions struct {
		ChatbotIDs    []string `json:"chatbotIDs"`
		StatusFilter  []string `json:"statusFilter"`
		RefreshRate   int      `json:"refreshRate"`
		ShowRefresh   bool     `json:"showRefresh"`
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
		EditFields:       []BlockField{{}},
		InlineEditFields: []BlockField{{}},
		SelectionButtons: []SelectionButton{{}},
		FilterPresets:    []RecordListFilterPreset{{}},
		Summaries:        []RecordListSummary{{}},
		SearchableFields: []string{},
	},
	"Record": RecordBlockOptions{
		Fields:          []BlockField{{}},
		FieldConditions: []FieldCondition{{}},
	},
	"Chart": ChartBlockOptions{
		DrillDown: DrillDown{RecordListOptions: DrillDownRecordList{Fields: []BlockField{{}}}},
	},
	"Automation": AutomationBlockOptions{Buttons: []AutomationButton{{}}},
	"Content":    ContentBlockOptions{},
	"SocialFeed": SocialFeedBlockOptions{},
	"Metric": MetricBlockOptions{Metrics: []MetricItem{{
		DrillDown: DrillDown{RecordListOptions: DrillDownRecordList{Fields: []BlockField{{}}}},
	}}},
	"Progress": ProgressBlockOptions{
		Display: ProgressDisplayOptions{Thresholds: []ProgressThreshold{{}}},
	},
	"Comment":         CommentBlockOptions{},
	"Calendar":        CalendarBlockOptions{Feeds: []CalendarFeed{{}}, Header: CalendarHeader{Views: []string{}}},
	"RecordOrganizer": RecordOrganizerBlockOptions{},
	"ChatbotInbox":    ChatbotInboxBlockOptions{ChatbotIDs: []string{}, StatusFilter: []string{}},
	"IFrame":          IFrameBlockOptions{},
	"File":            FileBlockOptions{Attachments: []string{}},
	"Navigation": NavigationBlockOptions{NavigationItems: []NavigationItem{{
		Options: NavigationItemOptions{
			Item: NavigationItemDestination{
				Dropdown: NavigationDropdown{Items: []NavigationDropdownItem{{}}},
			},
		},
	}}},
	"Tabs":            TabsBlockOptions{Tabs: []Tab{{}}},
	"RecordRevisions": RecordRevisionsBlockOptions{DisplayedFields: []string{}},
	"Geometry": GeometryBlockOptions{
		Center: []float64{},
		Feeds:  []GeometryFeed{{}},
	},
	"AgentChat": AgentChatBlockOptions{AllowedAgentIDs: []string{}},
}
