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
	ConfiguredConnection struct {
		ID           uint64 `json:"connectionID,string"`
		ConnectionID uint64 `json:"connectionID,string"`

		Name   string `json:"name"`
		Status string `json:"status"`

		Connection Connection                 `json:"connection"`
		Config     ConfiguredConnectionConfig `json:"config"`

		Labels map[string]labelTypes.LabelValue `json:"labels,omitempty"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
	}

	ConfiguredConnectionConfig struct {
		CredentialID uint64                      `json:"credentialID,string"`
		NamespaceID  uint64                      `json:"namespaceID,string"`
		Params       []ConfiguredConnectionParam `json:"params,omitempty"`
	}

	ConfiguredConnectionParam struct {
		Scope []string `json:"scope"`
		Name  string   `json:"name"`
		Value string   `json:"value"`
	}

	ConfiguredConnectionFilter struct {
		ConnectionID uint64   `json:"connectionID,string"`
		Status       []string `json:"status"`
		Query        string   `json:"query"`

		Deleted filter.State `json:"deleted"`

		LabeledIDs []uint64 `json:"-"`

		Check func(*ConfiguredConnection) (bool, error) `json:"-"`

		filter.Paging
		filter.Sorting
	}
)

func ParseConfiguredConnectionConfig(ss []string) (m ConfiguredConnectionConfig, err error) {
	if len(ss) == 0 {
		return
	}
	err = json.Unmarshal([]byte(ss[0]), &m)
	return
}

func (m *ConfiguredConnectionConfig) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m ConfiguredConnectionConfig) Value() (driver.Value, error) { return json.Marshal(m) }
func (m *Connection) Scan(src any) error                          { return sql.ParseJSON(src, m) }
func (m Connection) Value() (driver.Value, error)                 { return json.Marshal(m) }
