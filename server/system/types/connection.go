package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	labelTypes "github.com/cortezaproject/corteza/server/pkg/label/types"
	"github.com/cortezaproject/corteza/server/pkg/sql"
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
		Short       string `json:"short"`
		Description string `json:"description"`
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
		Handle string                     `json:"handle"`
		Meta   ConnectionResourceMeta     `json:"meta"`
		HTTP   ConnectionHTTPAction       `json:"http"`
		Input  []ConnectionOperationField `json:"input,omitempty"`
		Output []ConnectionOperationField `json:"output,omitempty"`
	}

	ConnectionOperationField struct {
		Name     string         `json:"name"`
		Type     string         `json:"type"`
		Required bool           `json:"required,omitempty"`
		Selector []string       `json:"selector,omitempty"`
		Meta     map[string]any `json:"meta,omitempty"`
	}

	ConnectionDerivedParam struct {
		Name        string   `json:"name"`
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
