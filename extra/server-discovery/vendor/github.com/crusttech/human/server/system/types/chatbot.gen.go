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
	"github.com/crusttech/human/server/pkg/revisions"
	"github.com/crusttech/human/server/pkg/sql"
	"reflect"
	"time"
)

type (
	Chatbot struct {
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

	ChatbotHandoff struct {
		Enabled    bool                     `json:"enabled"`
		Automation ChatbotHandoffAutomation `json:"automation"`
	}

	ChatbotHandoffAutomation struct {
		OnRequested ChatbotAutomationHook `json:"onRequested"`
		OnAccepted  ChatbotAutomationHook `json:"onAccepted"`
	}

	ChatbotAutomationHook struct {
		Automation string               `json:"automation"`
		Async      bool                 `json:"async"`
		Mappings   []ChatbotHookMapping `json:"mappings,omitempty"`
	}

	ChatbotHookMapping struct {
		StateExpression ChatbotStateExpression `json:"stateExpression"`
		TriggerParam    string                 `json:"triggerParam"`
	}

	ChatbotStateExpression struct {
		Source string `json:"source,omitempty"`
		Expr   string `json:"expr,omitempty"`
		Value  any    `json:"value,omitempty"`
		Type   string `json:"type,omitempty"`
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

	ChatbotStyling struct {
		LogoAttachmentID uint64           `json:"logoAttachmentID,string,omitempty"`
		LogoURL          string           `json:"logoURL,omitempty"`
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
)

func (r Chatbot) Clone() *Chatbot {
	dup := r
	dup.Handoff = *r.Handoff.Clone()

	dup.Styling = *r.Styling.Clone()

	if r.UpdatedAt != nil {
		v := *r.UpdatedAt
		dup.UpdatedAt = &v
	}

	if r.DeletedAt != nil {
		v := *r.DeletedAt
		dup.DeletedAt = &v
	}

	if r.Labels != nil {
		dup.Labels = make(map[string]labelTypes.LabelValue, len(r.Labels))
		for k, v := range r.Labels {
			dup.Labels[k] = v
		}
	}
	return &dup
}

func (r Chatbot) Diff(cmp *Chatbot) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &Chatbot{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "chatbotID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.Handle != cmp.Handle {
		out = append(out, &revisions.Change{Key: "handle", Old: []any{cmp.Handle}, New: []any{r.Handle}})
	}

	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Enabled != cmp.Enabled {
		out = append(out, &revisions.Change{Key: "enabled", Old: []any{cmp.Enabled}, New: []any{r.Enabled}})
	}

	if r.WidgetKey != cmp.WidgetKey {
		out = append(out, &revisions.Change{Key: "widgetKey", Old: []any{cmp.WidgetKey}, New: []any{r.WidgetKey}})
	}

	if !reflect.DeepEqual(r.AllowedOrigins, cmp.AllowedOrigins) {
		out = append(out, &revisions.Change{Key: "allowedOrigins", Old: []any{cmp.AllowedOrigins}, New: []any{r.AllowedOrigins}})
	}

	if r.SessionTTL != cmp.SessionTTL {
		out = append(out, &revisions.Change{Key: "sessionTTL", Old: []any{cmp.SessionTTL}, New: []any{r.SessionTTL}})
	}

	for _, c := range r.Handoff.Diff(&cmp.Handoff) {
		c.Key = "handoff." + c.Key
		out = append(out, c)
	}

	for _, c := range r.Styling.Diff(&cmp.Styling) {
		c.Key = "styling." + c.Key
		out = append(out, c)
	}

	if !reflect.DeepEqual(r.Scenarios, cmp.Scenarios) {
		out = append(out, &revisions.Change{Key: "scenarios", Old: []any{cmp.Scenarios}, New: []any{r.Scenarios}})
	}

	if !reflect.DeepEqual(r.CreatedAt, cmp.CreatedAt) {
		out = append(out, &revisions.Change{Key: "createdAt", Old: []any{cmp.CreatedAt}, New: []any{r.CreatedAt}})
	}

	if !reflect.DeepEqual(r.UpdatedAt, cmp.UpdatedAt) {
		out = append(out, &revisions.Change{Key: "updatedAt", Old: []any{cmp.UpdatedAt}, New: []any{r.UpdatedAt}})
	}

	if !reflect.DeepEqual(r.DeletedAt, cmp.DeletedAt) {
		out = append(out, &revisions.Change{Key: "deletedAt", Old: []any{cmp.DeletedAt}, New: []any{r.DeletedAt}})
	}

	if r.CreatedBy != cmp.CreatedBy {
		out = append(out, &revisions.Change{Key: "createdBy", Old: []any{cmp.CreatedBy}, New: []any{r.CreatedBy}})
	}

	if r.UpdatedBy != cmp.UpdatedBy {
		out = append(out, &revisions.Change{Key: "updatedBy", Old: []any{cmp.UpdatedBy}, New: []any{r.UpdatedBy}})
	}

	if r.DeletedBy != cmp.DeletedBy {
		out = append(out, &revisions.Change{Key: "deletedBy", Old: []any{cmp.DeletedBy}, New: []any{r.DeletedBy}})
	}

	if !reflect.DeepEqual(r.Labels, cmp.Labels) {
		out = append(out, &revisions.Change{Key: "labels", Old: []any{cmp.Labels}, New: []any{r.Labels}})
	}
	return out
}

func (r ChatbotHandoff) Clone() *ChatbotHandoff {
	dup := r
	dup.Automation = *r.Automation.Clone()

	return &dup
}

func (r ChatbotHandoff) Diff(cmp *ChatbotHandoff) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ChatbotHandoff{}
	}
	if r.Enabled != cmp.Enabled {
		out = append(out, &revisions.Change{Key: "enabled", Old: []any{cmp.Enabled}, New: []any{r.Enabled}})
	}

	for _, c := range r.Automation.Diff(&cmp.Automation) {
		c.Key = "automation." + c.Key
		out = append(out, c)
	}

	return out
}

