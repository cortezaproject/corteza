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
	Connection struct {
		ID             uint64                           `json:"connectionID,string"`
		Handle         string                           `json:"handle"`
		Revision       int                              `json:"revision"`
		Status         string                           `json:"status"`
		Source         string                           `json:"source,omitempty"`
		Meta           ConnectionMeta                   `json:"meta"`
		Service        ConnectionService                `json:"service"`
		Resources      ConnectionResources              `json:"resources"`
		Operations     ConnectionOperations             `json:"operations"`
		DerivedParams  []ConnectionDerivedParam         `json:"derivedParams,omitempty"`
		CatalogID      string                           `json:"catalogID,omitempty"`
		InstalledCount int                              `json:"installedCount,omitempty"`
		CreatedAt      time.Time                        `json:"createdAt,omitempty"`
		UpdatedAt      *time.Time                       `json:"updatedAt,omitempty"`
		DeletedAt      *time.Time                       `json:"deletedAt,omitempty"`
		CreatedBy      uint64                           `json:"createdBy,string"`
		UpdatedBy      uint64                           `json:"updatedBy,string,omitempty"`
		DeletedBy      uint64                           `json:"deletedBy,string,omitempty"`
		Labels         map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	}

	ConnectionMeta struct {
		Short       string   `json:"short"`
		Description string   `json:"description"`
		Icon        string   `json:"icon"`
		Tags        []string `json:"tags"`
	}

	ConnectionService struct {
		BaseURL     ConnectionTemplate            `json:"baseURL"`
		Protocol    string                        `json:"protocol"`
		ContentType string                        `json:"contentType"`
		Headers     map[string]ConnectionTemplate `json:"headers,omitempty"`
		Auth        ConnectionAuth                `json:"auth"`
		Probe       *ConnectionProbe              `json:"probe,omitempty"`
		Params      []ConnectionPlaceholder       `json:"params,omitempty"`
	}

	ConnectionTemplate struct {
		Value        string                  `json:"value"`
		Placeholders []ConnectionPlaceholder `json:"placeholders,omitempty"`
	}

	ConnectionAuth struct {
		Method string                        `json:"method"`
		Params map[string]ConnectionTemplate `json:"params,omitempty"`
	}

	ConnectionProbe struct {
		Path           ConnectionTemplate `json:"path"`
		ExpectedStatus int                `json:"expectedStatus,omitempty"`
	}

	ConnectionPlaceholder struct {
		Name        string   `json:"name"`
		Label       string   `json:"label,omitempty"`
		Type        string   `json:"type"`
		Description string   `json:"description"`
		Required    bool     `json:"required"`
		Default     string   `json:"default"`
		Options     []string `json:"options,omitempty"`
	}

	ConnectionDerivedParam struct {
		Name        string   `json:"name"`
		Label       string   `json:"label"`
		Scope       []string `json:"scope"`
		Type        string   `json:"type"`
		Description string   `json:"description"`
		Required    bool     `json:"required"`
		Default     string   `json:"default"`
		Options     []string `json:"options,omitempty"`
	}

	ConnectionResource struct {
		Handle     string                       `json:"handle"`
		Meta       ConnectionResourceMeta       `json:"meta"`
		Endpoint   ConnectionTemplate           `json:"endpoint"`
		Operations ConnectionResourceOperations `json:"operations,omitempty"`
		Fields     []ConnectionResourceField    `json:"fields"`
		Webhooks   []ConnectionWebhook          `json:"webhooks,omitempty"`
	}

	ConnectionResourceMeta struct {
		Short       string          `json:"short"`
		Description string          `json:"description"`
		Icon        *ConnectionIcon `json:"icon,omitempty"`
	}

	ConnectionIcon struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}

	ConnectionResourceField struct {
		Name     string         `json:"name"`
		Type     string         `json:"type"`
		Selector []string       `json:"selector,omitempty"`
		Meta     map[string]any `json:"meta,omitempty"`
	}

	ConnectionWebhook struct {
		Event   string                   `json:"event"`
		Path    string                   `json:"path"`
		Payload []ConnectionWebhookField `json:"payload,omitempty"`
		Mapping map[string]string        `json:"mapping,omitempty"`
	}

	ConnectionWebhookField struct {
		Name     string         `json:"name"`
		Type     string         `json:"type"`
		Selector []string       `json:"selector,omitempty"`
		Meta     map[string]any `json:"meta,omitempty"`
	}

	ConnectionResourceOperations struct {
		List   *ConnectionHTTPAction `json:"list,omitempty"`
		Read   *ConnectionHTTPAction `json:"read,omitempty"`
		Create *ConnectionHTTPAction `json:"create,omitempty"`
		Update *ConnectionHTTPAction `json:"update,omitempty"`
		Delete *ConnectionHTTPAction `json:"delete,omitempty"`
	}

	ConnectionHTTPAction struct {
		Method       string                        `json:"method"`
		Path         ConnectionTemplate            `json:"path"`
		Headers      map[string]ConnectionTemplate `json:"headers,omitempty"`
		QueryParams  map[string]ConnectionTemplate `json:"queryParams,omitempty"`
		BodyTemplate ConnectionTemplate            `json:"bodyTemplate,omitempty"`
		ResponseMap  map[string]string             `json:"responseMap,omitempty"`
	}

	ConnectionOperation struct {
		Handle string                           `json:"handle"`
		Meta   ConnectionResourceMeta           `json:"meta"`
		Input  []ConnectionOperationInputField  `json:"input,omitempty"`
		Output []ConnectionOperationOutputField `json:"output,omitempty"`
		Steps  []ConnectionOperationStep        `json:"steps,omitempty"`
	}

	ConnectionOperationInputField struct {
		Name      string         `json:"name"`
		Type      string         `json:"type"`
		Required  bool           `json:"required,omitempty"`
		Aggregate bool           `json:"aggregate,omitempty"`
		Meta      map[string]any `json:"meta,omitempty"`
	}

	ConnectionOperationOutputField struct {
		Name     string         `json:"name"`
		Type     string         `json:"type"`
		Selector []string       `json:"selector,omitempty"`
		Meta     map[string]any `json:"meta,omitempty"`
	}

	ConnectionOperationStep struct {
		Type      string                     `json:"type"`
		HTTP      *ConnectionHTTPAction      `json:"http,omitempty"`
		MimeBuild *ConnectionMimeBuildAction `json:"mime_build,omitempty"`
	}

	ConnectionMimeBuildAction struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
		From    string `json:"from,omitempty"`
		Output  string `json:"output"`
	}
)

