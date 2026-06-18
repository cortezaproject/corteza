package types

import (
	"encoding/json"

	"github.com/crusttech/human/server/pkg/filter"
)

type (
	ConfiguredConnectionConfig struct {
		NamespaceID     uint64 `json:"namespaceID,string"`
		DalConnectionID uint64 `json:"dalConnectionID,string"`

		CredentialID uint64                      `json:"credentialID,string"`
		Params       []ConfiguredConnectionParam `json:"params,omitempty"`

		// Discovery holds cached resource lists fetched at registration time.
		// Keyed by discovery key (e.g. "spreadsheets", "tabs").
		Discovery map[string]json.RawMessage `json:"discovery,omitempty"`
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

	ConfiguredConnectionCheckResult struct {
		Connectivity ConfiguredConnectionCheckStatus  `json:"connectivity"`
		Auth         ConfiguredConnectionCheckStatus  `json:"auth"`
		Probe        *ConfiguredConnectionCheckStatus `json:"probe,omitempty"`
	}

	ConfiguredConnectionCheckStatus struct {
		OK      bool   `json:"ok"`
		Message string `json:"message,omitempty"`
	}
)
