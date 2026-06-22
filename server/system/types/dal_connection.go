package types

import (
	"github.com/crusttech/human/server/pkg/geolocation"

	"github.com/crusttech/human/server/pkg/filter"
)

type (
	// Meta ...................................................................

	DalConnectionMeta struct {
		Name       string                      `json:"name"`
		Ownership  string                      `json:"ownership"`
		Location   geolocation.Full            `json:"location"`
		Properties DalConnectionMetaProperties `json:"properties"`
	}

	DalConnectionMetaProperties struct {
		DataAtRestEncryption    DalConnectionMetaProperty `json:"dataAtRestEncryption"`
		DataAtRestProtection    DalConnectionMetaProperty `json:"dataAtRestProtection"`
		DataAtTransitEncryption DalConnectionMetaProperty `json:"dataAtTransitEncryption"`
		DataRestoration         DalConnectionMetaProperty `json:"dataRestoration"`
	}

	DalConnectionMetaProperty struct {
		Enabled bool   `json:"enabled"`
		Notes   string `json:"notes"`
	}

	// Config .................................................................

	DalConnectionConfig struct {
		// DAL configuration
		// using ptr to allow nil values (when dealing with access-controlled data)
		DAL *DalConnectionConfigDAL `json:"dal,omitempty"`

		// Privacy configuration
		Privacy DalConnectionConfigPrivacy `json:"privacy"`
	}

	DalConnectionConfigPrivacy struct {
		// Sets max-allowed data-sensitivity level for this connection
		//
		// Fields of the modules using this connection should have equal or
		// lower sensitivity level
		SensitivityLevelID uint64 `json:"sensitivityLevelID,string,omitempty"`
	}

	// DalConnectionConfigDAL a set of connection parameters
	// and model configuration
	DalConnectionConfigDAL struct {
		// type of connection
		Type string `json:"type"`

		// parameters for the connection
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
		DalConnectionID []string `json:"connectionID"`
		Handle          string   `json:"handle"`
		Type            string   `json:"type"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*DalConnection) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Paging
		filter.Sorting
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