func (r Connection) Clone() *Connection {
	dup := r
	dup.Meta = *r.Meta.Clone()

	dup.Service = *r.Service.Clone()

	if r.DerivedParams != nil {
		dup.DerivedParams = make([]ConnectionDerivedParam, len(r.DerivedParams))
		for i := range r.DerivedParams {
			dup.DerivedParams[i] = *r.DerivedParams[i].Clone()
		}
	}

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

func (r Connection) Diff(cmp *Connection) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &Connection{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "connectionID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.Handle != cmp.Handle {
		out = append(out, &revisions.Change{Key: "handle", Old: []any{cmp.Handle}, New: []any{r.Handle}})
	}

	if r.Revision != cmp.Revision {
		out = append(out, &revisions.Change{Key: "revision", Old: []any{cmp.Revision}, New: []any{r.Revision}})
	}

	if r.Status != cmp.Status {
		out = append(out, &revisions.Change{Key: "status", Old: []any{cmp.Status}, New: []any{r.Status}})
	}

	if r.Source != cmp.Source {
		out = append(out, &revisions.Change{Key: "source", Old: []any{cmp.Source}, New: []any{r.Source}})
	}

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
		out = append(out, c)
	}

	for _, c := range r.Service.Diff(&cmp.Service) {
		c.Key = "service." + c.Key
		out = append(out, c)
	}

	if !reflect.DeepEqual(r.Resources, cmp.Resources) {
		out = append(out, &revisions.Change{Key: "resources", Old: []any{cmp.Resources}, New: []any{r.Resources}})
	}

	if !reflect.DeepEqual(r.Operations, cmp.Operations) {
		out = append(out, &revisions.Change{Key: "operations", Old: []any{cmp.Operations}, New: []any{r.Operations}})
	}

	if !reflect.DeepEqual(r.DerivedParams, cmp.DerivedParams) {
		out = append(out, &revisions.Change{Key: "derivedParams", Old: []any{cmp.DerivedParams}, New: []any{r.DerivedParams}})
	}

	if r.CatalogID != cmp.CatalogID {
		out = append(out, &revisions.Change{Key: "catalogID", Old: []any{cmp.CatalogID}, New: []any{r.CatalogID}})
	}

	if r.InstalledCount != cmp.InstalledCount {
		out = append(out, &revisions.Change{Key: "installedCount", Old: []any{cmp.InstalledCount}, New: []any{r.InstalledCount}})
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

func (r ConnectionMeta) Clone() *ConnectionMeta {
	dup := r
	if r.Tags != nil {
		dup.Tags = make([]string, len(r.Tags))
		copy(dup.Tags, r.Tags)
	}

	return &dup
}

func (r ConnectionMeta) Diff(cmp *ConnectionMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionMeta{}
	}
	if r.Short != cmp.Short {
		out = append(out, &revisions.Change{Key: "short", Old: []any{cmp.Short}, New: []any{r.Short}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if r.Icon != cmp.Icon {
		out = append(out, &revisions.Change{Key: "icon", Old: []any{cmp.Icon}, New: []any{r.Icon}})
	}

	if !reflect.DeepEqual(r.Tags, cmp.Tags) {
		out = append(out, &revisions.Change{Key: "tags", Old: []any{cmp.Tags}, New: []any{r.Tags}})
	}

	return out
}

func (r *ConnectionMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConnectionMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ConnectionService) Clone() *ConnectionService {
	dup := r
	dup.BaseURL = *r.BaseURL.Clone()

	if r.Headers != nil {
		dup.Headers = make(map[string]ConnectionTemplate, len(r.Headers))
		for k, v := range r.Headers {
			dup.Headers[k] = v
		}
	}

	dup.Auth = *r.Auth.Clone()

	if r.Probe != nil {
		dup.Probe = r.Probe.Clone()
	}

	if r.Params != nil {
		dup.Params = make([]ConnectionPlaceholder, len(r.Params))
		for i := range r.Params {
			dup.Params[i] = *r.Params[i].Clone()
		}
	}

	return &dup
}

func (r ConnectionService) Diff(cmp *ConnectionService) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionService{}
	}
	for _, c := range r.BaseURL.Diff(&cmp.BaseURL) {
		c.Key = "baseURL." + c.Key
		out = append(out, c)
	}

	if r.Protocol != cmp.Protocol {
		out = append(out, &revisions.Change{Key: "protocol", Old: []any{cmp.Protocol}, New: []any{r.Protocol}})
	}

	if r.ContentType != cmp.ContentType {
		out = append(out, &revisions.Change{Key: "contentType", Old: []any{cmp.ContentType}, New: []any{r.ContentType}})
	}

	if !reflect.DeepEqual(r.Headers, cmp.Headers) {
		out = append(out, &revisions.Change{Key: "headers", Old: []any{cmp.Headers}, New: []any{r.Headers}})
	}

	for _, c := range r.Auth.Diff(&cmp.Auth) {
		c.Key = "auth." + c.Key
		out = append(out, c)
	}

	if (r.Probe == nil) != (cmp.Probe == nil) {
		out = append(out, &revisions.Change{Key: "probe", Old: []any{cmp.Probe}, New: []any{r.Probe}})
	} else if r.Probe != nil {
		for _, c := range r.Probe.Diff(cmp.Probe) {
			c.Key = "probe." + c.Key
			out = append(out, c)
		}
	}

	if !reflect.DeepEqual(r.Params, cmp.Params) {
		out = append(out, &revisions.Change{Key: "params", Old: []any{cmp.Params}, New: []any{r.Params}})
	}

	return out
}

func (r *ConnectionService) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConnectionService) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ConnectionTemplate) Clone() *ConnectionTemplate {
	dup := r
	if r.Placeholders != nil {
		dup.Placeholders = make([]ConnectionPlaceholder, len(r.Placeholders))
		for i := range r.Placeholders {
			dup.Placeholders[i] = *r.Placeholders[i].Clone()
		}
	}

	return &dup
}

