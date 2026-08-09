package types

import (
	"encoding/json"

	"github.com/crusttech/human/server/pkg/filter"
)

type (
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
	aux := struct {
		*Alias
		Token      *ConnectionTemplate `json:"token"`
		HeaderName *ConnectionTemplate `json:"headerName"`
	}{Alias: (*Alias)(ct)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Catalog credential params may use a nested shape:
	//   { "token": { "value": "{{x}}", "placeholders": [...] },
	//     "headerName": { "value": "Authorization" } }
	// Flatten it so the token value and placeholders sit on the template and
	// the header name is preserved for the runtime auth layer.
	if ct.Value == "" && aux.Token != nil {
		ct.Value = aux.Token.Value
		ct.Placeholders = aux.Token.Placeholders
	}
	if aux.HeaderName != nil && ct.HeaderName == "" {
		ct.HeaderName = aux.HeaderName.Value
	}

	return nil
}