func (r *ChatbotHandoff) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ChatbotHandoff) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ChatbotHandoffAutomation) Clone() *ChatbotHandoffAutomation {
	dup := r
	dup.OnRequested = *r.OnRequested.Clone()

	dup.OnAccepted = *r.OnAccepted.Clone()

	return &dup
}

func (r ChatbotHandoffAutomation) Diff(cmp *ChatbotHandoffAutomation) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ChatbotHandoffAutomation{}
	}
	for _, c := range r.OnRequested.Diff(&cmp.OnRequested) {
		c.Key = "onRequested." + c.Key
		out = append(out, c)
	}

	for _, c := range r.OnAccepted.Diff(&cmp.OnAccepted) {
		c.Key = "onAccepted." + c.Key
		out = append(out, c)
	}

	return out
}

func (r *ChatbotHandoffAutomation) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ChatbotHandoffAutomation) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ChatbotAutomationHook) Clone() *ChatbotAutomationHook {
	dup := r
	if r.Mappings != nil {
		dup.Mappings = make([]ChatbotHookMapping, len(r.Mappings))
		for i := range r.Mappings {
			dup.Mappings[i] = *r.Mappings[i].Clone()
		}
	}

	return &dup
}

func (r ChatbotAutomationHook) Diff(cmp *ChatbotAutomationHook) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ChatbotAutomationHook{}
	}
	if r.Automation != cmp.Automation {
		out = append(out, &revisions.Change{Key: "automation", Old: []any{cmp.Automation}, New: []any{r.Automation}})
	}

	if r.Async != cmp.Async {
		out = append(out, &revisions.Change{Key: "async", Old: []any{cmp.Async}, New: []any{r.Async}})
	}

	if !reflect.DeepEqual(r.Mappings, cmp.Mappings) {
		out = append(out, &revisions.Change{Key: "mappings", Old: []any{cmp.Mappings}, New: []any{r.Mappings}})
	}

	return out
}

func (r *ChatbotAutomationHook) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ChatbotAutomationHook) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ChatbotHookMapping) Clone() *ChatbotHookMapping {
	dup := r
	dup.StateExpression = *r.StateExpression.Clone()

	return &dup
}

func (r ChatbotHookMapping) Diff(cmp *ChatbotHookMapping) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ChatbotHookMapping{}
	}
	for _, c := range r.StateExpression.Diff(&cmp.StateExpression) {
		c.Key = "stateExpression." + c.Key
		out = append(out, c)
	}

	if r.TriggerParam != cmp.TriggerParam {
		out = append(out, &revisions.Change{Key: "triggerParam", Old: []any{cmp.TriggerParam}, New: []any{r.TriggerParam}})
	}

	return out
}

func (r *ChatbotHookMapping) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ChatbotHookMapping) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ChatbotStateExpression) Clone() *ChatbotStateExpression {
	dup := r
	return &dup
}