func (r ConnectionTemplate) Diff(cmp *ConnectionTemplate) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionTemplate{}
	}
	if r.Value != cmp.Value {
		out = append(out, &revisions.Change{Key: "value", Old: []any{cmp.Value}, New: []any{r.Value}})
	}

	if !reflect.DeepEqual(r.Placeholders, cmp.Placeholders) {
		out = append(out, &revisions.Change{Key: "placeholders", Old: []any{cmp.Placeholders}, New: []any{r.Placeholders}})
	}

	return out
}

func (r ConnectionAuth) Clone() *ConnectionAuth {
	dup := r
	if r.Params != nil {
		dup.Params = make(map[string]ConnectionTemplate, len(r.Params))
		for k, v := range r.Params {
			dup.Params[k] = v
		}
	}

	return &dup
}

func (r ConnectionAuth) Diff(cmp *ConnectionAuth) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionAuth{}
	}
	if r.Method != cmp.Method {
		out = append(out, &revisions.Change{Key: "method", Old: []any{cmp.Method}, New: []any{r.Method}})
	}

	if !reflect.DeepEqual(r.Params, cmp.Params) {
		out = append(out, &revisions.Change{Key: "params", Old: []any{cmp.Params}, New: []any{r.Params}})
	}

	return out
}

