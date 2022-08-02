package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/geolocation"
	"github.com/cortezaproject/corteza-server/pkg/sql"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	DalConnection struct {
		ID     uint64 `json:"connectionID,string"`
		Handle string `json:"handle"`
		Type   string `json:"type"`

		Meta   ConnectionMeta   `json:"meta"`
		Config ConnectionConfig `json:"config"`

		Issues []string `json:"issues,omitempty" db:"-"`

		Labels map[string]string `json:"labels,omitempty"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		CreatedBy uint64     `json:"createdBy,string" `
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		UpdatedBy uint64     `json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
		DeletedBy uint64     `json:"deletedBy,string,omitempty" `
	}

	ConnectionConfig struct {
		Connection dal.ConnectionParams    `json:"connection"`
		Privacy    ConnectionConfigPrivacy `json:"privacy"`
		DAL        ConnectionConfigDAL     `json:"dal"`
	}

	ConnectionProperties struct {
		DataAtRestEncryption    ConnectionPropertyMeta `json:"dataAtRestEncryption"`
		DataAtRestProtection    ConnectionPropertyMeta `json:"dataAtRestProtection"`
		DataAtTransitEncryption ConnectionPropertyMeta `json:"dataAtTransitEncryption"`
		DataRestoration         ConnectionPropertyMeta `json:"dataRestoration"`
	}

	ConnectionPropertyMeta struct {
		Enabled bool   `json:"enabled"`
		Notes   string `json:"notes"`
	}

	ConnectionMeta struct {
		Location  geolocation.Full `json:"location"`
		Ownership string           `json:"ownership"`
		Name      string           `json:"name"`
	}

	ConnectionConfigPrivacy struct {
		SensitivityLevelID uint64 `json:"sensitivityLevelID,string,omitempty"`
	}

	ConnectionConfigDAL struct {
		Properties ConnectionProperties `json:"properties"`
		// @note operations, for now, will only be available on connections
		//       with a fallback on modules
		Operations dal.OperationSet `json:"operations"`

		ModelIdent     string `json:"modelIdent"`
		AttributeIdent string `json:"attributeIdent"`

		PartitionFormat         string `json:"partitionFormat"`
		PartitionIdentValidator string `json:"partitionIdentValidator"`
	}

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