func (r ChatbotStateExpression) Diff(cmp *ChatbotStateExpression) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ChatbotStateExpression{}
	}
	if r.Source != cmp.Source {
		out = append(out, &revisions.Change{Key: "source", Old: []any{cmp.Source}, New: []any{r.Source}})
	}

	if r.Expr != cmp.Expr {
		out = append(out, &revisions.Change{Key: "expr", Old: []any{cmp.Expr}, New: []any{r.Expr}})
	}

	if !reflect.DeepEqual(r.Value, cmp.Value) {
		out = append(out, &revisions.Change{Key: "value", Old: []any{cmp.Value}, New: []any{r.Value}})
	}

	if r.Type != cmp.Type {
		out = append(out, &revisions.Change{Key: "type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	return out
}

func (r ChatbotScenario) Clone() *ChatbotScenario {
	dup := r
	dup.Automation = *r.Automation.Clone()

	return &dup
}

func (r ChatbotScenario) Diff(cmp *ChatbotScenario) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ChatbotScenario{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "id", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Type != cmp.Type {
		out = append(out, &revisions.Change{Key: "type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	if r.AgentID != cmp.AgentID {
		out = append(out, &revisions.Change{Key: "agentID", Old: []any{cmp.AgentID}, New: []any{r.AgentID}})
	}

	if !reflect.DeepEqual(r.Config, cmp.Config) {
		out = append(out, &revisions.Change{Key: "config", Old: []any{cmp.Config}, New: []any{r.Config}})
	}

	for _, c := range r.Automation.Diff(&cmp.Automation) {
		c.Key = "automation." + c.Key
		out = append(out, c)
	}

	return out
}

func (r *ChatbotScenario) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ChatbotScenario) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ChatbotScenarioAutomation) Clone() *ChatbotScenarioAutomation {
	dup := r
	dup.Before = *r.Before.Clone()

	dup.After = *r.After.Clone()

	return &dup
}

func (r ChatbotScenarioAutomation) Diff(cmp *ChatbotScenarioAutomation) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ChatbotScenarioAutomation{}
	}
	for _, c := range r.Before.Diff(&cmp.Before) {
		c.Key = "before." + c.Key
		out = append(out, c)
	}

	for _, c := range r.After.Diff(&cmp.After) {
		c.Key = "after." + c.Key
		out = append(out, c)
	}

	return out
}

func (r *ChatbotScenarioAutomation) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ChatbotScenarioAutomation) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ChatbotStyling) Clone() *ChatbotStyling {
	dup := r
	dup.FontSizes = *r.FontSizes.Clone()

	dup.Colors = *r.Colors.Clone()

	dup.Launcher = *r.Launcher.Clone()

	return &dup
}

func (r ChatbotStyling) Diff(cmp *ChatbotStyling) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ChatbotStyling{}
	}
	if r.LogoAttachmentID != cmp.LogoAttachmentID {
		out = append(out, &revisions.Change{Key: "logoAttachmentID", Old: []any{cmp.LogoAttachmentID}, New: []any{r.LogoAttachmentID}})
	}

	if r.LogoURL != cmp.LogoURL {
		out = append(out, &revisions.Change{Key: "logoURL", Old: []any{cmp.LogoURL}, New: []any{r.LogoURL}})
	}

	if r.FontFamily != cmp.FontFamily {
		out = append(out, &revisions.Change{Key: "fontFamily", Old: []any{cmp.FontFamily}, New: []any{r.FontFamily}})
	}

	for _, c := range r.FontSizes.Diff(&cmp.FontSizes) {
		c.Key = "fontSizes." + c.Key
		out = append(out, c)
	}

	for _, c := range r.Colors.Diff(&cmp.Colors) {
		c.Key = "colors." + c.Key
		out = append(out, c)
	}

	for _, c := range r.Launcher.Diff(&cmp.Launcher) {
		c.Key = "launcher." + c.Key
		out = append(out, c)
	}

	return out
}

func (r ChatbotFontSizes) Clone() *ChatbotFontSizes {
	dup := r
	return &dup
}