func (r *ConnectionAuth) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConnectionAuth) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ConnectionProbe) Clone() *ConnectionProbe {
	dup := r
	dup.Path = *r.Path.Clone()

	return &dup
}

func (r ConnectionProbe) Diff(cmp *ConnectionProbe) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionProbe{}
	}
	for _, c := range r.Path.Diff(&cmp.Path) {
		c.Key = "path." + c.Key
		out = append(out, c)
	}

	if r.ExpectedStatus != cmp.ExpectedStatus {
		out = append(out, &revisions.Change{Key: "expectedStatus", Old: []any{cmp.ExpectedStatus}, New: []any{r.ExpectedStatus}})
	}

	return out
}

func (r *ConnectionProbe) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConnectionProbe) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ConnectionPlaceholder) Clone() *ConnectionPlaceholder {
	dup := r
	if r.Options != nil {
		dup.Options = make([]string, len(r.Options))
		copy(dup.Options, r.Options)
	}

	return &dup
}

func (r ConnectionPlaceholder) Diff(cmp *ConnectionPlaceholder) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionPlaceholder{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Label != cmp.Label {
		out = append(out, &revisions.Change{Key: "label", Old: []any{cmp.Label}, New: []any{r.Label}})
	}

	if r.Type != cmp.Type {
		out = append(out, &revisions.Change{Key: "type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if r.Required != cmp.Required {
		out = append(out, &revisions.Change{Key: "required", Old: []any{cmp.Required}, New: []any{r.Required}})
	}

	if r.Default != cmp.Default {
		out = append(out, &revisions.Change{Key: "default", Old: []any{cmp.Default}, New: []any{r.Default}})
	}

	if !reflect.DeepEqual(r.Options, cmp.Options) {
		out = append(out, &revisions.Change{Key: "options", Old: []any{cmp.Options}, New: []any{r.Options}})
	}

	return out
}

func (r *ConnectionPlaceholder) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConnectionPlaceholder) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ConnectionDerivedParam) Clone() *ConnectionDerivedParam {
	dup := r
	if r.Scope != nil {
		dup.Scope = make([]string, len(r.Scope))
		copy(dup.Scope, r.Scope)
	}

	if r.Options != nil {
		dup.Options = make([]string, len(r.Options))
		copy(dup.Options, r.Options)
	}

	return &dup
}

func (r ConnectionDerivedParam) Diff(cmp *ConnectionDerivedParam) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionDerivedParam{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Label != cmp.Label {
		out = append(out, &revisions.Change{Key: "label", Old: []any{cmp.Label}, New: []any{r.Label}})
	}

	if !reflect.DeepEqual(r.Scope, cmp.Scope) {
		out = append(out, &revisions.Change{Key: "scope", Old: []any{cmp.Scope}, New: []any{r.Scope}})
	}

	if r.Type != cmp.Type {
		out = append(out, &revisions.Change{Key: "type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if r.Required != cmp.Required {
		out = append(out, &revisions.Change{Key: "required", Old: []any{cmp.Required}, New: []any{r.Required}})
	}

	if r.Default != cmp.Default {
		out = append(out, &revisions.Change{Key: "default", Old: []any{cmp.Default}, New: []any{r.Default}})
	}

	if !reflect.DeepEqual(r.Options, cmp.Options) {
		out = append(out, &revisions.Change{Key: "options", Old: []any{cmp.Options}, New: []any{r.Options}})
	}

	return out
}

func (r *ConnectionDerivedParam) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConnectionDerivedParam) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ConnectionResource) Clone() *ConnectionResource {
	dup := r
	dup.Meta = *r.Meta.Clone()

	dup.Endpoint = *r.Endpoint.Clone()

	dup.Operations = *r.Operations.Clone()

	if r.Fields != nil {
		dup.Fields = make([]ConnectionResourceField, len(r.Fields))
		for i := range r.Fields {
			dup.Fields[i] = *r.Fields[i].Clone()
		}
	}

	if r.Webhooks != nil {
		dup.Webhooks = make([]ConnectionWebhook, len(r.Webhooks))
		for i := range r.Webhooks {
			dup.Webhooks[i] = *r.Webhooks[i].Clone()
		}
	}

	return &dup
}

