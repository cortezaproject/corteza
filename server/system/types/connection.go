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
	Connection struct {
		ID       uint64 `json:"connectionID,string"`
		Handle   string `json:"handle"`
		Revision int    `json:"revision"`
		Status   string `json:"status"`

		Meta    ConnectionMeta    `json:"meta"`
		Service ConnectionService `json:"service"`

		Resources     ConnectionResources      `json:"resources"`
		Operations    ConnectionOperations     `json:"operations"`
		DerivedParams []ConnectionDerivedParam `json:"derivedParams,omitempty"`

		Labels map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		// Runtime-only — not persisted in DB. Populated during Search.
		Source         string `json:"source,omitempty"`         // "catalog" | "local"
		CatalogID      string `json:"catalogID,omitempty"`      // appstore connection ID
		InstalledCount int    `json:"installedCount,omitempty"` // number of ConfiguredConnections

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
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
		// Params are plain connection-level inputs (e.g. impersonation subject)
		// that are collected from the user at configuration time, stored in cc.Config.Params,
		// and available for template substitution under the "service" scope.
		Params []ConnectionPlaceholder `json:"params,omitempty"`
	}

	ConnectionProbe struct {
		Path ConnectionTemplate `json:"path"`
		// Default expected status is 200
		ExpectedStatus int `json:"expectedStatus,omitempty"`
	}

	ConnectionAuth struct {
		Method string                        `json:"method"` // none | api_token | basic | oauth2_client_credentials
		Params map[string]ConnectionTemplate `json:"params,omitempty"`
	}

	ConnectionTemplate struct {
		Value        string                  `json:"value"`
		Placeholders []ConnectionPlaceholder `json:"placeholders,omitempty"`
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

	ConnectionHTTPAction struct {
		Method       string                        `json:"method"`
		Path         ConnectionTemplate            `json:"path"`
		Headers      map[string]ConnectionTemplate `json:"headers,omitempty"`
		QueryParams  map[string]ConnectionTemplate `json:"queryParams,omitempty"`
		BodyTemplate ConnectionTemplate            `json:"bodyTemplate,omitempty"`
		ResponseMap  map[string]string             `json:"responseMap,omitempty"`
	}

	ConnectionResourceOperations struct {
		List   *ConnectionHTTPAction `json:"list,omitempty"`
		Read   *ConnectionHTTPAction `json:"read,omitempty"`
		Create *ConnectionHTTPAction `json:"create,omitempty"`
		Update *ConnectionHTTPAction `json:"update,omitempty"`
		Delete *ConnectionHTTPAction `json:"delete,omitempty"`
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

	ConnectionMimeBuildAction struct {
		// Template strings referencing input vars via {{varName}}
		To      string `json:"to"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
		From    string `json:"from,omitempty"`
		// Output is the variable name the base64url-encoded raw message is stored under
		Output string `json:"output"`
	}

	ConnectionOperationStep struct {
		Type      string                     `json:"type"` // http | mime_build
		HTTP      *ConnectionHTTPAction      `json:"http,omitempty"`
		MimeBuild *ConnectionMimeBuildAction `json:"mime_build,omitempty"`
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

	ConnectionFilter struct {
		Handle string   `json:"handle"`
		Status []string `json:"status"`
		Query  string   `json:"query"`
		Tags   []string `json:"tags"`
		Source string   `json:"source"` // "catalog", "local", or "" for all

		Deleted filter.State `json:"deleted"`

		LabeledIDs []uint64 `json:"-"`

		Check func(*Connection) (bool, error) `json:"-"`

		filter.Paging
		filter.Sorting
	}

	// Slice types for DB JSON serialization
	ConnectionResources  []ConnectionResource
	ConnectionOperations []ConnectionOperation
)

func (ct *ConnectionTemplate) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		ct.Value = s
		return nil
	}
	type Alias ConnectionTemplate
	return json.Unmarshal(data, (*Alias)(ct))
}

func ParseConnectionMeta(ss []string) (m ConnectionMeta, err error) {
	if len(ss) == 0 {
		return
	}
	err = json.Unmarshal([]byte(ss[0]), &m)
	return
}

func ParseConnectionService(ss []string) (m ConnectionService, err error) {
	if len(ss) == 0 {
		return
	}
	err = json.Unmarshal([]byte(ss[0]), &m)
	return
}

func (m *ConnectionMeta) Scan(src any) error                        { return sql.ParseJSON(src, m) }
func (m ConnectionMeta) Value() (driver.Value, error)               { return json.Marshal(m) }
func (m *ConnectionService) Scan(src any) error                     { return sql.ParseJSON(src, m) }
func (m ConnectionService) Value() (driver.Value, error)            { return json.Marshal(m) }
func (m *ConnectionResources) Scan(src any) error                   { return sql.ParseJSON(src, m) }
func (m ConnectionResources) Value() (driver.Value, error)          { return json.Marshal(m) }
func (m *ConnectionResourceOperations) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ConnectionResourceOperations) Value() (driver.Value, error) { return json.Marshal(m) }
func (m *ConnectionOperations) Scan(src any) error                  { return sql.ParseJSON(src, m) }
func (m ConnectionOperations) Value() (driver.Value, error)         { return json.Marshal(m) }