func (r ChatbotFontSizes) Diff(cmp *ChatbotFontSizes) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ChatbotFontSizes{}
	}
	if r.Base != cmp.Base {
		out = append(out, &revisions.Change{Key: "base", Old: []any{cmp.Base}, New: []any{r.Base}})
	}

	if r.Small != cmp.Small {
		out = append(out, &revisions.Change{Key: "small", Old: []any{cmp.Small}, New: []any{r.Small}})
	}

	if r.Heading != cmp.Heading {
		out = append(out, &revisions.Change{Key: "heading", Old: []any{cmp.Heading}, New: []any{r.Heading}})
	}

	return out
}

func (r *ChatbotFontSizes) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ChatbotFontSizes) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ChatbotColors) Clone() *ChatbotColors {
	dup := r
	return &dup
}

func (r ChatbotColors) Diff(cmp *ChatbotColors) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ChatbotColors{}
	}
	if r.Primary != cmp.Primary {
		out = append(out, &revisions.Change{Key: "primary", Old: []any{cmp.Primary}, New: []any{r.Primary}})
	}

	if r.PrimaryText != cmp.PrimaryText {
		out = append(out, &revisions.Change{Key: "primaryText", Old: []any{cmp.PrimaryText}, New: []any{r.PrimaryText}})
	}

	if r.Header != cmp.Header {
		out = append(out, &revisions.Change{Key: "header", Old: []any{cmp.Header}, New: []any{r.Header}})
	}

	if r.HeaderText != cmp.HeaderText {
		out = append(out, &revisions.Change{Key: "headerText", Old: []any{cmp.HeaderText}, New: []any{r.HeaderText}})
	}

	if r.Background != cmp.Background {
		out = append(out, &revisions.Change{Key: "background", Old: []any{cmp.Background}, New: []any{r.Background}})
	}

	if r.Text != cmp.Text {
		out = append(out, &revisions.Change{Key: "text", Old: []any{cmp.Text}, New: []any{r.Text}})
	}

	if r.UserBubble != cmp.UserBubble {
		out = append(out, &revisions.Change{Key: "userBubble", Old: []any{cmp.UserBubble}, New: []any{r.UserBubble}})
	}

	if r.AgentBubble != cmp.AgentBubble {
		out = append(out, &revisions.Change{Key: "agentBubble", Old: []any{cmp.AgentBubble}, New: []any{r.AgentBubble}})
	}

	return out
}

func (r *ChatbotColors) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ChatbotColors) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ChatbotLauncher) Clone() *ChatbotLauncher {
	dup := r
	return &dup
}

func (r ChatbotLauncher) Diff(cmp *ChatbotLauncher) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ChatbotLauncher{}
	}
	if r.IconURL != cmp.IconURL {
		out = append(out, &revisions.Change{Key: "iconURL", Old: []any{cmp.IconURL}, New: []any{r.IconURL}})
	}

	if r.IconAttachmentID != cmp.IconAttachmentID {
		out = append(out, &revisions.Change{Key: "iconAttachmentID", Old: []any{cmp.IconAttachmentID}, New: []any{r.IconAttachmentID}})
	}

	if r.IconVisible != cmp.IconVisible {
		out = append(out, &revisions.Change{Key: "iconVisible", Old: []any{cmp.IconVisible}, New: []any{r.IconVisible}})
	}

	if r.Label != cmp.Label {
		out = append(out, &revisions.Change{Key: "label", Old: []any{cmp.Label}, New: []any{r.Label}})
	}

	if r.ButtonLabel != cmp.ButtonLabel {
		out = append(out, &revisions.Change{Key: "buttonLabel", Old: []any{cmp.ButtonLabel}, New: []any{r.ButtonLabel}})
	}

	if r.Size != cmp.Size {
		out = append(out, &revisions.Change{Key: "size", Old: []any{cmp.Size}, New: []any{r.Size}})
	}

	if r.Shape != cmp.Shape {
		out = append(out, &revisions.Change{Key: "shape", Old: []any{cmp.Shape}, New: []any{r.Shape}})
	}

	if r.Position != cmp.Position {
		out = append(out, &revisions.Change{Key: "position", Old: []any{cmp.Position}, New: []any{r.Position}})
	}

	if r.StartOpen != cmp.StartOpen {
		out = append(out, &revisions.Change{Key: "startOpen", Old: []any{cmp.StartOpen}, New: []any{r.StartOpen}})
	}

	return out
}

func (r *ChatbotLauncher) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ChatbotLauncher) Value() (driver.Value, error) { return json.Marshal(r) }

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

func (m *ChatbotScenarios) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ChatbotScenarios) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseChatbotScenarios(ss []string) (p ChatbotScenarios, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