func (r ConnectionResource) Diff(cmp *ConnectionResource) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionResource{}
	}
	if r.Handle != cmp.Handle {
		out = append(out, &revisions.Change{Key: "handle", Old: []any{cmp.Handle}, New: []any{r.Handle}})
	}

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
		out = append(out, c)
	}

	for _, c := range r.Endpoint.Diff(&cmp.Endpoint) {
		c.Key = "endpoint." + c.Key
		out = append(out, c)
	}

	for _, c := range r.Operations.Diff(&cmp.Operations) {
		c.Key = "operations." + c.Key
		out = append(out, c)
	}

	if !reflect.DeepEqual(r.Fields, cmp.Fields) {
		out = append(out, &revisions.Change{Key: "fields", Old: []any{cmp.Fields}, New: []any{r.Fields}})
	}

	if !reflect.DeepEqual(r.Webhooks, cmp.Webhooks) {
		out = append(out, &revisions.Change{Key: "webhooks", Old: []any{cmp.Webhooks}, New: []any{r.Webhooks}})
	}

	return out
}

func (r *ConnectionResource) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConnectionResource) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ConnectionResourceMeta) Clone() *ConnectionResourceMeta {
	dup := r
	if r.Icon != nil {
		dup.Icon = r.Icon.Clone()
	}

	return &dup
}

func (r ConnectionResourceMeta) Diff(cmp *ConnectionResourceMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionResourceMeta{}
	}
	if r.Short != cmp.Short {
		out = append(out, &revisions.Change{Key: "short", Old: []any{cmp.Short}, New: []any{r.Short}})
	}

	if r.Description != cmp.Description {
		out = append(out, &revisions.Change{Key: "description", Old: []any{cmp.Description}, New: []any{r.Description}})
	}

	if (r.Icon == nil) != (cmp.Icon == nil) {
		out = append(out, &revisions.Change{Key: "icon", Old: []any{cmp.Icon}, New: []any{r.Icon}})
	} else if r.Icon != nil {
		for _, c := range r.Icon.Diff(cmp.Icon) {
			c.Key = "icon." + c.Key
			out = append(out, c)
		}
	}

	return out
}

func (r *ConnectionResourceMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConnectionResourceMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ConnectionIcon) Clone() *ConnectionIcon {
	dup := r
	return &dup
}

func (r ConnectionIcon) Diff(cmp *ConnectionIcon) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionIcon{}
	}
	if r.Type != cmp.Type {
		out = append(out, &revisions.Change{Key: "type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	if r.Value != cmp.Value {
		out = append(out, &revisions.Change{Key: "value", Old: []any{cmp.Value}, New: []any{r.Value}})
	}

	return out
}

func (r ConnectionResourceField) Clone() *ConnectionResourceField {
	dup := r
	if r.Selector != nil {
		dup.Selector = make([]string, len(r.Selector))
		copy(dup.Selector, r.Selector)
	}

	if r.Meta != nil {
		dup.Meta = make(map[string]any, len(r.Meta))
		for k, v := range r.Meta {
			dup.Meta[k] = v
		}
	}

	return &dup
}

func (r ConnectionResourceField) Diff(cmp *ConnectionResourceField) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionResourceField{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Type != cmp.Type {
		out = append(out, &revisions.Change{Key: "type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	if !reflect.DeepEqual(r.Selector, cmp.Selector) {
		out = append(out, &revisions.Change{Key: "selector", Old: []any{cmp.Selector}, New: []any{r.Selector}})
	}

	if !reflect.DeepEqual(r.Meta, cmp.Meta) {
		out = append(out, &revisions.Change{Key: "meta", Old: []any{cmp.Meta}, New: []any{r.Meta}})
	}

	return out
}

func (r *ConnectionResourceField) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConnectionResourceField) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ConnectionWebhook) Clone() *ConnectionWebhook {
	dup := r
	if r.Payload != nil {
		dup.Payload = make([]ConnectionWebhookField, len(r.Payload))
		for i := range r.Payload {
			dup.Payload[i] = *r.Payload[i].Clone()
		}
	}

	if r.Mapping != nil {
		dup.Mapping = make(map[string]string, len(r.Mapping))
		for k, v := range r.Mapping {
			dup.Mapping[k] = v
		}
	}

	return &dup
}

func (r ConnectionWebhook) Diff(cmp *ConnectionWebhook) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionWebhook{}
	}
	if r.Event != cmp.Event {
		out = append(out, &revisions.Change{Key: "event", Old: []any{cmp.Event}, New: []any{r.Event}})
	}

	if r.Path != cmp.Path {
		out = append(out, &revisions.Change{Key: "path", Old: []any{cmp.Path}, New: []any{r.Path}})
	}

	if !reflect.DeepEqual(r.Payload, cmp.Payload) {
		out = append(out, &revisions.Change{Key: "payload", Old: []any{cmp.Payload}, New: []any{r.Payload}})
	}

	if !reflect.DeepEqual(r.Mapping, cmp.Mapping) {
		out = append(out, &revisions.Change{Key: "mapping", Old: []any{cmp.Mapping}, New: []any{r.Mapping}})
	}

	return out
}

