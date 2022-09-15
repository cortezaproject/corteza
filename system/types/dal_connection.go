package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/geolocation"
	"github.com/cortezaproject/corteza-server/pkg/sql"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	DalConnection struct {
		ID     uint64 `json:"connectionID,string"`
		Handle string `json:"handle"`
		Type   string `json:"type"`

		// descriptions, notes, and other user-provided meta-data
		Meta ConnectionMeta `json:"meta"`

		// collection of configurations for various subsystems that
		// use this connection and how it affects their behaviour
		Config ConnectionConfig `json:"config"`

		Issues []string `json:"issues,omitempty"`

		Labels map[string]string `json:"labels,omitempty"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty"`
	}

	// Meta ...................................................................

	ConnectionMeta struct {
		Name       string                   `json:"name"`
		Ownership  string                   `json:"ownership"`
		Location   geolocation.Full         `json:"location"`
		Properties ConnectionMetaProperties `json:"properties"`
	}

	ConnectionMetaProperties struct {
		DataAtRestEncryption    ConnectionMetaProperty `json:"dataAtRestEncryption"`
		DataAtRestProtection    ConnectionMetaProperty `json:"dataAtRestProtection"`
		DataAtTransitEncryption ConnectionMetaProperty `json:"dataAtTransitEncryption"`
		DataRestoration         ConnectionMetaProperty `json:"dataRestoration"`
	}

	ConnectionMetaProperty struct {
		Enabled bool   `json:"enabled"`
		Notes   string `json:"notes"`
	}

	// Config .................................................................

	ConnectionConfig struct {
		// DAL configuration
		// using ptr to allow nil values (when dealing with access-controlled data)
		DAL *ConnectionConfigDAL `json:"dal,omitempty"`

		// Privacy configuration
		Privacy ConnectionConfigPrivacy `json:"privacy"`
	}

	ConnectionConfigPrivacy struct {
		// Sets max-allowed data-sensitivity level for this connection
		//
		// Fields of the modules using this connection should have equal or
		// lower sensitivity level
		SensitivityLevelID uint64 `json:"sensitivityLevelID,string,omitempty"`
	}

	// ConnectionConfigDAL a set of connection parameters
	// and model configuration
	ConnectionConfigDAL struct {
		// type of connection
		Type string `json:"type"`

		// parameters for th connection
		Params map[string]any `json:"params"`

		// ident to be used when generating models from modules using this connection
		// it can use {{module}} and {{namespace}} as placeholders
		ModelIdent string `json:"modelIdent"`

		// set of regular-expression strings that will be used to match against
		// generated model identifiers
		ModelIdentCheck []string `json:"modelIdentCheck"`
	}

	// ........................................................................

	DalConnectionFilter struct {
		ConnectionID []uint64 `json:"connectionID,string"`
		Handle       string   `json:"handle"`
		Type         string   `json:"type"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*DalConnection) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Paging
	}
)

var (
	// Used to identify the primary DAL connection instead of an extra flag
	DalPrimaryConnectionResourceType = "corteza::system:primary-dal-connection"
	DalPrimaryConnectionHandle       = "primary-database"
)

func (c DalConnection) HasIssues() bool {
	return len(c.Issues) > 0
}

func ParseConnectionConfig(ss []string) (m ConnectionConfig, err error) {
	if len(ss) == 0 {
		return
	}

	err = json.Unmarshal([]byte(ss[0]), &m)
	return
}

func ParseConnectionMeta(ss []string) (m ConnectionMeta, err error) {
	if len(ss) == 0 {
		return
	}

	err = json.Unmarshal([]byte(ss[0]), &m)
	return
}

func (nm *ConnectionConfig) Scan(src any) error          { return sql.ParseJSON(src, nm) }
func (nm ConnectionConfig) Value() (driver.Value, error) { return json.Marshal(nm) }

func (nm *ConnectionMeta) Scan(src any) error          { return sql.ParseJSON(src, nm) }
func (nm ConnectionMeta) Value() (driver.Value, error) { return json.Marshal(nm) }