func (r *ConnectionWebhook) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConnectionWebhook) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ConnectionWebhookField) Clone() *ConnectionWebhookField {
	dup := r
	if r.Selector != nil {
		dup.Selector = make([]string, len(r.Selector))
		copy(dup.Selector, r.Selector)
	}

	if r.Meta != nil {
		dup.Meta = make(map[string]any, len(r.Meta))
		for k, v := range r.Meta {
			dup.Meta[k] = v
		}
	}

	return &dup
}

func (r ConnectionWebhookField) Diff(cmp *ConnectionWebhookField) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionWebhookField{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Type != cmp.Type {
		out = append(out, &revisions.Change{Key: "type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	if !reflect.DeepEqual(r.Selector, cmp.Selector) {
		out = append(out, &revisions.Change{Key: "selector", Old: []any{cmp.Selector}, New: []any{r.Selector}})
	}

	if !reflect.DeepEqual(r.Meta, cmp.Meta) {
		out = append(out, &revisions.Change{Key: "meta", Old: []any{cmp.Meta}, New: []any{r.Meta}})
	}

	return out
}

func (r *ConnectionWebhookField) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConnectionWebhookField) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ConnectionResourceOperations) Clone() *ConnectionResourceOperations {
	dup := r
	if r.List != nil {
		dup.List = r.List.Clone()
	}

	if r.Read != nil {
		dup.Read = r.Read.Clone()
	}

	if r.Create != nil {
		dup.Create = r.Create.Clone()
	}

	if r.Update != nil {
		dup.Update = r.Update.Clone()
	}

	if r.Delete != nil {
		dup.Delete = r.Delete.Clone()
	}

	return &dup
}

func (r ConnectionResourceOperations) Diff(cmp *ConnectionResourceOperations) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionResourceOperations{}
	}
	if (r.List == nil) != (cmp.List == nil) {
		out = append(out, &revisions.Change{Key: "list", Old: []any{cmp.List}, New: []any{r.List}})
	} else if r.List != nil {
		for _, c := range r.List.Diff(cmp.List) {
			c.Key = "list." + c.Key
			out = append(out, c)
		}
	}

	if (r.Read == nil) != (cmp.Read == nil) {
		out = append(out, &revisions.Change{Key: "read", Old: []any{cmp.Read}, New: []any{r.Read}})
	} else if r.Read != nil {
		for _, c := range r.Read.Diff(cmp.Read) {
			c.Key = "read." + c.Key
			out = append(out, c)
		}
	}

	if (r.Create == nil) != (cmp.Create == nil) {
		out = append(out, &revisions.Change{Key: "create", Old: []any{cmp.Create}, New: []any{r.Create}})
	} else if r.Create != nil {
		for _, c := range r.Create.Diff(cmp.Create) {
			c.Key = "create." + c.Key
			out = append(out, c)
		}
	}

	if (r.Update == nil) != (cmp.Update == nil) {
		out = append(out, &revisions.Change{Key: "update", Old: []any{cmp.Update}, New: []any{r.Update}})
	} else if r.Update != nil {
		for _, c := range r.Update.Diff(cmp.Update) {
			c.Key = "update." + c.Key
			out = append(out, c)
		}
	}

	if (r.Delete == nil) != (cmp.Delete == nil) {
		out = append(out, &revisions.Change{Key: "delete", Old: []any{cmp.Delete}, New: []any{r.Delete}})
	} else if r.Delete != nil {
		for _, c := range r.Delete.Diff(cmp.Delete) {
			c.Key = "delete." + c.Key
			out = append(out, c)
		}
	}

	return out
}

func (r *ConnectionResourceOperations) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConnectionResourceOperations) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ConnectionHTTPAction) Clone() *ConnectionHTTPAction {
	dup := r
	dup.Path = *r.Path.Clone()

	if r.Headers != nil {
		dup.Headers = make(map[string]ConnectionTemplate, len(r.Headers))
		for k, v := range r.Headers {
			dup.Headers[k] = v
		}
	}

	if r.QueryParams != nil {
		dup.QueryParams = make(map[string]ConnectionTemplate, len(r.QueryParams))
		for k, v := range r.QueryParams {
			dup.QueryParams[k] = v
		}
	}

	dup.BodyTemplate = *r.BodyTemplate.Clone()

	if r.ResponseMap != nil {
		dup.ResponseMap = make(map[string]string, len(r.ResponseMap))
		for k, v := range r.ResponseMap {
			dup.ResponseMap[k] = v
		}
	}

	return &dup
}

func (r ConnectionHTTPAction) Diff(cmp *ConnectionHTTPAction) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionHTTPAction{}
	}
	if r.Method != cmp.Method {
		out = append(out, &revisions.Change{Key: "method", Old: []any{cmp.Method}, New: []any{r.Method}})
	}

	for _, c := range r.Path.Diff(&cmp.Path) {
		c.Key = "path." + c.Key
		out = append(out, c)
	}

	if !reflect.DeepEqual(r.Headers, cmp.Headers) {
		out = append(out, &revisions.Change{Key: "headers", Old: []any{cmp.Headers}, New: []any{r.Headers}})
	}

	if !reflect.DeepEqual(r.QueryParams, cmp.QueryParams) {
		out = append(out, &revisions.Change{Key: "queryParams", Old: []any{cmp.QueryParams}, New: []any{r.QueryParams}})
	}

	for _, c := range r.BodyTemplate.Diff(&cmp.BodyTemplate) {
		c.Key = "bodyTemplate." + c.Key
		out = append(out, c)
	}

	if !reflect.DeepEqual(r.ResponseMap, cmp.ResponseMap) {
		out = append(out, &revisions.Change{Key: "responseMap", Old: []any{cmp.ResponseMap}, New: []any{r.ResponseMap}})
	}

	return out
}

func (r *ConnectionHTTPAction) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConnectionHTTPAction) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ConnectionOperation) Clone() *ConnectionOperation {
	dup := r
	dup.Meta = *r.Meta.Clone()

	if r.Input != nil {
		dup.Input = make([]ConnectionOperationInputField, len(r.Input))
		for i := range r.Input {
			dup.Input[i] = *r.Input[i].Clone()
		}
	}

	if r.Output != nil {
		dup.Output = make([]ConnectionOperationOutputField, len(r.Output))
		for i := range r.Output {
			dup.Output[i] = *r.Output[i].Clone()
		}
	}

	if r.Steps != nil {
		dup.Steps = make([]ConnectionOperationStep, len(r.Steps))
		for i := range r.Steps {
			dup.Steps[i] = *r.Steps[i].Clone()
		}
	}

	return &dup
}

func (r ConnectionOperation) Diff(cmp *ConnectionOperation) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionOperation{}
	}
	if r.Handle != cmp.Handle {
		out = append(out, &revisions.Change{Key: "handle", Old: []any{cmp.Handle}, New: []any{r.Handle}})
	}

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
		out = append(out, c)
	}

	if !reflect.DeepEqual(r.Input, cmp.Input) {
		out = append(out, &revisions.Change{Key: "input", Old: []any{cmp.Input}, New: []any{r.Input}})
	}

	if !reflect.DeepEqual(r.Output, cmp.Output) {
		out = append(out, &revisions.Change{Key: "output", Old: []any{cmp.Output}, New: []any{r.Output}})
	}

	if !reflect.DeepEqual(r.Steps, cmp.Steps) {
		out = append(out, &revisions.Change{Key: "steps", Old: []any{cmp.Steps}, New: []any{r.Steps}})
	}

	return out
}

func (r *ConnectionOperation) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConnectionOperation) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ConnectionOperationInputField) Clone() *ConnectionOperationInputField {
	dup := r
	if r.Meta != nil {
		dup.Meta = make(map[string]any, len(r.Meta))
		for k, v := range r.Meta {
			dup.Meta[k] = v
		}
	}

	return &dup
}

func (r ConnectionOperationInputField) Diff(cmp *ConnectionOperationInputField) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionOperationInputField{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Type != cmp.Type {
		out = append(out, &revisions.Change{Key: "type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	if r.Required != cmp.Required {
		out = append(out, &revisions.Change{Key: "required", Old: []any{cmp.Required}, New: []any{r.Required}})
	}

	if r.Aggregate != cmp.Aggregate {
		out = append(out, &revisions.Change{Key: "aggregate", Old: []any{cmp.Aggregate}, New: []any{r.Aggregate}})
	}

	if !reflect.DeepEqual(r.Meta, cmp.Meta) {
		out = append(out, &revisions.Change{Key: "meta", Old: []any{cmp.Meta}, New: []any{r.Meta}})
	}

	return out
}

func (r *ConnectionOperationInputField) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConnectionOperationInputField) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ConnectionOperationOutputField) Clone() *ConnectionOperationOutputField {
	dup := r
	if r.Selector != nil {
		dup.Selector = make([]string, len(r.Selector))
		copy(dup.Selector, r.Selector)
	}

	if r.Meta != nil {
		dup.Meta = make(map[string]any, len(r.Meta))
		for k, v := range r.Meta {
			dup.Meta[k] = v
		}
	}

	return &dup
}

func (r ConnectionOperationOutputField) Diff(cmp *ConnectionOperationOutputField) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionOperationOutputField{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Type != cmp.Type {
		out = append(out, &revisions.Change{Key: "type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	if !reflect.DeepEqual(r.Selector, cmp.Selector) {
		out = append(out, &revisions.Change{Key: "selector", Old: []any{cmp.Selector}, New: []any{r.Selector}})
	}

	if !reflect.DeepEqual(r.Meta, cmp.Meta) {
		out = append(out, &revisions.Change{Key: "meta", Old: []any{cmp.Meta}, New: []any{r.Meta}})
	}

	return out
}

func (r *ConnectionOperationOutputField) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConnectionOperationOutputField) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ConnectionOperationStep) Clone() *ConnectionOperationStep {
	dup := r
	if r.HTTP != nil {
		dup.HTTP = r.HTTP.Clone()
	}

	if r.MimeBuild != nil {
		dup.MimeBuild = r.MimeBuild.Clone()
	}

	return &dup
}

func (r ConnectionOperationStep) Diff(cmp *ConnectionOperationStep) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionOperationStep{}
	}
	if r.Type != cmp.Type {
		out = append(out, &revisions.Change{Key: "type", Old: []any{cmp.Type}, New: []any{r.Type}})
	}

	if (r.HTTP == nil) != (cmp.HTTP == nil) {
		out = append(out, &revisions.Change{Key: "http", Old: []any{cmp.HTTP}, New: []any{r.HTTP}})
	} else if r.HTTP != nil {
		for _, c := range r.HTTP.Diff(cmp.HTTP) {
			c.Key = "http." + c.Key
			out = append(out, c)
		}
	}

	if (r.MimeBuild == nil) != (cmp.MimeBuild == nil) {
		out = append(out, &revisions.Change{Key: "mime_build", Old: []any{cmp.MimeBuild}, New: []any{r.MimeBuild}})
	} else if r.MimeBuild != nil {
		for _, c := range r.MimeBuild.Diff(cmp.MimeBuild) {
			c.Key = "mime_build." + c.Key
			out = append(out, c)
		}
	}

	return out
}

func (r *ConnectionOperationStep) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConnectionOperationStep) Value() (driver.Value, error) { return json.Marshal(r) }

func (r ConnectionMimeBuildAction) Clone() *ConnectionMimeBuildAction {
	dup := r
	return &dup
}

func (r ConnectionMimeBuildAction) Diff(cmp *ConnectionMimeBuildAction) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &ConnectionMimeBuildAction{}
	}
	if r.To != cmp.To {
		out = append(out, &revisions.Change{Key: "to", Old: []any{cmp.To}, New: []any{r.To}})
	}

	if r.Subject != cmp.Subject {
		out = append(out, &revisions.Change{Key: "subject", Old: []any{cmp.Subject}, New: []any{r.Subject}})
	}

	if r.Body != cmp.Body {
		out = append(out, &revisions.Change{Key: "body", Old: []any{cmp.Body}, New: []any{r.Body}})
	}

	if r.From != cmp.From {
		out = append(out, &revisions.Change{Key: "from", Old: []any{cmp.From}, New: []any{r.From}})
	}

	if r.Output != cmp.Output {
		out = append(out, &revisions.Change{Key: "output", Old: []any{cmp.Output}, New: []any{r.Output}})
	}

	return out
}

func (r *ConnectionMimeBuildAction) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r ConnectionMimeBuildAction) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseConnectionMeta(ss []string) (p ConnectionMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func ParseConnectionService(ss []string) (p ConnectionService, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *ConnectionResources) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ConnectionResources) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseConnectionResources(ss []string) (p ConnectionResources, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}

func (m *ConnectionOperations) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ConnectionOperations) Value() (driver.Value, error) { return json.Marshal(m) }

func ParseConnectionOperations(ss []string) (p ConnectionOperations, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
